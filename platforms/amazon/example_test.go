package amazon

import (
	"context"
	"testing"
	"time"
)

// TestAmazonProductDetail 测试获取亚马逊商品详情
func TestAmazonProductDetail(t *testing.T) {
	// 创建亚马逊爬虫实例
	spider := NewAmazonSpider()

	// 设置测试URL (Sony WH-1000XM5 耳机)
	testURL := "https://www.amazon.com/Sony-WH-1000XM5-Canceling-Headphones-Hands-Free/dp/B09XS7JWHH"

	// 设置超时上下文
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// 获取商品详情
	result, err := spider.FetchProductDetail(ctx, testURL)
	if err != nil {
		t.Fatalf("❌ 获取商品详情失败: %v", err)
	}

	// 验证结果
	if result.Title == "" {
		t.Error("❌ 商品标题为空")
	} else {
		t.Logf("✅ 商品标题: %s", result.Title)
	}

	if len(result.Images) == 0 {
		t.Error("❌ 没有找到商品图片")
	} else {
		t.Logf("✅ 找到 %d 张商品图片", len(result.Images))
		t.Logf("  第一张图片: %s", result.Images[0])
	}

	if result.Price == nil {
		t.Log("⚠️ 没有找到商品价格")
	} else {
		t.Logf("✅ 商品价格: %s", *result.Price)
	}

	t.Logf("✅ 商品描述: %s", result.Desc)
	t.Logf("✅ 商品语言: %s", result.Language)

	// 输出完整结果
	priceStr := "无"
	if result.Price != nil {
		priceStr = *result.Price
	}
	discountStr := "无"
	if result.Discount != nil {
		discountStr = *result.Discount + "%"
	}

	t.Logf("\n🛍️ 商品详情结果:\n"+
		"  标题: %s\n"+
		"  描述: %s\n"+
		"  价格: %s\n"+
		"  折扣: %s\n"+
		"  图片数量: %d\n"+
		"  视频数量: %d\n",
		result.Title,
		truncateString(result.Desc, 100),
		priceStr,
		discountStr,
		len(result.Images),
		len(result.Videos))
}

// 辅助函数: 截断字符串
func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}
