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

type KVBBANKParser struct{}

var _ BankParser = (*KVBBANKParser)(nil)

func (p *KVBBANKParser) Parse(fileContent string, filePath string, holderAccount string) ([]Transaction, error) {
	ext := strings.ToLower(filepath.Ext(filePath))
	if ext != ".csv" {
		return nil, fmt.Errorf("unsupported file format for KVB. Expected: .csv (CSV file), but got: %s", ext)
	}

	records, err := p.ParseCSVToRecords(filePath)
	if err != nil {
		return nil, err
	}

	var transactions []Transaction

	for _, record := range records {
		var amount float64
		var transType string
		if record.TransactionDate == "" || record.Description == "" {
			continue
		}
		if strings.Contains(record.Description, "BALANCE") ||
			strings.Contains(record.Description, "TRANSACTION CHARGES") {
			continue
		}

		if record.Debit != "" && record.Debit != "0" && record.Debit != "0.00" {
			cleanedAmount := p.normalizeAmount(record.Debit)
			amount, err = strconv.ParseFloat(cleanedAmount, 64)
			if err != nil {
				log.Printf("Failed to parse the debit amount: %s, err: %v", record.Debit, err)
				continue
			}
			transType = "DEBIT"
			amount = -amount
		} else if record.Credit != "" && record.Credit != "0" && record.Credit != "0.00" {
			cleanedAmount := p.normalizeAmount(record.Credit)
			amount, err = strconv.ParseFloat(cleanedAmount, 64)
			if err != nil {
				log.Printf("Failed to parse the credit amount: %s, err: %v", record.Credit, err)
				continue
			}
			transType = "CREDIT"
		} else {
			continue
		}

		transactions = append(transactions, Transaction{
			Date:        p.formatDate(record.TransactionDate),
			Description: record.Description,
			Amount:      amount,
			Type:        transType,
			Account:     holderAccount,
		})
	}

	log.Printf("Successfully parsed %d KV bank transactions", len(transactions))
	return transactions, nil
}

func (p *KVBBANKParser) ConvertToTransInfo(transactions []Transaction, holderAccount string) []TransInfo {
	var transInfos []TransInfo

	for _, transaction := range transactions {
		transType := "TYPE_IN"
		if transaction.Type == "DEBIT" {
			transType = "TYPE_OUT"
		}

		// 设置 FundFlow 值
		var fundFlow int32
		if transaction.Type == "DEBIT" {
			fundFlow = 1 // debit去向
		} else if transaction.Type == "CREDIT" {
			fundFlow = 2 // credit来源
		}

		transName := p.extractName(transaction.Description)
		transUpistr := p.extractUpi(transaction.Description)
		bankTxnId := p.extractTxnId(transaction.Description)

		transAmount := fmt.Sprintf("%.2f", transaction.Amount)
		if transaction.Amount < 0 {
			transAmount = fmt.Sprintf("%.2f", -transaction.Amount)
		}
		transDate := FormatDateWithIndianTime(transaction.Date)

		transInfos = append(transInfos, TransInfo{
			TransType:    transType,
			TransName:    transName,
			TransAccount: holderAccount,
			TransUpistr:  transUpistr,
			TransAmount:  transAmount,
			BankTxnId:    bankTxnId,
			TransDate:    transDate,
			TransStatus:  "SUCCESS",
			FundFlow:     fundFlow, // 添加 FundFlow 字段
		})
	}

	return transInfos
}

func (p *KVBBANKParser) extractName(description string) string {
	// UPI 交易名称提取
	if strings.Contains(description, "UPI-") {
		parts := strings.Split(description, "-")
		if len(parts) >= 4 {
			return strings.TrimSpace(parts[3])
		}
	}

	// IMPS 交易名称提取
	if strings.Contains(description, "IMPS-") {
		parts := strings.Split(description, "-")
		if len(parts) >= 3 {
			namePart := strings.Split(parts[2], " ")[0]
			return strings.TrimSpace(namePart)
		}
	}

	// NEFT 交易名称提取
	if strings.Contains(description, "NEFT-") {
		parts := strings.Split(description, "-")
		if len(parts) >= 2 {
			return "NEFT_TRANSACTION"
		}
	}

	// CASH DEP 交易
	if strings.Contains(description, "CASH DEP-") {
		return "CASH_DEPOSIT"
	}

	// KVBLH 交易名称提取
	if strings.Contains(description, "KVBLH") {
		parts := strings.Split(description, "-")
		if len(parts) >= 3 {
			return strings.TrimSpace(parts[2])
		}
	}

	// 尝试从描述中提取名称
	if strings.Contains(description, "-") {
		parts := strings.Split(description, "-")
		lastPart := parts[len(parts)-1]
		if !strings.Contains(lastPart, "BANK") && !strings.Contains(lastPart, "UPI") {
			return strings.TrimSpace(lastPart)
		}
	}

	return ""
}

func (p *KVBBANKParser) extractUpi(description string) string {
	if strings.Contains(description, "UPI-") {
		return "UPI_TRANSACTION"
	}
	if strings.Contains(description, "IMPS-") {
		return "IMPS_TRANSACTION"
	}
	if strings.Contains(description, "NEFT-") {
		return "NEFT_TRANSACTION"
	}
	if strings.Contains(description, "CASH DEP-") {
		return "CASH_TRANSACTION"
	}
	return ""
}

func (p *KVBBANKParser) extractTxnId(description string) string {
	if strings.Contains(description, "UPI-") {
		parts := strings.Split(description, "-")
		if len(parts) >= 3 {
			upiId := strings.TrimSpace(parts[2])
			if p.isPureDigits(upiId) {
				return upiId
			}
		}
	}
	if strings.Contains(description, "IMPS-") {
		re := regexp.MustCompile(`\d{8,}`)
		matches := re.FindAllString(description, -1)
		if len(matches) > 0 {
			impsId := matches[0]
			if p.isPureDigits(impsId) {
				parts := strings.Split(description, "-")
				if len(parts) >= 2 {
					lastPart := strings.TrimSpace(parts[len(parts)-1])
					if lastPart != "" && !p.isBankCode(lastPart) && !p.isPureDigits(lastPart) {
						return impsId + "|+|" + lastPart
					}
				}
				return impsId
			}
		}
		parts := strings.Split(description, "-")
		if len(parts) >= 2 {
			candidate := strings.TrimSpace(parts[1])
			if p.isPureDigits(candidate) {
				return candidate
			}
		}
	}

	if strings.Contains(description, "NEFT-") {
		parts := strings.Split(description, "-")
		if len(parts) >= 2 {
			neftId := strings.TrimSpace(parts[1])
			if p.isPureDigits(neftId) {
				return neftId
			}
		}
	}

	// KVBLH 交易ID提取
	if strings.Contains(description, "KVBLH") {
		parts := strings.Split(description, "-")
		if len(parts) >= 2 {
			kvbId := strings.TrimSpace(parts[1])
			if p.isPureDigits(kvbId) {
				return kvbId
			}
		}
	}

	// 支票号提取
	if strings.Contains(description, "Cheque:") {
		parts := strings.Split(description, "Cheque:")
		if len(parts) >= 2 {
			chequeNo := strings.TrimSpace(parts[1])
			if p.isPureDigits(chequeNo) {
				return chequeNo
			}
		}
	}

	return ""
}

func (p *KVBBANKParser) isPureDigits(str string) bool {
	if str == "" {
		return false
	}
	for _, char := range str {
		if char < '0' || char > '9' {
			return false
		}
	}
	return true
}

func (p *KVBBANKParser) isBankCode(str string) bool {
	bankCodes := []string{"SBIN", "CNRB", "FDRL", "ICIC", "HDFC", "UTIB", "BDBL", "CRGB", "INDB"}
	for _, code := range bankCodes {
		if strings.Contains(strings.ToUpper(str), code) {
			return true
		}
	}
	return false
}

func (p *KVBBANKParser) normalizeAmount(amountStr string) string {
	cleaned := strings.ReplaceAll(amountStr, ",", "")
	cleaned = strings.ReplaceAll(cleaned, " ", "")
	cleaned = strings.ReplaceAll(cleaned, "\t", "")
	cleaned = strings.TrimSpace(cleaned)
	if strings.Contains(cleaned, `"`) {
		cleaned = strings.ReplaceAll(cleaned, `"`, "")
	}
	if strings.HasPrefix(cleaned, "-") {
		cleaned = strings.TrimPrefix(cleaned, "-")
	}

	return cleaned
}

func (p *KVBBANKParser) ParseCSVToRecords(filePath string) ([]KVBRecord, error) {
	csvContent, err := p.processKVBANKCSVFile(filePath)
	if err != nil {
		return nil, err
	}
	return p.parseKVBANKCSVToRecords(csvContent)
}

func (p *KVBBANKParser) processKVBANKCSVFile(filePath string) (string, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("Failed to read the file: %v", err)
	}

	lines := strings.Split(string(content), "\n")
	var transactionLines []string
	inTransactionSection := false
	headerFound := false

	for _, line := range lines {
		trimmedLine := strings.TrimSpace(line)
		if strings.HasPrefix(trimmedLine, "Transaction Date,Value Date,Branch,Cheque No.,Description,Debit,Credit,Balance") {
			inTransactionSection = true
			headerFound = true
			transactionLines = append(transactionLines, trimmedLine)
			continue
		}

		if inTransactionSection && headerFound && trimmedLine != "" {

			if strings.Contains(trimmedLine, ",,,") && len(strings.Split(trimmedLine, ",")) < 5 {
				continue
			}
			transactionLines = append(transactionLines, trimmedLine)
		}

		if strings.Contains(trimmedLine, "Account Number:") {
			continue
		}
	}

	if len(transactionLines) <= 1 {
		return "", fmt.Errorf("No transaction data found")
	}

	var result strings.Builder
	result.WriteString("Transaction Date,Value Date,Branch,Cheque No.,Description,Debit,Credit,Balance\n")

	for i := 1; i < len(transactionLines); i++ {
		line := transactionLines[i]
		if line != "" {
			result.WriteString(line + "\n")
		}
	}

	processedContent := result.String()
	return processedContent, nil
}

type KVBRecord struct {
	TransactionDate string
	ValueDate       string
	Branch          string
	ChequeNo        string
	Description     string
	Debit           string
	Credit          string
	Balance         string
}

func (p *KVBBANKParser) parseKVBANKCSVToRecords(csvContent string) ([]KVBRecord, error) {
	reader := csv.NewReader(strings.NewReader(csvContent))
	reader.FieldsPerRecord = -1
	reader.LazyQuotes = true
	reader.TrimLeadingSpace = true

	records, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("Failed to read the CSV content: %v", err)
	}

	var result []KVBRecord
	for i := 0; i < len(records); i++ {
		if len(records[i]) < 8 {
			continue
		}

		record := KVBRecord{
			TransactionDate: strings.TrimSpace(records[i][0]),
			ValueDate:       strings.TrimSpace(records[i][1]),
			Branch:          strings.TrimSpace(records[i][2]),
			ChequeNo:        strings.TrimSpace(records[i][3]),
			Description:     strings.TrimSpace(records[i][4]),
			Debit:           strings.TrimSpace(records[i][5]),
			Credit:          strings.TrimSpace(records[i][6]),
			Balance:         strings.TrimSpace(records[i][7]),
		}

		if record.TransactionDate == "" || record.Description == "" {
			continue
		}
		if strings.Contains(record.Description, "Opening Balance") ||
			strings.Contains(record.Description, "Closing Balance") {
			continue
		}

		result = append(result, record)
	}
	return result, nil
}

func (p *KVBBANKParser) formatDate(dateStr string) string {
	if dateStr == "" {
		return ""
	}
	if strings.Contains(dateStr, " ") {
		dateParts := strings.Split(dateStr, " ")
		if len(dateParts) > 0 {
			dateStr = dateParts[0]
		}
	}
	parsedDate, err := time.Parse("02-01-2006", dateStr)
	if err != nil {
		parsedDate, err = time.Parse("02/01/2006", dateStr)
		if err != nil {
			return dateStr
		}
	}
	return parsedDate.Format("2006-01-02")
}

func ParseKVBANKFile(filePath string, holderAccount string) (*BankResponse, error) {
	parser := &KVBBANKParser{}
	transactions, err := parser.Parse("", filePath, holderAccount)
	if err != nil {
		return nil, fmt.Errorf("Failed to parse the KVB bank file: %v", err)
	}
	transInfos := parser.ConvertToTransInfo(transactions, holderAccount)
	response := &BankResponse{
		HolderAccount: holderAccount,
		TransInfo:     transInfos,
	}

	return response, nil
}
