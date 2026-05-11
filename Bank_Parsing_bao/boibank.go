package bankparsing

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type BOIBANKParser struct{}

func (p *BOIBANKParser) Parse(_ string, filePath string, content string) ([]Transaction, error) {
	ext := strings.ToLower(filepath.Ext(filePath))
	if content != "" {
		if p.isFixedFormat(content) {
			return p.parseFixedFormatContent(content)
		}
		if p.isCSVFormat(content) {
			return p.parseCSVContent(content)
		}
	}
	switch ext {
	case ".txt":
		return p.parseTXTFile(filePath)
	case ".csv":
		return p.parseCSVFile(filePath)
	default:
		return nil, fmt.Errorf("unsupported file format for BOI. Expected: .txt or .csv, but got: %s", ext)
	}
}

func (p *BOIBANKParser) parseTXTFile(filePath string) ([]Transaction, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("读取BOI文件失败: %v", err)
	}

	return p.parseBOIContent(string(data))
}

func (p *BOIBANKParser) parseCSVFile(filePath string) ([]Transaction, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("读取BOI CSV文件失败: %v", err)
	}

	return p.parseCSVContent(string(data))
}

func (p *BOIBANKParser) parseBOIContent(content string) ([]Transaction, error) {
	if p.isFixedFormat(content) {
		return p.parseFixedFormatContent(content)
	}
	if p.isCSVFormat(content) {
		return p.parseCSVContent(content)
	}
	return p.parseTextContent(content)
}

func (p *BOIBANKParser) parseFixedFormatContent(content string) ([]Transaction, error) {
	lines := strings.Split(content, "\n")
	var transactions []Transaction

	i := 0
	recordCount := 0
	for i < len(lines) {
		line := strings.TrimSpace(lines[i])
		if strings.Contains(line, "=== 第") && strings.Contains(line, "次获取 ===") {
			i++
			break
		}
		if strings.Contains(line, "Account Holder") ||
			strings.Contains(line, "CustomerId") ||
			strings.Contains(line, "IFSC") ||
			strings.Contains(line, "Account Number") {
			i++
			continue
		}
		i++
	}

	for i < len(lines) {
		for i < len(lines) && strings.TrimSpace(lines[i]) == "" {
			i++
		}
		if i >= len(lines) {
			break
		}

		recordCount++
		dateLine := strings.TrimSpace(lines[i])
		date, err := p.parseFixedFormatDate(dateLine)
		if err != nil {
			i++
			continue
		}
		i++
		for i < len(lines) && strings.TrimSpace(lines[i]) == "" {
			i++
		}
		if i >= len(lines) {
			break
		}
		descLine := strings.TrimSpace(lines[i])
		i++
		for i < len(lines) && strings.TrimSpace(lines[i]) == "" {
			i++
		}
		if i >= len(lines) {
			break
		}
		amountLine := strings.TrimSpace(lines[i])
		amount := p.cleanAmountString(amountLine)
		i++
		for i < len(lines) && strings.TrimSpace(lines[i]) == "" {
			i++
		}
		if i >= len(lines) {
			transaction := p.createFixedFormatTransaction(date, descLine, amount, "0")
			if transaction != nil {
				transactions = append(transactions, *transaction)
			}
			break
		}
		balanceLine := strings.TrimSpace(lines[i])
		balance := p.cleanAmountString(balanceLine)
		i++
		transaction := p.createFixedFormatTransaction(date, descLine, amount, balance)
		if transaction != nil {
			transactions = append(transactions, *transaction)
		}
	}

	if len(transactions) == 0 {
		return nil, fmt.Errorf("未解析到BOI固定格式交易数据")
	}

	return transactions, nil
}

func (p *BOIBANKParser) createFixedFormatTransaction(date, descLine, amount, balance string) *Transaction {
	transType := p.determineTransTypeFromDesc(descLine)
	amountFloat, err := strconv.ParseFloat(amount, 64)
	if err != nil {
		amountFloat = 0
	}

	bankTxnId := p.extractTransactionIdFromDesc(descLine)
	if bankTxnId == "" {
		bankTxnId = p.generateTransactionID(descLine, date, amount, "")
	}
	formattedDate := p.formatTransactionDate(date)

	return &Transaction{
		Date:        formattedDate,
		Description: descLine,
		Amount:      amountFloat,
		Type:        transType,
		BankTxnId:   bankTxnId,
		Balance:     balance,
	}
}

func (p *BOIBANKParser) extractTransactionIdFromDesc(desc string) string {
	desc = strings.TrimSpace(desc)

	if desc == "" {
		return ""
	}
	re12digit := regexp.MustCompile(`\b\d{12}\b`)
	if match := re12digit.FindString(desc); match != "" {
		return match
	}

	patterns := []string{
		`UPI/(\d{12})/`,
		`IMPS/(\d{12})`,
		`NEFT.*?(\d{12})`,
		`RTGS.*?(\d{12})`,
		`TXN.*?(\d{12})`,
		`REF.*?(\d{12})`,
		`ID.*?(\d{12})`,
		`/(\d{12})/`,
	}

	for _, pattern := range patterns {
		re := regexp.MustCompile(pattern)
		if matches := re.FindStringSubmatch(desc); len(matches) > 1 {
			id := matches[1]
			if len(id) == 12 && p.isPureDigits(id) {
				return id
			}
		}
	}

	return ""
}

func (p *BOIBANKParser) isPureDigits(str string) bool {
	if str == "" {
		return false
	}
	for _, ch := range str {
		if ch < '0' || ch > '9' {
			return false
		}
	}
	return true
}

func (p *BOIBANKParser) determineTransTypeFromDesc(desc string) string {
	descUpper := strings.ToUpper(desc)
	if strings.Contains(descUpper, "/CR/") ||
		strings.Contains(descUpper, " CR/") ||
		strings.Contains(descUpper, "BY CASH") ||
		strings.Contains(descUpper, "UPI/") && strings.Contains(descUpper, "/CR/") {
		return "TYPE_IN"
	}
	if strings.Contains(descUpper, "/DR/") ||
		strings.Contains(descUpper, " DR/") ||
		strings.Contains(descUpper, "CHARGES") ||
		strings.Contains(descUpper, "CWDR") ||
		strings.Contains(descUpper, "IMPSUAIB") ||
		strings.Contains(descUpper, "MIBAL") ||
		strings.Contains(descUpper, "NON-MAINT") ||
		strings.Contains(descUpper, "CARD ISSUANCE") {
		return "TYPE_OUT"
	}

	return ""
}

func (p *BOIBANKParser) parseFixedFormatDate(dateStr string) (string, error) {
	dateStr = strings.TrimSpace(dateStr)

	layouts := []string{
		"2 Jan 2006",
		"02 Jan 2006",
		"2-Jan-2006",
		"02-Jan-2006",
		"2 Jan 06",
		"02 Jan 06",
		"Jan 2, 2006",
		"January 2, 2006",
		"2006-01-02",
		"01/02/2006",
		"02/01/2006",
	}

	for _, layout := range layouts {
		if t, err := time.Parse(layout, dateStr); err == nil {
			return t.Format("2006-01-02"), nil
		}
	}
	if strings.Contains(dateStr, "UPI/") ||
		strings.Contains(dateStr, "CWDR//") ||
		strings.Contains(dateStr, "Non-Maint") ||
		strings.Contains(dateStr, "IMPSUAIB") ||
		strings.Contains(dateStr, "MIBAL") ||
		strings.Contains(dateStr, "₹") {
		return "", fmt.Errorf("不是日期格式: %s", dateStr)
	}

	return "", fmt.Errorf("无法解析日期: %s", dateStr)
}

func (p *BOIBANKParser) isFixedFormat(content string) bool {
	lines := strings.Split(content, "\n")
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.Contains(trimmed, "=== 第") && strings.Contains(trimmed, "次获取 ===") {
			return true
		}
	}
	groupCount := 0
	for i := 0; i < len(lines); i++ {
		line := strings.TrimSpace(lines[i])
		if line == "" {
			continue
		}

		if _, err := p.parseFixedFormatDate(line); err == nil {
			if i+1 < len(lines) {
				nextLine := strings.TrimSpace(lines[i+1])
				if strings.Contains(nextLine, "UPI/") || strings.Contains(nextLine, "CWDR//") {
					groupCount++
				}
			}
		}
	}

	return groupCount > 0
}

func (p *BOIBANKParser) isCSVFormat(content string) bool {
	lines := strings.Split(content, "\n")

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.Contains(trimmed, `"Sr No","Date","Remarks","Debit","Credit","Balance Amount"`) {
			return true
		}
		if strings.HasPrefix(trimmed, `"`) && strings.Count(trimmed, `","`) >= 5 {
			parts := strings.Split(trimmed, `","`)
			if len(parts) >= 6 {
				firstField := strings.Trim(parts[0], `"`)
				if _, err := strconv.Atoi(firstField); err == nil {
					return true
				}
			}
		}
	}

	return false
}

func (p *BOIBANKParser) parseCSVContent(content string) ([]Transaction, error) {
	lines := strings.Split(content, "\n")
	var transactions []Transaction
	dataStart := -1
	for i, line := range lines {
		line = strings.TrimSpace(line)
		if line == `"Sr No","Date","Remarks","Debit","Credit","Balance Amount"` {
			dataStart = i + 1
			break
		}
	}

	if dataStart == -1 {
		for i, line := range lines {
			line = strings.TrimSpace(line)
			if line == "" {
				continue
			}
			if strings.Contains(line, "Account Holder") ||
				strings.Contains(line, "CustomerId") ||
				strings.Contains(line, "IFSC") ||
				strings.Contains(line, "Account Number") {
				continue
			}
			if strings.Contains(line, "Note :") {
				break
			}
			if strings.HasPrefix(line, `"`) && strings.Count(line, `","`) >= 5 {
				dataStart = i
				break
			}
		}
	}

	if dataStart == -1 || dataStart >= len(lines) {
		return nil, fmt.Errorf("未找到BOI CSV交易数据")
	}
	for i := dataStart; i < len(lines); i++ {
		line := strings.TrimSpace(lines[i])
		if line == "" {
			continue
		}
		if strings.Contains(line, "Note :") {
			break
		}
		transaction := p.parseCSVLine(line)
		if transaction != nil {
			transactions = append(transactions, *transaction)
		}
	}

	if len(transactions) == 0 {
		return nil, fmt.Errorf("未解析到BOI CSV交易数据")
	}

	return transactions, nil
}

func (p *BOIBANKParser) parseCSVLine(line string) *Transaction {
	parts := strings.Split(line, `","`)
	if len(parts) < 6 {
		return nil
	}
	for i := range parts {
		parts[i] = strings.Trim(parts[i], `"`)
	}
	dateStr := strings.TrimSpace(parts[1])
	if !p.isValidBOIDate(dateStr) && !p.isValidDate(dateStr) {
		return nil
	}
	remarks := strings.TrimSpace(parts[2])
	debitStr := strings.TrimSpace(parts[3])
	if debitStr == "" {
		debitStr = "0"
	}
	debitStr = p.cleanAmountString(debitStr)
	creditStr := strings.TrimSpace(parts[4])
	if creditStr == "" {
		creditStr = "0"
	}
	creditStr = p.cleanAmountString(creditStr)
	balanceStr := strings.TrimSpace(parts[5])
	if balanceStr == "" {
		balanceStr = "0"
	}
	balanceStr = p.cleanAmountString(balanceStr)
	var amountStr string
	var transType string

	debit, _ := strconv.ParseFloat(debitStr, 64)
	credit, _ := strconv.ParseFloat(creditStr, 64)

	if debit > 0 {
		amountStr = debitStr
		transType = "TYPE_OUT"
	} else if credit > 0 {
		amountStr = creditStr
		transType = "TYPE_IN"
	} else {
		transType = p.determineTransTypeFromDesc(remarks)
		amountStr = "0"
	}
	amountFloat, err := strconv.ParseFloat(amountStr, 64)
	if err != nil {
		amountFloat = 0
	}
	bankTxnId := p.extractTransactionIdFromDesc(remarks)
	if bankTxnId == "" {
		bankTxnId = p.generateTransactionID(remarks, dateStr, amountStr, "")
	}
	formattedDate := p.formatTransactionDate(dateStr)

	return &Transaction{
		Date:        formattedDate,
		Description: remarks,
		Amount:      amountFloat,
		Type:        transType,
		BankTxnId:   bankTxnId,
		Balance:     balanceStr,
	}
}

func (p *BOIBANKParser) parseTextContent(content string) ([]Transaction, error) {
	lines := strings.Split(content, "\n")

	var transactions []Transaction
	inDataSection := false
	skipNextLine := false

	for _, line := range lines {
		line = strings.TrimSpace(line)

		if strings.Contains(line, "Date       Description") {
			skipNextLine = true
			continue
		}

		if skipNextLine {
			skipNextLine = false
			continue
		}

		if !inDataSection && p.isTransactionLine(line) {
			inDataSection = true
		}

		if strings.Contains(line, "***Thank you for taking your Bank") ||
			strings.Contains(line, "Any discrepancy in this document") {
			inDataSection = false
			break
		}

		if !inDataSection {
			continue
		}

		if line == "" || strings.HasPrefix(line, "---") {
			continue
		}

		transaction := p.parseBOITransactionLine(line)
		if transaction != nil {
			transactions = append(transactions, *transaction)
		}
	}

	if len(transactions) == 0 {
		return nil, fmt.Errorf("未解析到BOI交易数据")
	}

	return transactions, nil
}

func (p *BOIBANKParser) isValidDate(dateStr string) bool {
	dateStr = strings.TrimSpace(dateStr)
	if dateStr == "" {
		return false
	}

	formats := []string{
		"02-01-2006",
		"02/01/2006",
		"2006-01-02",
		"02-Jan-2006",
		"02-Jan-06",
	}

	for _, format := range formats {
		if _, err := time.Parse(format, dateStr); err == nil {
			return true
		}
	}

	return false
}

func (p *BOIBANKParser) cleanAmountString(amount string) string {
	if amount == "" {
		return "0"
	}

	cleaned := strings.TrimSpace(amount)
	cleaned = strings.ReplaceAll(cleaned, ",", "")
	cleaned = strings.ReplaceAll(cleaned, "₹", "")
	cleaned = strings.ReplaceAll(cleaned, "$", "")
	cleaned = strings.ReplaceAll(cleaned, "£", "")
	cleaned = strings.ReplaceAll(cleaned, "€", "")
	cleaned = strings.Trim(cleaned, `"`)
	cleaned = strings.Trim(cleaned, `'`)
	if _, err := strconv.ParseFloat(cleaned, 64); err != nil {
		return "0"
	}

	return cleaned
}

// 格式化交易日期
func (p *BOIBANKParser) formatTransactionDate(dateStr string) string {
	dateStr = strings.TrimSpace(dateStr)
	if dateStr == "" {
		return time.Now().Format("2006-01-02 15:04:05")
	}
	formats := []string{
		"02-01-2006",
		"02/01/2006",
		"2006-01-02",
		"02-Jan-2006",
		"02-Jan-06",
		"2 Jan 2006",
		"02 Jan 2006",
		"2-Jan-2006",
		"02-Jan-2006",
	}

	for _, format := range formats {
		if t, err := time.Parse(format, dateStr); err == nil {
			now := time.Now()
			return time.Date(t.Year(), t.Month(), t.Day(),
				now.Hour(), now.Minute(), now.Second(), 0, now.Location()).Format("2006-01-02 15:04:05")
		}
	}

	return time.Now().Format("2006-01-02 15:04:05")
}

func (p *BOIBANKParser) isTransactionLine(line string) bool {
	if len(line) < 10 {
		return false
	}
	dateStr := strings.TrimSpace(line[0:10])
	return p.isValidBOIDate(dateStr)
}

func (p *BOIBANKParser) parseBOITransactionLine(line string) *Transaction {
	line = strings.TrimSpace(line)

	if len(line) < 20 {
		return nil
	}

	date := strings.TrimSpace(line[0:10])
	if !p.isValidBOIDate(date) {
		return nil
	}
	pattern := `^(\d{2}-\d{2}-\d{4})\s+(.+?)\s+(Dr|Cr)\s+([\d,.]+)\s+([\d,.]+)$`
	re := regexp.MustCompile(pattern)

	matches := re.FindStringSubmatch(line)
	if matches != nil && len(matches) == 6 {
		return p.createTransactionFromMatches(matches)
	}
	return p.manualParseTransactionLine(line, date)
}

func (p *BOIBANKParser) createTransactionFromMatches(matches []string) *Transaction {
	date := matches[1]
	description := strings.TrimSpace(matches[2])
	crDr := matches[3]
	amount := strings.ReplaceAll(matches[4], ",", "")
	amountFloat, err := strconv.ParseFloat(amount, 64)
	if err != nil {
		return nil
	}
	transType := "TYPE_OUT"
	if crDr == "Cr" {
		transType = "TYPE_IN"
	}
	formattedDate := p.formatBOIDate(date)
	bankTxnId := p.generateTransactionID(description, date, amount, crDr)

	return &Transaction{
		Date:        formattedDate,
		Description: description,
		Amount:      amountFloat,
		Type:        transType,
		BankTxnId:   bankTxnId,
	}
}

func (p *BOIBANKParser) manualParseTransactionLine(line, date string) *Transaction {
	line = strings.TrimSpace(line)
	var amountStr, balanceStr, crDr string
	parts := strings.Fields(line)
	if len(parts) < 4 {
		return nil
	}
	balanceStr = parts[len(parts)-1]
	balanceStr = strings.ReplaceAll(balanceStr, ",", "")
	amountStr = parts[len(parts)-2]
	amountStr = strings.ReplaceAll(amountStr, ",", "")

	if len(parts) >= 3 {
		crDrField := parts[len(parts)-3]
		if crDrField == "Dr" || crDrField == "Cr" {
			crDr = crDrField
		}
	}

	if crDr == "" {
		for i := len(parts) - 4; i >= 2; i-- {
			if parts[i] == "Dr" || parts[i] == "Cr" {
				crDr = parts[i]
				break
			}
		}
	}

	description := ""
	if crDr != "" {
		crDrIndex := -1
		for i, part := range parts {
			if part == crDr {
				crDrIndex = i
				break
			}
		}

		if crDrIndex > 1 {
			descParts := parts[1:crDrIndex]
			description = strings.Join(descParts, " ")
		}
	}

	if date == "" || amountStr == "" || crDr == "" || description == "" {
		return nil
	}

	amountFloat, err := strconv.ParseFloat(amountStr, 64)
	if err != nil {
		return nil
	}

	transType := "TYPE_OUT"
	if crDr == "Cr" {
		transType = "TYPE_IN"
	}
	formattedDate := p.formatBOIDate(date)
	bankTxnId := p.generateTransactionID(description, date, amountStr, crDr)

	return &Transaction{
		Date:        formattedDate,
		Description: description,
		Amount:      amountFloat,
		Type:        transType,
		BankTxnId:   bankTxnId,
	}
}

func (p *BOIBANKParser) isValidBOIDate(dateStr string) bool {
	if dateStr == "" {
		return false
	}
	pattern := `^\d{2}-\d{2}-\d{4}$`
	matched, _ := regexp.MatchString(pattern, dateStr)
	if !matched {
		return false
	}
	_, err := time.Parse("02-01-2006", dateStr)
	return err == nil
}

func (p *BOIBANKParser) formatBOIDate(dateStr string) string {
	dateStr = strings.TrimSpace(dateStr)
	if dateStr == "" {
		return time.Now().Format("2006-01-02 15:04:05")
	}

	if t, err := time.Parse("02-01-2006", dateStr); err == nil {
		now := time.Now()
		return time.Date(t.Year(), t.Month(), t.Day(),
			now.Hour(), now.Minute(), now.Second(), 0, now.Location()).Format("2006-01-02 15:04:05")
	}

	return time.Now().Format("2006-01-02 15:04:05")
}

func (p *BOIBANKParser) generateTransactionID(description, date, amount, crDr string) string {
	upiRef, _ := p.extractUPIInfo(description)
	if upiRef != "" && len(upiRef) >= 10 && len(upiRef) <= 12 {
		return upiRef
	}
	data := fmt.Sprintf("%s|%s|%s|%s|%d", description, date, amount, crDr, time.Now().UnixNano())
	hash := md5.Sum([]byte(data))
	return "BOI_" + hex.EncodeToString(hash[:])[:12]
}

func (p *BOIBANKParser) extractUPIInfo(description string) (string, string) {
	upiRef := ""
	upiTime := ""
	re := regexp.MustCompile(`(\d{10,12})`)
	matches := re.FindAllStringSubmatch(description, -1)

	if matches != nil {
		for _, match := range matches {
			if len(match) > 0 {
				numStr := match[0]
				if strings.Contains(description, "UPI/"+numStr) ||
					(strings.Contains(description, "UPI") && strings.Contains(description, numStr)) {
					upiRef = numStr
					break
				}
			}
		}
		if upiRef == "" && len(matches) > 0 && len(matches[0]) > 0 {
			numStr := matches[0][0]
			if len(numStr) >= 10 && len(numStr) <= 12 {
				upiRef = numStr
			}
		}
	}

	return upiRef, upiTime
}

func (p *BOIBANKParser) extractTransactionName(description string) string {
	if description == "" {
		return ""
	}
	if strings.Contains(description, "UPI") {
		re := regexp.MustCompile(`/CR/([^/]+)/`)
		matches := re.FindStringSubmatch(description)
		if len(matches) > 1 {
			name := strings.TrimSpace(matches[1])
			if name != "" {
				return name
			}
		}
		return "UPI Transaction"
	} else if strings.Contains(description, "CWDR") {
		return "Cash Withdrawal"
	} else if strings.Contains(description, "Charges") {
		return "Bank Charges"
	} else if strings.Contains(description, "IBNEFT") {
		parts := strings.Split(description, "/")
		if len(parts) >= 2 {
			name := strings.TrimSpace(parts[len(parts)-1])
			if name != "" {
				return name
			}
		}
		return ""
	}

	return ""
}

func (p *BOIBANKParser) extractAccountInfo(description string) string {
	if description == "" {
		return ""
	}
	parts := strings.Split(description, "/")
	if len(parts) > 1 {
		lastPart := strings.TrimSpace(parts[len(parts)-1])
		if lastPart != "" {
			if len(lastPart) >= 2 && len(lastPart) <= 4 &&
				!strings.Contains(lastPart, " ") &&
				!strings.Contains(lastPart, ".") {
				return lastPart
			}
		}
	}

	return ""
}

func (p *BOIBANKParser) ConvertToTransInfo(transactions []Transaction, holderAccount string) []TransInfo {
	var transInfos []TransInfo

	for _, txn := range transactions {
		upiRef, upiTime := p.extractUPIInfo(txn.Description)
		transDate := FormatTransactionDate(txn.Date)

		if upiTime != "" {
			datePart := extractDatePart(txn.Date)
			transDate = datePart + " " + upiTime
		}

		var fundFlow int32
		if txn.Type == "TYPE_OUT" {
			fundFlow = 1
		} else if txn.Type == "TYPE_IN" {
			fundFlow = 2
		}

		transName := p.extractTransactionName(txn.Description)
		transAccount := p.extractAccountInfo(txn.Description)
		bankTxnId := upiRef
		if bankTxnId == "" {
			bankTxnId = txn.BankTxnId
		}

		transInfo := TransInfo{
			TransType:    txn.Type,
			TransName:    transName,
			TransAccount: transAccount,
			TransUpistr:  txn.Description,
			TransAmount:  fmt.Sprintf("%.2f", txn.Amount),
			BankTxnId:    bankTxnId,
			TransDate:    transDate,
			TransStatus:  "SUCCESS",
			FundFlow:     fundFlow,
		}

		transInfos = append(transInfos, transInfo)
	}

	return transInfos
}

func (p *BOIBANKParser) GetBankName() string {
	return "BANK_OF_INDIA"
}
