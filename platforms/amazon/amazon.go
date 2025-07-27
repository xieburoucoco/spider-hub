package amazon

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
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

// 🔍 SearchProductsByImageURL 通过在线图片URL搜索相关商品
//
// 这个接口通过图片URL搜索亚马逊上的相关商品
//
// 参数:
//   - ctx: 上下文，用于控制请求的生命周期
//   - imageURL: 图片的URL地址
//   - proxies: 代理配置（可选）
//
// 返回:
//   - []ImageSearchProduct: 搜索到的商品列表
//   - error: 错误信息，如果没有错误则为nil
func (s *AmazonSpider) SearchProductsByImageURL(ctx context.Context, imageURL string, proxies map[string]string) ([]ImageSearchProduct, error) {
	for attempt := 0; attempt < MaxRetries; attempt++ {
		// 获取stylesnap值
		stylesnapValue, err := s.getStylesnapValue(ctx, proxies)
		if err != nil {
			if attempt == MaxRetries-1 {
				return nil, fmt.Errorf("❌ 获取stylesnap值失败: %w", err)
			}
			continue
		}

		// 下载图片
		imageData, err := s.downloadImage(ctx, imageURL)
		if err != nil {
			if attempt == MaxRetries-1 {
				return nil, fmt.Errorf("❌ 下载图片失败: %w", err)
			}
			continue
		}

		// 上传图片并搜索
		products, err := s.uploadImageAndSearch(ctx, imageData, stylesnapValue, proxies)
		if err != nil {
			if attempt == MaxRetries-1 {
				return nil, fmt.Errorf("❌ 图片搜索失败: %w", err)
			}
			continue
		}

		return products, nil
	}

	return nil, fmt.Errorf("❌ 达到最大重试次数，搜索失败")
}

// 🔍 SearchProductsByImageData 通过本地图片数据搜索相关商品
//
// 这个接口通过本地图片文件数据搜索亚马逊上的相关商品
//
// 参数:
//   - ctx: 上下文，用于控制请求的生命周期
//   - imageData: 图片的二进制数据
//   - proxies: 代理配置（可选）
//
// 返回:
//   - []ImageSearchProduct: 搜索到的商品列表
//   - error: 错误信息，如果没有错误则为nil
func (s *AmazonSpider) SearchProductsByImageData(ctx context.Context, imageData []byte, proxies map[string]string) ([]ImageSearchProduct, error) {
	for attempt := 0; attempt < MaxRetries; attempt++ {
		// 获取stylesnap值
		stylesnapValue, err := s.getStylesnapValue(ctx, proxies)
		if err != nil {
			if attempt == MaxRetries-1 {
				return nil, fmt.Errorf("❌ 获取stylesnap值失败: %w", err)
			}
			continue
		}

		// 上传图片并搜索
		products, err := s.uploadImageAndSearch(ctx, imageData, stylesnapValue, proxies)
		if err != nil {
			if attempt == MaxRetries-1 {
				return nil, fmt.Errorf("❌ 图片搜索失败: %w", err)
			}
			continue
		}

		return products, nil
	}

	return nil, fmt.Errorf("❌ 达到最大重试次数，搜索失败")
}

// 🔧 getStylesnapValue 获取stylesnap值用于图片上传请求
func (s *AmazonSpider) getStylesnapValue(ctx context.Context, proxies map[string]string) (string, error) {
	// 构建请求
	req := s.client.R().SetContext(ctx)

	// 设置代理
	if len(proxies) > 0 {
		for scheme, proxy := range proxies {
			if scheme == "http" || scheme == "https" {
				s.client.SetProxy(proxy)
				break
			}
		}
	}

	// 发送请求
	resp, err := req.Get(AmazonShopLookURL)
	if err != nil {
		return "", fmt.Errorf("请求失败: %w", err)
	}

	if resp.StatusCode() != 200 {
		return "", fmt.Errorf("状态码: %d", resp.StatusCode())
	}

	// 解析HTML获取stylesnap值
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(string(resp.Body())))
	if err != nil {
		return "", fmt.Errorf("解析HTML失败: %w", err)
	}

	var stylesnapValue string
	doc.Find("body > input[name=stylesnap]").Each(func(i int, s *goquery.Selection) {
		if val, exists := s.Attr("value"); exists {
			stylesnapValue = val
		}
	})

	if stylesnapValue == "" {
		return "", fmt.Errorf("找不到stylesnap，请重试")
	}

	fmt.Printf("请求头加密参数：【%s】\n", stylesnapValue)
	return stylesnapValue, nil
}

// 🔧 downloadImage 下载图片并返回二进制数据
func (s *AmazonSpider) downloadImage(ctx context.Context, imageURL string) ([]byte, error) {
	resp, err := s.client.R().
		SetContext(ctx).
		Get(imageURL)

	if err != nil {
		return nil, fmt.Errorf("下载图片失败: %w", err)
	}

	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("无法下载图片，状态码: %d", resp.StatusCode())
	}

	return resp.Body(), nil
}

// 🔧 uploadImageAndSearch 上传图片并获取搜索结果
func (s *AmazonSpider) uploadImageAndSearch(ctx context.Context, imageData []byte, stylesnapValue string, proxies map[string]string) ([]ImageSearchProduct, error) {
	// 设置代理
	if len(proxies) > 0 {
		for scheme, proxy := range proxies {
			if scheme == "http" || scheme == "https" {
				s.client.SetProxy(proxy)
				break
			}
		}
	}

	// 构建请求
	req := s.client.R().
		SetContext(ctx).
		SetQueryParam("stylesnapToken", stylesnapValue).
		SetFileReader("explore-looks.jpg", "explore-looks.jpg", bytes.NewReader(imageData))

	// 发送POST请求
	resp, err := req.Post(AmazonStyleSnapUploadURL)
	if err != nil {
		return nil, fmt.Errorf("上传图片失败: %w", err)
	}

	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("图片上传失败，状态码: %d", resp.StatusCode())
	}

	// 解析响应
	var styleSnapResp StyleSnapResponse
	if err := json.Unmarshal(resp.Body(), &styleSnapResp); err != nil {
		return nil, fmt.Errorf("解析响应失败: %w", err)
	}

	// 提取搜索结果
	var products []ImageSearchProduct
	if len(styleSnapResp.SearchResults) > 0 {
		products = styleSnapResp.SearchResults[0].BBXASINMetadataList
	}

	// 添加商品链接
	for i := range products {
		products[i].LinkURL = AmazonProductDetailPrefixURL + products[i].ASIN
	}

	return products, nil
}
