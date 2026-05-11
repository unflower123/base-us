package bankparsing

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/extrame/xls"
	"github.com/xuri/excelize/v2"
	"log"
	"net/url"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"
)

type IDBIBANKParser struct{}

func (p *IDBIBANKParser) Parse(fileContent string, filePath string, holderAccount string) ([]Transaction, error) {
	ext := strings.ToLower(filepath.Ext(filePath))
	if ext != ".xls" && ext != ".xlsx" {
		return nil, fmt.Errorf("unsupported file format for IDBI. Expected: .xls or .xlsx (Excel files), but got: %s", ext)
	}

	records, err := p.ParseXLSXToRecords(filePath)
	if err != nil {
		return nil, err
	}

	var transactions []Transaction
	parser := &XLSXBalanceParser{}

	for _, record := range records {
		cleanedAmount := parser.normalizeAmount(record.Amount)
		amount, err := strconv.ParseFloat(cleanedAmount, 64)
		if err != nil {
			log.Printf("解析金额失败: %s, 错误: %v", record.Amount, err)
			continue
		}

		transType := "CREDIT"
		if strings.Contains(strings.ToUpper(record.CRDR), "DR") {
			transType = "DEBIT"
			amount = -amount
		}
		bankTxnId := parser.extractTxnId(record.Description)
		if bankTxnId == "" {
			continue
		}

		transactions = append(transactions, Transaction{
			Date:        parser.formatDate(record.TxnDate),
			Description: record.Description,
			Amount:      amount,
			Type:        transType,
			Account:     holderAccount,
		})
	}

	return transactions, nil
}

func (p *IDBIBANKParser) ParseXLSXToRecords(filePath string) ([]XLSXRecord, error) {
	ext := strings.ToLower(filepath.Ext(filePath))

	if ext == ".xls" {
		return p.parseXLSToRecords(filePath)
	} else if ext == ".xlsx" {
		return p.parseXLSXToRecordsOriginal(filePath)
	}

	return nil, fmt.Errorf("不支持的文件格式: %s", ext)
}

func (p *IDBIBANKParser) parseXLSXToRecordsOriginal(xlsxPath string) ([]XLSXRecord, error) {
	var records []XLSXRecord
	f, err := excelize.OpenFile(xlsxPath)
	if err != nil {
		return nil, fmt.Errorf("打开XLSX文件失败: %v", err)
	}
	defer f.Close()

	sheets := f.GetSheetList()
	if len(sheets) == 0 {
		return nil, fmt.Errorf("XLSX文件中没有工作表")
	}

	sheetName := sheets[0]
	rows, err := f.GetRows(sheetName)
	if err != nil {
		return nil, fmt.Errorf("读取工作表数据失败: %v", err)
	}

	if len(rows) == 0 {
		return nil, fmt.Errorf("工作表为空")
	}

	headerRow := -1
	var srlCol, txnDateCol, valueDateCol, descCol, chequeNoCol, crdrCol, amountCol, balanceCol int = -1, -1, -1, -1, -1, -1, -1, -1
	for i, row := range rows {
		if len(row) == 0 {
			continue
		}
		rowText := strings.Join(row, " ")
		if strings.Contains(rowText, "Srl") &&
			(strings.Contains(rowText, "Txn Date") || strings.Contains(rowText, "Transaction Date")) &&
			(strings.Contains(rowText, "Value Date") || strings.Contains(rowText, "Val Date")) {
			headerRow = i

			for j, cell := range row {
				cell = strings.TrimSpace(cell)
				switch {
				case strings.Contains(cell, "Srl") || strings.Contains(cell, "Sl No"):
					srlCol = j
				case strings.Contains(cell, "Txn Date") || strings.Contains(cell, "Transaction Date"):
					txnDateCol = j
				case strings.Contains(cell, "Value Date") || strings.Contains(cell, "Val Date"):
					valueDateCol = j
				case strings.Contains(cell, "Description") || strings.Contains(cell, "Narration"):
					descCol = j
				case strings.Contains(cell, "Cheque No") || strings.Contains(cell, "Chq No"):
					chequeNoCol = j
				case strings.Contains(cell, "CR/DR") || strings.Contains(cell, "Dr/Cr"):
					crdrCol = j
				case strings.Contains(cell, "Amount"):
					amountCol = j
				case strings.Contains(cell, "Balance"):
					balanceCol = j
				}
			}
			break
		}
	}

	if headerRow == -1 {
		return nil, nil
	}
	for i := headerRow + 1; i < len(rows); i++ {
		row := rows[i]
		if len(row) == 0 {
			continue
		}

		record := XLSXRecord{}
		if srlCol >= 0 && srlCol < len(row) {
			record.Srl = strings.TrimSpace(row[srlCol])
		}
		if txnDateCol >= 0 && txnDateCol < len(row) {
			record.TxnDate = strings.TrimSpace(row[txnDateCol])
		}
		if valueDateCol >= 0 && valueDateCol < len(row) {
			record.ValueDate = strings.TrimSpace(row[valueDateCol])
		}
		if descCol >= 0 && descCol < len(row) {
			record.Description = strings.TrimSpace(row[descCol])
		}
		if chequeNoCol >= 0 && chequeNoCol < len(row) {
			record.ChequeNo = strings.TrimSpace(row[chequeNoCol])
		}
		if crdrCol >= 0 && crdrCol < len(row) {
			record.CRDR = strings.TrimSpace(row[crdrCol])
		}
		if amountCol >= 0 && amountCol < len(row) {
			record.Amount = strings.TrimSpace(row[amountCol])
		}
		if balanceCol >= 0 && balanceCol < len(row) {
			record.Balance = strings.TrimSpace(row[balanceCol])
		}

		if record.TxnDate == "" || record.Amount == "" {
			continue
		}

		records = append(records, record)
	}

	return records, nil
}

func (p *IDBIBANKParser) parseXLSToRecords(xlsPath string) ([]XLSXRecord, error) {
	var records []XLSXRecord
	file, err := xls.Open(xlsPath, "utf-8")
	if err != nil {
		return nil, fmt.Errorf("打开XLS文件失败: %v", err)
	}

	sheet := file.GetSheet(0)
	if sheet == nil {
		return nil, fmt.Errorf("XLS文件中没有工作表")
	}

	headerRow := -1
	var srlCol, txnDateCol, valueDateCol, descCol, chequeNoCol, crdrCol, amountCol, balanceCol int = -1, -1, -1, -1, -1, -1, -1, -1
	for i := 0; i <= int(sheet.MaxRow); i++ {
		row := sheet.Row(i)
		if row == nil {
			continue
		}

		var rowCells []string
		for j := 0; j < row.LastCol(); j++ {
			cell := row.Col(j)
			rowCells = append(rowCells, strings.TrimSpace(cell))
		}

		if len(rowCells) == 0 {
			continue
		}

		rowText := strings.Join(rowCells, " ")

		if strings.Contains(rowText, "Srl") &&
			(strings.Contains(rowText, "Txn Date") || strings.Contains(rowText, "Transaction Date")) {
			headerRow = i
			for j, cell := range rowCells {
				cell = strings.TrimSpace(cell)
				switch {
				case strings.Contains(cell, "Srl") || strings.Contains(cell, "Sl No"):
					srlCol = j
				case strings.Contains(cell, "Txn Date") || strings.Contains(cell, "Transaction Date"):
					txnDateCol = j
				case strings.Contains(cell, "Value Date") || strings.Contains(cell, "Val Date"):
					valueDateCol = j
				case strings.Contains(cell, "Description") || strings.Contains(cell, "Narration"):
					descCol = j
				case strings.Contains(cell, "Cheque No") || strings.Contains(cell, "Chq No"):
					chequeNoCol = j
				case strings.Contains(cell, "CR/DR") || strings.Contains(cell, "Dr/Cr"):
					crdrCol = j
				case strings.Contains(cell, "Amount"):
					amountCol = j
				case strings.Contains(cell, "Balance"):
					balanceCol = j
				}
			}
			break
		}
	}

	if headerRow == -1 {
		return nil, fmt.Errorf("未找到表头行")
	}

	log.Printf("表头行: %d, 列索引: Srl=%d, TxnDate=%d, ValueDate=%d, Desc=%d, ChequeNo=%d, CRDR=%d, Amount=%d, Balance=%d",
		headerRow+1, srlCol+1, txnDateCol+1, valueDateCol+1, descCol+1, chequeNoCol+1, crdrCol+1, amountCol+1, balanceCol+1)

	dataRowCount := 0
	for i := headerRow + 1; i <= int(sheet.MaxRow); i++ {
		row := sheet.Row(i)
		if row == nil {
			continue
		}

		var rowCells []string
		for j := 0; j < row.LastCol(); j++ {
			cell := row.Col(j)
			rowCells = append(rowCells, strings.TrimSpace(cell))
		}

		if len(rowCells) == 0 {
			continue
		}

		record := XLSXRecord{}
		if srlCol >= 0 && srlCol < len(rowCells) {
			record.Srl = rowCells[srlCol]
		}
		if txnDateCol >= 0 && txnDateCol < len(rowCells) {
			record.TxnDate = rowCells[txnDateCol]
		}
		if valueDateCol >= 0 && valueDateCol < len(rowCells) {
			record.ValueDate = rowCells[valueDateCol]
		}
		if descCol >= 0 && descCol < len(rowCells) {
			record.Description = rowCells[descCol]
		}
		if chequeNoCol >= 0 && chequeNoCol < len(rowCells) {
			record.ChequeNo = rowCells[chequeNoCol]
		}
		if crdrCol >= 0 && crdrCol < len(rowCells) {
			record.CRDR = rowCells[crdrCol]
		}
		if amountCol >= 0 && amountCol < len(rowCells) {
			record.Amount = rowCells[amountCol]
		}
		if balanceCol >= 0 && balanceCol < len(rowCells) {
			record.Balance = rowCells[balanceCol]
		}

		if record.TxnDate == "" || record.Amount == "" {
			continue
		}

		records = append(records, record)
		dataRowCount++
	}
	return records, nil
}

func (p *XLSXBalanceParser) extractName(desc string) string {
	parts := strings.Split(desc, ",")
	if len(parts) > 0 {
		return strings.TrimSpace(parts[0])
	}
	return desc
}

func (p *XLSXBalanceParser) extractUpi(desc string) string {
	descLower := strings.ToLower(desc)
	if strings.Contains(descLower, "@") {
		start := strings.Index(descLower, "@")
		left := start
		for left > 0 && desc[left-1] != ' ' {
			left--
		}
		right := start
		for right < len(desc) && desc[right] != ' ' {
			right++
		}
		return desc[left:right]
	}
	return ""
}

func (p *XLSXBalanceParser) extractTxnId(desc string) string {
	if strings.Contains(desc, "/") {
		parts := strings.Split(desc, "/")
		if len(parts) >= 3 && parts[0] == "UPI" {
			return strings.TrimSpace(parts[1])
		}
	}
	return ""
}

func (p *XLSXBalanceParser) formatDate(dateStr string) string {
	for _, layout := range []string{"02-01-2006", "01/02/2006", "2006-01-02"} {
		if t, err := time.Parse(layout, dateStr); err == nil {
			return t.Format("2006-01-02")
		}
	}
	return dateStr
}

func (p *XLSXBalanceParser) normalizeAmount(amount string) string {
	cleaned := strings.ReplaceAll(amount, ",", "")
	cleaned = strings.ReplaceAll(cleaned, "INR ", "")
	cleaned = strings.ReplaceAll(cleaned, " ", "")
	return cleaned
}

func (p *XLSXBalanceParser) MarshalStructToSortedString(v any) (string, string, error) {
	paramMap := make(map[string]any)

	jsonBytes, err := json.Marshal(v)
	if err != nil {
		return "", "", fmt.Errorf("failed to marshal input to JSON: %w", err)
	}

	decoder := json.NewDecoder(bytes.NewReader(jsonBytes))
	decoder.UseNumber()

	if err := decoder.Decode(&paramMap); err != nil {
		return "", "", fmt.Errorf("failed to decode JSON to map: %w", err)
	}

	paramSign, ok := paramMap["sign"].(string)
	if !ok {
		paramSign = ""
	}

	flatMap := make(map[string]any)
	p.flattenMap("", paramMap, flatMap)
	delete(flatMap, "sign")

	var keys []string
	for k := range flatMap {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var builder strings.Builder
	for i, k := range keys {
		valueStr, err := p.valueToString(flatMap[k])
		if err != nil {
			return "", "", fmt.Errorf("failed to convert value for key %s: %w", k, err)
		}
		encodedKey := url.QueryEscape(k)
		encodedValue := url.QueryEscape(valueStr)
		if i > 0 {
			builder.WriteByte('&')
		}
		builder.WriteString(encodedKey)
		builder.WriteByte('=')
		builder.WriteString(encodedValue)
	}

	return builder.String(), paramSign, nil
}

// flattenMap 扁平化 map
func (p *XLSXBalanceParser) flattenMap(prefix string, input map[string]interface{}, output map[string]interface{}) {
	for k, v := range input {
		key := k
		if prefix != "" {
			key = prefix + "." + k
		}
		switch val := v.(type) {
		case map[string]interface{}:
			p.flattenMap(key, val, output)
		default:
			output[key] = v
		}
	}
}

func (p *IDBIBANKParser) ConvertToTransInfo(transactions []Transaction, holderAccount string) []TransInfo {
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

		transName := parser.extractName(transaction.Description)
		transUpistr := parser.extractUpi(transaction.Description)
		bankTxnId := parser.extractTxnId(transaction.Description)

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

// valueToString 将值转换为字符串
func (p *XLSXBalanceParser) valueToString(v interface{}) (string, error) {
	switch val := v.(type) {
	case string:
		return val, nil
	case float64:
		return strconv.FormatFloat(val, 'f', -1, 64), nil
	case float32:
		return strconv.FormatFloat(float64(val), 'f', -1, 32), nil
	case int64:
		return strconv.FormatInt(val, 10), nil
	case int32:
		return strconv.FormatInt(int64(val), 10), nil
	case int:
		return strconv.Itoa(val), nil
	case uint64:
		return strconv.FormatUint(val, 10), nil
	case uint32:
		return strconv.FormatUint(uint64(val), 10), nil
	case uint:
		return strconv.FormatUint(uint64(val), 10), nil
	case bool:
		return strconv.FormatBool(val), nil
	case []interface{}:
		var strs []string
		for _, item := range val {
			str, err := p.valueToString(item)
			if err != nil {
				return "", err
			}
			strs = append(strs, str)
		}
		return strings.Join(strs, ","), nil
	default:
		jsonVal, err := json.Marshal(val)
		if err != nil {
			return "", fmt.Errorf("failed to serialize value: %w", err)
		}
		return string(jsonVal), nil
	}
}
