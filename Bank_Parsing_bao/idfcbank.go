package bankparsing

import (
	"fmt"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/xuri/excelize/v2"
)

type IDFCBANKParser struct{}

var _ BankParser = (*IDFCBANKParser)(nil)

func (p *IDFCBANKParser) Parse(fileContent string, filePath string, holderAccount string) ([]Transaction, error) {
	ext := strings.ToLower(filepath.Ext(filePath))
	if ext != ".xlsx" {
		return nil, fmt.Errorf("unsupported file format for IDFC. Expected: .xlsx (Excel 2007+ format), but got: %s", ext)
	}

	records, err := p.ParseXLSXToRecords(filePath)
	if err != nil {
		return nil, err
	}
	var transactions []Transaction
	for _, record := range records {
		var amount float64
		var transType string

		if record.CRDR == "TYPE_IN" || record.CRDR == "TYPE_OUT" {
			cleanedAmount := p.normalizeAmount(record.Amount)
			amount, err = strconv.ParseFloat(cleanedAmount, 64)
			if err != nil {
				continue
			}
			transType = record.CRDR
		} else {
			continue
		}

		transactions = append(transactions, Transaction{
			Date:        p.formatDate(record.TxnDate),
			ValueDate:   p.formatDate(record.ValueDate),
			Description: record.Description,
			Amount:      amount,
			Type:        transType,
			Account:     holderAccount,
			ChequeNo:    record.ChequeNo,
		})
	}
	return transactions, nil
}

func (p *IDFCBANKParser) ParseXLSXToRecords(filePath string) ([]XLSXRecord, error) {
	rows, err := p.readXLSXFile(filePath)
	if err != nil {
		return nil, err
	}
	return p.extractIDFCTransactions(rows)
}

func (p *IDFCBANKParser) readXLSXFile(filePath string) ([][]string, error) {
	f, err := excelize.OpenFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open the XLSX file: %v", err)
	}
	defer f.Close()
	sheetName := f.GetSheetName(0)
	if sheetName == "" {
		return nil, fmt.Errorf("worksheet not found")
	}
	rows, err := f.GetRows(sheetName)
	if err != nil {
		return nil, fmt.Errorf("failed to read the worksheet: %v", err)
	}
	var nonEmptyRows [][]string
	for _, row := range rows {
		if !p.isEmptyRow(row) {
			nonEmptyRows = append(nonEmptyRows, row)
		}
	}

	if len(nonEmptyRows) == 0 {
		return nil, fmt.Errorf("no valid data available")
	}

	return nonEmptyRows, nil
}

func (p *IDFCBANKParser) isEmptyRow(cells []string) bool {
	for _, cell := range cells {
		if strings.TrimSpace(cell) != "" {
			return false
		}
	}
	return true
}

func (p *IDFCBANKParser) extractIDFCTransactions(rows [][]string) ([]XLSXRecord, error) {
	var records []XLSXRecord
	headerRowIndex := p.findIDFCHeaderRow(rows)
	if headerRowIndex == -1 {
		return nil, fmt.Errorf("header row not found")
	}
	startRow := headerRowIndex + 1
	transactionCount := 0
	for i := startRow; i < len(rows); i++ {
		row := rows[i]
		if len(row) < 6 {
			fmt.Printf("Row %d has insufficient columns: %d\n", i, len(row))
			continue
		}
		if p.isFooterRow(row) {
			fmt.Printf("Footer found at row %d, stopping\n", i)
			break
		}
		if p.isEmptyDataRow(row) {
			fmt.Printf("Row %d is empty data row\n", i)
			continue
		}

		record, err := p.parseIDFCTransactionRow(row)
		if err != nil {
			fmt.Printf("Failed to parse row %d: %v\n", i, err)
			continue
		}

		if record != nil {
			records = append(records, *record)
			transactionCount++
		}
	}
	return records, nil
}

func (p *IDFCBANKParser) findIDFCHeaderRow(rows [][]string) int {
	for i, row := range rows {
		if len(row) < 6 {
			continue
		}
		rowText := strings.ToUpper(strings.Join(row, "|"))
		hasTransactionDate := strings.Contains(rowText, "TRANSACTION DATE") ||
			strings.Contains(rowText, "TRAN DATE") ||
			strings.Contains(rowText, "DATE")
		hasValueDate := strings.Contains(rowText, "VALUE DATE") ||
			strings.Contains(rowText, "VAL DATE")
		hasParticulars := strings.Contains(rowText, "PARTICULARS") ||
			strings.Contains(rowText, "DESCRIPTION") ||
			strings.Contains(rowText, "NARRATION")
		hasDebit := strings.Contains(rowText, "DEBIT")
		hasCredit := strings.Contains(rowText, "CREDIT")
		if hasTransactionDate && hasValueDate && hasParticulars && (hasDebit || hasCredit) {
			fmt.Printf("Found header at row %d\n", i)
			return i
		}
	}
	return -1
}

func (p *IDFCBANKParser) parseIDFCTransactionRow(row []string) (*XLSXRecord, error) {
	if len(row) < 7 {
		return nil, fmt.Errorf("insufficient row data，只有 %d 列", len(row))
	}

	var txnDate, valueDate, description, chequeNo, debit, credit string
	txnDate = strings.TrimSpace(row[0])
	valueDate = strings.TrimSpace(row[1])
	description = strings.TrimSpace(row[2])

	if len(row) > 3 {
		chequeNo = strings.TrimSpace(row[3])
	}
	if len(row) > 4 {
		debit = strings.TrimSpace(row[4])
	}
	if len(row) > 5 {
		credit = strings.TrimSpace(row[5])
	}

	if txnDate == "" && description == "" && debit == "" && credit == "" {
		return nil, fmt.Errorf("empty row")
	}

	if txnDate == "" || !p.isValidDate(txnDate) {
		return nil, fmt.Errorf("invalid transaction date: '%s'", txnDate)
	}
	var amountStr string
	var transType string

	if debit != "" && debit != "0" && debit != "0.00" && debit != "-" {
		amountStr = debit
		transType = "TYPE_OUT"
	} else if credit != "" && credit != "0" && credit != "0.00" && credit != "-" {
		amountStr = credit
		transType = "TYPE_IN"
	} else {
		return nil, fmt.Errorf("no valid amount found, debit: '%s', credit: '%s'", debit, credit)
	}

	amountStr = p.normalizeAmount(amountStr)
	_, err := strconv.ParseFloat(amountStr, 64)
	if err != nil {
		return nil, fmt.Errorf("failed to parse the amount '%s': %v", amountStr, err)
	}

	record := &XLSXRecord{
		TxnDate:     txnDate,
		ValueDate:   valueDate,
		Description: description,
		ChequeNo:    chequeNo,
		CRDR:        transType,
		Amount:      amountStr,
		Debit:       debit,
		Credit:      credit,
	}

	return record, nil
}

func (p *IDFCBANKParser) isEmptyDataRow(row []string) bool {
	if len(row) == 0 {
		return true
	}
	hasDate := len(row) > 0 && strings.TrimSpace(row[0]) != ""
	hasDescription := len(row) > 2 && strings.TrimSpace(row[2]) != ""
	hasDebit := len(row) > 4 && strings.TrimSpace(row[4]) != ""
	hasCredit := len(row) > 5 && strings.TrimSpace(row[5]) != ""

	return !hasDate && !hasDescription && !hasDebit && !hasCredit
}

func (p *IDFCBANKParser) ConvertToTransInfo(transactions []Transaction, holderAccount string) []TransInfo {
	var transInfos []TransInfo
	for _, txn := range transactions {
		upiRef, _ := p.extractUPIInfo(txn.Description)
		transDate := FormatTransactionDate(txn.Date)

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

func (p *IDFCBANKParser) extractUPIInfo(narration string) (string, string) {
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

func (p *IDFCBANKParser) extractNameFromDescription(description string) string {
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

	return "Bank Transaction"
}

func (p *IDFCBANKParser) extractAccountFromNarration(narration string) string {
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

func (p *IDFCBANKParser) normalizeAmount(amountStr string) string {
	cleaned := strings.ReplaceAll(amountStr, ",", "")
	cleaned = strings.ReplaceAll(cleaned, " ", "")
	cleaned = strings.TrimSpace(cleaned)
	return cleaned
}

func (p *IDFCBANKParser) formatDate(dateStr string) string {
	if dateStr == "" {
		return ""
	}

	formats := []string{
		"02-Jan-2006",
		"02-January-2006",
		"02/01/2006",
		"2006-01-02",
		"02-01-2006",
	}

	for _, format := range formats {
		if t, err := time.Parse(format, dateStr); err == nil {
			return t.Format("2006-01-02")
		}
	}
	return dateStr
}

func (p *IDFCBANKParser) isValidDate(dateStr string) bool {
	if dateStr == "" {
		return false
	}

	formats := []string{
		"02-Jan-2006",
		"02-January-2006",
		"02/01/2006",
		"2006-01-02",
		"02-01-2006",
	}

	for _, format := range formats {
		if _, err := time.Parse(format, dateStr); err == nil {
			return true
		}
	}

	return false
}

func (p *IDFCBANKParser) isNumeric(s string) bool {
	_, err := strconv.ParseFloat(s, 64)
	return err == nil
}

func (p *IDFCBANKParser) isFooterRow(row []string) bool {
	if len(row) == 0 {
		return true
	}

	rowText := strings.Join(row, " ")
	return strings.Contains(strings.ToUpper(rowText), "TOTAL") ||
		strings.Contains(strings.ToUpper(rowText), "END OF STATEMENT") ||
		strings.Contains(strings.ToUpper(rowText), "TOTAL NUMBER OF") ||
		strings.Contains(rowText, "Total number of") ||
		strings.Contains(rowText, "End of the Statement")
}

func (p *IDFCBANKParser) GetBankName() string {
	return "IDFC"
}

func (p *IDFCBANKParser) SupportsFile(filePath string) bool {
	return strings.HasSuffix(strings.ToLower(filePath), ".xlsx")
}

func ParseIDFCBANKFile(filePath string, holderAccount string) (*BankResponse, error) {
	parser := &IDFCBANKParser{}

	transactions, err := parser.Parse("", filePath, holderAccount)
	if err != nil {
		return nil, fmt.Errorf("failed to parse the IDFC Bank file: %v", err)
	}

	transInfos := parser.ConvertToTransInfo(transactions, holderAccount)
	response := &BankResponse{
		HolderAccount: holderAccount,
		TransInfo:     transInfos,
	}
	return response, nil
}
