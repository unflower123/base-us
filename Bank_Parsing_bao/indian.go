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

	"github.com/xuri/excelize/v2"
)

type INDIANBANKParser struct{}

func (p *INDIANBANKParser) Parse(fileContent string, filePath string, holderAccount string) ([]Transaction, error) {
	ext := strings.ToLower(filepath.Ext(filePath))
	if ext == ".xlsx" {
		return p.parseXLSXFile(filePath)
	} else if ext == ".txt" {
		return p.parseTXTFile(filePath)
	}
	return nil, fmt.Errorf("unsupported file format for INDIAN. Expected: .xlsx or .txt, but got: %s", ext)
}

func (p *INDIANBANKParser) parseXLSXFile(filePath string) ([]Transaction, error) {
	f, err := excelize.OpenFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open XLSX file: %v", err)
	}
	defer f.Close()

	// 获取第一个工作表
	sheets := f.GetSheetList()
	if len(sheets) == 0 {
		return nil, fmt.Errorf("no sheets found in XLSX file")
	}

	rows, err := f.GetRows(sheets[0])
	if err != nil {
		return nil, fmt.Errorf("failed to get rows from sheet: %v", err)
	}
	return p.extractINDIANTransactions(rows)
}

// 解析 TXT 格式文件 - 从 api.go 移植过来
func (p *INDIANBANKParser) parseTXTFile(filePath string) ([]Transaction, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("读取文件失败: %v", err)
	}

	content := string(data)
	lines := strings.Split(content, "\n")

	var records []TXTRecord
	inDataSection := false
	headerFound := false

	for _, line := range lines {
		line = strings.TrimSpace(line)

		// 查找数据开始标记
		if strings.Contains(line, "Value Date|Post Date |") &&
			strings.Contains(line, "Remitter Branch") &&
			strings.Contains(line, "Description") {
			inDataSection = true
			headerFound = true
			continue
		}

		// 查找数据结束标记
		if strings.Contains(line, "END OF STATEMENT") ||
			strings.Contains(line, "Total") ||
			strings.Contains(line, "DOWNLOAD LIMIT") {
			break
		}

		if !inDataSection || !headerFound {
			continue
		}

		// 跳过分隔线
		if strings.HasPrefix(line, "---") || strings.Trim(line, "-") == "" {
			continue
		}

		// 解析数据行
		record := p.parseIndianTxtLine(line)
		if record != nil {
			records = append(records, *record)
		}
	}

	if len(records) == 0 {
		return nil, fmt.Errorf("未解析到交易数据")
	}

	// 转换为 Transaction 格式
	return p.convertTXTRecordsToTransactions(records), nil
}

// TXTRecord 结构体 - 从 api.go 移植过来
type TXTRecord struct {
	ValueDate        string `json:"ValueDate"`        // Value Date
	PostDate         string `json:"PostDate"`         // Post Date
	RemitterBranch   string `json:"RemitterBranch"`   // Remitter Branch
	Description      string `json:"Description"`      // Description
	ChequeNo         string `json:"ChequeNo"`         // Cheque No
	DebitAmount      string `json:"DebitAmount"`      // Debit Amount
	CreditAmount     string `json:"CreditAmount"`     // Credit Amount
	Balance          string `json:"Balance"`          // Balance
	TransactionRefNo string `json:"TransactionRefNo"` // 交易参考号
}

// 解析印度银行文本格式的单行数据 - 从 api.go 移植过来
func (p *INDIANBANKParser) parseIndianTxtLine(line string) *TXTRecord {
	// 使用管道符号分割列
	parts := strings.Split(line, "|")
	if len(parts) < 8 {
		return nil
	}

	// 清理每个字段的空格
	for i := range parts {
		parts[i] = strings.TrimSpace(parts[i])
	}

	valueDate := parts[0]
	postDate := parts[1]
	remitterBranch := parts[2]
	description := parts[3]
	chequeNo := parts[4]
	debitAmount := parts[5]
	creditAmount := parts[6]
	balance := parts[7]

	// 跳过表头行和汇总行
	if valueDate == "" || description == "" ||
		strings.Contains(strings.ToUpper(valueDate), "VALUE DATE") ||
		strings.Contains(strings.ToUpper(description), "BALANCE B/F") ||
		strings.Contains(strings.ToUpper(description), "TOTAL") {
		return nil
	}

	// 处理空余额情况
	if balance == "" {
		balance = "0.00"
	}

	normalizedDebit := p.normalizeAmount(debitAmount)
	normalizedCredit := p.normalizeAmount(creditAmount)
	cleanedBalance := p.normalizeIndianBankBalance(balance)

	// 生成交易参考号
	transactionRefNo := p.generateIndianBankTransactionRefNo(description, valueDate, normalizedDebit, normalizedCredit)

	record := &TXTRecord{
		ValueDate:        valueDate,
		PostDate:         postDate,
		RemitterBranch:   remitterBranch,
		Description:      description,
		ChequeNo:         chequeNo,
		DebitAmount:      normalizedDebit,
		CreditAmount:     normalizedCredit,
		Balance:          cleanedBalance,
		TransactionRefNo: transactionRefNo,
	}

	return record
}

// 金额标准化 - 从 api.go 移植过来
func (p *INDIANBANKParser) normalizeAmount(amount string) string {
	if amount == "" {
		return "0.00"
	}

	cleaned := strings.TrimSpace(amount)
	cleaned = strings.ReplaceAll(cleaned, ",", "")
	cleaned = strings.ReplaceAll(cleaned, " ", "")

	// 检查是否是数字
	if cleaned == "" || cleaned == "-" || cleaned == "0" {
		return "0.00"
	}
	if _, err := strconv.ParseFloat(cleaned, 64); err != nil {
		return "0.00"
	}

	return cleaned
}

// 印度银行余额标准化 - 从 api.go 移植过来
func (p *INDIANBANKParser) normalizeIndianBankBalance(balance string) string {
	if balance == "" {
		return "0.00"
	}

	cleaned := strings.TrimSpace(balance)
	cleaned = strings.ReplaceAll(cleaned, ",", "")
	cleaned = strings.ReplaceAll(cleaned, " ", "")

	// 处理 CR 后缀（贷方余额）
	if strings.Contains(strings.ToUpper(cleaned), "CR") {
		cleaned = strings.ReplaceAll(strings.ToUpper(cleaned), "CR", "")
	}

	cleaned = strings.TrimSpace(cleaned)

	if cleaned == "" {
		return "0.00"
	}

	// 验证是否为有效数字
	if _, err := strconv.ParseFloat(cleaned, 64); err != nil {
		return "0.00"
	}

	return cleaned
}

// 生成印度银行交易参考号 - 从 api.go 移植过来
func (p *INDIANBANKParser) generateIndianBankTransactionRefNo(description, date, debit, credit string) string {
	// 使用描述、日期和金额生成唯一标识
	data := fmt.Sprintf("%s|%s|%s|%s|%d", description, date, debit, credit, time.Now().UnixNano())
	hash := md5.Sum([]byte(data))
	return "INB_" + hex.EncodeToString(hash[:])[:12]
}

// 将 TXTRecord 转换为 Transaction - 从 api.go 移植过来并适配
func (p *INDIANBANKParser) convertTXTRecordsToTransactions(records []TXTRecord) []Transaction {
	var transactions []Transaction

	for _, record := range records {
		narration := record.Description
		transAmount := p.calculateTransAmount(record)
		transType := p.determineTransType(record)
		bankTxnId := p.extractBankTxnId(narration)

		amount, _ := strconv.ParseFloat(transAmount, 64)

		transaction := Transaction{
			Date:        p.formatTransactionDate(record.PostDate),
			Description: narration,
			Amount:      amount,
			Type:        transType,
			BankTxnId:   bankTxnId,
		}

		transactions = append(transactions, transaction)
	}

	return transactions
}

// 计算交易金额 - 从 api.go 移植过来
func (p *INDIANBANKParser) calculateTransAmount(record TXTRecord) string {
	debit := p.normalizeAmount(record.DebitAmount)
	credit := p.normalizeAmount(record.CreditAmount)

	debitFloat, _ := strconv.ParseFloat(debit, 64)
	creditFloat, _ := strconv.ParseFloat(credit, 64)

	if debitFloat > 0 {
		return debit
	} else if creditFloat > 0 {
		return credit
	}
	return "0.00"
}

// 确定交易类型 - 从 api.go 移植过来
func (p *INDIANBANKParser) determineTransType(record TXTRecord) string {
	debit := p.normalizeAmount(record.DebitAmount)
	credit := p.normalizeAmount(record.CreditAmount)

	debitFloat, _ := strconv.ParseFloat(debit, 64)
	creditFloat, _ := strconv.ParseFloat(credit, 64)

	if debitFloat > 0 {
		return "TYPE_OUT"
	} else if creditFloat > 0 {
		return "TYPE_IN"
	}

	return "Parsing failed"
}

func (p *INDIANBANKParser) extractBankTxnId(narration string) string {
	re := regexp.MustCompile(`\d{11}`)
	matches := re.FindStringSubmatch(narration)

	if len(matches) > 0 {
		return matches[0]
	}

	return "unknown"
}

func (p *INDIANBANKParser) formatTransactionDate(dateStr string) string {
	dateStr = strings.TrimSpace(dateStr)
	if dateStr == "" {
		return time.Now().Format("2006-01-02 15:04:05")
	}
	formats := []string{
		"02/01/2006",
		"02-01-2006",
		"2006-01-02",
		"01/02/2006",
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

// 原有的 XLSX 解析方法保持不变
func (p *INDIANBANKParser) extractINDIANTransactions(rows [][]string) ([]Transaction, error) {
	var transactions []Transaction
	headerRowIndex := p.findINDIANHeaderRow(rows)
	if headerRowIndex == -1 {
		return nil, nil
	}

	startRow := headerRowIndex + 1
	for i := startRow; i < len(rows); i++ {
		row := rows[i]
		if len(row) == 0 {
			continue
		}
		if p.isFooterRow(row) {
			break
		}

		transaction, err := p.parseINDIANTransactionRow(row)
		if err != nil {
			continue
		}

		if transaction != nil {
			transactions = append(transactions, *transaction)
		}
	}

	if len(transactions) == 0 {
		return nil, fmt.Errorf("未解析到任何交易数据")
	}

	return transactions, nil
}

func (p *INDIANBANKParser) findINDIANHeaderRow(rows [][]string) int {
	for i, row := range rows {
		if len(row) < 5 {
			continue
		}
		rowText := strings.ToUpper(strings.Join(row, "|"))
		hasDate := strings.Contains(rowText, "DATE")
		hasTransactionDetails := strings.Contains(rowText, "TRANSACTION DETAILS") || strings.Contains(rowText, "TRANSACTIONDETAILS")
		hasDebit := strings.Contains(rowText, "DEBIT") || strings.Contains(rowText, "DEBITS")
		hasCredit := strings.Contains(rowText, "CREDIT") || strings.Contains(rowText, "CREDITS")

		if hasDate && hasTransactionDetails && (hasDebit || hasCredit) {
			return i
		}
	}
	return -1
}

func (p *INDIANBANKParser) parseINDIANTransactionRow(row []string) (*Transaction, error) {
	if len(row) < 5 {
		return nil, fmt.Errorf("行数据不足")
	}
	date, description, debit, credit := p.analyzeINDIANRowData(row)
	if date == "" || !p.isValidDate(date) {
		return nil, fmt.Errorf("无效的交易日期: '%s'", date)
	}
	amountStr := ""
	transType := ""
	if debit != "" && debit != "0" && debit != "-" && debit != "0.00" {
		amountStr = debit
		transType = "TYPE_OUT"
	} else if credit != "" && credit != "0" && credit != "-" && credit != "0.00" {
		amountStr = credit
		transType = "TYPE_IN"
	} else {
		return nil, fmt.Errorf("未找到有效金额")
	}
	amountStr = p.cleanAmountString(amountStr)
	amount, err := p.parseFloat(amountStr)
	if err != nil {
		return nil, fmt.Errorf("解析金额失败 '%s': %v", amountStr, err)
	}
	transaction := &Transaction{
		Date:        date,
		Description: description,
		Amount:      amount,
		Type:        transType,
	}

	return transaction, nil
}

func (p *INDIANBANKParser) analyzeINDIANRowData(row []string) (string, string, string, string) {
	var date, description, debit, credit string

	for i, cell := range row {
		cell = strings.TrimSpace(cell)
		if cell == "" {
			continue
		}
		if p.isValidDate(cell) && date == "" && i >= 1 {
			date = cell
			continue
		}
		if p.isAmount(cell) {
			if debit == "" && i >= 8 && i <= 9 {
				debit = cell
			} else if credit == "" && i >= 10 && i <= 12 {
				credit = cell
			}
			continue
		}
		if description == "" && len(cell) > 10 && i >= 2 && i <= 7 {
			description = cell
		}
	}

	return date, description, debit, credit
}

func (p *INDIANBANKParser) isAmount(s string) bool {
	s = strings.TrimSpace(s)
	if s == "" || s == "-" {
		return false
	}
	clean := strings.ReplaceAll(s, ",", "")
	clean = strings.ReplaceAll(clean, " ", "")
	clean = strings.ReplaceAll(clean, "INR", "")
	clean = strings.ReplaceAll(clean, "Cr", "")
	clean = strings.ReplaceAll(clean, "Dr", "")
	_, err := strconv.ParseFloat(clean, 64)
	return err == nil
}

func (p *INDIANBANKParser) isFooterRow(row []string) bool {
	if len(row) == 0 {
		return true
	}
	rowText := strings.Join(row, " ")
	return strings.Contains(strings.ToUpper(rowText), "PAGE") ||
		strings.Contains(rowText, "This is computer-generated") ||
		strings.Contains(rowText, "Account Activity") ||
		len(rowText) < 5
}

func (p *INDIANBANKParser) cleanAmountString(amountStr string) string {
	amountStr = strings.TrimSpace(amountStr)
	amountStr = strings.ReplaceAll(amountStr, ",", "")
	amountStr = strings.ReplaceAll(amountStr, "INR", "")
	amountStr = strings.ReplaceAll(amountStr, "Cr", "")
	amountStr = strings.ReplaceAll(amountStr, "Dr", "")
	amountStr = strings.ReplaceAll(amountStr, " ", "")
	return amountStr
}

func (p *INDIANBANKParser) isValidDate(dateStr string) bool {
	if dateStr == "" {
		return false
	}
	formats := []string{
		"02 Jan 2006",
		"02-Jan-2006",
		"02/01/2006",
		"2006-01-02",
		"02 Jan 2006 15:04:05",
	}

	for _, format := range formats {
		if _, err := time.Parse(format, dateStr); err == nil {
			return true
		}
	}

	return false
}

func (p *INDIANBANKParser) parseFloat(s string) (float64, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return 0, fmt.Errorf("empty string")
	}

	var result float64
	_, err := fmt.Sscanf(s, "%f", &result)
	if err != nil {
		if floatResult, err := strconv.ParseFloat(s, 64); err == nil {
			return floatResult, nil
		}
		return 0, err
	}
	return result, nil
}

func (p *INDIANBANKParser) ConvertToTransInfo(transactions []Transaction, holderAccount string) []TransInfo {
	var transInfos []TransInfo
	for _, txn := range transactions {
		upiRef, upiTime := p.extractUPIInfo(txn.Description)
		transDate := FormatTransactionDate(txn.Date)
		if upiTime != "" {
			datePart := extractDatePart(txn.Date)
			transDate = datePart + " " + upiTime
		}

		// 设置 FundFlow 值
		var fundFlow int32
		if txn.Type == "TYPE_OUT" {
			fundFlow = 1 // debit去向
		} else if txn.Type == "TYPE_IN" {
			fundFlow = 2 // credit来源
		}

		transInfo := TransInfo{
			TransType:    txn.Type,
			TransName:    p.extractNameFromDescription(txn.Description),
			TransAccount: p.extractAccountFromNarration(txn.Description),
			TransUpistr:  txn.Description,
			TransAmount:  fmt.Sprintf("%.2f", txn.Amount),
			BankTxnId:    upiRef,
			TransDate:    transDate,
			TransStatus:  "SUCCESS",
			FundFlow:     fundFlow, // 添加 FundFlow 字段
		}
		transInfos = append(transInfos, transInfo)
	}

	return transInfos
}

func (p *INDIANBANKParser) extractUPIInfo(narration string) (string, string) {
	if !strings.Contains(narration, "UPI") {
		return "", ""
	}
	parts := strings.Split(narration, "/")
	upiRef := ""
	upiTime := ""

	for _, part := range parts {
		part = strings.TrimSpace(part)

		if upiRef == "" && len(part) >= 10 && p.isNumeric(part) {
			upiRef = part
		}
		if upiTime == "" && p.isValidTimeFormat(part) {
			upiTime = part
		}

		if upiRef != "" && upiTime != "" {
			break
		}
	}

	return upiRef, upiTime
}

func (p *INDIANBANKParser) isValidTimeFormat(timeStr string) bool {
	if len(timeStr) != 8 {
		return false
	}

	if timeStr[2] != ':' || timeStr[5] != ':' {
		return false
	}
	hourStr := timeStr[0:2]
	minuteStr := timeStr[3:5]
	secondStr := timeStr[6:8]

	if !p.isNumeric(hourStr) || !p.isNumeric(minuteStr) || !p.isNumeric(secondStr) {
		return false
	}
	hour, _ := strconv.Atoi(hourStr)
	minute, _ := strconv.Atoi(minuteStr)
	second, _ := strconv.Atoi(secondStr)

	return hour >= 0 && hour <= 23 && minute >= 0 && minute <= 59 && second >= 0 && second <= 59
}

func (p *INDIANBANKParser) extractNameFromDescription(description string) string {
	if description == "" {
		return "Bank Transaction"
	}
	if strings.Contains(description, "UPI") {
		return "UPI Transaction"
	}
	return "Bank Transaction"
}

func (p *INDIANBANKParser) extractAccountFromNarration(narration string) string {
	if narration == "" {
		return "N/A"
	}
	if strings.Contains(narration, "@") {
		parts := strings.Split(narration, "/")
		for _, part := range parts {
			part = strings.TrimSpace(part)
			if strings.Contains(part, "@") {
				return part
			}
		}
	}
	return "N/A"
}

func (p *INDIANBANKParser) isNumeric(s string) bool {
	_, err := strconv.ParseFloat(s, 64)
	return err == nil
}

func (p *INDIANBANKParser) GetBankName() string {
	return "INDIANBANK"
}
