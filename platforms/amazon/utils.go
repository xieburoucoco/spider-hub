package amazon

import (
	"regexp"
	"strconv"
	"strings"
)

type AmazonUtil struct{}

// GetAllValues 根据给定的键名，从 CountriesInfo 列表中取出相应的值
func (a *AmazonUtil) GetAllValues(keyName string) map[string]bool {
	result := make(map[string]bool)
	switch keyName {
	case "country":
		for _, country := range CountriesInfo {
			result[country.Country] = true
		}
	case "currency_symbol":
		for _, country := range CountriesInfo {
			result[country.CurrencySymbol] = true
		}
	case "language_code":
		for _, country := range CountriesInfo {
			result[country.LanguageCode] = true
		}
	case "whisper_language_code":
		for _, country := range CountriesInfo {
			result[country.WhisperLanguageCode] = true
		}
	case "dify_language":
		for _, country := range CountriesInfo {
			result[country.DifyLanguage] = true
		}
	case "tts_voice_name":
		for _, country := range CountriesInfo {
			result[country.TTSVoiceName] = true
		}
	default:
		for _, country := range CountriesInfo {
			result[country.LanguageCode] = true
		}
	}
	return result
}

// FindValue 在给定的国家信息列表中搜索与给定的搜索值匹配的字典，并返回与给定的键名对应的值
func (a *AmazonUtil) FindValue(searchValue, searchKey, valueKey string) string {
	searchValue = strings.ToUpper(searchValue)
	for _, country := range CountriesInfo {
		var compareValue string
		switch searchKey {
		case "country":
			compareValue = strings.ToUpper(country.Country)
		case "currency_symbol":
			compareValue = strings.ToUpper(country.CurrencySymbol)
		case "language_code":
			compareValue = strings.ToUpper(country.LanguageCode)
		case "whisper_language_code":
			compareValue = strings.ToUpper(country.WhisperLanguageCode)
		case "dify_language":
			compareValue = strings.ToUpper(country.DifyLanguage)
		case "tts_voice_name":
			compareValue = strings.ToUpper(country.TTSVoiceName)
		}
		if compareValue == searchValue {
			switch valueKey {
			case "country":
				return country.Country
			case "currency_symbol":
				return country.CurrencySymbol
			case "language_code":
				return country.LanguageCode
			case "whisper_language_code":
				return country.WhisperLanguageCode
			case "dify_language":
				return country.DifyLanguage
			case "tts_voice_name":
				return country.TTSVoiceName
			}
		}
	}
	return CountriesInfo[1].LanguageCode
}

// MatchCurrency 在给定的文本中使用正则表达式搜索货币单位
func (a *AmazonUtil) MatchCurrency(text string) string {
	currencyRegex := regexp.MustCompile(`(\$|￥|€|₹|Rp|¥|₩|RM|रू|฿|₫)`)
	matches := currencyRegex.FindStringSubmatch(text)
	if len(matches) > 1 {
		return matches[1]
	}
	return "$"
}

// FindValueForCurrency 先用 MatchCurrency 找货币单位，再用 FindValue 找对应值
func (a *AmazonUtil) FindValueForCurrency(text, valueKey string) string {
	currency := a.MatchCurrency(text)
	return a.FindValue(currency, "currency_symbol", valueKey)
}

// ExtractPrice 从字符串中提取价格并转为浮点数
func (a *AmazonUtil) ExtractPrice(priceStr string) *float64 {
	pattern := regexp.MustCompile(`\d+(,\d{3})*(\.\d+)?`)
	match := pattern.FindString(priceStr)
	if match != "" {
		priceNoComma := strings.ReplaceAll(match, ",", "")
		if price, err := strconv.ParseFloat(priceNoComma, 64); err == nil {
			return &price
		}
	}
	return nil
}

// URLCheck 检查URL列表中的Amazon链接
func (a *AmazonUtil) URLCheck(urlList []string) []string {
	pattern := regexp.MustCompile(AmazonPattern)
	var checkURLList []string
	for _, url := range urlList {
		if pattern.MatchString(url) {
			checkURLList = append(checkURLList, url)
		}
	}
	return checkURLList
}
