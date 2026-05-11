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

type BOMBANKParser struct{}

var _ BankParser = (*BOMBANKParser)(nil)

func (p *BOMBANKParser) Parse(fileContent string, filePath string, holderAccount string) ([]Transaction, error) {
	ext := strings.ToLower(filepath.Ext(filePath))
	if ext != ".xls" {
		return nil, fmt.Errorf("unsupported file format for BOM. Expected: .xls (Excel 97-2003 format), but got: %s", ext)
	}

	transactions, err := p.parseXLSFile(filePath)
	if err != nil {
		return nil, err
	}
	p.sortTransactionsByDate(transactions)

	return transactions, nil
}

func (p *BOMBANKParser) parseXLSFile(filePath string) ([]Transaction, error) {
	file, err := xls.Open(filePath, "utf-8")
	if err != nil {
		return nil, fmt.Errorf("failed to open XLS file: %v", err)
	}

	sheet := file.GetSheet(0)
	if sheet == nil {
		return nil, fmt.Errorf("no sheets found in XLS file")
	}

	var transactions []Transaction
	targetRows := []int{26, 27, 28, 29, 30, 31, 32, 33}

	for _, rowIndex := range targetRows {
		func() {
			defer func() {
				if r := recover(); r != nil {
					fmt.Printf("Recovered from panic at row %d: %v\n", rowIndex, r)
				}
			}()

			row := sheet.Row(rowIndex)
			if row == nil {
				return
			}

			var cells []string
			for j := 0; j < 8; j++ {
				cell := strings.TrimSpace(row.Col(j))
				cells = append(cells, cell)
			}

			// 解析交易行
			if len(cells) >= 7 {
				transaction, err := p.parseBOMTransactionRow(cells)
				if err == nil && transaction != nil {
					transactions = append(transactions, *transaction)
				}
			}
		}()
	}

	if len(transactions) == 0 {
		return nil, fmt.Errorf("no transactions found in expected rows")
	}

	return transactions, nil
}

func (p *BOMBANKParser) extractBOMTransactions(rows [][]string) ([]Transaction, error) {
	var transactions []Transaction
	headerRowIndex := p.findBOMHeaderRow(rows)

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
			fmt.Printf("Found footer at row %d, stopping\n", i)
			break
		}
		transaction, err := p.parseBOMTransactionRow(row)
		if err != nil {
			fmt.Printf("Failed to parse row %d: %v\n", i, err)
			continue
		}

		if transaction != nil {
			fmt.Printf("Found transaction at row %d: %s %s %.2f\n", i, transaction.Date, transaction.Type, transaction.Amount)
			transactions = append(transactions, *transaction)
		}
	}

	if len(transactions) == 0 {
		return nil, fmt.Errorf("no transaction data found")
	}

	fmt.Printf("Successfully parsed %d transactions\n", len(transactions))
	return transactions, nil
}

func (p *BOMBANKParser) findBOMHeaderRow(rows [][]string) int {
	for i, row := range rows {
		if len(row) < 7 {
			continue
		}

		// 检查是否包含BOM银行对账单的表头列
		rowText := strings.ToUpper(strings.Join(row, " "))
		hasDate := strings.Contains(rowText, "DATE")
		hasType := strings.Contains(rowText, "TYPE")
		hasParticulars := strings.Contains(rowText, "PARTICULARS")
		hasDebit := strings.Contains(rowText, "DEBIT")
		hasCredit := strings.Contains(rowText, "CREDIT")
		headerCount := 0
		if hasDate {
			headerCount++
		}
		if hasType {
			headerCount++
		}
		if hasParticulars {
			headerCount++
		}
		if hasDebit {
			headerCount++
		}
		if hasCredit {
			headerCount++
		}

		if headerCount >= 4 {
			fmt.Printf("Found BOM header at row %d: %s\n", i, rowText)
			return i
		}
	}
	for i, row := range rows {
		if len(row) >= 7 {
			dateCell := strings.TrimSpace(row[0])
			if p.isValidDate(dateCell) {
				fmt.Printf("Found transaction data starting at row %d\n", i)
				return i - 1
			}
		}
	}

	return -1
}

func (p *BOMBANKParser) parseBOMTransactionRow(cells []string) (*Transaction, error) {
	if len(cells) < 7 {
		return nil, fmt.Errorf("insufficient columns: %d", len(cells))
	}
	date := strings.TrimSpace(cells[0])
	transType := strings.TrimSpace(cells[1])
	particulars := strings.TrimSpace(cells[2])
	chequeRef := strings.TrimSpace(cells[3])
	debit := strings.TrimSpace(cells[4])
	credit := strings.TrimSpace(cells[5])
	balance := strings.TrimSpace(cells[6])
	channel := ""
	if len(cells) > 7 {
		channel = strings.TrimSpace(cells[7])
	}

	if date == "" || !p.isValidDate(date) {
		return nil, fmt.Errorf("invalid date: %s", date)
	}
	amountStr := ""
	transDirection := ""

	if debit != "" && debit != "0" && debit != "-" && debit != "0.00" {
		amountStr = debit
		transDirection = "TYPE_OUT"
	} else if credit != "" && credit != "0" && credit != "-" && credit != "0.00" {
		amountStr = credit
		transDirection = "TYPE_IN"
	} else {
		return nil, fmt.Errorf("no valid amount")
	}
	amountStr = p.cleanAmountString(amountStr)
	amount, err := p.parseFloat(amountStr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse amount: %v", err)
	}

	description := particulars
	if transType != "" && !strings.Contains(description, transType) {
		description = transType + " - " + description
	}

	transaction := &Transaction{
		Date:        date,
		Description: description,
		Amount:      amount,
		Type:        transDirection,
		ChequeNo:    chequeRef,
		Balance:     balance,
		Channel:     channel,
		BankTxnId:   chequeRef,
	}

	return transaction, nil
}

func (p *BOMBANKParser) isFooterRow(cells []string) bool {
	if len(cells) == 0 {
		return false
	}

	rowText := strings.Join(cells, " ")
	return strings.Contains(strings.ToUpper(rowText), "ALL THE AMOUNTS") ||
		strings.Contains(rowText, "Summary for Account") ||
		strings.Contains(rowText, "Total Transaction Count") ||
		strings.Contains(rowText, "Opening Balance") ||
		strings.Contains(rowText, "Closing Balance") ||
		strings.Contains(rowText, "This is a System Generated Statement")
}

func (p *BOMBANKParser) isAmount(s string) bool {
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

func (p *BOMBANKParser) cleanAmountString(amountStr string) string {
	amountStr = strings.TrimSpace(amountStr)
	amountStr = strings.ReplaceAll(amountStr, ",", "")
	amountStr = strings.ReplaceAll(amountStr, "Cr", "")
	amountStr = strings.ReplaceAll(amountStr, "Dr", "")
	amountStr = strings.ReplaceAll(amountStr, " ", "")
	return amountStr
}

func (p *BOMBANKParser) isValidDate(dateStr string) bool {
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

func (p *BOMBANKParser) parseFloat(s string) (float64, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return 0, fmt.Errorf("empty string")
	}

	s = strings.ReplaceAll(s, ",", "")

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

func (p *BOMBANKParser) sortTransactionsByDate(transactions []Transaction) {
	sort.Slice(transactions, func(i, j int) bool {
		dateI, errI := p.parseDateWithIndianTimezone(transactions[i].Date)
		dateJ, errJ := p.parseDateWithIndianTimezone(transactions[j].Date)

		if errI != nil || errJ != nil {
			return i > j
		}

		return dateI.After(dateJ)
	})
}

func (p *BOMBANKParser) parseDateWithIndianTimezone(dateStr string) (time.Time, error) {
	loc, err := time.LoadLocation("Asia/Kolkata")
	if err != nil {
		loc = time.FixedZone("IST", 5*60*60+30*60)
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

func (p *BOMBANKParser) ConvertToTransInfo(transactions []Transaction, holderAccount string) []TransInfo {
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
		bankTxnId := txn.BankTxnId
		if bankTxnId == "" {
			bankTxnId = upiRef
		}

		transInfo := TransInfo{
			TransType:    txn.Type,
			TransName:    p.extractNameFromDescription(txn.Description),
			TransAccount: p.extractAccountFromNarration(txn.Description),
			TransUpistr:  "",
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

func (p *BOMBANKParser) formatTransactionDate(dateStr string) string {
	dateStr = strings.TrimSpace(dateStr)
	if dateStr == "" {
		return FormatDateWithIndianTime("")
	}

	formats := []string{
		"02/01/2006",
		"02-01-2006",
		"2006-01-02",
		"02/01/06",
	}

	for _, format := range formats {
		if t, err := time.Parse(format, dateStr); err == nil {
			return FormatDateWithIndianTime(t.Format("2006-01-02"))
		}
	}
	return FormatDateWithIndianTime("")
}

func (p *BOMBANKParser) extractDatePart(dateStr string) string {
	return extractDatePart(dateStr)
}

func (p *BOMBANKParser) extractUPIInfo(narration string) (string, string) {
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

func (p *BOMBANKParser) isValidTimeFormat(timeStr string) bool {
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

func (p *BOMBANKParser) extractNameFromDescription(description string) string {
	if description == "" {
		return "Bank Transaction"
	}

	if strings.Contains(description, "UPI") {
		parts := strings.Split(description, "/")
		for i, part := range parts {
			part = strings.TrimSpace(part)
			if strings.HasPrefix(part, "MAHB") || strings.HasPrefix(part, "KKBK") {
				if i+1 < len(parts) {
					namePart := strings.TrimSpace(parts[i+1])
					if namePart != "" && !strings.Contains(namePart, "@") && !p.isNumeric(namePart) {
						namePart = strings.TrimPrefix(namePart, "Mr ")
						namePart = strings.TrimPrefix(namePart, "Mrs ")
						namePart = strings.TrimPrefix(namePart, "Ms ")
						namePart = strings.TrimPrefix(namePart, "Shri ")
						return namePart
					}
				}
			}
		}
		return "UPI Transaction"
	}

	if len(description) > 30 {
		return description[:30] + "..."
	}
	return description
}

func (p *BOMBANKParser) extractAccountFromNarration(narration string) string {
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

func (p *BOMBANKParser) isNumeric(s string) bool {
	_, err := strconv.ParseFloat(s, 64)
	return err == nil
}

func (p *BOMBANKParser) GetBankName() string {
	return "BOMBANK"
}

func (p *BOMBANKParser) SupportsFile(filePath string) bool {
	return strings.HasSuffix(strings.ToLower(filePath), ".xls")
}
