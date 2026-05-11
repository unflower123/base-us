package bankparsing

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/xuri/excelize/v2"
)

type SBIBANKParser struct{}

var _ BankParser = (*SBIBANKParser)(nil)

type SBIRecord struct {
	Date               string `json:"Date"`
	TransactionDetails string `json:"TransactionDetails"`
	Debits             string `json:"Debits"`
	Credits            string `json:"Credits"`
	Balance            string `json:"Balance"`
}

func (p *SBIBANKParser) Parse(_ string, filePath string, holderAccount string) ([]Transaction, error) {
	ext := strings.ToLower(filepath.Ext(filePath))
	if ext != ".xls" && ext != ".xlsx" && ext != ".xlsm" && ext != ".xltx" && ext != ".xltm" {
		return nil, fmt.Errorf("unsupported file format for SBI. Expected: xls file, but got: %s", ext)
	}

	rows, err := p.readExcelFile(filePath)
	if err != nil {
		return nil, err
	}

	records, err := p.parseXLSXToRecords(rows)
	if err != nil {
		return nil, err
	}

	var transactions []Transaction

	for _, record := range records {
		if record.Date == "" || record.TransactionDetails == "" {
			continue
		}

		upperDesc := strings.ToUpper(record.TransactionDetails)
		if strings.Contains(upperDesc, "MIN BAL CHGS") ||
			strings.Contains(upperDesc, "CHARGES") ||
			strings.Contains(upperDesc, "SERVICE") {
			continue
		}

		var amount float64
		var transType string

		debitCleaned := p.cleanAmount(record.Debits)
		creditCleaned := p.cleanAmount(record.Credits)

		debit, _ := strconv.ParseFloat(debitCleaned, 64)
		credit, _ := strconv.ParseFloat(creditCleaned, 64)

		if debit > 0 {
			transType = "DEBIT"
			amount = -debit
		} else if credit > 0 {
			transType = "CREDIT"
			amount = credit
		} else {
			continue
		}

		transactions = append(transactions, Transaction{
			Date:        p.formatDate(record.Date),
			Description: record.TransactionDetails,
			Amount:      amount,
			Type:        transType,
			Account:     holderAccount,
		})
	}

	log.Printf("Successfully parsed %d SBI bank transactions", len(transactions))
	return transactions, nil
}

func (p *SBIBANKParser) ConvertToTransInfo(transactions []Transaction, holderAccount string) []TransInfo {
	var transInfos []TransInfo

	for _, transaction := range transactions {
		transType := "TYPE_IN"
		if transaction.Type == "DEBIT" {
			transType = "TYPE_OUT"
		}

		var fundFlow int32
		if transaction.Type == "DEBIT" {
			fundFlow = 1
		} else if transaction.Type == "CREDIT" {
			fundFlow = 2
		}

		transName := p.extractName(transaction.Description)
		transUpistr := p.extractUpi(transaction.Description)
		bankTxnId := p.extractTxnId(transaction.Description)
		transAccount := p.extractAccount(transaction.Description)

		transAmount := fmt.Sprintf("%.2f", transaction.Amount)
		if transaction.Amount < 0 {
			transAmount = fmt.Sprintf("%.2f", -transaction.Amount)
		}

		transDate := FormatDateWithIndianTime(transaction.Date)

		transInfos = append(transInfos, TransInfo{
			TransType:    transType,
			TransName:    transName,
			TransAccount: transAccount,
			TransUpistr:  transUpistr,
			TransAmount:  transAmount,
			BankTxnId:    bankTxnId,
			TransDate:    transDate,
			TransStatus:  "SUCCESS",
			FundFlow:     fundFlow,
		})
	}

	return transInfos
}

func (p *SBIBANKParser) extractName(description string) string {
	parts := strings.Split(description, "/")

	for i, part := range parts {
		part = strings.TrimSpace(part)
		if p.isPureDigits(part) && i+1 < len(parts) {
			nextPart := strings.TrimSpace(parts[i+1])
			if !strings.Contains(strings.ToUpper(nextPart), "BKID") &&
				!strings.Contains(strings.ToUpper(nextPart), "PSIB") &&
				!strings.Contains(strings.ToUpper(nextPart), "CNRB") &&
				!strings.Contains(strings.ToUpper(nextPart), "PAYME") &&
				!p.isPureDigits(nextPart) {
				nameParts := strings.Fields(nextPart)
				if len(nameParts) > 0 {
					name := strings.TrimSpace(nameParts[0])
					if len(name) > 1 && !p.isPureDigits(name) {
						return name
					}
				}
			}
		}
	}

	words := strings.Fields(description)
	for _, word := range words {
		word = strings.TrimSpace(word)
		if len(word) > 2 && !p.isPureDigits(word) &&
			!strings.Contains(strings.ToUpper(word), "UPI") &&
			!strings.Contains(strings.ToUpper(word), "BKID") &&
			!strings.Contains(strings.ToUpper(word), "PSIB") &&
			!strings.Contains(strings.ToUpper(word), "CNRB") &&
			!strings.Contains(strings.ToUpper(word), "PAYME") &&
			!strings.Contains(strings.ToUpper(word), "TRANSFER") &&
			!strings.Contains(strings.ToUpper(word), "BY") &&
			!strings.Contains(strings.ToUpper(word), "CR") {
			return word
		}
	}

	return "Unknown"
}

func (p *SBIBANKParser) extractUpi(description string) string {
	if strings.Contains(description, "UPI") {
		re := regexp.MustCompile(`UPI/CR/(\d+)`)
		match := re.FindStringSubmatch(description)
		if len(match) > 1 {
			return fmt.Sprintf("upi://%s", match[1])
		}
	}
	return ""
}

func (p *SBIBANKParser) extractTxnId(description string) string {
	re := regexp.MustCompile(`UPI/CR/(\d{10,15})/`)
	match := re.FindStringSubmatch(description)
	if len(match) > 1 {
		return match[1]
	}

	re2 := regexp.MustCompile(`TRANSFER FROM (\d+)`)
	match2 := re2.FindStringSubmatch(description)
	if len(match2) > 1 {
		return match2[1]
	}

	re3 := regexp.MustCompile(`\b(\d{10,20})\b`)
	matches := re3.FindAllStringSubmatch(description, -1)
	for _, match := range matches {
		if len(match) > 1 {
			return match[1]
		}
	}

	return ""
}

func (p *SBIBANKParser) extractAccount(description string) string {
	re := regexp.MustCompile(`TRANSFER FROM (\d+)`)
	match := re.FindStringSubmatch(description)
	if len(match) > 1 {
		return match[1]
	}

	parts := strings.Split(description, "/")
	for _, part := range parts {
		trimmedPart := strings.TrimSpace(part)
		if matched, _ := regexp.MatchString(`^\d{10,20}$`, trimmedPart); matched {
			return trimmedPart
		}
	}
	return ""
}

func (p *SBIBANKParser) isPureDigits(str string) bool {
	matched, _ := regexp.MatchString(`^\d+$`, str)
	return matched
}

func (p *SBIBANKParser) normalizeAmount(amountStr string) string {
	cleaned := strings.ReplaceAll(amountStr, ",", "")
	cleaned = strings.ReplaceAll(cleaned, " ", "")
	cleaned = strings.ReplaceAll(cleaned, "\t", "")
	cleaned = strings.TrimSpace(cleaned)

	if strings.HasPrefix(cleaned, "INR") || strings.HasPrefix(cleaned, "₹") {
		cleaned = strings.TrimPrefix(cleaned, "INR")
		cleaned = strings.TrimPrefix(cleaned, "₹")
		cleaned = strings.TrimSpace(cleaned)
	}

	cleaned = strings.ReplaceAll(cleaned, "(", "")
	cleaned = strings.ReplaceAll(cleaned, ")", "")

	return cleaned
}

func (p *SBIBANKParser) formatDate(dateStr string) string {
	if dateStr == "" {
		return ""
	}

	dateStr = strings.TrimSpace(dateStr)
	if strings.Contains(dateStr, " ") {
		parts := strings.Split(dateStr, " ")
		dateStr = parts[0]
	}

	layouts := []string{
		"2 Jan 2006",
		"02 Jan 2006",
		"2-Jan-2006",
		"02-Jan-2006",
		"2 Jan 06",
		"02-Jan-06",
		"2006-01-02",
		"02/01/2006",
		"01/02/2006",
		"2 January 2006",
	}

	for _, layout := range layouts {
		parsedDate, err := time.Parse(layout, dateStr)
		if err == nil {
			return parsedDate.Format("2006-01-02")
		}
	}

	return dateStr
}

func (p *SBIBANKParser) readExcelFile(filePath string) ([][]string, error) {
	ext := strings.ToLower(filepath.Ext(filePath))
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return nil, fmt.Errorf("无法获取文件信息: %v", err)
	}

	if fileInfo.Size() == 0 {
		return nil, fmt.Errorf("文件为空")
	}

	if ext == ".xls" {
		return p.processBankStatementFile(filePath)
	}

	fileType, err := p.detectExcelFileType(filePath)
	if err != nil {
		return nil, fmt.Errorf("检测文件类型失败: %v", err)
	}

	switch fileType {
	case "excel_xls_biff8":
		return p.readXLSFileWithLibxls(filePath)
	case "excel_xlsx":
		return p.readXLSXFile(filePath)
	case "text_tab_delimited":
		return p.processBankStatementFile(filePath)
	default:
		return p.processBankStatementFile(filePath)
	}
}

func (p *SBIBANKParser) readXLSXFile(filePath string) ([][]string, error) {
	f, err := excelize.OpenFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("打开XLSX文件失败: %v", err)
	}
	defer f.Close()

	sheets := f.GetSheetList()
	if len(sheets) == 0 {
		return nil, fmt.Errorf("XLSX文件中没有工作表")
	}

	rows, err := f.GetRows(sheets[0])
	if err != nil {
		return nil, fmt.Errorf("读取工作表失败: %v", err)
	}
	return rows, nil
}

func (p *SBIBANKParser) detectExcelFileType(filePath string) (string, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}

	if len(content) == 0 {
		return "empty", nil
	}

	if len(content) >= 8 {
		if content[0] == 0xD0 && content[1] == 0xCF && content[2] == 0x11 && content[3] == 0xE0 {
			return "excel_xls_biff8", nil
		}

		if content[0] == 0x50 && content[1] == 0x4B && content[2] == 0x03 && content[3] == 0x04 {
			return "excel_xlsx", nil
		}
	}

	text := string(content)
	tabCount := strings.Count(text, "\t")
	lineCount := strings.Count(text, "\n") + 1
	if tabCount > 0 && lineCount > 0 {
		avgTabsPerLine := float64(tabCount) / float64(lineCount)
		if avgTabsPerLine > 1.0 {
			return "text_tab_delimited", nil
		}
	}

	if strings.Contains(strings.ToUpper(text), "ACCOUNT") ||
		strings.Contains(strings.ToUpper(text), "BALANCE") ||
		strings.Contains(strings.ToUpper(text), "TRANSACTION") {
		return "text_tab_delimited", nil
	}
	return "text_tab_delimited", nil
}

func (p *SBIBANKParser) processBankStatementFile(filePath string) ([][]string, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	text := string(content)
	lines := strings.Split(text, "\n")
	var rows [][]string

	for _, line := range lines {
		line = strings.TrimRight(line, "\r\n")
		line = strings.TrimSpace(line)

		if line == "" {
			continue
		}
		if strings.Contains(line, "\t") {
			fields := strings.Split(line, "\t")
			rows = append(rows, fields)
		} else {
			re := regexp.MustCompile(`\s{2,}`)
			fields := re.Split(line, -1)
			cleanedFields := make([]string, len(fields))
			for j, field := range fields {
				cleanedFields[j] = strings.TrimSpace(field)
			}
			rows = append(rows, cleanedFields)
		}
	}
	return rows, nil
}

func (p *SBIBANKParser) readXLSFileWithLibxls(filePath string) ([][]string, error) {
	os.Setenv("EXCELIZE_ENABLE_XLS", "true")
	f, err := excelize.OpenFile(filePath, excelize.Options{
		RawCellValue: true,
	})

	if err != nil {
		if strings.Contains(err.Error(), "unsupported workbook file format") ||
			strings.Contains(err.Error(), "operation not supported") {
		}
		return nil, fmt.Errorf("无法读取XLS文件: %v", err)
	}
	defer f.Close()
	sheets := f.GetSheetList()
	if len(sheets) == 0 {
		return nil, fmt.Errorf("XLS文件中没有工作表")
	}

	sheetName := sheets[0]
	rows, err := f.GetRows(sheetName)
	if err != nil {
		rawRows, err := f.Rows(sheetName)
		if err != nil {
			return nil, fmt.Errorf("获取行迭代器失败: %v", err)
		}

		var manualRows [][]string
		for rawRows.Next() {
			row, err := rawRows.Columns()
			if err != nil {
				continue
			}
			manualRows = append(manualRows, row)
		}

		return manualRows, nil
	}
	return rows, nil
}

func (p *SBIBANKParser) parseXLSXToRecords(rows [][]string) ([]SBIRecord, error) {
	var result []SBIRecord
	if len(rows) == 0 {
		return nil, fmt.Errorf("输入数据为空")
	}
	headerRow := -1
	for i, row := range rows {
		for _, cell := range row {
			if strings.Contains(strings.ToUpper(cell), "TXN DATE") {
				headerRow = i
				break
			}
		}
		if headerRow != -1 {
			break
		}
	}
	if headerRow == -1 {
	}
	dateCol := 0
	descCol := 2
	debitCol := 5
	creditCol := 6
	balanceCol := 7

	if headerRow < len(rows) {
		headerRowData := rows[headerRow]
		for i, cell := range headerRowData {
			cellUpper := strings.ToUpper(cell)
			if strings.Contains(cellUpper, "TXN DATE") || strings.Contains(cellUpper, "DATE") {
				dateCol = i
			}
			if strings.Contains(cellUpper, "DESCRIPTION") {
				descCol = i
			}
			if strings.Contains(cellUpper, "DEBIT") {
				debitCol = i
			}
			if strings.Contains(cellUpper, "CREDIT") {
				creditCol = i
			}
			if strings.Contains(cellUpper, "BALANCE") {
				balanceCol = i
			}
		}
	}
	dataStartRow := headerRow + 1

	for i := dataStartRow; i < len(rows); i++ {
		row := rows[i]

		if len(row) < 3 {
			continue
		}

		if len(row) > 0 && strings.Contains(strings.ToUpper(row[0]), "THIS IS A COMPUTER") {
			continue
		}

		txnDate := ""
		if dateCol < len(row) {
			txnDate = p.cleanDateField(row[dateCol])
		}

		description := ""
		if descCol < len(row) {
			description = strings.TrimSpace(row[descCol])
		}

		debit := "0"
		if debitCol < len(row) {
			debit = strings.TrimSpace(row[debitCol])
		}

		credit := "0"
		if creditCol < len(row) {
			credit = strings.TrimSpace(row[creditCol])
		}

		balance := ""
		if balanceCol < len(row) {
			balance = strings.TrimSpace(row[balanceCol])
		}

		txnDate = p.cleanDateField(txnDate)
		debit = p.cleanAmount(debit)
		credit = p.cleanAmount(credit)
		balance = p.cleanBalance(balance)

		if txnDate == "" || description == "" {
			continue
		}

		if !p.isValidDateForSBI(txnDate) {
			if p.containsMonthAndYear(txnDate) {
				continue
			}
		}

		if debit == "0" && credit == "0" && balance == "0" {
			continue
		}
		result = append(result, SBIRecord{
			Date:               txnDate,
			TransactionDetails: description,
			Debits:             debit,
			Credits:            credit,
			Balance:            balance,
		})
	}

	return result, nil
}

func (p *SBIBANKParser) containsMonthAndYear(str string) bool {
	str = strings.ToUpper(str)

	months := []string{"JAN", "FEB", "MAR", "APR", "MAY", "JUN",
		"JUL", "AUG", "SEP", "OCT", "NOV", "DEC"}

	hasMonth := false
	for _, month := range months {
		if strings.Contains(str, month) {
			hasMonth = true
			break
		}
	}

	hasYear := strings.Contains(str, "2025") ||
		strings.Contains(str, "2024") ||
		strings.Contains(str, "2023") ||
		strings.Contains(str, "2026")

	return hasMonth && hasYear
}

func (p *SBIBANKParser) isValidDateForSBI(dateStr string) bool {
	dateStr = strings.TrimSpace(dateStr)
	if dateStr == "" {
		return false
	}
	layouts := []string{
		"2 Jan 2006",
		"02 Jan 2006",
		"2-Jan-2006",
		"02-Jan-2006",
		"2 Jan 06",
		"02-Jan-06",
		"2006-01-02",
		"02/01/2006",
		"01/02/2006",
		"2 January 2006",
	}

	for _, layout := range layouts {
		if _, err := time.Parse(layout, dateStr); err == nil {
			return true
		}
	}
	hasMonth := false
	months := []string{"Jan", "Feb", "Mar", "Apr", "May", "Jun",
		"Jul", "Aug", "Sep", "Oct", "Nov", "Dec"}

	for _, month := range months {
		if strings.Contains(dateStr, month) {
			hasMonth = true
			break
		}
	}

	hasYear := strings.Contains(dateStr, "2025") ||
		strings.Contains(dateStr, "2024") ||
		strings.Contains(dateStr, "2023") ||
		strings.Contains(dateStr, "2026")

	if hasMonth && hasYear {
		return true
	}

	re := regexp.MustCompile(`^\d{1,2}[-/]\d{1,2}[-/]\d{4}$`)
	if re.MatchString(dateStr) {
		return true
	}
	return false
}

func (p *SBIBANKParser) cleanDateField(dateStr string) string {
	dateStr = strings.TrimSpace(dateStr)
	dateStr = strings.Join(strings.Fields(dateStr), " ")
	return dateStr
}

func (p *SBIBANKParser) cleanBalance(balance string) string {
	balance = strings.TrimSpace(balance)
	if balance == "" || balance == "-" {
		return "0"
	}
	balance = strings.ReplaceAll(balance, "INR ", "")
	balance = strings.ReplaceAll(balance, "Rs.", "")
	balance = strings.ReplaceAll(balance, "₹", "")
	balance = strings.ReplaceAll(balance, ",", "")

	return balance
}

func (p *SBIBANKParser) cleanAmount(amount string) string {
	if amount == "" || amount == "-" || amount == "." {
		return "0"
	}

	amount = strings.TrimSpace(amount)
	amount = strings.ReplaceAll(amount, ",", "")
	amount = strings.ReplaceAll(amount, "INR", "")
	amount = strings.ReplaceAll(amount, "Rs.", "")
	amount = strings.ReplaceAll(amount, "₹", "")
	amount = strings.ReplaceAll(amount, "$", "")
	amount = strings.TrimSpace(amount)

	if amount == "" {
		return "0"
	}

	if strings.HasSuffix(amount, ".") {
		amount = amount + "0"
	}

	return amount
}

func ParseSBIBANKFile(filePath string, holderAccount string) (*BankResponse, error) {
	parser := &SBIBANKParser{}
	transactions, err := parser.Parse("", filePath, holderAccount)
	if err != nil {
		return nil, fmt.Errorf("Failed to parse the SBI bank file: %v", err)
	}
	transInfos := parser.ConvertToTransInfo(transactions, holderAccount)
	response := &BankResponse{
		HolderAccount: holderAccount,
		TransInfo:     transInfos,
	}

	return response, nil
}

func (p *SBIBANKParser) MarshalStructToSortedString(v any) (string, string, error) {
	paramMap := make(map[string]any)
	jsonBytes, err := json.Marshal(v)
	if err != nil {
		return "", "", fmt.Errorf("failed to marshal input to JSON: %w", err)
	}
	decoder := json.NewDecoder(strings.NewReader(string(jsonBytes)))
	decoder.UseNumber()
	if err := decoder.Decode(&paramMap); err != nil {
		return "", "", fmt.Errorf("failed to decode JSON to map: %w", err)
	}
	paramSign, ok := paramMap["sign"].(string)
	if !ok {
		paramSign = ""
	}
	flatMap := make(map[string]any)
	p.flattenMap("", paramMap, flatMap)
	delete(flatMap, "sign")
	var keys []string
	for k := range flatMap {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	var builder strings.Builder
	for i, k := range keys {
		valueStr, err := p.valueToString(flatMap[k])
		if err != nil {
			return "", "", fmt.Errorf("failed to convert value for key %s: %w", k, err)
		}
		if i > 0 {
			builder.WriteByte('&')
		}
		builder.WriteString(k)
		builder.WriteByte('=')
		builder.WriteString(valueStr)
	}
	return builder.String(), paramSign, nil
}

func (p *SBIBANKParser) flattenMap(prefix string, input map[string]interface{}, output map[string]interface{}) {
	for k, v := range input {
		key := k
		if prefix != "" {
			key = prefix + "." + k
		}
		switch val := v.(type) {
		case map[string]interface{}:
			p.flattenMap(key, val, output)
		default:
			output[key] = v
		}
	}
}

func (p *SBIBANKParser) valueToString(v interface{}) (string, error) {
	switch val := v.(type) {
	case string:
		return val, nil
	case float64:
		return strconv.FormatFloat(val, 'f', -1, 64), nil
	case float32:
		return strconv.FormatFloat(float64(val), 'f', -1, 32), nil
	case int64:
		return strconv.FormatInt(val, 10), nil
	case int32:
		return strconv.FormatInt(int64(val), 10), nil
	case int:
		return strconv.Itoa(val), nil
	case uint64:
		return strconv.FormatUint(val, 10), nil
	case uint32:
		return strconv.FormatUint(uint64(val), 10), nil
	case uint:
		return strconv.FormatUint(uint64(val), 10), nil
	case bool:
		return strconv.FormatBool(val), nil
	case []interface{}:
		var strs []string
		for _, item := range val {
			str, err := p.valueToString(item)
			if err != nil {
				return "", err
			}
			strs = append(strs, str)
		}
		return strings.Join(strs, ","), nil
	default:
		jsonVal, err := json.Marshal(val)
		if err != nil {
			return "", fmt.Errorf("failed to serialize value: %w", err)
		}
		return string(jsonVal), nil
	}
}
