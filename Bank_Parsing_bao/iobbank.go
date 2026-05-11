package bankparsing

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type IOBBANKParser struct{}

var _ BankParser = (*IOBBANKParser)(nil)

func (p *IOBBANKParser) Parse(fileContent string, filePath string, holderAccount string) ([]Transaction, error) {
	ext := strings.ToLower(filepath.Ext(filePath))
	if ext != ".csv" {
		return nil, fmt.Errorf("unsupported file format for IOB. Expected: .csv (CSV file), but got: %s", ext)
	}

	records, err := p.ParseCSVToRecords(filePath)
	if err != nil {
		return nil, err
	}
	var transactions []Transaction
	parser := &XLSXBalanceParser{}

	for _, record := range records {
		var amount float64
		var transType string
		if record.Debit != "" && record.Debit != "0" {
			cleanedAmount := parser.normalizeAmount(record.Debit)
			amount, err = strconv.ParseFloat(cleanedAmount, 64)
			if err != nil {
				log.Printf("Failed to parse debit amount: %s, error: %v", record.Debit, err)
				continue
			}
			transType = "DEBIT"
			amount = -amount
		} else if record.Credit != "" && record.Credit != "0" {
			cleanedAmount := parser.normalizeAmount(record.Credit)
			amount, err = strconv.ParseFloat(cleanedAmount, 64)
			if err != nil {
				log.Printf("Failed to parse credit amount: %s, error: %v", record.Credit, err)
				continue
			}
			transType = "CREDIT"
		} else {
			continue
		}
		transactions = append(transactions, Transaction{
			Date:        p.formatTransactionDate(record.Date),
			Description: record.Narration,
			Amount:      amount,
			Type:        transType,
			Account:     holderAccount,
		})
	}

	return transactions, nil
}

func (p *IOBBANKParser) ParseCSVToRecords(filePath string) ([]CsvRecord, error) {
	csvContent, err := p.processCSVFile(filePath)
	if err != nil {
		return nil, err
	}
	return p.parseCSVToRecords(csvContent)
}

func (p *IOBBANKParser) ConvertToTransInfo(transactions []Transaction, currency string) []TransInfo {
	var transInfos []TransInfo
	parser := &XLSXBalanceParser{}

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

		bankTxnId := p.extractIOBTxnId(transaction.Description)
		transName := p.extractIOBName(transaction.Description)

		if bankTxnId == "" {
			bankTxnId = parser.extractTxnId(transaction.Description)
		}
		if transName == "" {
			transName = parser.extractName(transaction.Description)
		}
		if bankTxnId == "" {
			continue
		}
		transAmount := fmt.Sprintf("%.2f", transaction.Amount)
		if transaction.Amount < 0 {
			transAmount = fmt.Sprintf("%.2f", -transaction.Amount)
		}
		transDate := FormatDateWithIndianTime(transaction.Date)

		transInfos = append(transInfos, TransInfo{
			TransType:    transType,
			TransName:    transName,
			TransAccount: transaction.Account,
			TransUpistr:  parser.extractUpi(transaction.Description),
			TransAmount:  transAmount,
			BankTxnId:    bankTxnId,
			TransDate:    transDate,
			TransStatus:  "SUCCESS",
			FundFlow:     fundFlow,
		})
	}

	return transInfos
}

func (p *IOBBANKParser) formatTransactionDate(dateStr string) string {
	return FormatDateWithIndianTime(dateStr)
}

func (p *IOBBANKParser) processCSVFile(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to open file: %v", err)
	}
	defer file.Close()

	content, err := os.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to read file: %v", err)
	}

	cleanedContent := p.cleanBankCSVContent(string(content))
	reader := csv.NewReader(strings.NewReader(cleanedContent))
	reader.FieldsPerRecord = -1
	reader.TrimLeadingSpace = true

	records, err := reader.ReadAll()
	if err != nil {
		return p.processBankCSVWithManualParsing(cleanedContent)
	}

	if len(records) <= 1 {
		return "", fmt.Errorf("insufficient CSV file content")
	}

	processedContent, _, err := p.ProcessBankRecords(records)
	if err != nil {
		return "", err
	}

	return processedContent, nil
}

func (p *IOBBANKParser) cleanBankCSVContent(content string) string {
	lines := strings.Split(content, "\n")
	var cleanedLines []string

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		line = strings.TrimSuffix(line, ",")
		cleanedLines = append(cleanedLines, line)
	}

	return strings.Join(cleanedLines, "\n")
}

func (p *IOBBANKParser) processBankCSVWithManualParsing(content string) (string, error) {
	lines := strings.Split(content, "\n")
	if len(lines) <= 1 {
		return "", fmt.Errorf("insufficient CSV file content")
	}

	var records [][]string

	for i, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		fields := p.parseCSVLine(line)

		if i == 0 {
			continue
		}

		if len(fields) >= 8 {
			records = append(records, fields[:8])
		}
	}

	processedContent, _, err := p.ProcessBankRecords(records)
	if err != nil {
		return "", err
	}

	return processedContent, nil
}

func (p *IOBBANKParser) parseCSVLine(line string) []string {
	var fields []string
	var currentField strings.Builder
	inQuotes := false

	for _, char := range line {
		switch char {
		case '"':
			inQuotes = !inQuotes
		case ',':
			if inQuotes {
				currentField.WriteRune(char)
			} else {
				fields = append(fields, strings.TrimSpace(currentField.String()))
				currentField.Reset()
			}
		default:
			currentField.WriteRune(char)
		}
	}
	if currentField.Len() > 0 {
		fields = append(fields, strings.TrimSpace(currentField.String()))
	}

	return fields
}

func (p *IOBBANKParser) ProcessBankRecords(records [][]string) (string, bool, error) {
	header := []string{"Date", "Value Date", "Chq No", "Narration", "Cod", "Debit", "Credit", "Balance"}
	var newRecords [][]string
	newRecords = append(newRecords, header)

	for i := 1; i < len(records); i++ {
		record := records[i]
		if len(record) >= 8 {
			dateKey := strings.TrimSpace(record[0])
			if dateKey != "" && dateKey != "Date" {
				normalizedRecord := make([]string, 8)
				for j := 0; j < 8; j++ {
					if j < len(record) {
						normalizedRecord[j] = strings.TrimSpace(record[j])
					}
				}
				newRecords = append(newRecords, normalizedRecord)
			}
		}
	}

	if len(newRecords) <= 1 {
		return "", false, fmt.Errorf("no valid transaction records")
	}

	var result strings.Builder
	writer := csv.NewWriter(&result)
	err := writer.WriteAll(newRecords)
	if err != nil {
		return "", false, fmt.Errorf("failed to write CSV: %v", err)
	}
	writer.Flush()
	return result.String(), true, nil
}

func (p *IOBBANKParser) parseCSVToRecords(csvContent string) ([]CsvRecord, error) {
	reader := csv.NewReader(strings.NewReader(csvContent))
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	var result []CsvRecord
	for i := 1; i < len(records); i++ {
		if len(records[i]) < 8 {
			continue
		}
		record := CsvRecord{
			Date:      strings.TrimSpace(records[i][0]),
			ValueDate: strings.TrimSpace(records[i][1]),
			ChqNo:     strings.TrimSpace(records[i][2]),
			Narration: strings.TrimSpace(records[i][3]),
			Cod:       strings.TrimSpace(records[i][4]),
			Debit:     strings.TrimSpace(records[i][5]),
			Credit:    strings.TrimSpace(records[i][6]),
			Balance:   strings.TrimSpace(records[i][7]),
		}

		if record.Date == "" || (record.Debit == "" && record.Credit == "") {
			continue
		}

		result = append(result, record)
	}

	log.Printf("Successfully parsed %d valid transaction records from CSV", len(result))
	return result, nil
}

func (p *IOBBANKParser) extractIOBTxnId(description string) string {
	if strings.Contains(description, "UPI/") {
		parts := strings.Split(description, "/")
		if len(parts) >= 3 {
			txnId := strings.TrimSpace(parts[1])
			code := strings.TrimSpace(parts[len(parts)-1])
			return fmt.Sprintf("%s|+|%s", txnId, code)
		}
	}
	return ""
}

func (p *IOBBANKParser) extractIOBName(description string) string {
	if strings.Contains(description, "UPI/") {
		parts := strings.Split(description, "/")
		if len(parts) >= 5 {
			name := strings.TrimSpace(parts[3])
			return strings.Join(strings.Fields(name), " ")
		}
	}
	return ""
}

func ParseIOBBANKFile(filePath string, holderAccount string, currency string) (*BankResponse, error) {
	parser := &IOBBANKParser{}
	transactions, err := parser.Parse("", filePath, holderAccount)
	if err != nil {
		return nil, fmt.Errorf("failed to parse IOB bank file: %v", err)
	}
	transInfos := parser.ConvertToTransInfo(transactions, currency)

	response := &BankResponse{
		HolderAccount: holderAccount,
		TransInfo:     transInfos,
	}

	return response, nil
}
