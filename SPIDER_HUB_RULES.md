# 🕷️ Spider-Hub 项目开发规则与规范

## 📋 项目概述

Spider-Hub 是一个专注于各平台爬虫逆向技术的Go语言项目集合，通过统一的架构设计和开发规范，为不同平台提供高效、稳定的数据采集解决方案。

## 🏗️ 项目结构规范

### 总体架构
```
spider-hub/
├── README.md                    # 项目总览文档
├── LICENSE                      # 开源协议
├── go.mod                       # Go模块文件
├── go.sum                       # 依赖校验文件
├── .gitignore                   # Git忽略文件
├── SPIDER_HUB_RULES.md         # 项目规则规范（本文档）
├── platforms/                   # 各平台爬虫实现目录
│   ├── coinglass/              # Coinglass平台实现
│   ├── amazon/                 # Amazon平台实现
│   └── [platform_name]/       # 其他平台实现
├── cmd/                        # 命令行工具
│   └── spider-hub/            # 主程序入口
└── internal/                   # 内部工具包
    └── util/                  # 通用工具函数
```

## 📁 Platforms目录规范

### 规则1: 平台目录结构

每个平台必须遵循以下目录结构：

#### 🔄 统一接口类型平台（如coinglass）
当平台的所有接口返回数据格式相同，可以封装为统一的函数集时：

```
platforms/[platform_name]/
├── spider.go              # 主要爬虫实现，包含所有接口
├── types.go               # 数据结构定义
├── constants.go           # 常量配置
├── utils.go               # 工具函数（可选）
├── example_test.go        # 测试用例
└── README.md              # 平台说明文档（唯一README）
```

**示例：Coinglass平台**
```go
// spider.go - 统一的接口实现
func (c *CoinglassClient) GetLongShortRatio(params map[string]interface{}) (*Response, error)
func (c *CoinglassClient) GetOpenInterest(params map[string]interface{}) (*Response, error)
func (c *CoinglassClient) GetLiquidation(params map[string]interface{}) (*Response, error)
// 所有接口返回统一的Response结构
```

#### 🔀 多样化接口类型平台（如amazon）
当平台的不同接口返回数据格式差异较大，需要按功能模块分离时：

```
platforms/[platform_name]/
├── [platform_name].go     # 主要爬虫实现（商品详情等核心功能）
├── image_search.go         # 图片搜索功能（可选，如需要）
├── review_crawler.go       # 评论爬取功能（可选，如需要）
├── types.go               # 所有数据结构定义
├── constants.go           # 常量配置
├── utils.go               # 工具函数
├── extractor.go           # 数据提取器（可选）
├── example_test.go        # 测试用例
└── README.md              # 平台说明文档（唯一README）
```

**示例：Amazon平台**
```go
// amazon.go - 商品详情接口
func (s *AmazonSpider) FetchProductDetail(ctx context.Context, productURL string) (ProductResult, error)

// amazon.go - 图片搜索接口（集成在同一文件中）
func (s *AmazonSpider) SearchProductsByImageURL(ctx context.Context, imageURL string, proxies map[string]string) ([]ImageSearchProduct, error)
func (s *AmazonSpider) SearchProductsByImageData(ctx context.Context, imageData []byte, proxies map[string]string) ([]ImageSearchProduct, error)
```

### 规则2: 文件命名规范

#### 核心文件（必需）
- `[platform_name].go` - 主要爬虫实现，以平台名命名
- `types.go` - 数据结构定义
- `constants.go` - 常量和配置
- `README.md` - **唯一的平台说明文档**
- `example_test.go` - 测试用例

#### 扩展文件（可选）
- `utils.go` - 工具函数
- `extractor.go` - 数据提取器
- `[feature_name].go` - 特定功能模块（如需要）

#### 🚫 禁止的文件
- ❌ 多个README文件（如`FEATURE_README.md`等）
- ❌ 单独的测试目录（测试应集中在`example_test.go`）
- ❌ 配置文件分散（应统一在`constants.go`）

### 规则3: README文档规范

每个平台**有且仅有一个**`README.md`文件，必须包含以下章节：

```markdown
# 🛒 [Platform Name]爬虫

> 平台简介和功能描述

## 📋 功能特性
- ✅ 功能1描述
- ✅ 功能2描述
- ✅ 功能N描述

## 🚀 快速开始
### 安装依赖
### 基本使用
### 高级用法

## 📊 返回数据结构
### 数据结构1
### 数据结构2

## 🧪 测试和示例
```bash
# 运行测试命令
```

## ⚠️ 注意事项
## 📝 更新日志
```

## 🔧 代码规范

### 规则4: 结构体命名

#### 客户端/爬虫结构体
```go
// 统一接口类型
type [Platform]Client struct {
    // 字段定义
}

// 多样化接口类型
type [Platform]Spider struct {
    // 字段定义
}
```

#### 构造函数
```go
// 统一接口类型
func New[Platform]Client() *[Platform]Client

// 多样化接口类型  
func New[Platform]Spider() *[Platform]Spider
```

### 规则5: 接口设计

#### 统一接口类型（如coinglass）
```go
// 所有接口返回统一格式
type Response struct {
    Success bool        `json:"success"`
    Data    interface{} `json:"data"`
    Message string      `json:"message"`
}

// 接口方法
func (c *CoinglassClient) Get[Feature](params map[string]interface{}) (*Response, error)
```

#### 多样化接口类型（如amazon）
```go
// 不同功能返回不同结构
func (s *AmazonSpider) FetchProductDetail(ctx context.Context, productURL string) (ProductResult, error)
func (s *AmazonSpider) SearchProductsByImageURL(ctx context.Context, imageURL string, proxies map[string]string) ([]ImageSearchProduct, error)
```

### 规则6: 错误处理

```go
// 统一错误格式
return nil, fmt.Errorf("❌ 操作失败: %w", err)

// 日志输出格式
fmt.Printf("✅ 操作成功: %s\n", result)
fmt.Printf("⚠️ 警告信息: %s\n", warning)
fmt.Printf("❌ 错误信息: %s\n", error)
```

### 规则7: 测试规范

#### 测试文件结构
```go
// example_test.go
package [platform]

import (
    "context"
    "fmt"
    "testing"
    "time"
)

// 功能测试1
func Test[Feature1](t *testing.T) {
    // 测试实现
}

// 功能测试2  
func Test[Feature2](t *testing.T) {
    // 测试实现
}

// 辅助函数
func helperFunction() {
    // 辅助实现
}
```

## 📦 平台集成规范

### 规则8: 新平台集成流程

1. **创建平台目录**
   ```bash
   mkdir platforms/[platform_name]
   cd platforms/[platform_name]
   ```

2. **确定平台类型**
   - 🔄 **统一接口类型**：所有API返回格式相同 → 使用coinglass模式
   - 🔀 **多样化接口类型**：不同API返回格式差异大 → 使用amazon模式

3. **创建核心文件**
   ```bash
   touch [platform_name].go types.go constants.go README.md example_test.go
   ```

4. **实现核心功能**
   - 定义数据结构（`types.go`）
   - 配置常量（`constants.go`）
   - 实现爬虫逻辑（`[platform_name].go`）
   - 编写测试用例（`example_test.go`）
   - 完善文档（`README.md`）

5. **质量检查**
   ```bash
   go build -o /dev/null .     # 编译检查
   go test -v                  # 测试检查
   go fmt ./...               # 格式检查
   ```

### 规则9: 代码质量标准

#### 必须满足的条件
- ✅ 代码能够正常编译
- ✅ 所有测试用例通过
- ✅ 遵循Go语言规范（gofmt）
- ✅ 包含详细的注释文档
- ✅ 错误处理完善
- ✅ 只有一个README文件

#### 推荐的最佳实践
- 🎯 接口设计简洁明了
- 🔄 支持上下文控制（context.Context）
- 🛡️ 完善的错误处理和重试机制
- 📝 清晰的日志输出
- 🧪 充分的测试coverage

## 🚨 平台开发注意事项

### 规则10: 合规性要求

1. **法律合规**
   - 遵守目标网站的robots.txt协议
   - 遵守相关法律法规
   - 不进行恶意爬取

2. **技术合规**
   - 合理控制请求频率
   - 使用适当的User-Agent
   - 支持代理配置
   - 实现请求重试机制

3. **文档合规**
   - 每个平台只能有一个README文件
   - 必须包含免责声明
   - 说明使用限制和注意事项

## 📊 项目维护规范

### 规则11: 版本管理

- 使用语义化版本控制
- 及时更新go.mod依赖
- 维护更新日志

### 规则12: 文档维护

- 保持README文档的时效性
- 及时更新API变化
- 记录重要的配置变更

## 🎯 总结

Spider-Hub项目通过以上规范，确保：

1. **结构清晰**：统一的目录结构和命名规范
2. **文档统一**：每个平台只有一个综合性README
3. **代码规范**：遵循Go语言最佳实践
4. **易于维护**：清晰的分层和模块化设计
5. **扩展性强**：支持不同类型平台的灵活集成

遵循这些规则，可以确保项目的一致性、可维护性和可扩展性，为各平台爬虫开发提供标准化的框架支持。 