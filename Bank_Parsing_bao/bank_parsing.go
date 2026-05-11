package bankparsing

import (
	"fmt"
	"strings"
	"time"
)

// 基础结构体定义
type XLSXRecord struct {
	Srl         string
	TxnDate     string
	ValueDate   string
	Description string
	ChequeNo    string
	CRDR        string
	Amount      string
	Balance     string
	Debit       string
	Credit      string
}

type XLSXBalanceParser struct{}

type CsvRecord struct {
	Date      string `json:"Date"`      // 交易日期
	ValueDate string `json:"ValueDate"` // 生效日期
	ChqNo     string `json:"ChqNo"`     // 支票号码
	Narration string `json:"Narration"` // 交易叙述
	Cod       string `json:"Cod"`       // 代码
	Debit     string `json:"Debit"`     // 借方金额（支出）
	Credit    string `json:"Credit"`    // 贷方金额（收入）
	Balance   string `json:"Balance"`   // 交易后的账户余额
}

type BalanceReportRequest struct {
	ClientID      string `json:"clientID"`
	DeviceID      string `json:"deviceID"`
	Timestamp     string `json:"timestamp"`
	HolderName    string `json:"holderName"`
	HolderAccount string `json:"holderAccount"`
	HolderEmail   string `json:"holderEmail"`
	Balance       string `json:"balance"`
}

// 交易信息结构体
type TransInfo struct {
	TransType    string `json:"transType"`
	TransName    string `json:"transName"`
	TransAccount string `json:"transAccount"`
	TransUpistr  string `json:"transUpistr"`
	TransAmount  string `json:"transAmount"`
	BankTxnId    string `json:"bankTxnId"`
	TransDate    string `json:"transDate"`
	TransStatus  string `json:"transStatus"`
	FundFlow     int32  `json:"fundFlow,optional"`
}

// 银行响应结构体
type BankResponse struct {
	HolderAccount string      `json:"holderAccount"`
	TransInfo     []TransInfo `json:"transInfo"`
}

// 交易数据结构
type Transaction struct {
	Date        string
	Description string
	Amount      float64
	Type        string
	Account     string
	ChequeNo    string
	ValueDate   string
	BankTxnId   string
	Balance     string
	FundFlow    int32
	Channel     string
}

// 银行解析接口
type BankParser interface {
	Parse(fileContent string, filePath string, holderAccount string) ([]Transaction, error)
}

type BankParserFactory struct{}

func (f *BankParserFactory) GetParser(bankName string) (BankParser, error) {
	bankType := strings.ToUpper(strings.TrimSpace(bankName))
	switch bankType {
	case "BOBBANK":
		return &BOBBANKParser{}, nil
	case "IDBIBANK":
		return &IDBIBANKParser{}, nil
	case "IOBBANK":
		return &IOBBANKParser{}, nil
	case "JKBANK":
		return &JKBANKParser{}, nil
	case "INDIANBANK":
		return &INDIANBANKParser{}, nil
	case "CANARABANK":
		return &CANARABANKParser{}, nil
	case "CSBBANK":
		return &CSBBANKParser{}, nil
	case "CUBBANK":
		return &CUBBANKParser{}, nil
	case "NEOBANK":
		return &NEOBANKParser{}, nil
	case "KVBBANK":
		return &KVBBANKParser{}, nil
	case "BOIBANK":
		return &BOIBANKParser{}, nil
	case "IDFCBANK":
		return &IDFCBANKParser{}, nil
	case "JSFBBANK":
		return &JSFBBANKParser{}, nil
	case "SBIBANK":
		return &SBIBANKParser{}, nil
	case "BOMBANK":
		return &BOMBANKParser{}, nil
	default:
		return nil, fmt.Errorf("400001:Unsupported bank type: %s", bankName)
	}
}

func IsParseBank(bankName string, importType int32) (err error) {
	factory := &BankParserFactory{}
	_, err = factory.GetParser(bankName)
	//if importType == 1 {
	//} else {
	//	// TODO payout适配
	//	err = fmt.Errorf("400002:Unsupported payout type: %d", importType)
	//}

	return
}

// 主要的解析调用函数
func ParseBankFile(bankName string, filePath string, holderAccount string) (*BankResponse, error) {
	factory := &BankParserFactory{}
	parser, err := factory.GetParser(bankName)
	if err != nil {
		return nil, err
	}
	transactions, err := parser.Parse("", filePath, holderAccount)
	if err != nil {
		return nil, err
	}
	if customConverter, ok := parser.(interface {
		ConvertToTransInfo([]Transaction, string) []TransInfo
	}); ok {
		transInfos := customConverter.ConvertToTransInfo(transactions, holderAccount)
		var filteredTransInfos []TransInfo
		for _, transInfo := range transInfos {
			if transInfo.BankTxnId != "" {
				filteredTransInfos = append(filteredTransInfos, transInfo)
			} else {
			}
		}

		return &BankResponse{
			HolderAccount: holderAccount,
			TransInfo:     filteredTransInfos,
		}, nil
	}

	return nil, fmt.Errorf("parser does not implement ConvertToTransInfo")
}

func GetIndianTimeString() string {
	ist, err := time.LoadLocation("Asia/Kolkata")
	if err != nil {
		now := time.Now().UTC()
		return fmt.Sprintf("%02d:%02d:%02d", now.Hour(), now.Minute(), now.Second())
	}
	now := time.Now().In(ist)
	return fmt.Sprintf("%02d:%02d:%02d", now.Hour(), now.Minute(), now.Second())
}

func extractDatePart(dateStr string) string {
	formats := []string{
		"02-Jan-2006",
		"2006-01-02",
		"02/01/2006",
		"01/02/2006",
		"02-Jan-2006 15:04:05",
		"2006-01-02 15:04:05",
		"02/01/2006 15:04:05",
	}

	for _, format := range formats {
		if t, err := time.Parse(format, dateStr); err == nil {
			return t.Format("2006-01-02")
		}
	}
	ist, err := time.LoadLocation("Asia/Kolkata")
	if err != nil {
		return time.Now().UTC().Format("2006-01-02")
	}
	return time.Now().In(ist).Format("2006-01-02")
}

func FormatDateWithIndianTime(dateStr string) string {
	datePart := extractDatePart(dateStr)
	indianTime := GetIndianTimeString()
	return datePart + " " + indianTime
}

func FormatTransactionDate(dateStr string) string {
	return FormatDateWithIndianTime(dateStr)
}
