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

type CSBBANKParser struct{}

var _ BankParser = (*CSBBANKParser)(nil)

func (p *CSBBANKParser) Parse(fileContent string, filePath string, holderAccount string) ([]Transaction, error) {
	ext := strings.ToLower(filepath.Ext(filePath))
	if ext != ".xls" {
		return nil, fmt.Errorf("unsupported file format for CSB. Expected: .xls (Excel 97-2003 format), but got: %s", ext)
	}

	transactions, err := p.parseXLSFile(filePath)
	if err != nil {
		return nil, err
	}
	p.sortTransactionsByDate(transactions)

	return transactions, nil
}

func (p *CSBBANKParser) sortTransactionsByDate(transactions []Transaction) {
	sort.Slice(transactions, func(i, j int) bool {
		dateI, errI := p.parseDateWithIndianTimezone(transactions[i].Date)
		dateJ, errJ := p.parseDateWithIndianTimezone(transactions[j].Date)

		if errI != nil || errJ != nil {
			return i > j
		}

		return dateI.After(dateJ)
	})
}

func (p *CSBBANKParser) parseDateWithIndianTimezone(dateStr string) (time.Time, error) {
	loc, err := time.LoadLocation("Asia/Kolkata")
	if err != nil {
		loc = time.FixedZone("IST", 5*60*60+30*60)
	}

	formats := []string{
		"02-01-2006",
		"02/01/2006",
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

func (p *CSBBANKParser) parseXLSFile(filePath string) ([]Transaction, error) {
	file, err := xls.Open(filePath, "utf-8")
	if err != nil {
		return nil, fmt.Errorf("failed to open XLS file: %v", err)
	}
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
	return p.extractCSBTransactions(rows)
}

func (p *CSBBANKParser) extractCSBTransactions(rows [][]string) ([]Transaction, error) {
	var transactions []Transaction
	headerRowIndex := p.findCSBHeaderRow(rows)
	if headerRowIndex == -1 {
		return nil, fmt.Errorf("header row not found")
	}

	fmt.Printf("Found header at row %d\n", headerRowIndex)

	startRow := headerRowIndex + 1
	for i := startRow; i < len(rows); i++ {
		row := rows[i]
		if len(row) == 0 {
			continue
		}
		if p.isFooterRow(row) {
			fmt.Printf("Found footer at row %d, stopping parsing\n", i)
			break
		}

		transaction, err := p.parseCSBTransactionRow(row)
		if err != nil {
			fmt.Printf("Failed to parse row %d: %v\n", i, err)
			continue
		}

		if transaction != nil {
			transactions = append(transactions, *transaction)
			fmt.Printf("Successfully parsed transaction: %+v\n", *transaction)
		}
	}

	if len(transactions) == 0 {
		return nil, fmt.Errorf("no transaction data parsed")
	}

	fmt.Printf("Total transactions parsed: %d\n", len(transactions))
	return transactions, nil
}

func (p *CSBBANKParser) findCSBHeaderRow(rows [][]string) int {
	for i, row := range rows {
		if len(row) < 5 {
			continue
		}
		rowText := strings.ToUpper(strings.Join(row, "|"))
		hasDate := strings.Contains(rowText, "DATE")
		hasParticulars := strings.Contains(rowText, "PARTICULARS")
		hasChqRef := strings.Contains(rowText, "CHQ") || strings.Contains(rowText, "REF")
		hasDebit := strings.Contains(rowText, "DEBIT")
		hasCredit := strings.Contains(rowText, "CREDIT")

		fmt.Printf("Row %d: Date=%t, Particulars=%t, ChqRef=%t, Debit=%t, Credit=%t\n",
			i, hasDate, hasParticulars, hasChqRef, hasDebit, hasCredit)

		if hasDate && hasParticulars && hasChqRef && (hasDebit || hasCredit) {
			fmt.Printf("Found header row at index %d: %v\n", i, row)
			return i
		}
	}
	return -1
}

func (p *CSBBANKParser) parseCSBTransactionRow(row []string) (*Transaction, error) {
	if len(row) < 10 {
		return nil, fmt.Errorf("insufficient row data: only %d columns", len(row))
	}

	fmt.Printf("Parsing row with %d columns: %v\n", len(row), row)
	date, particulars, debit, credit, chequeNo := p.analyzeCSBRowData(row)

	if date == "" || !p.isValidDate(date) {
		return nil, fmt.Errorf("invalid or missing date: '%s'", date)
	}

	if particulars == "" {
		return nil, fmt.Errorf("missing description")
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
		return nil, fmt.Errorf("no valid amount found (debit: '%s', credit: '%s')", debit, credit)
	}

	amountStr = p.cleanAmountString(amountStr)
	amount, err := p.parseFloat(amountStr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse the amount '%s': %v", amountStr, err)
	}

	transaction := &Transaction{
		Date:        date,
		ValueDate:   date,
		Description: particulars,
		Amount:      amount,
		Type:        transType,
		ChequeNo:    chequeNo,
	}

	fmt.Printf("Parsed transaction: Date=%s, Desc=%s, Amount=%.2f, Type=%s, ChequeNo=%s\n",
		date, particulars, amount, transType, chequeNo)

	return transaction, nil
}

func (p *CSBBANKParser) analyzeCSBRowData(row []string) (string, string, string, string, string) {
	var date, particulars, debit, credit, chequeNo string

	fmt.Printf("Analyzing CSB row data (%d columns): %v\n", len(row), row)
	if len(row) > 0 {
		date = strings.TrimSpace(row[0])
		if p.isValidDate(date) {
			fmt.Printf("Found date at column 0: '%s'\n", date)
		}
	}
	if len(row) > 3 {
		particulars = strings.TrimSpace(row[3])
		if particulars != "" {
			fmt.Printf("Found particulars at column 3: '%s'\n", particulars)
		}
	}
	if len(row) > 10 {
		chequeNo = strings.TrimSpace(row[10])
		if chequeNo != "" && chequeNo != "0" {
			fmt.Printf("Found chequeNo at column 10: '%s'\n", chequeNo)
		}
	}
	if len(row) > 16 {
		debit = strings.TrimSpace(row[16])
		if p.isAmount(debit) && debit != "0" {
			fmt.Printf("Found debit at column 16: '%s'\n", debit)
		} else {
			debit = ""
		}
	}
	if len(row) > 20 {
		credit = strings.TrimSpace(row[20])
		if p.isAmount(credit) && credit != "0" {
			fmt.Printf("Found credit at column 20: '%s'\n", credit)
		} else {
			credit = ""
		}
	}
	if date == "" || !p.isValidDate(date) {
		for i, cell := range row {
			if p.isValidDate(cell) {
				date = cell
				fmt.Printf("Found valid date at column %d: '%s'\n", i, date)
				break
			}
		}
	}

	if particulars == "" {
		for i := 1; i < len(row) && i < 10; i++ {
			cell := strings.TrimSpace(row[i])
			if cell != "" && !p.isValidDate(cell) && !p.isAmount(cell) && cell != "0" {
				particulars = cell
				fmt.Printf("Found particulars at column %d: '%s'\n", i, particulars)
				break
			}
		}
	}

	return date, particulars, debit, credit, chequeNo
}

func (p *CSBBANKParser) isAmount(s string) bool {
	s = strings.TrimSpace(s)
	if s == "" || s == "-" {
		return false
	}
	clean := strings.ReplaceAll(s, ",", "")
	clean = strings.ReplaceAll(clean, " ", "")
	clean = strings.ReplaceAll(clean, "Cr", "")
	clean = strings.ReplaceAll(clean, "Dr", "")
	_, err := strconv.ParseFloat(clean, 64)
	return err == nil
}

func (p *CSBBANKParser) getCellValue(row []string, index int) string {
	if index < 0 || index >= len(row) {
		return ""
	}
	return strings.TrimSpace(row[index])
}

func (p *CSBBANKParser) isFooterRow(row []string) bool {
	if len(row) == 0 {
		return true
	}

	rowText := strings.Join(row, " ")
	upperText := strings.ToUpper(rowText)

	return strings.Contains(upperText, "TOTAL") ||
		strings.Contains(upperText, "CLOSING BALANCE") ||
		strings.Contains(upperText, "PAGE") ||
		(len(row) > 0 && strings.TrimSpace(row[0]) == "Total") ||
		(len(row) > 0 && strings.TrimSpace(row[0]) == "Closing Balance")
}

func (p *CSBBANKParser) cleanAmountString(amountStr string) string {
	amountStr = strings.TrimSpace(amountStr)
	amountStr = strings.ReplaceAll(amountStr, ",", "")
	amountStr = strings.ReplaceAll(amountStr, "Cr", "")
	amountStr = strings.ReplaceAll(amountStr, "Dr", "")
	amountStr = strings.ReplaceAll(amountStr, " ", "")
	return amountStr
}

func (p *CSBBANKParser) isValidDate(dateStr string) bool {
	if dateStr == "" {
		return false
	}

	formats := []string{
		"02-01-2006",
		"02/01/2006",
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

func (p *CSBBANKParser) parseFloat(s string) (float64, error) {
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

func (p *CSBBANKParser) ConvertToTransInfo(transactions []Transaction, holderAccount string) []TransInfo {
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

func (p *CSBBANKParser) extractUPIInfo(narration string) (string, string) {
	if !strings.Contains(narration, "UPI") {
		return "", ""
	}
	parts := strings.Split(narration, "/")
	upiRef := ""
	upiTime := ""

	fmt.Printf("Extracting UPI info from: %v\n", parts)

	for i, part := range parts {
		part = strings.TrimSpace(part)
		if upiRef == "" && i == 0 && p.isNumeric(part) && len(part) >= 10 {
			upiRef = part
			fmt.Printf("Found UPI ref: %s\n", upiRef)
		}
		if p.isValidDate(part) {
			if strings.Contains(part, " ") {
				dateTimeParts := strings.Split(part, " ")
				if len(dateTimeParts) > 1 && p.isValidTimeFormat(dateTimeParts[1]) {
					upiTime = dateTimeParts[1]
				}
			}
		}

		if upiRef != "" {
			break
		}
	}

	return upiRef, upiTime
}

func (p *CSBBANKParser) isValidTimeFormat(timeStr string) bool {
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

func (p *CSBBANKParser) extractNameFromDescription(description string) string {
	if description == "" {
		return "Bank Transaction"
	}

	if strings.Contains(description, "UPI") {
		parts := strings.Split(description, "/")
		fmt.Printf("UPI description parts: %v\n", parts)
		for i := len(parts) - 1; i >= 0; i-- {
			namePart := strings.TrimSpace(parts[i])
			if namePart != "" &&
				!strings.Contains(namePart, "UPI") &&
				!strings.Contains(namePart, "CR") &&
				!strings.Contains(namePart, "DR") &&
				!strings.Contains(namePart, "QR") &&
				!p.isValidDate(namePart) &&
				!p.isNumeric(namePart) {
				fmt.Printf("Extracted name: %s\n", namePart)
				return namePart
			}
		}
		return "UPI Transaction"
	}
	if len(description) > 30 {
		return description[:30] + "..."
	}
	return description
}

func (p *CSBBANKParser) extractAccountFromNarration(narration string) string {
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
	name := p.extractNameFromDescription(narration)
	if name != "Bank Transaction" && name != "UPI Transaction" {
		return name
	}

	return "N/A"
}

func (p *CSBBANKParser) isNumeric(s string) bool {
	s = strings.TrimSpace(s)
	if s == "" {
		return false
	}
	_, err := strconv.ParseFloat(s, 64)
	return err == nil
}

func (p *CSBBANKParser) ParseXLSXToRecords(filePath string) ([]XLSXRecord, error) {
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

func ParseCSBBANKFile(filePath string, holderAccount string) (*BankResponse, error) {
	parser := &CSBBANKParser{}

	transactions, err := parser.Parse("", filePath, holderAccount)
	if err != nil {
		return nil, fmt.Errorf("failed to parse the CSB Bank file: %v", err)
	}

	transInfos := parser.ConvertToTransInfo(transactions, holderAccount)
	response := &BankResponse{
		HolderAccount: holderAccount,
		TransInfo:     transInfos,
	}
	return response, nil
}

func (p *CSBBANKParser) GetBankName() string {
	return "CSB"
}

func (p *CSBBANKParser) SupportsFile(filePath string) bool {
	return strings.HasSuffix(strings.ToLower(filePath), ".xls")
}
