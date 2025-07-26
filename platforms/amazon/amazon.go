package amazon

import (
	"context"
	"fmt"
	"time"

	"github.com/go-resty/resty/v2"
)

// 🕷️ AmazonSpider 亚马逊爬虫结构体
// 负责爬取亚马逊商品详情页面并解析数据
type AmazonSpider struct {
	client    *resty.Client    // HTTP客户端
	extractor *AmazonExtractor // 数据提取器
	util      *AmazonUtil      // 工具函数
}

// 🏭 NewAmazonSpider 创建新的亚马逊爬虫实例
//
// 返回一个配置好的亚马逊爬虫，包含:
// - 配置了超时和重试的resty客户端
// - 数据提取器
// - 工具函数集
func NewAmazonSpider() *AmazonSpider {
	// 创建resty客户端
	client := resty.New().
		SetTimeout(time.Duration(RequestTimeout) * time.Second).
		SetRetryCount(3).
		SetRetryWaitTime(5 * time.Second).
		SetRetryMaxWaitTime(20 * time.Second)

	// 设置默认请求头
	for key, value := range AmazonHeaders {
		client.SetHeader(key, value)
	}

	return &AmazonSpider{
		client:    client,
		extractor: NewAmazonExtractor(),
		util:      &AmazonUtil{},
	}
}

// 🔍 FetchProductDetail 获取亚马逊商品详情
//
// 这是暴露给外部调用的主要接口，通过商品URL获取商品的详细信息
//
// 参数:
//   - ctx: 上下文，用于控制请求的生命周期
//   - productURL: 亚马逊商品URL
//
// 返回:
//   - ProductResult: 商品详情结果
//   - error: 错误信息，如果没有错误则为nil
func (s *AmazonSpider) FetchProductDetail(ctx context.Context, productURL string) (ProductResult, error) {
	// ✅ 验证URL是否为有效的亚马逊商品链接
	if validURLs := s.util.URLCheck([]string{productURL}); len(validURLs) == 0 {
		return ProductResult{}, fmt.Errorf("❌ 无效的亚马逊URL: %s", productURL)
	}

	// 📡 发送HTTP请求获取页面内容
	resp, err := s.client.R().
		SetContext(ctx).
		Get(productURL)

	if err != nil {
		return ProductResult{}, fmt.Errorf("❌ 请求失败: %w", err)
	}

	if resp.StatusCode() != 200 {
		return ProductResult{}, fmt.Errorf("❌ HTTP请求失败，状态码: %d", resp.StatusCode())
	}

	// 📝 提取商品详情
	result := s.extractor.GetProductDetail(productURL, string(resp.Body()))
	return result, nil
}
