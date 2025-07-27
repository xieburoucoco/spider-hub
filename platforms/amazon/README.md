# 🛒 Amazon爬虫

> 亚马逊商品详情爬虫，支持提取商品标题、描述、价格、图片和视频等信息

## 📋 功能特性

- ✅ **商品详情提取**：标题、描述、价格、折扣等
- 🖼️ **图片提取**：高清商品图片
- 🎬 **视频提取**：商品相关视频
- 🌐 **多语言支持**：自动检测商品语言
- 🔄 **批量处理**：支持批量爬取多个商品
- 🔍 **图片搜索**：通过图片URL或本地图片搜索相似商品
- 🌐 **代理支持**：支持HTTP/HTTPS代理配置
- 🔄 **自动重试**：内置重试机制，提高搜索成功率

## 🚀 快速开始

### 安装依赖

```bash
go get github.com/PuerkitoBio/goquery
go get github.com/go-resty/resty/v2
```

### 基本使用

```go
package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/xieburoucoco/spider-hub/platforms/amazon"
)

func main() {
	// 创建爬虫实例
	spider := amazon.NewAmazonSpider()

	// 设置商品URL
	productURL := "https://www.amazon.com/Sony-WH-1000XM5-Canceling-Headphones-Hands-Free/dp/B09XS7JWHH"

	// 设置超时上下文
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// 获取商品详情
	result, err := spider.FetchProductDetail(ctx, productURL)
	if err != nil {
		log.Fatalf("获取商品详情失败: %v", err)
	}

	// 输出结果
	fmt.Printf("商品标题: %s\n", result.Title)
	fmt.Printf("商品描述: %s\n", result.Desc)
	fmt.Printf("商品价格: %v\n", result.Price)
	fmt.Printf("商品图片数量: %d\n", len(result.Images))
	fmt.Printf("商品视频数量: %d\n", len(result.Videos))
}
```

### 图片搜索功能

#### 通过在线图片URL搜索商品

```go
package main

import (
	"context"
	"fmt"
	"log"
	"time"
	
	"github.com/xieburoucoco/spider-hub/platforms/amazon"
)

func main() {
	// 创建爬虫实例
	spider := amazon.NewAmazonSpider()
	
	// 设置超时上下文
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	
	// 图片URL
	imageURL := "https://m.media-amazon.com/images/I/71c-jiE2IcL._AC_SX679_.jpg"
	
	// 搜索商品
	products, err := spider.SearchProductsByImageURL(ctx, imageURL, nil)
	if err != nil {
		log.Fatalf("搜索失败: %v", err)
	}
	
	// 输出结果
	fmt.Printf("找到 %d 个相关商品:\n", len(products))
	for i, product := range products {
		fmt.Printf("\n商品 %d:\n", i+1)
		fmt.Printf("  标题: %s\n", product.Title)
		fmt.Printf("  品牌: %s\n", product.ByLine)
		fmt.Printf("  价格: %s\n", product.Price)
		fmt.Printf("  评分: %.1f\n", product.AverageOverallRating)
		fmt.Printf("  评论数: %s\n", product.TotalReviewCount)
		fmt.Printf("  链接: %s\n", product.LinkURL)
	}
}
```

#### 通过本地图片文件搜索商品

```go
package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"time"
	
	"github.com/xieburoucoco/spider-hub/platforms/amazon"
)

func main() {
	// 创建爬虫实例
	spider := amazon.NewAmazonSpider()
	
	// 设置超时上下文
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	
	// 读取本地图片文件
	imageData, err := readImageFile("camera.jpg")
	if err != nil {
		log.Fatalf("读取图片失败: %v", err)
	}
	
	// 搜索商品
	products, err := spider.SearchProductsByImageData(ctx, imageData, nil)
	if err != nil {
		log.Fatalf("搜索失败: %v", err)
	}
	
	// 输出结果
	fmt.Printf("找到 %d 个相关商品:\n", len(products))
	for i, product := range products {
		fmt.Printf("\n商品 %d:\n", i+1)
		fmt.Printf("  标题: %s\n", product.Title)
		fmt.Printf("  品牌: %s\n", product.ByLine)
		fmt.Printf("  价格: %s\n", product.Price)
		fmt.Printf("  评分: %.1f\n", product.AverageOverallRating)
		fmt.Printf("  评论数: %s\n", product.TotalReviewCount)
		fmt.Printf("  链接: %s\n", product.LinkURL)
	}
}

// 读取本地图片文件
func readImageFile(filePath string) ([]byte, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	
	return io.ReadAll(file)
}
```

#### 使用代理配置

```go
// 配置代理
proxies := map[string]string{
	"http":  "http://proxy.example.com:8080",
	"https": "https://proxy.example.com:8080",
}

// 使用代理搜索
products, err := spider.SearchProductsByImageURL(ctx, imageURL, proxies)
```

### 批量处理

```go
// 设置多个商品URL
productURLs := []string{
	"https://www.amazon.com/Sony-WH-1000XM5-Canceling-Headphones-Hands-Free/dp/B09XS7JWHH",
	"https://www.amazon.com/Apple-MacBook-13-inch-256GB-Storage/dp/B08N5LLDSG",
}

// 批量获取商品详情
results, errors := spider.FetchProductDetails(ctx, productURLs)

// 处理结果
for i, result := range results {
	fmt.Printf("商品 %d: %s\n", i+1, result.Title)
}
```

## 📊 返回数据结构

### 商品详情结构

```go
type ProductResult struct {
	LinkURL  string   `json:"link_url"`   // 商品链接
	Title    string   `json:"title"`      // 商品标题
	Desc     string   `json:"desc"`       // 商品描述
	Language string   `json:"language"`   // 商品语言
	Images   []string `json:"images"`     // 商品图片URL列表
	Videos   []string `json:"videos"`     // 商品视频URL列表
	Price    *string  `json:"price"`      // 商品价格
	Discount *string  `json:"discount"`   // 商品折扣
}
```

### 图片搜索结果结构

```go
type ImageSearchProduct struct {
	GLProductGroup               string                    `json:"glProductGroup"`              // 商品分组
	ByLine                      string                    `json:"byLine"`                      // 品牌/制造商
	Price                       string                    `json:"price"`                       // 当前价格
	ListPrice                   *string                   `json:"listPrice"`                   // 原价
	CurrencyPriceRange          *string                   `json:"currencyPriceRange"`          // 价格区间
	VariationalSomePrimeEligible *string                   `json:"variationalSomePrimeEligible"` // Prime资格
	ImageURL                    string                    `json:"imageUrl"`                    // 商品图片URL
	ASIN                        string                    `json:"asin"`                        // 亚马逊商品ID
	Availability                string                    `json:"availability"`                // 库存状态
	Title                       string                    `json:"title"`                       // 商品标题
	IsAdultProduct              string                    `json:"isAdultProduct"`              // 是否成人商品
	IsEligibleForPrimeShipping  *string                   `json:"isEligibleForPrimeShipping"`  // Prime配送资格
	AverageOverallRating        float64                   `json:"averageOverallRating"`        // 平均评分
	TotalReviewCount            string                    `json:"totalReviewCount"`            // 评论总数
	ColorSwatches               []interface{}             `json:"colorSwatches"`               // 颜色选项
	TwisterVariations           []TwisterVariation        `json:"twisterVariations"`           // 变体信息
	LinkURL                     string                    `json:"link_url"`                    // 商品链接
}

type TwisterVariation struct {
	ASIN     string `json:"asin"`
	ImageURL string `json:"imageUrl"`
}
```

## ⚠️ 注意事项

- 请遵守亚马逊的robots.txt规则
- 过于频繁的请求可能会导致IP被封禁
- 建议使用代理IP轮换请求
- 图片搜索功能基于亚马逊StyleSnap技术
- 建议在请求之间添加适当延时，避免被反爬虫机制阻止
- 使用高质量、清晰的图片能获得更好的搜索结果

## 🧪 测试和示例

### 运行测试

```bash
# 运行图片搜索测试
cd platforms/amazon

# 测试通过URL搜索商品
go test -v -run TestSearchProductsByImageURL

# 测试通过本地图片搜索商品
go test -v -run TestSearchProductsByImageData

# 运行所有测试
go test -v
```

### 测试输出示例

以下是实际运行测试时的输出日志：

```
=== RUN   TestSearchProductsByImageURL
🖼️ 亚马逊图片搜索功能演示
================================================================================

📡 示例1: 通过在线图片URL搜索商品
--------------------------------------------------
🔍 搜索图片: https://m.media-amazon.com/images/I/71c-jiE2IcL._AC_SX679_.jpg

请求头加密参数：【hFYtEMYspiQsknjSuiY6h7r9EwnPbLhIcukjGWUrPDbJAAAAAGiGP8sAAAAB】
✅ 找到 16 个相关商品:

📦 商品 1:
   🏷️  标题: AiTechny Digital Camera for Kids, 1080P FHD 44MP Point an...
   🏭 品牌: AiTechny
   💰 价格: $36.79
   ⭐ 评分: 4.2 (471条评论)
   📦 库存: 有库存 ✅
   🔗 链接: https://www.amazon.com/dp/B0DFW4RR1H
   🎨 变体: 3个选项

📦 商品 2:
   🏷️  标题: Digital Camera, Upgraded FHD 1080P Point and Shoot Kids C...
   🏭 品牌: Lecran
   💰 价格: $36.79 (原价: $45.99)
   ⭐ 评分: 4.6 (627条评论)
   📦 库存: 有库存 ✅
   🔗 链接: https://www.amazon.com/dp/B0DHKGWYVG
   🎨 变体: 4个选项

📦 商品 3:
   🏷️  标题: AiTechny Digital Camera, 1080P FHD Camera for Kids, 44MP ...
   🏭 品牌: AiTechny
   💰 价格: $36.79
   ⭐ 评分: 3.8 (113条评论)
   📦 库存: 库存紧张 ⚠️
   🔗 链接: https://www.amazon.com/dp/B0CM6MVW2B

📦 商品 4:
   🏷️  标题: Digital Camera,Autofocus 4K Vlogging Camera for Photograp...
   🏭 品牌: Lecnippy
   💰 价格: $49.99 (原价: $64.99)
   ⭐ 评分: 4.5 (1,145条评论)
   📦 库存: 有库存 ✅
   🔗 链接: https://www.amazon.com/dp/B0DNFGNG18
   🎨 变体: 3个选项

📦 商品 5:
   🏷️  标题: FHD 1080P Digital Camera for Kids with 32GB SD Card - Com...
   🏭 品牌: VAHOIALD
   💰 价格: $45.99
   ⭐ 评分: 4.0 (1,074条评论)
   📦 库存: 有库存 ✅
   🔗 链接: https://www.amazon.com/dp/B0D12TN394
   🎨 变体: 3个选项

... 还有 11 个商品未显示

================================================================================
🎉 URL图片搜索测试完成！
--- PASS: TestSearchProductsByImageURL (7.05s)
PASS
```

## 🔧 技术实现要点

### 核心功能
- ✅ **获取stylesnap值**: 从亚马逊页面提取认证参数
- ✅ **图片下载**: 支持从URL下载图片到内存
- ✅ **图片上传**: 将图片上传到亚马逊StyleSnap API
- ✅ **结果解析**: 解析JSON响应并构建Go结构体
- ✅ **自动重试**: 内置3次重试机制，提高成功率
- ✅ **代理支持**: 支持HTTP/HTTPS代理配置

### 简化的代码结构
- 🎯 **简洁的接口**: 只暴露两个核心方法，隐藏内部实现细节
- 🔄 **统一的重试逻辑**: 两个接口共享相同的重试机制
- 📦 **完整的数据结构**: 定义了详细的返回数据结构
- 🛡️ **错误处理**: 完善的错误处理和日志输出

## 🔍 与原Python代码的对比

### 保持的功能
- ✅ 获取stylesnap认证参数
- ✅ 图片下载和上传
- ✅ 失败重试机制
- ✅ 代理支持
- ✅ 完整的错误处理

### 改进的地方
- 🎯 **接口简化**: 从多个函数简化为2个核心接口
- 📦 **类型安全**: 使用Go的强类型系统，避免运行时错误
- 🔄 **统一重试**: 重试逻辑更加统一和可配置
- 📝 **更好的文档**: 详细的注释和使用说明
- 🏗️ **更清晰的结构**: 遵循Go的最佳实践

## 📝 更新日志

### v1.2.0 - 图片搜索功能
- ✅ 新增通过图片URL搜索商品功能
- ✅ 新增通过本地图片数据搜索商品功能
- ✅ 支持HTTP/HTTPS代理配置
- ✅ 内置3次重试机制
- ✅ 完善的错误处理和日志输出
- ✅ 详细的测试用例和使用示例

### v1.1.0 - 商品详情爬取
- ✅ 基础商品详情提取功能
- ✅ 支持图片和视频提取
- ✅ 多语言支持
- ✅ 批量处理功能 