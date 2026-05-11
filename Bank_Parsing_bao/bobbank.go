package bankparsing

import (
	"fmt"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/extrame/xls"
)

type BOBBANKParser struct{}

var _ BankParser = (*BOBBANKParser)(nil)

func (p *BOBBANKParser) Parse(fileContent string, filePath string, holderAccount string) ([]Transaction, error) {
	ext := strings.ToLower(filepath.Ext(filePath))
	if ext != ".xls" {
		return nil, fmt.Errorf("unsupported file format for BOB. Expected: .xls (Excel 97-2003 format), but got: %s", ext)
	}

	transactions, err := p.parseXLSFile(filePath)
	if err != nil {
		return nil, err
	}

	p.sortTransactionsByDate(transactions)

	return transactions, nil
}

func (p *BOBBANKParser) sortTransactionsByDate(transactions []Transaction) {
	sort.Slice(transactions, func(i, j int) bool {
		dateI, errI := p.parseDateWithIndianTimezone(transactions[i].Date)
		dateJ, errJ := p.parseDateWithIndianTimezone(transactions[j].Date)

		if errI != nil || errJ != nil {
			return i > j // 如果解析失败，保持原顺序
		}

		return dateI.After(dateJ)
	})
}

func (p *BOBBANKParser) parseDateWithIndianTimezone(dateStr string) (time.Time, error) {
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

func (p *BOBBANKParser) parseXLSFile(filePath string) ([]Transaction, error) {
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
	for i := 0; i <= int(sheet.MaxRow); i++ {
		row := sheet.Row(i)
		if row == nil {
			continue
		}

		var cells []string
		for j := 0; j < row.LastCol(); j++ {
			cell := row.Col(j)
			cells = append(cells, strings.TrimSpace(cell))
		}
		if len(cells) > 0 {
			rows = append(rows, cells)
		}
	}
	return p.extractBOBTransactions(rows)
}

func (p *BOBBANKParser) extractBOBTransactions(rows [][]string) ([]Transaction, error) {
	var transactions []Transaction
	headerRowIndex := p.findBOBHeaderRow(rows)
	if headerRowIndex == -1 {
		return nil, fmt.Errorf("header row not found")
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

		transaction, err := p.parseBOBTransactionRow(row)
		if err != nil {
			continue
		}

		if transaction != nil {
			transactions = append(transactions, *transaction)
		}
	}

	if len(transactions) == 0 {
		return nil, fmt.Errorf("no transaction data parsed")
	}

	return transactions, nil
}

// 查找表头行
func (p *BOBBANKParser) findBOBHeaderRow(rows [][]string) int {
	for i, row := range rows {
		if len(row) < 10 {
			continue
		}
		rowText := strings.ToUpper(strings.Join(row, "|"))
		hasTranDate := strings.Contains(rowText, "TRAN DATE") || strings.Contains(rowText, "TRANDATE")
		hasValueDate := strings.Contains(rowText, "VALUE DATE") || strings.Contains(rowText, "VALUEDATE")
		hasNarration := strings.Contains(rowText, "NARRATION")
		if hasTranDate && hasValueDate && hasNarration {
			return i
		}
	}
	return -1
}

func (p *BOBBANKParser) parseBOBTransactionRow(row []string) (*Transaction, error) {
	if len(row) < 10 {
		return nil, fmt.Errorf("insufficient row data")
	}
	tranDate, valueDate, narration, withdrawal, deposit, _ := p.analyzeRowData(row)
	if tranDate == "" || !p.isValidDate(tranDate) {
		return nil, fmt.Errorf("unable to parse the date: '%s'", tranDate)
	}
	amountStr := ""
	transType := ""

	if withdrawal != "" && withdrawal != "0" && withdrawal != "-" && withdrawal != "0.00" {
		amountStr = withdrawal
		transType = "TYPE_OUT"
	} else if deposit != "" && deposit != "0" && deposit != "-" && deposit != "0.00" {
		amountStr = deposit
		transType = "TYPE_IN"
	} else {
		return nil, fmt.Errorf("no valid amount found")
	}
	amountStr = p.cleanAmountString(amountStr)
	amount, err := p.parseFloat(amountStr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse the amount '%s': %v", amountStr, err)
	}

	transaction := &Transaction{
		Date:        tranDate,
		ValueDate:   valueDate,
		Description: narration,
		Amount:      amount,
		Type:        transType,
	}

	return transaction, nil
}

func (p *BOBBANKParser) analyzeRowData(row []string) (string, string, string, string, string, string) {
	var tranDate, valueDate, narration, withdrawal, deposit, balance string
	for i, cell := range row {
		cell = strings.TrimSpace(cell)
		if cell == "" {
			continue
		}
		if p.isValidDate(cell) {
			if tranDate == "" {
				tranDate = cell
			} else if valueDate == "" {
				valueDate = cell
			}
			continue
		}
		if p.isAmount(cell) {
			if withdrawal == "" && i < 15 {
				withdrawal = cell
			} else if deposit == "" {
				deposit = cell
			}
			continue
		}
		if strings.Contains(cell, "Cr") && balance == "" {
			balance = cell
			continue
		}
		if len(cell) > 20 && strings.Contains(cell, "UPI") && narration == "" {
			narration = cell
			continue
		}
	}

	return tranDate, valueDate, narration, withdrawal, deposit, balance
}

func (p *BOBBANKParser) isAmount(s string) bool {
	s = strings.TrimSpace(s)
	if s == "" {
		return false
	}
	clean := strings.ReplaceAll(s, ",", "")
	clean = strings.ReplaceAll(clean, " ", "")
	clean = strings.ReplaceAll(clean, "Cr", "")
	clean = strings.ReplaceAll(clean, "Dr", "")
	_, err := strconv.ParseFloat(clean, 64)
	return err == nil
}

func (p *BOBBANKParser) getCellValue(row []string, index int) string {
	if index < 0 || index >= len(row) {
		return ""
	}
	return strings.TrimSpace(row[index])
}

func (p *BOBBANKParser) isFooterRow(row []string) bool {
	if len(row) == 0 {
		return true
	}

	rowText := strings.Join(row, " ")
	return strings.Contains(strings.ToUpper(rowText), "PAGE") ||
		strings.Contains(rowText, "This is computer-generated") ||
		strings.Contains(rowText, "Contact-Us") ||
		strings.Contains(rowText, "18005700") ||
		strings.Contains(rowText, "Page 1 of")
}

func (p *BOBBANKParser) cleanAmountString(amountStr string) string {
	amountStr = strings.TrimSpace(amountStr)
	amountStr = strings.ReplaceAll(amountStr, ",", "")
	amountStr = strings.ReplaceAll(amountStr, "Cr", "")
	amountStr = strings.ReplaceAll(amountStr, "Dr", "")
	amountStr = strings.ReplaceAll(amountStr, " ", "")
	return amountStr
}

func (p *BOBBANKParser) isValidDate(dateStr string) bool {
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

func (p *BOBBANKParser) parseFloat(s string) (float64, error) {
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

func (p *BOBBANKParser) ConvertToTransInfo(transactions []Transaction, holderAccount string) []TransInfo {
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
			FundFlow:     fundFlow,
		}
		transInfos = append(transInfos, transInfo)
	}

	return transInfos
}

func (p *BOBBANKParser) extractUPIInfo(narration string) (string, string) {
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

func (p *BOBBANKParser) isValidTimeFormat(timeStr string) bool {
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

func (p *BOBBANKParser) extractNameFromDescription(description string) string {
	if description == "" {
		return "Bank Transaction"
	}
	if strings.Contains(description, "UPI") {
		return "UPI Transaction"
	}
	return "Bank Transaction"
}

func (p *BOBBANKParser) extractAccountFromNarration(narration string) string {
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

func (p *BOBBANKParser) isNumeric(s string) bool {
	_, err := strconv.ParseFloat(s, 64)
	return err == nil
}

func (p *BOBBANKParser) ParseXLSXToRecords(filePath string) ([]XLSXRecord, error) {
	transactions, err := p.Parse("", filePath, "")
	if err != nil {
		return nil, err
	}

	var records []XLSXRecord
	for _, txn := range transactions {
		record := XLSXRecord{
			TxnDate:     txn.Date,
			ValueDate:   txn.ValueDate,
			Description: txn.Description,
			ChequeNo:    txn.ChequeNo,
			CRDR:        txn.Type,
			Amount:      fmt.Sprintf("%.2f", txn.Amount),
		}
		records = append(records, record)
	}

	return records, nil
}

func ParseBOBBANKFile(filePath string, holderAccount string) (*BankResponse, error) {
	parser := &BOBBANKParser{}

	transactions, err := parser.Parse("", filePath, holderAccount)
	if err != nil {
		return nil, fmt.Errorf("failed to parse the BOB Bank file: %v", err)
	}

	transInfos := parser.ConvertToTransInfo(transactions, holderAccount)
	response := &BankResponse{
		HolderAccount: holderAccount,
		TransInfo:     transInfos,
	}
	return response, nil
}

func (p *BOBBANKParser) GetBankName() string {
	return "BOB"
}

func (p *BOBBANKParser) SupportsFile(filePath string) bool {
	return strings.HasSuffix(strings.ToLower(filePath), ".xls")
}
