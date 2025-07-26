package amazon

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// AmazonExtractor Amazon数据提取器
type AmazonExtractor struct {
	util *AmazonUtil
}

// NewAmazonExtractor 创建新的Amazon提取器
func NewAmazonExtractor() *AmazonExtractor {
	return &AmazonExtractor{
		util: &AmazonUtil{},
	}
}

// getInnerText 获取HTML元素的文本内容 (私有方法)
func (e *AmazonExtractor) getInnerText(doc *goquery.Document, selectorText string) string {
	text := doc.Find(selectorText).Text()
	return strings.TrimSpace(strings.ReplaceAll(text, "\t", " "))
}

// getProductTextDetail 获取产品详细信息 (私有方法)
func (e *AmazonExtractor) getProductTextDetail(doc *goquery.Document) ProductDetail {
	productTitle := e.getInnerText(doc, "#productTitle")

	byDesc := e.getInnerText(doc, ".a-unordered-list.a-vertical.a-spacing-mini")
	if byDesc == "" || len(strings.ReplaceAll(byDesc, " ", "")) == 0 {
		byDesc = e.getInnerText(doc, ".a-expander-content.a-expander-partial-collapse-content")
	}

	featureBullets := e.getInnerText(doc, "#feature-bullets")
	bucketdividerDesc := e.getInnerText(doc, ".aplus-v2.desktop.celwidget")
	spacingTopBase := e.getInnerText(doc, ".a-row.a-spacing-top-base")
	productDescription := e.getInnerText(doc, "#productDescription")
	productDiscount := e.getInnerText(doc, ".a-size-large.a-color-price.savingPriceOverride.aok-align-center.reinventPriceSavingsPercentageMargin.savingsPercentage")

	productPrice := strings.Split(e.getInnerText(doc, ".aok-offscreen"), " ")[0]
	symbolFlag := false
	currencySymbols := e.util.GetAllValues("currency_symbol")
	for symbol := range currencySymbols {
		if strings.Contains(productPrice, symbol) {
			symbolFlag = true
			break
		}
	}

	if !symbolFlag {
		// 兼容AMAZON FASHION网站
		productPrice = strings.Split(e.getInnerText(doc, ".a-offscreen"), " ")[0]
	}

	priceFloat := e.util.ExtractPrice(productPrice)
	language := e.util.FindValueForCurrency(productPrice, "language_code")

	var price *string
	if priceFloat != nil {
		currencySymbol := e.util.FindValue(language, "language_code", "currency_symbol")
		priceStr := currencySymbol + fmt.Sprintf("%.2f", *priceFloat)
		price = &priceStr
	}

	var discount *string
	if discountFloat := e.util.ExtractPrice(productDiscount); discountFloat != nil {
		discountStr := fmt.Sprintf("%.2f", *discountFloat)
		discount = &discountStr
	}

	return ProductDetail{
		Title:             productTitle,
		ByDesc:            byDesc,
		Feature:           featureBullets,
		BucketdividerDesc: bucketdividerDesc,
		ProductParameters: spacingTopBase,
		ProductDesc:       productDescription,
		ProductDiscount:   discount,
		ProductPrice:      price,
		Language:          language,
	}
}

// getImgSrc 获取产品图片URL列表 (私有方法)
func (e *AmazonExtractor) getImgSrc(content string) []string {
	// 图片JSON解析
	patternImgObj := regexp.MustCompile(`P\.when\('A'\)\.register\("ImageBlockATF", function\(A\)\{\s*var data = \{\s*'enableS2WithoutS1': [True|true|False|false]+\,\s*'notShowVideoCount': [True|true|False|false]+\,\s*'colorImages': \{ 'initial': (.*?)\}\,\s*'colorToAsin'`)
	matches := patternImgObj.FindStringSubmatch(content)

	if len(matches) > 1 {
		totalImageJSONStr := strings.ReplaceAll(matches[1], "\\", "")
		var totalImageJSON []map[string]interface{}

		if err := json.Unmarshal([]byte(totalImageJSONStr), &totalImageJSON); err == nil {
			var images []string
			for _, row := range totalImageJSON {
				if hiResImg, ok := row["hiRes"].(string); ok && hiResImg != "" {
					images = append(images, hiResImg)
				} else if largeImg, ok := row["large"].(string); ok && largeImg != "" {
					images = append(images, largeImg)
				}
			}
			return images
		}
	}

	// 备用图片提取方法
	imgPattern := regexp.MustCompile(`register.*?var data = (.*?)colorToAsin`)
	if cleanContentMatch := imgPattern.FindStringSubmatch(content); len(cleanContentMatch) > 1 {
		content = cleanContentMatch[1]
	}

	imgURLRegex := regexp.MustCompile(`(https://m\.media-amazon\.com/images/I/.*?)":\s?\[(\d+),\s?\d+\]`)
	imgMatches := imgURLRegex.FindAllStringSubmatch(content, -1)

	var srcs []string
	for _, match := range imgMatches {
		if len(match) >= 3 {
			url := match[1]
			name := strings.Split(filepath.Base(url), "_")[0]
			fix := strings.Split(filepath.Base(url), "_")
			if len(fix) > 0 {
				srcs = append(srcs, "https://m.media-amazon.com/images/I/"+name+fix[len(fix)-1])
			}
		}
	}

	// 去重
	uniqueSrcs := make(map[string]bool)
	var result []string
	for _, src := range srcs {
		if !uniqueSrcs[src] {
			uniqueSrcs[src] = true
			result = append(result, src)
		}
	}

	return result
}

// getVideos 获取产品视频URL列表 (私有方法)
func (e *AmazonExtractor) getVideos(text string, doc *goquery.Document, aplus string) []string {
	var videos []string

	// Aplus视频提取
	if aplus != "" {
		aplusPattern := regexp.MustCompile(`https://m\.media-amazon\.com/images/[A-Za-z0-9\-/_]+\.mp4`)
		matches := aplusPattern.FindAllString(aplus, -1)
		videos = append(videos, matches...)
	}

	// 视频JSON解析
	patternObj := regexp.MustCompile(`P\.when\('A'\)\.execute\('triggerVideoAjax', function\(A\)\{\nvar obj = A\.\$\.parseJSON\('(.*?)'\);\nA\.trigger\('enableS2WithoutS1Ajax`)
	jsonMatches := patternObj.FindAllStringSubmatch(text, -1)

	for _, match := range jsonMatches {
		if len(match) > 1 {
			jsonStr := match[1]
			if strings.HasPrefix(jsonStr, `{"dataInJson"`) && strings.HasSuffix(jsonStr, "}") {
				jsonStr = strings.ReplaceAll(jsonStr, "\\", "")
				var jsonInfo map[string]interface{}

				if err := json.Unmarshal([]byte(jsonStr), &jsonInfo); err == nil {
					if videosData, ok := jsonInfo["videos"].([]interface{}); ok {
						for _, video := range videosData {
							if videoMap, ok := video.(map[string]interface{}); ok {
								if url, ok := videoMap["url"].(string); ok && url != "" {
									videos = append(videos, url)
								}
							}
						}
					}
				}
			}
		}
	}

	// 轮播卡片视频提取
	doc.Find("li._vse-vw-dp-card_style_carouselElement__AVBU9").Each(func(i int, s *goquery.Selection) {
		if div := s.Find("div").First(); div.Length() > 0 {
			if videoURL, exists := div.Attr("data-video-url"); exists && videoURL != "" {
				videos = append(videos, videoURL)
			}
		}
	})

	// 评论视频提取
	doc.Find("div[id^='review-video-id-']").Each(func(i int, s *goquery.Selection) {
		if videoURL, exists := s.Attr("data-video-url"); exists && videoURL != "" {
			videos = append(videos, videoURL)
		}
	})

	// 去重
	uniqueVideos := make(map[string]bool)
	var result []string
	for _, video := range videos {
		if !uniqueVideos[video] {
			uniqueVideos[video] = true
			result = append(result, video)
		}
	}

	// 生成m3u8和mp4列表
	_, m3u8URLList, mp4List := e.generateIDM3u8Mp4List(result)
	return append(m3u8URLList, mp4List...)
}

// generateIDM3u8Mp4List 生成视频ID、m3u8和mp4 URL列表 (私有方法)
func (e *AmazonExtractor) generateIDM3u8Mp4List(videos []string) ([]string, []string, []string) {
	var idList, m3u8URLList, mp4URLList []string

	for _, video := range videos {
		if strings.HasSuffix(video, ".m3u8") {
			match := regexp.MustCompile(`-prod/(.*?)/`).FindStringSubmatch(video)
			if len(match) > 1 {
				idList = append(idList, match[1])
				m3u8URLList = append(m3u8URLList, video)
			}
		}
		if strings.HasSuffix(video, ".mp4") {
			mp4URLList = append(mp4URLList, video)
		}
	}

	return idList, m3u8URLList, mp4URLList
}

// GetProductDetail 获取产品详情 (公开方法)
// 这是唯一暴露给外部的方法，用于提取产品详情
func (e *AmazonExtractor) GetProductDetail(url, text string) ProductResult {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(text))
	if err != nil {
		return ProductResult{}
	}

	productDetail := e.getProductTextDetail(doc)
	srcs := e.getImgSrc(text)

	aplus := doc.Find("#aplus")
	var aplusHTML string
	if aplus.Length() > 0 {
		aplusHTML, _ = aplus.Html()
	}

	videos := e.getVideos(text, doc, aplusHTML)

	return ProductResult{
		LinkURL:  url,
		Title:    productDetail.Title,
		Desc:     productDetail.ByDesc,
		Language: productDetail.Language,
		Images:   srcs,
		Videos:   videos,
		Price:    productDetail.ProductPrice,
		Discount: productDetail.ProductDiscount,
	}
}
