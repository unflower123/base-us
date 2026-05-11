/**
 * @Author: darry
 * @Desc:
 * @Date: 2025/4/12 11:35
 */
package consts

import (
	"fmt"
	"time"
)

// Currency codes according to ISO 4217 standard
// The following includes currency codes for 100 major countries and regions
const (
	// AED - United Arab Emirates - UAE dirham
	AED = "AED"
	// AFN - Afghanistan - Afghan afghani
	AFN = "AFN"
	// ALL - Albania - Albanian lek
	ALL = "ALL"
	// AMD - Armenia - Armenian dram
	AMD = "AMD"
	// ANG - Curaçao/Sint Maarten - Netherlands Antillean guilder
	ANG = "ANG"
	// AOA - Angola - Angolan kwanza
	AOA = "AOA"
	// ARS - Argentina - Argentine peso
	ARS = "ARS"
	// AUD - Australia - Australian dollar
	AUD = "AUD"
	// AWG - Aruba - Aruban florin
	AWG = "AWG"
	// BAM - Bosnia and Herzegovina - Bosnian convertible mark
	BAM = "BAM"
	// BDT - Bangladesh - Bangladeshi taka
	BDT = "BDT"
	// BGN - Bulgaria - Bulgarian lev
	BGN = "BGN"
	// BHD - Bahrain - Bahraini dinar
	BHD = "BHD"
	// BIF - Burundi - Burundian franc
	BIF = "BIF"
	// BOB - Bolivia - Bolivian boliviano
	BOB = "BOB"
	// BRL - Brazil - Brazilian real
	BRL = "BRL"
	// BWP - Botswana - Botswana pula
	BWP = "BWP"
	// BYN - Belarus - Belarusian ruble
	BYN = "BYN"
	// CAD - Canada - Canadian dollar
	CAD = "CAD"
	// CDF - Congo (DRC) - Congolese franc
	CDF = "CDF"
	// CHF - Switzerland - Swiss franc
	CHF = "CHF"
	// CLP - Chile - Chilean peso
	CLP = "CLP"
	// CNY - China - Chinese yuan
	CNY = "CNY"
	// COP - Colombia - Colombian peso
	COP = "COP"
	// CRC - Costa Rica - Costa Rican colón
	CRC = "CRC"
	// CUP - Cuba - Cuban peso
	CUP = "CUP"
	// CZK - Czech Republic - Czech koruna
	CZK = "CZK"
	// DKK - Denmark - Danish krone
	DKK = "DKK"
	// DOP - Dominican Republic - Dominican peso
	DOP = "DOP"
	// DZD - Algeria - Algerian dinar
	DZD = "DZD"
	// EGP - Egypt - Egyptian pound
	EGP = "EGP"
	// ETB - Ethiopia - Ethiopian birr
	ETB = "ETB"
	// EUR - Eurozone (Austria, Belgium, etc.) - Euro
	EUR = "EUR"
	// GBP - United Kingdom - British pound
	GBP = "GBP"
	// GEL - Georgia - Georgian lari
	GEL = "GEL"
	// GHS - Ghana - Ghanaian cedi
	GHS = "GHS"
	// GTQ - Guatemala - Guatemalan quetzal
	GTQ = "GTQ"
	// HKD - Hong Kong - Hong Kong dollar
	HKD = "HKD"
	// HNL - Honduras - Honduran lempira
	HNL = "HNL"
	// HTG - Haiti - Haitian gourde
	HTG = "HTG"
	// HUF - Hungary - Hungarian forint
	HUF = "HUF"
	// IDR - Indonesia - Indonesian rupiah
	IDR = "IDR"
	// ILS - Israel - Israeli new shekel
	ILS = "ILS"
	// INR - India - Indian rupee
	INR = "INR"
	// IRR - Iran - Iranian rial
	IRR = "IRR"
	// JOD - Jordan - Jordanian dinar
	JOD = "JOD"
	// JPY - Japan - Japanese yen
	JPY = "JPY"
	// KES - Kenya - Kenyan shilling
	KES = "KES"
	// KGS - Kyrgyzstan - Kyrgyzstani som
	KGS = "KGS"
	// KHR - Cambodia - Cambodian riel
	KHR = "KHR"
	// KRW - South Korea - South Korean won
	KRW = "KRW"
	// KWD - Kuwait - Kuwaiti dinar
	KWD = "KWD"
	// KZT - Kazakhstan - Kazakhstani tenge
	KZT = "KZT"
	// LAK - Laos - Lao kip
	LAK = "LAK"
	// LBP - Lebanon - Lebanese pound
	LBP = "LBP"
	// LKR - Sri Lanka - Sri Lankan rupee
	LKR = "LKR"
	// LRD - Liberia - Liberian dollar
	LRD = "LRD"
	// LYD - Libya - Libyan dinar
	LYD = "LYD"
	// MAD - Morocco - Moroccan dirham
	MAD = "MAD"
	// MDL - Moldova - Moldovan leu
	MDL = "MDL"
	// MGA - Madagascar - Malagasy ariary
	MGA = "MGA"
	// MKD - North Macedonia - Macedonian denar
	MKD = "MKD"
	// MMK - Myanmar - Burmese kyat
	MMK = "MMK"
	// MNT - Mongolia - Mongolian tögrög
	MNT = "MNT"
	// MOP - Macao - Macanese pataca
	MOP = "MOP"
	// MXN - Mexico - Mexican peso
	MXN = "MXN"
	// MYR - Malaysia - Malaysian ringgit
	MYR = "MYR"
	// MWK - Malawi - Malawian kwacha
	MWK = "MWK"
	// NAD - Namibia - Namibian dollar
	NAD = "NAD"
	// NGN - Nigeria - Nigerian naira
	NGN = "NGN"
	// NIO - Nicaragua - Nicaraguan córdoba
	NIO = "NIO"
	// NOK - Norway - Norwegian krone
	NOK = "NOK"
	// NPR - Nepal - Nepalese rupee
	NPR = "NPR"
	// NZD - New Zealand - New Zealand dollar
	NZD = "NZD"
	// OMR - Oman - Omani rial
	OMR = "OMR"
	// PAB - Panama - Panamanian balboa
	PAB = "PAB"
	// PEN - Peru - Peruvian sol
	PEN = "PEN"
	// PGK - Papua New Guinea - Papua New Guinean kina
	PGK = "PGK"
	// PHP - Philippines - Philippine peso
	PHP = "PHP"
	// PKR - Pakistan - Pakistani rupee
	PKR = "PKR"
	// PLN - Poland - Polish złoty
	PLN = "PLN"
	// PYG - Paraguay - Paraguayan guaraní
	PYG = "PYG"
	// QAR - Qatar - Qatari riyal
	QAR = "QAR"
	// RON - Romania - Romanian leu
	RON = "RON"
	// RSD - Serbia - Serbian dinar
	RSD = "RSD"
	// RUB - Russia - Russian ruble
	RUB = "RUB"
	// SAR - Saudi Arabia - Saudi riyal
	SAR = "SAR"
	// SDG - Sudan - Sudanese pound
	SDG = "SDG"
	// SEK - Sweden - Swedish krona
	SEK = "SEK"
	// SGD - Singapore - Singapore dollar
	SGD = "SGD"
	// SRD - Suriname - Surinamese dollar
	SRD = "SRD"
	// TJS - Tajikistan - Tajikistani somoni
	TJS = "TJS"
	// THB - Thailand - Thai baht
	THB = "THB"
	// TMT - Turkmenistan - Turkmenistani manat
	TMT = "TMT"
	// TND - Tunisia - Tunisian dinar
	TND = "TND"
	// TOP - Tonga - Tongan paʻanga
	TOP = "TOP"
	// TRY - Turkey - Turkish lira
	TRY = "TRY"
	// TTD - Trinidad and Tobago - Trinidad and Tobago dollar
	TTD = "TTD"
	// TWD - Taiwan - New Taiwan dollar
	TWD = "TWD"
	// TZS - Tanzania - Tanzanian shilling
	TZS = "TZS"
	// UAH - Ukraine - Ukrainian hryvnia
	UAH = "UAH"
	// UGX - Uganda - Ugandan shilling
	UGX = "UGX"
	// USD - United States - United States dollar
	USD = "USD"
	// UYU - Uruguay - Uruguayan peso
	UYU = "UYU"
	// UZS - Uzbekistan - Uzbekistani soʻm
	UZS = "UZS"
	// VES - Venezuela - Venezuelan bolívar
	VES = "VES"
	// VND - Vietnam - Vietnamese đồng
	VND = "VND"
	// XAF - Central African CFA franc
	XAF = "XAF"
	// XOF - West African CFA franc
	XOF = "XOF"
	// YER - Yemen - Yemeni rial
	YER = "YER"
	// ZAR - South Africa - South African rand
	ZAR = "ZAR"
	// ZMW - Zambia - Zambian kwacha
	ZMW = "ZMW"
	// ZWL - Zimbabwe - Zimbabwean dollar
	ZWL = "ZWL"
)

// CurrencySymbols maps currency codes to currency symbols
var CurrencySymbols = map[string]string{
	AED: "د.إ",  // 阿联酋迪拉姆 - 阿拉伯联合酋长国
	AFN: "؋",    // 阿富汗尼 - 阿富汗
	ALL: "L",    // 列克 - 阿尔巴尼亚
	AMD: "֏",    // 德拉姆 - 亚美尼亚
	ANG: "ƒ",    // 荷属安的列斯盾 - 库拉索/荷属圣马丁
	AOA: "Kz",   // 宽扎 - 安哥拉
	ARS: "$",    // 比索 - 阿根廷
	AUD: "$",    // 澳元 - 澳大利亚
	AWG: "ƒ",    // 阿鲁巴弗罗林 - 阿鲁巴
	BAM: "KM",   // 可兑换马克 - 波斯尼亚和黑塞哥维那
	BDT: "৳",    // 塔卡 - 孟加拉国
	BGN: "лв",   // 列弗 - 保加利亚
	BHD: "د.ب",  // 第纳尔 - 巴林
	BIF: "Fr",   // 法郎 - 布隆迪
	BOB: "Bs",   // 玻利维亚诺 - 玻利维亚
	BRL: "R$",   // 雷亚尔 - 巴西
	BWP: "P",    // 普拉 - 博茨瓦纳
	BYN: "Br",   // 卢布 - 白俄罗斯
	CAD: "$",    // 加元 - 加拿大
	CDF: "FC",   // 刚果法郎 - 刚果（金）
	CHF: "Fr",   // 瑞士法郎 - 瑞士
	CLP: "$",    // 比索 - 智利
	CNY: "¥",    // 人民币 - 中国
	COP: "$",    // 比索 - 哥伦比亚
	CRC: "₡",    // 科朗 - 哥斯达黎加
	CUP: "₱",    // 比索 - 古巴
	CZK: "Kč",   // 克朗 - 捷克
	DKK: "kr",   // 克朗 - 丹麦
	DOP: "RD$",  // 比索 - 多米尼加共和国
	DZD: "د.ج",  // 第纳尔 - 阿尔及利亚
	EGP: "£",    // 埃及镑 - 埃及
	ETB: "Br",   // 比尔 - 埃塞俄比亚
	EUR: "€",    // 欧元 - 欧元区
	GBP: "£",    // 英镑 - 英国
	GEL: "₾",    // 拉里 - 格鲁吉亚
	GHS: "₵",    // 塞地 - 加纳
	GTQ: "Q",    // 格查尔 - 危地马拉
	HKD: "$",    // 港元 - 香港
	HNL: "L",    // 伦皮拉 - 洪都拉斯
	HTG: "G",    // 古德 - 海地
	HUF: "Ft",   // 福林 - 匈牙利
	IDR: "Rp",   // 卢比 - 印度尼西亚
	ILS: "₪",    // 新谢克尔 - 以色列
	INR: "₹",    // 卢比 - 印度
	IRR: "﷼",    // 里亚尔 - 伊朗
	JOD: "د.ا",  // 第纳尔 - 约旦
	JPY: "¥",    // 日元 - 日本
	KES: "KSh",  // 先令 - 肯尼亚
	KGS: "с",    // 索姆 - 吉尔吉斯斯坦
	KHR: "៛",    // 瑞尔 - 柬埔寨
	KRW: "₩",    // 韩元 - 韩国
	KWD: "د.ك",  // 第纳尔 - 科威特
	KZT: "₸",    // 坚戈 - 哈萨克斯坦
	LAK: "₭",    // 基普 - 老挝
	LBP: "ل.ل",  // 黎巴嫩镑 - 黎巴嫩
	LKR: "Rs",   // 卢比 - 斯里兰卡
	LRD: "$",    // 利比里亚元 - 利比里亚
	LYD: "ل.د",  // 第纳尔 - 利比亚
	MAD: "د.م.", // 迪拉姆 - 摩洛哥
	MDL: "L",    // 列伊 - 摩尔多瓦
	MGA: "Ar",   // 阿里亚里 - 马达加斯加
	MKD: "ден",  // 第纳尔 - 北马其顿
	MMK: "K",    // 缅元 - 缅甸
	MNT: "₮",    // 图格里克 - 蒙古
	MOP: "P",    // 澳门元 - 澳门
	MXN: "$",    // 比索 - 墨西哥
	MYR: "RM",   // 林吉特 - 马来西亚
	MWK: "MK",   // 克瓦查 - 马拉维
	NAD: "$",    // 纳米比亚元 - 纳米比亚
	NGN: "₦",    // 奈拉 - 尼日利亚
	NIO: "C$",   // 科尔多瓦 - 尼加拉瓜
	NOK: "kr",   // 克朗 - 挪威
	NPR: "₨",    // 卢比 - 尼泊尔
	NZD: "$",    // 纽元 - 新西兰
	OMR: "ر.ع.", // 里亚尔 - 阿曼
	PAB: "B/.",  // 巴波亚 - 巴拿马
	PEN: "S/.",  // 索尔 - 秘鲁
	PGK: "K",    // 基那 - 巴布亚新几内亚
	PHP: "₱",    // 比索 - 菲律宾
	PKR: "₨",    // 卢比 - 巴基斯坦
	PLN: "zł",   // 兹罗提 - 波兰
	PYG: "₲",    // 瓜拉尼 - 巴拉圭
	QAR: "ر.ق",  // 里亚尔 - 卡塔尔
	RON: "lei",  // 列伊 - 罗马尼亚
	RSD: "дин",  // 第纳尔 - 塞尔维亚
	RUB: "₽",    // 卢布 - 俄罗斯
	SAR: "﷼",    // 里亚尔 - 沙特阿拉伯
	SDG: "£",    // 镑 - 苏丹
	SEK: "kr",   // 克朗 - 瑞典
	SGD: "$",    // 新元 - 新加坡
	SRD: "$",    // 苏里南元 - 苏里南
	TJS: "ЅМ",   // 索莫尼 - 塔吉克斯坦
	THB: "฿",    // 泰铢 - 泰国
	TMT: "m",    // 马纳特 - 土库曼斯坦
	TND: "د.ت",  // 第纳尔 - 突尼斯
	TOP: "T$",   // 潘加 - 汤加
	TRY: "₺",    // 里拉 - 土耳其
	TTD: "TT$",  // 元 - 特立尼达和多巴哥
	TWD: "NT$",  // 新台币 - 台湾
	TZS: "Sh",   // 先令 - 坦桑尼亚
	UAH: "₴",    // 格里夫纳 - 乌克兰
	UGX: "USh",  // 先令 - 乌干达
	USD: "$",    // 美元 - 美国
	UYU: "$U",   // 比索 - 乌拉圭
	UZS: "so'm", // 苏姆 - 乌兹别克斯坦
	VES: "Bs.",  // 玻利瓦尔 - 委内瑞拉
	VND: "₫",    // 盾 - 越南
	XAF: "FCFA", // 中非法郎 - 中非经济共同体
	XOF: "CFA",  // 西非法郎 - 西非经济货币联盟
	YER: "﷼",    // 里亚尔 - 也门
	ZAR: "R",    // 兰特 - 南非
	ZMW: "ZK",   // 克瓦查 - 赞比亚
	ZWL: "$",    // 津巴布韦元 - 津巴布韦
}

// CurrencyMappingCountry maps currency codes to ISO 3166-1 alpha-2 country codes
var CurrencyMappingCountry = map[string]string{
	AED: "AE", // 阿拉伯联合酋长国
	AFN: "AF", // 阿富汗
	ALL: "AL", // 阿尔巴尼亚
	AMD: "AM", // 亚美尼亚
	ANG: "CW", // 库拉索
	AOA: "AO", // 安哥拉
	ARS: "AR", // 阿根廷
	AUD: "AU", // 澳大利亚
	AWG: "AW", // 阿鲁巴
	BAM: "BA", // 波斯尼亚和黑塞哥维那
	BDT: "BD", // 孟加拉国
	BGN: "BG", // 保加利亚
	BHD: "BH", // 巴林
	BIF: "BI", // 布隆迪
	BOB: "BO", // 玻利维亚
	BRL: "BR", // 巴西
	BWP: "BW", // 博茨瓦纳
	BYN: "BY", // 白俄罗斯
	CAD: "CA", // 加拿大
	CDF: "CD", // 刚果民主共和国
	CHF: "CH", // 瑞士
	CLP: "CL", // 智利
	CNY: "CN", // 中国
	COP: "CO", // 哥伦比亚
	CRC: "CR", // 哥斯达黎加
	CUP: "CU", // 古巴
	CZK: "CZ", // 捷克
	DKK: "DK", // 丹麦
	DOP: "DO", // 多米尼加共和国
	DZD: "DZ", // 阿尔及利亚
	EGP: "EG", // 埃及
	ETB: "ET", // 埃塞俄比亚
	EUR: "EU", // 欧元区
	GBP: "GB", // 英国
	GEL: "GE", // 格鲁吉亚
	GHS: "GH", // 加纳
	GTQ: "GT", // 危地马拉
	HKD: "HK", // 中国香港
	HNL: "HN", // 洪都拉斯
	HTG: "HT", // 海地
	HUF: "HU", // 匈牙利
	IDR: "ID", // 印度尼西亚
	ILS: "IL", // 以色列
	INR: "IN", // 印度
	IRR: "IR", // 伊朗
	JOD: "JO", // 约旦
	JPY: "JP", // 日本
	KES: "KE", // 肯尼亚
	KGS: "KG", // 吉尔吉斯斯坦
	KHR: "KH", // 柬埔寨
	KRW: "KR", // 韩国
	KWD: "KW", // 科威特
	KZT: "KZ", // 哈萨克斯坦
	LAK: "LA", // 老挝
	LBP: "LB", // 黎巴嫩
	LKR: "LK", // 斯里兰卡
	LRD: "LR", // 利比里亚
	LYD: "LY", // 利比亚
	MAD: "MA", // 摩洛哥
	MDL: "MD", // 摩尔多瓦
	MGA: "MG", // 马达加斯加
	MKD: "MK", // 北马其顿
	MMK: "MM", // 缅甸
	MNT: "MN", // 蒙古
	MOP: "MO", // 中国澳门
	MXN: "MX", // 墨西哥
	MYR: "MY", // 马来西亚
	MWK: "MW", // 马拉维
	NAD: "NA", // 纳米比亚
	NGN: "NG", // 尼日利亚
	NIO: "NI", // 尼加拉瓜
	NOK: "NO", // 挪威
	NPR: "NP", // 尼泊尔
	NZD: "NZ", // 新西兰
	OMR: "OM", // 阿曼
	PAB: "PA", // 巴拿马
	PEN: "PE", // 秘鲁
	PGK: "PG", // 巴布亚新几内亚
	PHP: "PH", // 菲律宾
	PKR: "PK", // 巴基斯坦
	PLN: "PL", // 波兰
	PYG: "PY", // 巴拉圭
	QAR: "QA", // 卡塔尔
	RON: "RO", // 罗马尼亚
	RSD: "RS", // 塞尔维亚
	RUB: "RU", // 俄罗斯
	SAR: "SA", // 沙特阿拉伯
	SDG: "SD", // 苏丹
	SEK: "SE", // 瑞典
	SGD: "SG", // 新加坡
	SRD: "SR", // 苏里南
	TJS: "TJ", // 塔吉克斯坦
	THB: "TH", // 泰国
	TMT: "TM", // 土库曼斯坦
	TND: "TN", // 突尼斯
	TOP: "TO", // 汤加
	TRY: "TR", // 土耳其
	TTD: "TT", // 特立尼达和多巴哥
	TWD: "TW", // 台湾
	TZS: "TZ", // 坦桑尼亚
	UAH: "UA", // 乌克兰
	UGX: "UG", // 乌干达
	USD: "US", // 美国
	UYU: "UY", // 乌拉圭
	UZS: "UZ", // 乌兹别克斯坦
	VES: "VE", // 委内瑞拉
	VND: "VN", // 越南
	XAF: "CM", // 喀麦隆（中非经济共同体成员）
	XOF: "CI", // 科特迪瓦（西非经济货币联盟成员）
	YER: "YE", // 也门
	ZAR: "ZA", // 南非
	ZMW: "ZM", // 赞比亚
	ZWL: "ZW", // 津巴布韦
}

// CurrencyTimeZoneMap maps currency codes to representative timezones
var CurrencyTimeZoneMap = map[string]string{
	AED: "Asia/Dubai",                     // 阿联酋迪拉姆 - 迪拜时间
	AFN: "Asia/Kabul",                     // 阿富汗尼 - 喀布尔时间
	ALL: "Europe/Tirane",                  // 列克 - 地拉那时间
	AMD: "Asia/Yerevan",                   // 德拉姆 - 埃里温时间
	ANG: "America/Curacao",                // 荷属安的列斯盾 - 库拉索时间
	AOA: "Africa/Luanda",                  // 宽扎 - 罗安达时间
	ARS: "America/Argentina/Buenos_Aires", // 比索 - 布宜诺斯艾利斯时间
	AUD: "Australia/Sydney",               // 澳元 - 悉尼时间
	AWG: "America/Aruba",                  // 阿鲁巴弗罗林 - 阿鲁巴时间
	BAM: "Europe/Sarajevo",                // 可兑换马克 - 萨拉热窝时间
	BDT: "Asia/Dhaka",                     // 塔卡 - 达卡时间
	BGN: "Europe/Sofia",                   // 列弗 - 索非亚时间
	BHD: "Asia/Bahrain",                   // 第纳尔 - 巴林时间
	BIF: "Africa/Bujumbura",               // 法郎 - 布琼布拉时间
	BOB: "America/La_Paz",                 // 玻利维亚诺 - 拉巴斯时间
	BRL: "America/Sao_Paulo",              // 雷亚尔 - 圣保罗时间
	BWP: "Africa/Gaborone",                // 普拉 - 哈博罗内时间
	BYN: "Europe/Minsk",                   // 卢布 - 明斯克时间
	CAD: "America/Toronto",                // 加元 - 多伦多时间
	CDF: "Africa/Kinshasa",                // 刚果法郎 - 金沙萨时间
	CHF: "Europe/Zurich",                  // 瑞士法郎 - 苏黎世时间
	CLP: "America/Santiago",               // 比索 - 圣地亚哥时间
	CNY: "Asia/Shanghai",                  // 人民币 - 上海时间
	COP: "America/Bogota",                 // 比索 - 波哥大时间
	CRC: "America/Costa_Rica",             // 科朗 - 圣何塞时间
	CUP: "America/Havana",                 // 比索 - 哈瓦那时间
	CZK: "Europe/Prague",                  // 克朗 - 布拉格时间
	DKK: "Europe/Copenhagen",              // 克朗 - 哥本哈根时间
	DOP: "America/Santo_Domingo",          // 比索 - 圣多明各时间
	DZD: "Africa/Algiers",                 // 第纳尔 - 阿尔及尔时间
	EGP: "Africa/Cairo",                   // 埃及镑 - 开罗时间
	ETB: "Africa/Addis_Ababa",             // 比尔 - 亚的斯亚贝巴时间
	EUR: "Europe/Berlin",                  // 欧元 - 柏林时间
	GBP: "Europe/London",                  // 英镑 - 伦敦时间
	GEL: "Asia/Tbilisi",                   // 拉里 - 第比利斯时间
	GHS: "Africa/Accra",                   // 塞地 - 阿克拉时间
	GTQ: "America/Guatemala",              // 格查尔 - 危地马拉城时间
	HKD: "Asia/Hong_Kong",                 // 港元 - 香港时间
	HNL: "America/Tegucigalpa",            // 伦皮拉 - 特古西加尔巴时间
	HTG: "America/Port-au-Prince",         // 古德 - 太子港时间
	HUF: "Europe/Budapest",                // 福林 - 布达佩斯时间
	IDR: "Asia/Jakarta",                   // 卢比 - 雅加达时间
	ILS: "Asia/Jerusalem",                 // 新谢克尔 - 耶路撒冷时间
	INR: "Asia/Kolkata",                   // 卢比 - 加尔各答时间
	IRR: "Asia/Tehran",                    // 里亚尔 - 德黑兰时间
	JOD: "Asia/Amman",                     // 第纳尔 - 安曼时间
	JPY: "Asia/Tokyo",                     // 日元 - 东京时间
	KES: "Africa/Nairobi",                 // 先令 - 内罗毕时间
	KGS: "Asia/Bishkek",                   // 索姆 - 比什凯克时间
	KHR: "Asia/Phnom_Penh",                // 瑞尔 - 金边时间
	KRW: "Asia/Seoul",                     // 韩元 - 首尔时间
	KWD: "Asia/Kuwait",                    // 第纳尔 - 科威特时间
	KZT: "Asia/Almaty",                    // 坚戈 - 阿拉木图时间
	LAK: "Asia/Vientiane",                 // 基普 - 万象时间
	LBP: "Asia/Beirut",                    // 黎巴嫩镑 - 贝鲁特时间
	LKR: "Asia/Colombo",                   // 卢比 - 科伦坡时间
	LRD: "Africa/Monrovia",                // 利比里亚元 - 蒙罗维亚时间
	LYD: "Africa/Tripoli",                 // 第纳尔 - 的黎波里时间
	MAD: "Africa/Casablanca",              // 迪拉姆 - 卡萨布兰卡时间
	MDL: "Europe/Chisinau",                // 列伊 - 基希讷乌时间
	MGA: "Indian/Antananarivo",            // 阿里亚里 - 塔那那利佛时间
	MKD: "Europe/Skopje",                  // 第纳尔 - 斯科普里时间
	MMK: "Asia/Yangon",                    // 缅元 - 仰光时间
	MNT: "Asia/Ulaanbaatar",               // 图格里克 - 乌兰巴托时间
	MOP: "Asia/Macau",                     // 澳门元 - 澳门时间
	MXN: "America/Mexico_City",            // 比索 - 墨西哥城时间
	MYR: "Asia/Kuala_Lumpur",              // 林吉特 - 吉隆坡时间
	MWK: "Africa/Blantyre",                // 克瓦查 - 布兰太尔时间
	NAD: "Africa/Windhoek",                // 纳米比亚元 - 温得和克时间
	NGN: "Africa/Lagos",                   // 奈拉 - 拉各斯时间
	NIO: "America/Managua",                // 科尔多瓦 - 马那瓜时间
	NOK: "Europe/Oslo",                    // 克朗 - 奥斯陆时间
	NPR: "Asia/Kathmandu",                 // 卢比 - 加德满都时间
	NZD: "Pacific/Auckland",               // 纽元 - 奥克兰时间
	OMR: "Asia/Muscat",                    // 里亚尔 - 马斯喀特时间
	PAB: "America/Panama",                 // 巴波亚 - 巴拿马城时间
	PEN: "America/Lima",                   // 索尔 - 利马时间
	PGK: "Pacific/Port_Moresby",           // 基那 - 莫尔兹比港时间
	PHP: "Asia/Manila",                    // 比索 - 马尼拉时间
	PKR: "Asia/Karachi",                   // 卢比 - 卡拉奇时间
	PLN: "Europe/Warsaw",                  // 兹罗提 - 华沙时间
	PYG: "America/Asuncion",               // 瓜拉尼 - 亚松森时间
	QAR: "Asia/Qatar",                     // 里亚尔 - 多哈时间
	RON: "Europe/Bucharest",               // 列伊 - 布加勒斯特时间
	RSD: "Europe/Belgrade",                // 第纳尔 - 贝尔格莱德时间
	RUB: "Europe/Moscow",                  // 卢布 - 莫斯科时间
	SAR: "Asia/Riyadh",                    // 里亚尔 - 利雅得时间
	SDG: "Africa/Khartoum",                // 镑 - 喀土穆时间
	SEK: "Europe/Stockholm",               // 克朗 - 斯德哥尔摩时间
	SGD: "Asia/Singapore",                 // 新元 - 新加坡时间
	SRD: "America/Paramaribo",             // 苏里南元 - 帕拉马里博时间
	TJS: "Asia/Dushanbe",                  // 索莫尼 - 杜尚别时间
	THB: "Asia/Bangkok",                   // 泰铢 - 曼谷时间
	TMT: "Asia/Ashgabat",                  // 马纳特 - 阿什哈巴德时间
	TND: "Africa/Tunis",                   // 第纳尔 - 突尼斯时间
	TOP: "Pacific/Tongatapu",              // 潘加 - 努库阿洛法时间
	TRY: "Europe/Istanbul",                // 里拉 - 伊斯坦布尔时间
	TTD: "America/Port_of_Spain",          // 元 - 西班牙港时间
	TWD: "Asia/Taipei",                    // 新台币 - 台北时间
	TZS: "Africa/Dar_es_Salaam",           // 先令 - 达累斯萨拉姆时间
	UAH: "Europe/Kyiv",                    // 格里夫纳 - 基辅时间
	UGX: "Africa/Kampala",                 // 先令 - 坎帕拉时间
	USD: "America/New_York",               // 美元 - 纽约时间
	UYU: "America/Montevideo",             // 比索 - 蒙得维的亚时间
	UZS: "Asia/Tashkent",                  // 苏姆 - 塔什干时间
	VES: "America/Caracas",                // 玻利瓦尔 - 加拉加斯时间
	VND: "Asia/Ho_Chi_Minh",               // 盾 - 胡志明市时间
	XAF: "Africa/Douala",                  // 中非法郎 - 杜阿拉时间
	XOF: "Africa/Abidjan",                 // 西非法郎 - 阿比让时间
	YER: "Asia/Aden",                      // 里亚尔 - 亚丁时间
	ZAR: "Africa/Johannesburg",            // 兰特 - 约翰内斯堡时间
	ZMW: "Africa/Lusaka",                  // 克瓦查 - 卢萨卡时间
	ZWL: "Africa/Harare",                  // 津巴布韦元 - 哈拉雷时间
}

// CountryMappingCurrency maps ISO 3166-1 alpha-2 country codes to their primary currency codes.
// Note: For countries with multiple currencies, this map represents one common or primary currency.
// Special cases like the Eurozone (EU) are handled based on the original mapping.
var CountryMappingCurrency = map[string]string{
	"AE": AED, // United Arab Emirates
	"AF": AFN, // Afghanistan
	"AL": ALL, // Albania
	"AM": AMD, // Armenia
	"AO": AOA, // Angola
	"AR": ARS, // Argentina
	"AU": AUD, // Australia
	"AW": AWG, // Aruba
	"BA": BAM, // Bosnia and Herzegovina
	"BD": BDT, // Bangladesh
	"BG": BGN, // Bulgaria
	"BH": BHD, // Bahrain
	"BI": BIF, // Burundi
	"BO": BOB, // Bolivia
	"BR": BRL, // Brazil
	"BW": BWP, // Botswana
	"BY": BYN, // Belarus
	"CA": CAD, // Canada
	"CD": CDF, // Democratic Republic of the Congo
	"CH": CHF, // Switzerland
	"CL": CLP, // Chile
	"CN": CNY, // China
	"CO": COP, // Colombia
	"CR": CRC, // Costa Rica
	"CU": CUP, // Cuba
	"CZ": CZK, // Czech Republic
	"DK": DKK, // Denmark
	"DO": DOP, // Dominican Republic
	"DZ": DZD, // Algeria
	"EG": EGP, // Egypt
	"ET": ETB, // Ethiopia
	"EU": EUR, // Eurozone - Special case, not a country but a region
	"GB": GBP, // United Kingdom
	"GE": GEL, // Georgia
	"GH": GHS, // Ghana
	"GT": GTQ, // Guatemala
	"HK": HKD, // Hong Kong SAR China
	"HN": HNL, // Honduras
	"HT": HTG, // Haiti
	"HU": HUF, // Hungary
	"ID": IDR, // Indonesia
	"IL": ILS, // Israel
	"IN": INR, // India
	"IR": IRR, // Iran
	"JO": JOD, // Jordan
	"JP": JPY, // Japan
	"KE": KES, // Kenya
	"KG": KGS, // Kyrgyzstan
	"KH": KHR, // Cambodia
	"KR": KRW, // South Korea
	"KW": KWD, // Kuwait
	"KZ": KZT, // Kazakhstan
	"LA": LAK, // Laos
	"LB": LBP, // Lebanon
	"LK": LKR, // Sri Lanka
	"LR": LRD, // Liberia
	"LY": LYD, // Libya
	"MA": MAD, // Morocco
	"MD": MDL, // Moldova
	"MG": MGA, // Madagascar
	"MK": MKD, // North Macedonia
	"MM": MMK, // Myanmar (Burma)
	"MN": MNT, // Mongolia
	"MO": MOP, // Macao SAR China
	"MX": MXN, // Mexico
	"MY": MYR, // Malaysia
	"MW": MWK, // Malawi
	"NA": NAD, // Namibia
	"NG": NGN, // Nigeria
	"NI": NIO, // Nicaragua
	"NO": NOK, // Norway
	"NP": NPR, // Nepal
	"NZ": NZD, // New Zealand
	"OM": OMR, // Oman
	"PA": PAB, // Panama
	"PE": PEN, // Peru
	"PG": PGK, // Papua New Guinea
	"PH": PHP, // Philippines
	"PK": PKR, // Pakistan
	"PL": PLN, // Poland
	"PY": PYG, // Paraguay
	"QA": QAR, // Qatar
	"RO": RON, // Romania
	"RS": RSD, // Serbia
	"RU": RUB, // Russia
	"SA": SAR, // Saudi Arabia
	"SD": SDG, // Sudan
	"SE": SEK, // Sweden
	"SG": SGD, // Singapore
	"SR": SRD, // Suriname
	"TJ": TJS, // Tajikistan
	"TH": THB, // Thailand
	"TM": TMT, // Turkmenistan
	"TN": TND, // Tunisia
	"TO": TOP, // Tonga
	"TR": TRY, // Turkey
	"TT": TTD, // Trinidad and Tobago
	"TW": TWD, // Taiwan
	"TZ": TZS, // Tanzania
	"UA": UAH, // Ukraine
	"UG": UGX, // Uganda
	"US": USD, // United States
	"UY": UYU, // Uruguay
	"UZ": UZS, // Uzbekistan
	"VE": VES, // Venezuela
	"VN": VND, // Vietnam
	"CM": XAF, // Cameroon (representing XAF region)
	"CI": XOF, // Côte d'Ivoire (representing XOF region)
	"YE": YER, // Yemen
	"ZA": ZAR, // South Africa
	"ZM": ZMW, // Zambia
	"ZW": ZWL, // Zimbabwe
	// Note: ANG (Netherlands Antillean Guilder) is mapped to CW (Curaçao) in the original.
	"CW": ANG, // Curaçao
}

// GetDayRangeByCurrency 根据币种代码获取当天的开始和结束时间。
// 如果找不到对应的时区，则返回错误。
// GetDayRangeByCurrency gets the start and end time of the current day based on the currency code.
// It returns an error if the corresponding timezone is not found.
func GetDayRangeByCurrency(currency string) (startTime, endTime time.Time, err error) {
	// 1. 根据币种获取时区名称
	// 1. Get the timezone name based on the currency
	tzName, ok := CurrencyTimeZoneMap[currency]
	if !ok {
		return time.Time{}, time.Time{}, fmt.Errorf("未知的币种: %s", currency)
	}

	// 2. 加载时区
	// 2. Load the timezone
	loc, err := time.LoadLocation(tzName)
	if err != nil {
		return time.Time{}, time.Time{}, fmt.Errorf("加载时区失败: %w", err)
	}

	// 3. 获取该时区下的当前时间
	// 3. Get the current time in that timezone
	now := time.Now().In(loc)

	// 4. 计算当天的开始时间（零点）
	// 4. Calculate the start time of the day (midnight)
	startTime = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, loc)

	// 5. 计算当天的结束时间（下一天零点前一纳秒）
	// 5. Calculate the end time of the day (one nanosecond before the next day's midnight)
	endTime = startTime.Add(24 * time.Hour).Add(-time.Nanosecond)

	return startTime, endTime, nil
}
