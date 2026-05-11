package bankparsing

import (
	"fmt"
	"github.com/xuri/excelize/v2"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/tealeg/xlsx/v3"
)

type JSFBBANKParser struct{}

var _ BankParser = (*JSFBBANKParser)(nil)

func (p *JSFBBANKParser) Parse(fileContent string, filePath string, holderAccount string) ([]Transaction, error) {
	ext := strings.ToLower(filepath.Ext(filePath))
	if ext != ".xls" && ext != ".xlsx" {
		return nil, fmt.Errorf("unsupported file format for JSFB. Expected: .xls or .xlsx (Excel files), but got: %s", ext)
	}

	transactions, err := p.parseExcelFile(filePath)
	if err != nil {
		return nil, err
	}
	p.sortTransactionsByDate(transactions)

	return transactions, nil
}

func (p *JSFBBANKParser) sortTransactionsByDate(transactions []Transaction) {
	if len(transactions) == 0 {
		return
	}
	sort.Slice(transactions, func(i, j int) bool {
		dateI, errI := p.parseDateWithIndianTimezone(transactions[i].Date)
		dateJ, errJ := p.parseDateWithIndianTimezone(transactions[j].Date)

		if errI != nil || errJ != nil {
			return i > j
		}
		return dateI.After(dateJ)
	})
}

func (p *JSFBBANKParser) parseDateWithIndianTimezone(dateStr string) (time.Time, error) {
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

func (p *JSFBBANKParser) parseExcelFile(filePath string) ([]Transaction, error) {
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return nil, fmt.Errorf("file does not exist or cannot access: %v", err)
	}

	if fileInfo.Size() == 0 {
		return nil, fmt.Errorf("file is empty")
	}

	ext := strings.ToLower(filepath.Ext(filePath))
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("cannot open file: %v", err)
	}
	defer file.Close()

	header := make([]byte, 8)
	n, err := file.Read(header)
	if err != nil {
		return nil, fmt.Errorf("cannot read file header: %v", err)
	}
	fileType := p.detectFileType(header, n, ext)
	if fileType != "XLS" && fileType != "XLSX" {
		return nil, fmt.Errorf("unsupported file format: %s", fileType)
	}

	xlFile, err := xlsx.OpenFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open Excel file: %v", err)
	}
	sheetCount := len(xlFile.Sheets)
	if sheetCount == 0 {
		return nil, fmt.Errorf("no sheets found in file")
	}

	sheet := xlFile.Sheets[0]
	if sheet == nil {
		return nil, fmt.Errorf("failed to get first sheet")
	}
	var rows [][]string
	validRowCount := 0

	for rowIndex := 0; rowIndex <= sheet.MaxRow && rowIndex < 100; rowIndex++ {
		row, err := sheet.Row(rowIndex)
		if err != nil {
			continue
		}

		var cells []string
		hasData := false
		for colIndex := 0; colIndex < 50; colIndex++ {
			cell := row.GetCell(colIndex)
			cellValue := p.getCellValueSafely(cell)
			cells = append(cells, cellValue)
			if cellValue != "" {
				hasData = true
			}
		}

		if hasData {
			rows = append(rows, cells)
			validRowCount++
		}
	}

	if len(rows) == 0 {
		return p.parseWithExcelize(filePath)
	}
	return p.extractJSFBTransactions(rows)
}

func (p *JSFBBANKParser) parseWithExcelize(filePath string) ([]Transaction, error) {
	f, _ := excelize.OpenFile(filePath)
	defer f.Close()
	sheets := f.GetSheetList()

	sheetName := sheets[0]
	rows, _ := f.GetRows(sheetName)

	return p.extractJSFBTransactions(rows)
}

func (p *JSFBBANKParser) getCellValueSafely(cell *xlsx.Cell) string {
	if cell == nil {
		return ""
	}

	return strings.TrimSpace(cell.Value)
}

func (p *JSFBBANKParser) parseExcelFileAlternative(filePath string) ([]Transaction, error) {
	xlFile, err := xlsx.OpenFile(filePath)
	if err != nil {
		return nil, err
	}

	sheet := xlFile.Sheets[0]
	var rows [][]string

	err = sheet.ForEachRow(func(row *xlsx.Row) error {
		var cells []string

		err := row.ForEachCell(func(cell *xlsx.Cell) error {
			value, err := cell.FormattedValue()
			if err != nil {
				value = cell.Value
			}
			cells = append(cells, strings.TrimSpace(value))
			return nil
		})

		if err == nil && len(cells) > 0 {
			hasData := false
			for _, cell := range cells {
				if cell != "" {
					hasData = true
					break
				}
			}
			if hasData {
				rows = append(rows, cells)
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	if len(rows) == 0 {
		return nil, fmt.Errorf("no data found in file using alternative method")
	}

	return p.extractJSFBTransactions(rows)
}

func (p *JSFBBANKParser) findJSFBHeaderRow(rows [][]string) int {
	for i, row := range rows {
		if len(row) < 25 {
			continue
		}
		rowText := strings.ToUpper(strings.Join(row, " "))
		hasTxnDate := strings.Contains(rowText, "TXN DATE") || strings.Contains(rowText, "DATE")
		hasNarration := strings.Contains(rowText, "NARRATION")
		hasReference := strings.Contains(rowText, "REFERENCE")
		hasDeposit := strings.Contains(rowText, "DEPOSIT")
		hasWithdrawal := strings.Contains(rowText, "WITHDRAWAL")
		hasBalance := strings.Contains(rowText, "BALANCE")

		if hasTxnDate && hasNarration && hasReference && hasDeposit && hasWithdrawal && hasBalance {
			return i
		}
	}

	return -1
}

func (p *JSFBBANKParser) detectFileType(header []byte, headerLen int, ext string) string {
	if headerLen >= 8 {
		if header[0] == 0xD0 && header[1] == 0xCF && header[2] == 0x11 && header[3] == 0xE0 {
			return "XLS"
		}
		if header[0] == 0x50 && header[1] == 0x4B && header[2] == 0x03 && header[3] == 0x04 {
			return "XLSX"
		}
	}
	switch ext {
	case ".xls":
		return "XLS"
	case ".xlsx":
		return "XLSX"
	default:
		return "UNKNOWN"
	}
}

func (p *JSFBBANKParser) extractJSFBTransactions(rows [][]string) ([]Transaction, error) {
	var transactions []Transaction
	headerRowIndex := p.findJSFBHeaderRow(rows)
	if headerRowIndex == -1 {
		headerRowIndex = 49
	}

	startRow := headerRowIndex + 1
	foundTransactions := false

	for i := startRow; i < len(rows); i++ {
		row := rows[i]
		if len(row) < 5 {
			continue
		}
		if p.isFooterRow(row) {
			break
		}
		transaction, err := p.parseJSFBTransactionRow(row)
		if err != nil {
			continue
		}

		if transaction != nil {
			transactions = append(transactions, *transaction)
			foundTransactions = true
		}
	}

	if !foundTransactions {
		for _, row := range rows {
			if len(row) > 3 {
				dateStr := strings.TrimSpace(row[3])
				if p.isValidDate(dateStr) {
					transaction, err := p.parseJSFBTransactionRow(row)
					if err == nil && transaction != nil {
						transactions = append(transactions, *transaction)
					}
				}
			}
		}
	}

	if len(transactions) == 0 {
		return nil, fmt.Errorf("no transaction data parsed")
	}
	return transactions, nil
}

//func (p *JSFBBANKParser) findJSFBHeaderRow(rows [][]string) int {
//	for i, row := range rows {
//		if len(row) < 20 {
//			continue
//		}
//		rowText := strings.ToUpper(strings.Join(row, "|"))
//
//		// 根据实际表头查找
//		hasTxnDate := strings.Contains(rowText, "TXN DATE") ||
//			(strings.Contains(rowText, "DATE") && strings.Contains(rowText, "NARRATION"))
//		hasNarration := strings.Contains(rowText, "NARRATION")
//		hasDeposit := strings.Contains(rowText, "DEPOSITS")
//		hasWithdrawal := strings.Contains(rowText, "WITHDRAWAL")
//		hasBalance := strings.Contains(rowText, "BALANCE")
//
//		if hasTxnDate && hasNarration && hasDeposit && hasWithdrawal && hasBalance {
//			fmt.Printf("找到JSFB表头行: %d, 内容: %s\n", i, rowText)
//			return i
//		}
//	}
//	return -1
//}

func (p *JSFBBANKParser) parseJSFBTransactionRow(row []string) (*Transaction, error) {
	if len(row) < 24 {
		return nil, fmt.Errorf("insufficient row data: %d columns", len(row))
	}
	txnDate := ""
	if len(row) > 4 && p.isValidDate(strings.TrimSpace(row[4])) {
		txnDate = strings.TrimSpace(row[4])
	} else if len(row) > 5 && p.isValidDate(strings.TrimSpace(row[5])) {
		txnDate = strings.TrimSpace(row[5])
	}
	narration := ""
	for i := 6; i <= 10 && i < len(row); i++ {
		part := strings.TrimSpace(row[i])
		if part != "" {
			if narration != "" {
				narration += " "
			}
			narration += part
		}
	}
	reference := ""
	for i := 11; i <= 14 && i < len(row); i++ {
		part := strings.TrimSpace(row[i])
		if part != "" && strings.HasSuffix(part, "D") && p.isNumeric(strings.TrimSuffix(part, "D")) {
			reference = part
			break
		}
	}
	deposit := ""
	for i := 15; i <= 17 && i < len(row); i++ {
		part := strings.TrimSpace(row[i])
		if part != "" && p.isAmount(part) {
			deposit = part
			break
		}
	}
	if deposit == "" {
		deposit = "0.00"
	}
	withdrawal := ""
	for i := 18; i <= 22 && i < len(row); i++ {
		part := strings.TrimSpace(row[i])
		if part != "" && p.isAmount(part) {
			withdrawal = part
			break
		}
	}
	balance := ""
	if len(row) > 23 {
		balance = strings.TrimSpace(row[23])
	}
	if txnDate == "" || !p.isValidDate(txnDate) {
		return nil, fmt.Errorf("invalid or empty date: '%s'", txnDate)
	}
	bankTxnId := reference
	if bankTxnId != "" && len(bankTxnId) > 1 && strings.HasSuffix(bankTxnId, "D") {
		bankTxnId = bankTxnId[:len(bankTxnId)-1]
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
		return nil, fmt.Errorf("failed to parse amount '%s': %v", amountStr, err)
	}
	p.formatTransactionDateForOutput(txnDate)

	transaction := &Transaction{
		Date:        txnDate,
		ValueDate:   txnDate,
		Description: narration,
		Amount:      amount,
		Type:        transType,
		Balance:     balance,
		BankTxnId:   bankTxnId,
	}
	return transaction, nil
}

func (p *JSFBBANKParser) formatTransactionDateForOutput(dateStr string) string {
	return FormatTransactionDate(dateStr)
}

func (p *JSFBBANKParser) isAmount(s string) bool {
	s = strings.TrimSpace(s)
	if s == "" || s == "-" || s == "0" || s == "0.00" {
		return false
	}

	clean := strings.ReplaceAll(s, ",", "")
	clean = strings.ReplaceAll(clean, " ", "")
	clean = strings.ReplaceAll(clean, "₹", "")

	_, err := strconv.ParseFloat(clean, 64)
	return err == nil
}

func (p *JSFBBANKParser) analyzeJSFBRowData(row []string) (string, string, string, string) {
	txnDate := ""
	narration := ""
	withdrawal := ""
	deposit := ""

	if len(row) > 3 {
		txnDate = strings.TrimSpace(row[3])
	}
	if len(row) > 6 {
		narration = strings.TrimSpace(row[6])
	}
	if len(row) > 15 {
		deposit = strings.TrimSpace(row[15])
	}
	if len(row) > 17 {
		withdrawal = strings.TrimSpace(row[17])
	}

	return txnDate, narration, withdrawal, deposit
}

func (p *JSFBBANKParser) isFooterRow(row []string) bool {
	if len(row) == 0 {
		return true
	}

	rowText := strings.Join(row, " ")
	upperText := strings.ToUpper(rowText)

	return strings.Contains(upperText, "REGISTERED OFFICE") ||
		strings.Contains(upperText, "JANA SMALL FINANCE BANK") ||
		strings.Contains(rowText, "This is a computer-generated") ||
		strings.Contains(rowText, "customercare@jana.bank.in") ||
		strings.Contains(rowText, "1800-2080") ||
		strings.Contains(upperText, "COMMONLY USED NARRATIONS") ||
		strings.Contains(upperText, "AS PER RBI CIRCULAR") ||
		strings.Contains(upperText, "COMPUTER-GENERATED ADVICE")
}

func (p *JSFBBANKParser) cleanAmountString(amountStr string) string {
	amountStr = strings.TrimSpace(amountStr)
	amountStr = strings.ReplaceAll(amountStr, ",", "")
	amountStr = strings.ReplaceAll(amountStr, "Cr", "")
	amountStr = strings.ReplaceAll(amountStr, "Dr", "")
	amountStr = strings.ReplaceAll(amountStr, " ", "")
	amountStr = strings.ReplaceAll(amountStr, "₹", "")
	return amountStr
}

func (p *JSFBBANKParser) isValidDate(dateStr string) bool {
	if dateStr == "" {
		return false
	}

	dateStr = strings.TrimSpace(dateStr)

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

func (p *JSFBBANKParser) parseFloat(s string) (float64, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return 0, fmt.Errorf("empty string")
	}
	s = strings.ReplaceAll(s, ",", "")
	s = strings.ReplaceAll(s, "₹", "")

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

func (p *JSFBBANKParser) ConvertToTransInfo(transactions []Transaction, holderAccount string) []TransInfo {
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
		transName := p.extractNameFromDescription(txn.Description)
		transAccount := p.extractAccountFromNarration(txn.Description)
		bankTxnId := txn.BankTxnId
		if bankTxnId == "" {
			bankTxnId = upiRef
		}

		transInfo := TransInfo{
			TransType:    txn.Type,
			TransName:    transName,
			TransAccount: transAccount,
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

func (p *JSFBBANKParser) extractUPIInfo(narration string) (string, string) {
	if !strings.Contains(narration, "UPI") {
		return "", ""
	}

	parts := strings.Split(narration, "/")
	upiRef := ""
	upiTime := ""

	for _, part := range parts {
		part = strings.TrimSpace(part)
		if upiRef == "" && len(part) >= 10 && len(part) <= 12 && p.isNumeric(part) {
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

func (p *JSFBBANKParser) isValidTimeFormat(timeStr string) bool {
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

func (p *JSFBBANKParser) extractNameFromDescription(description string) string {
	if description == "" {
		return "Bank Transaction"
	}

	if strings.Contains(description, "UPI") {
		parts := strings.Split(description, "/")
		if len(parts) >= 4 {
			name := parts[3]
			if name != "" && !p.isNumeric(name) {
				return name
			}
		}
		return "UPI Transaction"
	}

	return "Bank Transaction"
}

func (p *JSFBBANKParser) extractAccountFromNarration(narration string) string {
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
	return ""
}

func (p *JSFBBANKParser) isNumeric(s string) bool {
	_, err := strconv.ParseFloat(s, 64)
	return err == nil
}

func (p *JSFBBANKParser) GetBankName() string {
	return "JSFBBANK"
}

func (p *JSFBBANKParser) extractDatePart(dateStr string) string {
	return extractDatePart(dateStr)
}
