package bankparsing

import (
	"fmt"
	"log"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/extrame/xls"
)

type CUBBANKParser struct{}

var _ BankParser = (*CUBBANKParser)(nil)

func (p *CUBBANKParser) Parse(fileContent string, filePath string, holderAccount string) ([]Transaction, error) {
	ext := strings.ToLower(filepath.Ext(filePath))
	if ext != ".xls" {
		return nil, fmt.Errorf("unsupported file format for CUB. Expected: .xls (Excel 97-2003 format), but got: %s", ext)
	}

	transactions, err := p.parseXLSFile(filePath)
	if err != nil {
		return nil, err
	}
	p.sortTransactionsByDate(transactions)

	return transactions, nil
}

func (p *CUBBANKParser) sortTransactionsByDate(transactions []Transaction) {
	sort.Slice(transactions, func(i, j int) bool {
		dateI, errI := p.parseDateWithIndianTimezone(transactions[i].Date)
		dateJ, errJ := p.parseDateWithIndianTimezone(transactions[j].Date)

		if errI != nil || errJ != nil {
			return i > j // 如果解析失败，保持原顺序
		}

		return dateI.After(dateJ)
	})
}

func (p *CUBBANKParser) parseDateWithIndianTimezone(dateStr string) (time.Time, error) {
	loc, err := time.LoadLocation("Asia/Kolkata")
	if err != nil {
		loc = time.FixedZone("IST", 5*60*60+30*60) // UTC+5:30
	}

	formats := []string{
		"02/01/2006",
		"02-01-2006",
		"2006-01-02",
		"02/01/06",
	}

	for _, format := range formats {
		if t, err := time.ParseInLocation(format, dateStr, loc); err == nil {
			return t, nil
		}
	}

	return time.Time{}, fmt.Errorf("unable to parse the date: %s", dateStr)
}

func (p *CUBBANKParser) parseXLSFile(filePath string) ([]Transaction, error) {
	file, err := xls.Open(filePath, "utf-8")
	if err != nil {
		return nil, fmt.Errorf("failed to open XLS file: %v", err)
	}

	// 获取第一个工作表
	sheet := file.GetSheet(0)
	if sheet == nil {
		return nil, fmt.Errorf("no sheets found in XLS file")
	}

	var rows [][]string
	maxRow := int(sheet.MaxRow)
	if maxRow <= 0 {
		maxRow = 100 // 设置合理上限
	}

	// 使用 panic 恢复机制保护整个解析过程
	for i := 0; i <= maxRow; i++ {
		func() {
			defer func() {
				if r := recover(); r != nil {
					log.Printf("读取第 %d 行时发生 panic，已恢复: %v", i+1, r)
				}
			}()

			row := sheet.Row(i)
			if row == nil {
				return
			}

			var cells []string
			colCount := row.LastCol()
			if colCount <= 0 {
				return
			}

			for j := 0; j < colCount; j++ {
				cell := row.Col(j)
				cells = append(cells, strings.TrimSpace(cell))
			}

			if len(cells) > 0 && p.hasValidData(cells) {
				rows = append(rows, cells)
			}
		}()
	}

	log.Printf("CUB银行文件读取完成，共 %d 行数据", len(rows))

	return p.extractCUBTransactions(rows)
}

func (p *CUBBANKParser) hasValidData(cells []string) bool {
	for _, cell := range cells {
		if strings.TrimSpace(cell) != "" {
			return true
		}
	}
	return false
}

func (p *CUBBANKParser) extractCUBTransactions(rows [][]string) ([]Transaction, error) {
	var transactions []Transaction
	headerRowIndex := p.findCUBHeaderRow(rows)
	if headerRowIndex == -1 {
		return nil, fmt.Errorf("CUB银行表头行未找到")
	}

	log.Printf("找到CUB银行表头行在第 %d 行", headerRowIndex+1)

	startRow := headerRowIndex + 1
	transactionCount := 0

	for i := startRow; i < len(rows); i++ {
		row := rows[i]
		if len(row) == 0 {
			continue
		}

		if p.isFooterRow(row) {
			log.Printf("在第 %d 行遇到页脚，停止解析", i+1)
			break
		}

		if p.isEmptyDataRow(row) {
			continue
		}

		transaction, err := p.parseCUBTransactionRow(row)
		if err != nil {
			log.Printf("解析第 %d 行失败: %v", i+1, err)
			continue
		}

		if transaction != nil {
			transactions = append(transactions, *transaction)
			transactionCount++
		}
	}

	log.Printf("CUB银行交易解析完成，共 %d 笔交易", transactionCount)

	if len(transactions) == 0 {
		return nil, fmt.Errorf("未解析到任何CUB银行交易数据")
	}

	return transactions, nil
}

// 查找CUB银行表头行
func (p *CUBBANKParser) findCUBHeaderRow(rows [][]string) int {
	for i, row := range rows {
		if len(row) < 6 {
			continue
		}
		rowText := strings.ToUpper(strings.Join(row, "|"))
		hasDate := strings.Contains(rowText, "DATE")
		hasDescription := strings.Contains(rowText, "DESCRIPTION") || strings.Contains(rowText, "PARTICULARS")
		hasDebit := strings.Contains(rowText, "DEBIT")
		hasCredit := strings.Contains(rowText, "CREDIT")
		if hasDate && hasDescription && (hasDebit || hasCredit) {
			log.Printf("匹配到CUB银行表头: %s", rowText)
			return i
		}
	}
	return -1
}

func (p *CUBBANKParser) parseCUBTransactionRow(row []string) (*Transaction, error) {
	if len(row) < 6 {
		return nil, fmt.Errorf("行数据不足，只有 %d 列", len(row))
	}
	date := strings.TrimSpace(row[0])
	description := strings.TrimSpace(row[1])
	debit := strings.TrimSpace(row[3])
	credit := strings.TrimSpace(row[4])
	if date == "" && description == "" && debit == "" && credit == "" {
		return nil, fmt.Errorf("空数据行")
	}

	if date == "" || !p.isValidDate(date) {
		return nil, fmt.Errorf("无效的交易日期: '%s'", date)
	}

	var amountStr string
	var transType string

	if debit != "" && debit != "0" && debit != "-" && debit != "0.00" {
		amountStr = debit
		transType = "TYPE_OUT"
	} else if credit != "" && credit != "0" && credit != "-" && credit != "0.00" {
		amountStr = credit
		transType = "TYPE_IN"
	} else {
		return nil, fmt.Errorf("no valid amount found (debit: '%s', credit: '%s')", debit, credit)
	}

	amountStr = p.cleanAmountString(amountStr)
	amount, err := p.parseFloat(amountStr)
	if err != nil {
		return nil, fmt.Errorf("解析金额失败 '%s': %v", amountStr, err)
	}

	transaction := &Transaction{
		Date:        p.formatDate(date),
		ValueDate:   p.formatDate(date),
		Description: description,
		Amount:      amount,
		Type:        transType,
	}

	return transaction, nil
}

func (p *CUBBANKParser) isEmptyDataRow(row []string) bool {
	if len(row) == 0 {
		return true
	}
	hasDate := len(row) > 0 && strings.TrimSpace(row[0]) != ""
	hasDescription := len(row) > 1 && strings.TrimSpace(row[1]) != ""
	hasDebit := len(row) > 3 && strings.TrimSpace(row[3]) != ""
	hasCredit := len(row) > 4 && strings.TrimSpace(row[4]) != ""

	return !hasDate && !hasDescription && !hasDebit && !hasCredit
}

func (p *CUBBANKParser) cleanAmountString(amountStr string) string {
	amountStr = strings.TrimSpace(amountStr)
	amountStr = strings.ReplaceAll(amountStr, ",", "")
	amountStr = strings.ReplaceAll(amountStr, " ", "")
	amountStr = strings.ReplaceAll(amountStr, "Cr", "")
	amountStr = strings.ReplaceAll(amountStr, "Dr", "")
	return amountStr
}

func (p *CUBBANKParser) formatDate(dateStr string) string {
	if dateStr == "" {
		return ""
	}

	formats := []string{
		"02/01/2006",
		"02-01-2006",
		"2006-01-02",
		"02/01/06",
	}

	for _, format := range formats {
		if t, err := time.Parse(format, dateStr); err == nil {
			return t.Format("2006-01-02")
		}
	}

	return dateStr
}

func (p *CUBBANKParser) isValidDate(dateStr string) bool {
	if dateStr == "" {
		return false
	}

	formats := []string{
		"02/01/2006",
		"02-01-2006",
		"2006-01-02",
		"02/01/06",
	}

	for _, format := range formats {
		if _, err := time.Parse(format, dateStr); err == nil {
			return true
		}
	}

	return false
}

func (p *CUBBANKParser) parseFloat(s string) (float64, error) {
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

func (p *CUBBANKParser) isFooterRow(row []string) bool {
	if len(row) == 0 {
		return true
	}

	rowText := strings.Join(row, " ")
	return strings.Contains(strings.ToUpper(rowText), "TOTAL") ||
		strings.Contains(strings.ToUpper(rowText), "END OF STATEMENT") ||
		strings.Contains(rowText, "Statement Downloaded") ||
		strings.Contains(rowText, "If any discrepancy") ||
		strings.Contains(rowText, "Regd. Office") ||
		strings.Contains(rowText, "Website:") ||
		strings.Contains(rowText, "CIN :") ||
		strings.Contains(rowText, "Page") ||
		strings.Contains(rowText, "This is computer-generated")
}

func (p *CUBBANKParser) ConvertToTransInfo(transactions []Transaction, holderAccount string) []TransInfo {
	var transInfos []TransInfo

	for _, txn := range transactions {
		upiRef, _ := p.extractUPIInfo(txn.Description)
		transDate := FormatTransactionDate(txn.Date)

		// 设置 FundFlow 值
		var fundFlow int32
		if txn.Type == "TYPE_OUT" {
			fundFlow = 1
		} else if txn.Type == "TYPE_IN" {
			fundFlow = 2
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
			FundFlow:     fundFlow,
		}
		transInfos = append(transInfos, transInfo)
	}

	return transInfos
}

func (p *CUBBANKParser) extractUPIInfo(narration string) (string, string) {
	if !strings.Contains(narration, "UPI") {
		return "", ""
	}

	parts := strings.Split(narration, "/")
	upiRef := ""

	for _, part := range parts {
		part = strings.TrimSpace(part)
		if upiRef == "" && len(part) >= 10 && p.isNumeric(part) {
			upiRef = part
			break
		}
	}

	return upiRef, ""
}

func (p *CUBBANKParser) extractNameFromDescription(description string) string {
	if description == "" {
		return "Bank Transaction"
	}

	if strings.Contains(description, "UPI") {
		parts := strings.Split(description, "/")
		if len(parts) >= 4 {
			return strings.TrimSpace(parts[3])
		}
		return "UPI Transaction"
	}

	if strings.Contains(description, "BY") {
		return "Credit Transaction"
	} else if strings.Contains(description, "TO") {
		return "Debit Transaction"
	}

	return "Bank Transaction"
}

func (p *CUBBANKParser) extractAccountFromNarration(narration string) string {
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

func (p *CUBBANKParser) isNumeric(s string) bool {
	_, err := strconv.ParseFloat(s, 64)
	return err == nil
}

func (p *CUBBANKParser) GetBankName() string {
	return "CUB"
}

func (p *CUBBANKParser) SupportsFile(filePath string) bool {
	return strings.HasSuffix(strings.ToLower(filePath), ".xls")
}

func ParseCUBBANKFile(filePath string, holderAccount string) (*BankResponse, error) {
	parser := &CUBBANKParser{}

	transactions, err := parser.Parse("", filePath, holderAccount)
	if err != nil {
		return nil, fmt.Errorf("解析CUB银行文件失败: %v", err)
	}

	transInfos := parser.ConvertToTransInfo(transactions, holderAccount)
	response := &BankResponse{
		HolderAccount: holderAccount,
		TransInfo:     transInfos,
	}
	return response, nil
}
