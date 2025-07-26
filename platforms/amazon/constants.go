package amazon

// 参数常量
const (
	Keyword  = "keyword"
	Page     = "page"
	PageSize = "page_size"
)

// 返回格式常量
const (
	LinkURL  = "link_url"
	Title    = "title"
	Desc     = "desc"
	Language = "language"
	Images   = "images"
	Videos   = "videos"
	Price    = "price"
	Discount = "discount"
)

// HTTP请求超时时间
const (
	HttpxRequestTimeout = 60
	RequestTimeout      = 30
)

// Amazon请求头
var AmazonHeaders = map[string]string{
	"accept":                     "*/*",
	"accept-language":            "en-US,en;q=0.6",
	"cache-control":              "no-cache",
	"pragma":                     "no-cache",
	"priority":                   "u=1, i",
	"referer":                    "https://www.amazon.com",
	"sec-ch-ua":                  `"Chromium";v="124", "Brave";v="124", "Not-A.Brand";v="99"`,
	"sec-ch-ua-mobile":           "?0",
	"sec-ch-ua-platform":         `"Windows"`,
	"sec-ch-ua-platform-version": `"10.0.0"`,
	"sec-fetch-dest":             "empty",
	"sec-fetch-mode":             "cors",
	"sec-fetch-site":             "same-origin",
	"sec-gpc":                    "1",
	"user-agent":                 "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/124.0.0.0 Safari/537.36",
}

// 链接校验正则表达式
const AmazonPattern = `https://(www|us)\.amazon.*?/.*?(d|dp)/.*?`

// 任务切块大小
const (
	ChunkSize       = 5
	VideoThresholds = 100
	ImageThresholds = 100
)

// 请求类型
const (
	Httpx    = "HTTPX"
	Requests = "REQUESTS"
)

// 国家信息结构
type CountryInfo struct {
	Country             string `json:"country"`
	CurrencySymbol      string `json:"currency_symbol"`
	Language            string `json:"language"`
	LanguageCode        string `json:"language_code"`
	WhisperLanguageCode string `json:"whisper_language_code"`
	DifyLanguage        string `json:"dify_language"`
	TTSVoiceName        string `json:"tts_voice_name"`
}

// 国家信息列表
var CountriesInfo = []CountryInfo{
	{Country: "中国", CurrencySymbol: "￥", Language: "中文", LanguageCode: "ZH", WhisperLanguageCode: "zh", DifyLanguage: "中文", TTSVoiceName: "zh-CN-XiaoxiaoNeural"},
	{Country: "美国", CurrencySymbol: "$", Language: "英文", LanguageCode: "EN", WhisperLanguageCode: "en", DifyLanguage: "英语", TTSVoiceName: "en-US-GuyNeural"},
	{Country: "法国", CurrencySymbol: "€", Language: "法语", LanguageCode: "FRA", WhisperLanguageCode: "fr", DifyLanguage: "法语", TTSVoiceName: "fr-FR-DeniseNeural"},
	{Country: "印度", CurrencySymbol: "₹", Language: "印地语", LanguageCode: "HI", WhisperLanguageCode: "hi", DifyLanguage: "印地语", TTSVoiceName: "hi-IN-MadhurNeural"},
	{Country: "印度尼西亚", CurrencySymbol: "Rp", Language: "印尼语", LanguageCode: "ID", WhisperLanguageCode: "id", DifyLanguage: "印尼语", TTSVoiceName: "id-ID-ArdiNeural"},
	{Country: "日本", CurrencySymbol: "¥", Language: "日语", LanguageCode: "JP", WhisperLanguageCode: "ja", DifyLanguage: "日语", TTSVoiceName: "ja-JP-NanamiNeural"},
	{Country: "韩国", CurrencySymbol: "₩", Language: "韩语", LanguageCode: "KOR", WhisperLanguageCode: "ko", DifyLanguage: "韩语", TTSVoiceName: "ko-KR-SunHiNeural"},
	{Country: "马来西亚", CurrencySymbol: "RM", Language: "马来语", LanguageCode: "MS", WhisperLanguageCode: "ms", DifyLanguage: "马来语", TTSVoiceName: "ms-MY-OsmanNeural"},
	{Country: "尼泊尔", CurrencySymbol: "रू", Language: "尼泊尔语", LanguageCode: "NE", WhisperLanguageCode: "ne", DifyLanguage: "尼泊尔语", TTSVoiceName: "ne-NP-SagarNeural"},
	{Country: "西班牙", CurrencySymbol: "€", Language: "西班牙语", LanguageCode: "ES", WhisperLanguageCode: "es", DifyLanguage: "西班牙语", TTSVoiceName: "es-ES-ElviraNeural"},
	{Country: "泰国", CurrencySymbol: "฿", Language: "泰语", LanguageCode: "TH", WhisperLanguageCode: "th", DifyLanguage: "泰语", TTSVoiceName: "th-TH-NiwatNeural"},
	{Country: "越南", CurrencySymbol: "₫", Language: "越南语", LanguageCode: "VIE", WhisperLanguageCode: "vi", DifyLanguage: "越南语", TTSVoiceName: "vi-VN-NamMinhNeural"},
}
