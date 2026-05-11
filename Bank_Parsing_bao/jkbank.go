package bankparsing

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

type JKBANKParser struct{}

var _ BankParser = (*JKBANKParser)(nil)

func (p *JKBANKParser) Parse(fileContent string, filePath string, holderAccount string) ([]Transaction, error) {
	ext := strings.ToLower(filepath.Ext(filePath))
	if ext != ".csv" {
		return nil, fmt.Errorf("unsupported file format for JK. Expected: .csv (CSV file), but got: %s", ext)
	}

	records, err := p.ParseCSVToRecords(filePath)
	if err != nil {
		return nil, err
	}

	var transactions []Transaction

	for _, record := range records {
		var amount float64
		var transType string
		if record.Debit != "" && record.Debit != "0" && record.Debit != "0.0" && record.Debit != "0.0  " {
			cleanedAmount := p.normalizeAmount(record.Debit)
			amount, err = strconv.ParseFloat(cleanedAmount, 64)
			if err != nil {
				log.Printf("解析借方金额失败: %s, 错误: %v", record.Debit, err)
				continue
			}
			transType = "DEBIT"
			amount = -amount // 支出为负数
		} else if record.Credit != "" && record.Credit != "0" && record.Credit != "0.0" && record.Credit != "0.0  " {
			cleanedAmount := p.normalizeAmount(record.Credit)
			amount, err = strconv.ParseFloat(cleanedAmount, 64)
			if err != nil {
				log.Printf("解析贷方金额失败: %s, 错误: %v", record.Credit, err)
				continue
			}
			transType = "CREDIT"
		} else {
			continue
		}

		transactions = append(transactions, Transaction{
			Date:        p.formatDate(record.Date),
			Description: record.Narration,
			Amount:      amount,
			Type:        transType,
			Account:     holderAccount,
		})
	}

	log.Printf("成功解析 %d 笔JK银行交易", len(transactions))
	return transactions, nil
}

func (p *JKBANKParser) ConvertToTransInfo(transactions []Transaction, holderAccount string) []TransInfo {
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
			FundFlow:     fundFlow,
		})
	}

	return transInfos
}

func (p *JKBANKParser) extractName(description string) string {
	if strings.Contains(description, "UPI") {
		parts := strings.Split(description, "/")
		if len(parts) >= 4 {
			return strings.TrimSpace(parts[3])
		}
	}
	if strings.Contains(description, "TO ") {
		parts := strings.Split(description, "TO ")
		if len(parts) >= 2 {
			return strings.TrimSpace(parts[1])
		}
	}

	return ""
}

func (p *JKBANKParser) extractUpi(description string) string {
	if strings.Contains(description, "UPI") {
		return "UPI_TRANSACTION"
	}
	return ""
}

func (p *JKBANKParser) extractTxnId(description string) string {
	if strings.Contains(description, "UPI") {
		parts := strings.Split(description, "/")
		if len(parts) >= 3 {
			return strings.TrimSpace(parts[2])
		}
	}
	if strings.Contains(description, "Ref No:") {
		parts := strings.Split(description, "Ref No:")
		if len(parts) >= 2 {
			return strings.TrimSpace(parts[1])
		}
	}

	return ""
}

func (p *JKBANKParser) normalizeAmount(amountStr string) string {
	cleaned := strings.ReplaceAll(amountStr, ",", "")
	cleaned = strings.ReplaceAll(cleaned, " ", "")
	cleaned = strings.TrimSpace(cleaned)
	if strings.Contains(cleaned, `"`) {
		cleaned = strings.ReplaceAll(cleaned, `"`, "")
	}
	if strings.HasPrefix(cleaned, "-") {
		cleaned = strings.TrimPrefix(cleaned, "-")
	}

	return cleaned
}

func (p *JKBANKParser) ParseCSVToRecords(filePath string) ([]CsvRecord, error) {
	csvContent, err := p.processJKBankCSVFile(filePath)
	if err != nil {
		return nil, err
	}
	return p.parseJKBankCSVToRecords(csvContent)
}

func (p *JKBANKParser) processJKBankCSVFile(filePath string) (string, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("读取文件失败: %v", err)
	}

	lines := strings.Split(string(content), "\n")
	var transactionLines []string
	inTransactionSection := false
	headerFound := false
	for _, line := range lines {
		trimmedLine := strings.TrimSpace(line)
		if strings.Contains(trimmedLine, "Transactions List") {
			inTransactionSection = true
			continue
		}
		if strings.Contains(trimmedLine, "Legends Used in Account Statement") {
			break
		}
		if inTransactionSection && trimmedLine != "" {
			if strings.Contains(trimmedLine, "Value Date") &&
				strings.Contains(trimmedLine, "Transaction Date") &&
				strings.Contains(trimmedLine, "Withdrawal") {
				headerFound = true
				continue
			}
			if headerFound && trimmedLine != "" {
				transactionLines = append(transactionLines, trimmedLine)
			}
		}
	}

	if len(transactionLines) == 0 {
		return "", fmt.Errorf("未找到交易数据")
	}
	var result strings.Builder
	result.WriteString("Date,Value Date,Chq No,Narration,Debit,Credit,Balance\n")

	for i, line := range transactionLines {
		fields := p.parseJKBankTransactionLine(line)
		if len(fields) >= 10 {
			debit := strings.TrimSpace(fields[7])
			credit := strings.TrimSpace(fields[8])
			balance := strings.TrimSpace(fields[9])
			csvLine := fmt.Sprintf("%s,%s,%s,%s,%s,%s,%s",
				strings.TrimSpace(fields[3]),
				strings.TrimSpace(fields[2]),
				strings.TrimSpace(fields[4]),
				strings.TrimSpace(fields[6]),
				debit,
				credit,
				balance,
			)
			result.WriteString(csvLine + "\n")
		} else {
			log.Printf("第%d行字段数量不足: %d, 跳过", i+1, len(fields))
		}
	}

	processedContent := result.String()
	return processedContent, nil
}

// parseJKBankTransactionLine 解析 JK Bank 交易行
func (p *JKBANKParser) parseJKBankTransactionLine(line string) []string {
	reader := csv.NewReader(strings.NewReader(line))
	reader.LazyQuotes = true
	reader.TrimLeadingSpace = true

	fields, err := reader.Read()
	if err != nil {
		return p.parseJKBankTransactionLineManual(line)
	}

	return fields
}

// parseJKBankTransactionLineManual 手动解析交易行
func (p *JKBANKParser) parseJKBankTransactionLineManual(line string) []string {
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

	// 添加最后一个字段
	if currentField.Len() > 0 {
		fields = append(fields, strings.TrimSpace(currentField.String()))
	}

	return fields
}

func (p *JKBANKParser) parseJKBankCSVToRecords(csvContent string) ([]CsvRecord, error) {
	reader := csv.NewReader(strings.NewReader(csvContent))
	reader.FieldsPerRecord = -1
	reader.LazyQuotes = true
	reader.TrimLeadingSpace = true

	records, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("读取CSV内容失败: %v", err)
	}

	var result []CsvRecord
	for i := 1; i < len(records); i++ {
		if len(records[i]) < 7 {
			continue
		}

		record := CsvRecord{
			Date:      strings.TrimSpace(records[i][0]),
			ValueDate: strings.TrimSpace(records[i][1]),
			ChqNo:     strings.TrimSpace(records[i][2]),
			Narration: strings.TrimSpace(records[i][3]),
			Debit:     strings.TrimSpace(records[i][4]),
			Credit:    strings.TrimSpace(records[i][5]),
			Balance:   strings.TrimSpace(records[i][6]),
		}
		if record.Date == "" {
			continue
		}
		result = append(result, record)
	}
	return result, nil
}

func (p *JKBANKParser) formatDate(dateStr string) string {
	if dateStr == "" {
		return ""
	}
	parsedDate, err := time.Parse("02/01/2006", dateStr)
	if err != nil {
		return dateStr
	}
	return parsedDate.Format("2006-01-02")
}

func ParseJKBANKFile(filePath string, holderAccount string) (*BankResponse, error) {
	parser := &JKBANKParser{}
	transactions, err := parser.Parse("", filePath, holderAccount)
	if err != nil {
		return nil, fmt.Errorf("解析JK银行文件失败: %v", err)
	}
	transInfos := parser.ConvertToTransInfo(transactions, holderAccount)
	response := &BankResponse{
		HolderAccount: holderAccount,
		TransInfo:     transInfos,
	}

	return response, nil
}
