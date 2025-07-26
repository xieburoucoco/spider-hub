# 🛒 Amazon爬虫

> 亚马逊商品详情爬虫，支持提取商品标题、描述、价格、图片和视频等信息

## 📋 功能特性

- ✅ **商品详情提取**：标题、描述、价格、折扣等
- 🖼️ **图片提取**：高清商品图片
- 🎬 **视频提取**：商品相关视频
- 🌐 **多语言支持**：自动检测商品语言
- 🔄 **批量处理**：支持批量爬取多个商品

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

## ⚠️ 注意事项

- 请遵守亚马逊的robots.txt规则
- 过于频繁的请求可能会导致IP被封禁
- 建议使用代理IP轮换请求

## 🧪 测试

```bash
# 运行测试
cd platforms/amazon
go test -v
``` 