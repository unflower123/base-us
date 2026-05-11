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

type NEOBANKParser struct{}

var _ BankParser = (*NEOBANKParser)(nil)

func (p *NEOBANKParser) Parse(fileContent string, filePath string, holderAccount string) ([]Transaction, error) {
	ext := strings.ToLower(filepath.Ext(filePath))
	if ext != ".csv" {
		return nil, fmt.Errorf("unsupported file format for NEO. Expected: .csv (CSV file), but got: %s", ext)
	}

	records, err := p.ParseCSVToRecords(filePath)
	if err != nil {
		return nil, err
	}

	var transactions []Transaction

	for _, record := range records {
		var amount float64
		var transType string
		if record.TransactionDate == "" || record.Particulars == "" {
			continue
		}

		if record.DebitCredit == "DR" {
			cleanedAmount := p.normalizeAmount(record.Amount)
			amount, err = strconv.ParseFloat(cleanedAmount, 64)
			if err != nil {
				log.Printf("Failed to parse the debit amount: %s, err: %v", record.Amount, err)
				continue
			}
			transType = "DEBIT"
			amount = -amount
		} else if record.DebitCredit == "CR" {
			cleanedAmount := p.normalizeAmount(record.Amount)
			amount, err = strconv.ParseFloat(cleanedAmount, 64)
			if err != nil {
				log.Printf("Failed to parse the credit amount: %s, err: %v", record.Amount, err)
				continue
			}
			transType = "CREDIT"
		} else {
			continue
		}

		transactions = append(transactions, Transaction{
			Date:        p.formatDate(record.TransactionDate),
			Description: record.Particulars,
			Amount:      amount,
			Type:        transType,
			Account:     holderAccount,
		})
	}

	log.Printf("Successfully parsed %d ENO bank transactions", len(transactions))
	return transactions, nil
}

func (p *NEOBANKParser) ConvertToTransInfo(transactions []Transaction, holderAccount string) []TransInfo {
	var transInfos []TransInfo

	for _, transaction := range transactions {
		transType := "TYPE_IN"
		if transaction.Type == "DEBIT" {
			transType = "TYPE_OUT"
			holderAccount = ""
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
		if bankTxnId != "" && transName != "" {
			bankTxnId = bankTxnId + "|+|" + transName
		}

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
			FundFlow:     fundFlow, // 修正 FundFlow 字段
		})
	}

	return transInfos
}

func (p *NEOBANKParser) extractName(description string) string {
	if strings.Contains(description, "IMPS/") {
		parts := strings.Split(description, "/")
		if len(parts) >= 4 {
			nameCandidate := strings.TrimSpace(parts[3])
			return nameCandidate
		}
	}
	if strings.Contains(description, "UPI/") {
		parts := strings.Split(description, "/")
		if len(parts) >= 4 {
			return strings.TrimSpace(parts[3])
		}
	}

	return ""
}

func (p *NEOBANKParser) extractUpi(description string) string {
	if strings.Contains(description, "UPI") {
		return "UPI_TRANSACTION"
	}
	if strings.Contains(description, "IMPS") {
		return "IMPS_TRANSACTION"
	}
	return ""
}

func (p *NEOBANKParser) extractTxnId(description string) string {
	if strings.Contains(description, "UPI/") {
		parts := strings.Split(description, "/")
		if len(parts) >= 3 {
			upiId := strings.TrimSpace(parts[2])
			if p.isPureDigits(upiId) {
				return upiId
			}
		}
	}
	if strings.Contains(description, "IMPS/") {
		re := regexp.MustCompile(`\d{8,}`)
		matches := re.FindAllString(description, -1)
		if len(matches) > 0 {
			impsId := matches[0]
			if p.isPureDigits(impsId) {
				return impsId
			}
		}
		parts := strings.Split(description, "/")
		if len(parts) >= 3 {
			candidate := strings.TrimSpace(parts[2])
			if p.isPureDigits(candidate) {
				return candidate
			}
		}
	}
	if strings.Contains(description, "NEFT/") {
		parts := strings.Split(description, "/")
		if len(parts) >= 3 {
			neftId := strings.TrimSpace(parts[2])
			if p.isPureDigits(neftId) {
				return neftId
			}
		}
	}
	re := regexp.MustCompile(`(Ref[:/]?\s*)(\d{8,})|(Txn[:/]?\s*)(\d{8,})|(ID[:/]?\s*)(\d{8,})`)
	matches := re.FindStringSubmatch(description)
	if len(matches) > 0 {
		for i := 2; i < len(matches); i += 2 {
			if matches[i] != "" && p.isPureDigits(matches[i]) {
				return matches[i]
			}
		}
	}

	reLongDigits := regexp.MustCompile(`\d{10,}`)
	longDigits := reLongDigits.FindString(description)
	if longDigits != "" && p.isPureDigits(longDigits) {
		return longDigits
	}

	reMediumDigits := regexp.MustCompile(`\d{8,}`)
	mediumDigits := reMediumDigits.FindString(description)
	if mediumDigits != "" && p.isPureDigits(mediumDigits) {
		return mediumDigits
	}

	if strings.Contains(description, "Chq/") || strings.Contains(description, "Cheque/") {
		parts := strings.Split(description, "/")
		for _, part := range parts {
			trimmedPart := strings.TrimSpace(part)
			if p.isPureDigits(trimmedPart) && len(trimmedPart) >= 6 {
				return trimmedPart
			}
		}
	}

	return ""
}

func (p *NEOBANKParser) isPureDigits(str string) bool {
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

func (p *NEOBANKParser) isBankCode(str string) bool {
	bankCodes := []string{"SBIN", "CNRB", "FDRL", "ICIC", "HDFC", "UTIB", "BDBL", "CRGB", "INDB", "YESB", "AXIS"}
	for _, code := range bankCodes {
		if strings.Contains(strings.ToUpper(str), code) {
			return true
		}
	}
	return false
}

func (p *NEOBANKParser) normalizeAmount(amountStr string) string {
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

func (p *NEOBANKParser) ParseCSVToRecords(filePath string) ([]NeoBankRecord, error) {
	csvContent, err := p.processNeoBankCSVFile(filePath)
	if err != nil {
		return nil, err
	}
	return p.parseNeoBankCSVToRecords(csvContent)
}

func (p *NEOBANKParser) processNeoBankCSVFile(filePath string) (string, error) {
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

		if strings.HasPrefix(trimmedLine, "S.No,Transaction Date") {
			inTransactionSection = true
			headerFound = true
			transactionLines = append(transactionLines, trimmedLine)
			continue
		}

		if strings.Contains(trimmedLine, "TRANSACTION TOTAL") ||
			strings.Contains(trimmedLine, "CLOSING BALANCE") ||
			strings.Contains(trimmedLine, "Unless the constituent") {
			break
		}

		if inTransactionSection && headerFound && trimmedLine != "" {
			if !strings.HasPrefix(trimmedLine, "S.No,") && trimmedLine != "" {
				transactionLines = append(transactionLines, trimmedLine)
			}
		}
	}

	if len(transactionLines) == 0 {
		return "", fmt.Errorf("No transaction data found")
	}
	var result strings.Builder
	for _, line := range transactionLines {
		result.WriteString(line + "\n")
	}

	processedContent := result.String()
	return processedContent, nil
}

type NeoBankRecord struct {
	SNo             string
	TransactionDate string
	ValueDate       string
	Particulars     string
	Amount          string
	DebitCredit     string
	Balance         string
	ChequeNumber    string
	BranchName      string
}

func (p *NEOBANKParser) parseNeoBankCSVToRecords(csvContent string) ([]NeoBankRecord, error) {
	reader := csv.NewReader(strings.NewReader(csvContent))
	reader.FieldsPerRecord = -1
	reader.LazyQuotes = true
	reader.TrimLeadingSpace = true

	records, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("Failed to read the CSV content: %v", err)
	}

	var result []NeoBankRecord
	for i := 1; i < len(records); i++ {
		if len(records[i]) < 9 {
			continue
		}

		record := NeoBankRecord{
			SNo:             strings.TrimSpace(records[i][0]),
			TransactionDate: strings.TrimSpace(records[i][1]),
			ValueDate:       strings.TrimSpace(records[i][2]),
			Particulars:     strings.TrimSpace(records[i][3]),
			Amount:          "",
			DebitCredit:     "",
			Balance:         strings.TrimSpace(records[i][6]),
			ChequeNumber:    strings.TrimSpace(records[i][7]),
			BranchName:      strings.TrimSpace(records[i][8]),
		}
		debitAmount := strings.TrimSpace(records[i][4])
		creditAmount := strings.TrimSpace(records[i][5])

		if debitAmount != "" && debitAmount != "0" {
			record.Amount = debitAmount
			record.DebitCredit = "DR"
		} else if creditAmount != "" && creditAmount != "0" {
			record.Amount = creditAmount
			record.DebitCredit = "CR"
		} else {
			continue
		}

		if record.SNo == "" || record.Particulars == "" ||
			strings.Contains(record.Particulars, "TRANSACTION TOTAL") ||
			strings.Contains(record.Particulars, "OPENING BALANCE") ||
			strings.Contains(record.Particulars, "CLOSING BALANCE") {
			continue
		}

		result = append(result, record)
	}
	return result, nil
}

func (p *NEOBANKParser) formatDate(dateStr string) string {
	if dateStr == "" {
		return ""
	}
	parsedDate, err := time.Parse("02/01/2006", dateStr)
	if err != nil {
		return dateStr
	}
	return parsedDate.Format("2006-01-02")
}

func ParseNEOBANKFile(filePath string, holderAccount string) (*BankResponse, error) {
	parser := &NEOBANKParser{}
	transactions, err := parser.Parse("", filePath, holderAccount)
	if err != nil {
		return nil, fmt.Errorf("Failed to parse the NEO bank file: %v", err)
	}
	transInfos := parser.ConvertToTransInfo(transactions, holderAccount)
	response := &BankResponse{
		HolderAccount: holderAccount,
		TransInfo:     transInfos,
	}

	return response, nil
}
