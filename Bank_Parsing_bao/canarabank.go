package bankparsing

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type CANARABANKParser struct{}

func (p *CANARABANKParser) Parse(fileContent string, filePath string, holderAccount string) ([]Transaction, error) {
	ext := strings.ToLower(filepath.Ext(filePath))
	if ext != ".csv" {
		return nil, fmt.Errorf("unsupported file format for CANARA. Expected: .csv (CSV file), but got: %s", ext)
	}

	records, err := p.ParseCSVToRecords(filePath)
	if err != nil {
		return nil, err
	}
	var transactions []Transaction
	for _, record := range records {
		var amount float64
		var transType string
		if record.Debit != "" && record.Debit != "0" && record.Debit != "0.00" {
			cleanedAmount := p.normalizeAmount(record.Debit)
			amount, err = strconv.ParseFloat(cleanedAmount, 64)
			if err != nil {
				continue
			}
			transType = "TYPE_OUT"
		} else if record.Credit != "" && record.Credit != "0" && record.Credit != "0.00" {
			cleanedAmount := p.normalizeAmount(record.Credit)
			amount, err = strconv.ParseFloat(cleanedAmount, 64)
			if err != nil {
				continue
			}
			transType = "TYPE_IN"
		} else {
			continue
		}
		transactions = append(transactions, Transaction{
			Date:        p.formatDate(record.Date),
			ValueDate:   p.formatValueDate(record.ValueDate),
			Description: record.Narration,
			Amount:      amount,
			Type:        transType,
			Account:     holderAccount,
			ChequeNo:    record.ChqNo,
		})
	}
	return transactions, nil
}

func (p *CANARABANKParser) formatDateForTransInfo(dateStr string) string {
	if dateStr == "" {
		return ""
	}
	dateStr = p.cleanExcelFormula(dateStr)
	formats := []string{
		"2006-01-02",
		"2006-01-02 15:04:05",
		"02-01-2006 15:04:05",
		"02-01-2006",
		"02/01/2006 15:04:05",
		"02/01/2006",
	}
	for _, format := range formats {
		if t, err := time.Parse(format, dateStr); err == nil {
			return t.Format("2006-01-02")
		}
	}
	return dateStr
}

func (p *CANARABANKParser) ConvertToTransInfo(transactions []Transaction, holderAccount string) []TransInfo {
	var transInfos []TransInfo
	for _, transaction := range transactions {
		upiRef, upiTime := p.extractUPIInfo(transaction.Description)
		datePart := p.formatDateForTransInfo(transaction.Date)
		transDate := datePart
		if upiTime != "" {
			transDate = datePart + " " + upiTime
		} else {
			transDate = datePart + " 00:00:00"
		}

		// 设置 FundFlow 值
		var fundFlow int32
		if transaction.Type == "TYPE_OUT" {
			fundFlow = 1 // debit去向
		} else if transaction.Type == "TYPE_IN" {
			fundFlow = 2 // credit来源
		}

		transInfo := TransInfo{
			TransType:    transaction.Type,
			TransName:    p.extractNameFromDescription(transaction.Description),
			TransAccount: p.extractAccountFromNarration(transaction.Description),
			TransUpistr:  transaction.Description,
			TransAmount:  fmt.Sprintf("%.2f", transaction.Amount),
			BankTxnId:    upiRef,
			TransDate:    transDate,
			TransStatus:  "SUCCESS",
			FundFlow:     fundFlow, // 添加 FundFlow 字段
		}
		transInfos = append(transInfos, transInfo)
	}

	return transInfos
}

func (p *CANARABANKParser) extractUPIInfo(narration string) (string, string) {
	if !strings.Contains(narration, "UPI") {
		return "", ""
	}

	parts := strings.Split(narration, "/")
	upiRef := ""
	upiTime := ""
	for i, part := range parts {
		part = strings.TrimSpace(part)
		if upiRef == "" && len(part) >= 10 && p.isNumeric(part) {
			upiRef = part
		}
		if upiTime == "" && p.isValidTimestampFormat(part) {
			if len(part) >= 19 {
				timePart := part[11:19]
				if p.isValidTimeFormat(timePart) {
					upiTime = timePart
					break
				}
			}
		}
		if upiTime == "" && i < len(parts)-1 {
			nextPart := strings.TrimSpace(parts[i+1])
			if p.isValidTimestampFormat(nextPart) {
				if len(nextPart) >= 19 {
					timePart := nextPart[11:19]
					if p.isValidTimeFormat(timePart) {
						upiTime = timePart
						break
					}
				}
			}
		}
	}
	if upiTime == "" {
		timePattern := `(\d{2}:\d{2}:\d{2})`
		if matches := regexp.MustCompile(timePattern).FindStringSubmatch(narration); len(matches) > 1 {
			timeStr := matches[1]
			if p.isValidTimeFormat(timeStr) {
				upiTime = timeStr
			}
		}
	}

	return upiRef, upiTime
}

func (p *CANARABANKParser) isValidTimestampFormat(timestampStr string) bool {
	if len(timestampStr) != 19 {
		return false
	}
	if timestampStr[2] != '/' || timestampStr[5] != '/' || timestampStr[10] != ' ' ||
		timestampStr[13] != ':' || timestampStr[16] != ':' {
		return false
	}
	datePart := timestampStr[0:10]
	timePart := timestampStr[11:19]
	return p.isValidDateFormat(datePart) && p.isValidTimeFormat(timePart)
}

func (p *CANARABANKParser) isValidDateFormat(dateStr string) bool {
	if len(dateStr) != 10 {
		return false
	}
	dayStr := dateStr[0:2]
	monthStr := dateStr[3:5]
	yearStr := dateStr[6:10]

	if !p.isNumeric(dayStr) || !p.isNumeric(monthStr) || !p.isNumeric(yearStr) {
		return false
	}
	day, _ := strconv.Atoi(dayStr)
	month, _ := strconv.Atoi(monthStr)
	year, _ := strconv.Atoi(yearStr)
	return day >= 1 && day <= 31 && month >= 1 && month <= 12 && year >= 2000 && year <= 2100
}

func (p *CANARABANKParser) isValidTimeFormat(timeStr string) bool {
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

func (p *CANARABANKParser) extractNameFromDescription(description string) string {
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
	if strings.Contains(description, "TO ") {
		parts := strings.Split(description, "TO ")
		if len(parts) >= 2 {
			return strings.TrimSpace(parts[1])
		}
	}
	return "Bank Transaction"
}

func (p *CANARABANKParser) extractAccountFromNarration(narration string) string {
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

func (p *CANARABANKParser) isNumeric(s string) bool {
	_, err := strconv.ParseFloat(s, 64)
	return err == nil
}

func (p *CANARABANKParser) normalizeAmount(amountStr string) string {
	cleaned := strings.ReplaceAll(amountStr, ",", "")
	cleaned = strings.ReplaceAll(cleaned, "Rs.", "")
	cleaned = strings.ReplaceAll(cleaned, " ", "")
	cleaned = strings.TrimSpace(cleaned)
	if strings.HasPrefix(cleaned, "=") && strings.Contains(cleaned, `"`) {
		cleaned = strings.ReplaceAll(cleaned, `"`, "")
		cleaned = strings.TrimPrefix(cleaned, "=")
	}
	return cleaned
}

func (p *CANARABANKParser) ParseCSVToRecords(filePath string) ([]CsvRecord, error) {
	csvContent, err := p.processCanaraBankCSVFile(filePath)
	if err != nil {
		return nil, err
	}
	return p.parseCanaraBankCSVToRecords(csvContent)
}

func (p *CANARABANKParser) processCanaraBankCSVFile(filePath string) (string, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("读取文件失败: %v", err)
	}
	lines := strings.Split(string(content), "\n")
	var transactionLines []string
	inTransactionSection := false
	for _, line := range lines {
		trimmedLine := strings.TrimSpace(line)
		if strings.Contains(trimmedLine, "Txn Date,Value Date,Cheque No.,Description,Branch Code,Debit,Credit,Balance,") {
			inTransactionSection = true

			transactionLines = append(transactionLines, trimmedLine)
			continue
		}
		if inTransactionSection {
			if trimmedLine == "" || strings.Contains(trimmedLine, "Legends Used") {
				break
			}
			transactionLines = append(transactionLines, trimmedLine)
		}
	}
	if len(transactionLines) == 0 {
		return "", fmt.Errorf("未找到交易数据")
	}
	var result strings.Builder
	for _, line := range transactionLines {
		result.WriteString(line + "\n")
	}

	processedContent := result.String()
	return processedContent, nil
}

func (p *CANARABANKParser) parseCanaraBankCSVToRecords(csvContent string) ([]CsvRecord, error) {
	reader := csv.NewReader(strings.NewReader(csvContent))
	reader.FieldsPerRecord = -1
	reader.LazyQuotes = true
	reader.TrimLeadingSpace = true

	records, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("读取CSV内容失败: %v", err)
	}

	var result []CsvRecord
	headerRow := -1
	for i, record := range records {
		if len(record) >= 8 {
			rowText := strings.ToUpper(strings.Join(record, " "))
			if strings.Contains(rowText, "TXN DATE") &&
				strings.Contains(rowText, "VALUE DATE") &&
				strings.Contains(rowText, "DEBIT") &&
				strings.Contains(rowText, "CREDIT") {
				headerRow = i
				break
			}
		}
	}

	if headerRow == -1 {
		return nil, fmt.Errorf("未找到表头行")
	}
	for i := headerRow + 1; i < len(records); i++ {
		if len(records[i]) < 8 {
			continue
		}
		record := CsvRecord{
			Date:      strings.TrimSpace(records[i][0]),
			ValueDate: strings.TrimSpace(records[i][1]),
			ChqNo:     strings.TrimSpace(records[i][2]),
			Narration: strings.TrimSpace(records[i][3]),
			Debit:     strings.TrimSpace(records[i][5]),
			Credit:    strings.TrimSpace(records[i][6]),
			Balance:   strings.TrimSpace(records[i][7]),
		}
		if record.Date == "" || (record.Debit == "" && record.Credit == "") {
			continue
		}
		record.Date = p.cleanExcelFormula(record.Date)
		record.ValueDate = p.cleanExcelFormula(record.ValueDate)
		record.ChqNo = p.cleanExcelFormula(record.ChqNo)
		record.Debit = p.cleanExcelFormula(record.Debit)
		record.Credit = p.cleanExcelFormula(record.Credit)

		result = append(result, record)
	}

	return result, nil
}

func (p *CANARABANKParser) cleanExcelFormula(s string) string {
	s = strings.TrimSpace(s)
	if strings.HasPrefix(s, "=") && strings.Contains(s, `"`) {
		s = strings.ReplaceAll(s, `"`, "")
		s = strings.TrimPrefix(s, "=")
	}
	return s
}

func (p *CANARABANKParser) formatDate(dateStr string) string {
	if dateStr == "" {
		return ""
	}

	dateStr = p.cleanExcelFormula(dateStr)
	formats := []string{
		"02-01-2006 15:04:05",
		"02-01-2006",
		"02/01/2006 15:04:05",
		"02/01/2006",
	}
	for _, format := range formats {
		if t, err := time.Parse(format, dateStr); err == nil {
			return t.Format("2006-01-02")
		}
	}

	return dateStr
}

func (p *CANARABANKParser) formatValueDate(dateStr string) string {
	if dateStr == "" {
		return ""
	}
	dateStr = p.cleanExcelFormula(dateStr)
	formats := []string{
		"02 Jan 2006",
		"02-Jan-2006",
		"02/01/2006",
	}

	for _, format := range formats {
		if t, err := time.Parse(format, dateStr); err == nil {
			return t.Format("2006-01-02")
		}
	}

	return dateStr
}

func (p *CANARABANKParser) GetBankName() string {
	return "CANARA"
}

func ParseCANARABANKFile(filePath string, holderAccount string) (*BankResponse, error) {
	parser := &CANARABANKParser{}
	transactions, err := parser.Parse("", filePath, holderAccount)
	if err != nil {
		return nil, fmt.Errorf("解析Canara银行文件失败: %v", err)
	}
	transInfos := parser.ConvertToTransInfo(transactions, holderAccount)
	response := &BankResponse{
		HolderAccount: holderAccount,
		TransInfo:     transInfos,
	}

	log.Printf("成功解析 %d 笔Canara银行交易", len(transInfos))
	return response, nil
}
