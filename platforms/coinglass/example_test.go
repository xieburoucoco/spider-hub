package coinglass

import (
	"fmt"
	"testing"
)

// 创建爬虫实例
var spider = NewSpider()

// TestCoinglassAPI 测试Coinglass API
func TestCoinglassAPI(t *testing.T) {

	// 测试币种市场数据API
	apiURL := "https://capi.coinglass.com/api/home/v2/coinMarkets?sort=&order=&keyword=&pageNum=1&pageSize=5&ex=all"

	fmt.Printf("🕷️ 测试 Coinglass API\n")
	fmt.Printf("URL: %s\n", apiURL)

	// 获取数据
	data, err := spider.GetData(apiURL)
	if err != nil {
		t.Errorf("获取数据失败: %v", err)
		return
	}

	preview := getDataPreview(data, 100)
	t.Logf("✅ 市场数据获取成功，数据长度: %d, 数据预览: %s", len(data), preview)
}

// 🔍 TestOpenInterestChart 测试持仓量图表数据接口
// 📊 获取指定币种的持仓量变化趋势数据
func TestOpenInterestChart(t *testing.T) {
	// 📋 构建查询参数
	params := map[string]string{
		"symbol":       "BTC", // 🪙 币种符号
		"timeType":     "0",   // ⏰ 时间类型
		"exchangeName": "",    // 🏢 交易所名称（空表示全部）
		"currency":     "USD", // 💵 计价货币
		"type":         "0",   // 📈 数据类型
	}

	// 🌐 构建完整的API URL
	apiURL := "https://capi.coinglass.com/api/openInterest/v3/chart"
	fullURL := buildURLWithParams(apiURL, params)

	// 🚀 发起请求获取数据
	result, err := spider.GetData(fullURL)
	if err != nil {
		t.Fatalf("❌ 获取持仓量图表数据失败: %v", err)
	}

	preview := getDataPreview(result, 100)
	t.Logf("✅ 持仓量图表数据获取成功，数据长度: %d, 数据预览: %s", len(result), preview)
}

// 📊 TestFuturesHomeStatistics 测试期货首页统计数据接口
// 🏠 获取期货市场的整体统计信息和概览数据
func TestFuturesHomeStatistics(t *testing.T) {
	// 🌐 构建API URL（无需额外参数）
	apiURL := "https://capi.coinglass.com/api/futures/home/statistics"

	// 🚀 发起请求获取数据
	result, err := spider.GetData(apiURL)
	if err != nil {
		t.Fatalf("❌ 获取期货首页统计数据失败: %v", err)
	}

	preview := getDataPreview(result, 100)
	t.Logf("✅ 期货首页统计数据获取成功，数据长度: %d, 数据预览: %s", len(result), preview)
}

// 🏢 TestDerivativeExchangeList 测试衍生品交易所列表接口
// 📋 获取支持衍生品交易的所有交易所信息
func TestDerivativeExchangeList(t *testing.T) {
	// 🌐 构建API URL（无需额外参数）
	apiURL := "https://capi.coinglass.com/api/derivative/exchange/list"

	// 🚀 发起请求获取数据
	result, err := spider.GetData(apiURL)
	if err != nil {
		t.Fatalf("❌ 获取衍生品交易所列表失败: %v", err)
	}

	preview := getDataPreview(result, 100)
	t.Logf("✅ 衍生品交易所列表获取成功，数据长度: %d, 数据预览: %s", len(result), preview)
}

// 💰 TestCoinMarkets 测试币种市场数据接口
// 📈 获取币种的市场排行和详细交易数据
func TestCoinMarkets(t *testing.T) {
	// 📋 构建查询参数
	params := map[string]string{
		"sort":     "",    // 🔄 排序字段（空表示默认排序）
		"order":    "",    // ⬆️ 排序方向（空表示默认方向）
		"keyword":  "",    // 🔍 搜索关键词（空表示不筛选）
		"pageNum":  "1",   // 📄 页码
		"pageSize": "20",  // 📊 每页数量
		"ex":       "all", // 🏢 交易所筛选（all表示全部）
	}

	// 🌐 构建完整的API URL
	apiURL := "https://capi.coinglass.com/api/home/v2/coinMarkets"
	fullURL := buildURLWithParams(apiURL, params)

	// 🚀 发起请求获取数据
	result, err := spider.GetData(fullURL)
	if err != nil {
		t.Fatalf("❌ 获取币种市场数据失败: %v", err)
	}

	preview := getDataPreview(result, 100)
	t.Logf("✅ 币种市场数据获取成功，数据长度: %d, 数据预览: %s", len(result), preview)
}

// 🔗 TestExchangeFuturesPairInfo 测试交易所期货交易对信息接口
// 📊 获取各大交易所的期货交易对详细配置信息
func TestExchangeFuturesPairInfo(t *testing.T) {
	// 🌐 构建API URL（无需额外参数）
	apiURL := "https://capi.coinglass.com/api/exchange/futures/pairInfo"

	// 🚀 发起请求获取数据
	result, err := spider.GetData(apiURL)
	if err != nil {
		t.Fatalf("❌ 获取交易所期货交易对信息失败: %v", err)
	}

	preview := getDataPreview(result, 100)
	t.Logf("✅ 交易所期货交易对信息获取成功，数据长度: %d, 数据预览: %s", len(result), preview)
}

// 🪙 TestSpotSupportCoin 测试现货支持币种接口
// 📋 获取平台支持的所有现货交易币种列表
func TestSpotSupportCoin(t *testing.T) {
	// 🌐 构建API URL（无需额外参数）
	apiURL := "https://capi.coinglass.com/api/spot/support/coin"

	// 🚀 发起请求获取数据
	result, err := spider.GetData(apiURL)
	if err != nil {
		t.Fatalf("❌ 获取现货支持币种失败: %v", err)
	}

	preview := getDataPreview(result, 100)
	t.Logf("✅ 现货支持币种获取成功，数据长度: %d, 数据预览: %s", len(result), preview)
}

// 🔧 buildURLWithParams 构建带参数的URL
// 📝 将参数映射转换为URL查询字符串并拼接到基础URL上
func buildURLWithParams(baseURL string, params map[string]string) string {
	if len(params) == 0 {
		return baseURL
	}

	url := baseURL + "?"
	first := true
	for key, value := range params {
		if !first {
			url += "&"
		}
		url += key + "=" + value
		first = false
	}
	return url
}

// 👀 getDataPreview 获取数据预览
// 📄 截取指定长度的数据内容用于预览，避免输出过长的数据
func getDataPreview(data string, maxLength int) string {
	if len(data) <= maxLength {
		return data
	}
	return data[:maxLength] + "..."
}
