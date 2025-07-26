# Coinglass 数据爬虫 🕷️

Coinglass 是一个专业的加密货币数据分析平台，提供期货、现货、持仓量等多维度的市场数据。

## 🏗️ 项目结构

## 功能特性 ✨

- 🔍 持仓量图表数据获取
- 📊 期货首页统计数据
- 🏢 衍生品交易所列表
- 💰 币种市场数据
- 🔗 交易所期货交易对信息
- 🪙 现货支持币种信息

## API 接口说明 📋

### 1. 持仓量图表数据 (`/api/openInterest/v3/chart`) 📈
获取指定币种的持仓量图表数据

**请求示例:**
```
https://capi.coinglass.com/api/openInterest/v3/chart?symbol=BTC&timeType=0&exchangeName=&currency=USD&type=0
```
**实际测试结果:**
```
✅ 持仓量图表数据获取成功，数据长度: 422718
📊 响应时间: 3.81s
📄 数据预览: {"dataMap":{"Binance":[161745430,181309262,180784983,182010652,176903553,176218564,183792561,1963354...
```

### 2. 期货首页统计 (`/api/futures/home/statistics`) 🏠
获取期货市场的整体统计数据，无需参数

**请求示例:**
```
https://capi.coinglass.com/api/futures/home/statistics
```

**实际测试结果:**
```
✅ 期货首页统计数据获取成功，数据长度: 223
📊 响应时间: 0.35s
📄 数据预览: {"volH24Chain":-26.5900,"shortRate":51.8900,"oiH24Chain":1.5900,"openInterest":202194926751,"longRat...
```


### 3. 衍生品交易所列表 (`/api/derivative/exchange/list`) 🏢
获取支持衍生品交易的交易所列表，无需参数

**请求示例:**
```
https://capi.coinglass.com/api/derivative/exchange/list
```

**实际测试结果:**
```
✅ 衍生品交易所列表获取成功，数据长度: 7824
📊 响应时间: 0.24s
📄 数据预览: [{"exchangeName":"Binance","liquidationVolUsd":95900704.33530666,"logo":"https://cdn.coinglasscdn.co...
```


### 4. 币种市场数据 (`/api/home/v2/coinMarkets`) 💰
获取币种的市场数据和排行信息

**请求示例:**
```
https://capi.coinglass.com/api/home/v2/coinMarkets?sort=&order=&keyword=&pageNum=1&pageSize=20&ex=all
```

**参数说明:**
- `sort`: 排序字段 🔄
- `order`: 排序方向 ⬆️
- `keyword`: 搜索关键词 🔍
- `pageNum`: 页码 (默认: 1) 📄
- `pageSize`: 每页数量 (默认: 20) 📊
- `ex`: 交易所筛选 (all: 全部) 🏢

**实际测试结果:**
```
✅ 币种市场数据获取成功，数据长度: 53169
📊 响应时间: 0.81s
📄 数据预览: {"total":664,"pageSize":20,"list":[{"avgFundingRateByOi":0.011645,"avgFundingRateByOiAPR":12.7513,"a...
```

### 5. 交易所期货交易对信息 (`/api/exchange/futures/pairInfo`) 🔗
获取各交易所的期货交易对详细信息，无需参数

**请求示例:**
```
https://capi.coinglass.com/api/exchange/futures/pairInfo
```

**实际测试结果:**
```
✅ 交易所期货交易对信息获取成功，数据长度: 73495
📊 响应时间: 1.01s
📄 数据预览: [{"symbol":"BTC","symbolLogo":"https://cdn.coinglasscdn.com/static/img/coins/bitcoin-BTC.png"},{"sym...
```

### 6. 现货支持币种 (`/api/spot/support/coin`) 🪙
获取平台支持的现货交易币种列表，无需参数

**请求示例:**
```
https://capi.coinglass.com/api/spot/support/coin
```

**实际测试结果:**
```
✅ 现货支持币种获取成功，数据长度: 5794
📊 响应时间: 0.25s
📄 数据预览: ["BTC","ETH","XRP","USDT","BNB","SOL","USDC","DOGE","TRX","ADA","HYPE","SUI","XLM","LINK","HBAR","BC...
```

### 7. 币种市场数据（小数据量测试） (`/api/home/v2/coinMarkets`) 📊
获取少量币种数据用于快速测试

**请求示例:**
```
https://capi.coinglass.com/api/home/v2/coinMarkets?sort=&order=&keyword=&pageNum=1&pageSize=5&ex=all
```

**实际测试结果:**
```
✅ 市场数据获取成功，数据长度: 13592
📊 响应时间: 0.22s
📄 数据预览: {"total":664,"pageSize":5,"list":[{"avgFundingRateByOi":0.011645,"avgFundingRateByOiAPR":12.7513,"av...
```

## ⚠️ 免责声明

### 法律声明 📋
本项目仅供学习和研究使用，请遵守相关网站的使用条款和robots.txt协议。使用者需要自行承担使用风险，作者不承担任何法律责任。

1. **🔍 学习目的**：本爬虫工具仅用于技术学习、数据分析研究等合法用途
2. **📚 教育用途**：适用于编程教学、算法研究、数据科学学习等教育场景
3. **🚫 禁止滥用**：严禁用于恶意攻击、数据盗取、商业牟利等非法行为
4. **⚖️ 法律责任**：使用者需自行承担因违法使用而产生的一切法律后果

### 使用限制 🚦

1. **📊 数据使用**：
   - 获取的数据仅供个人学习和研究使用
   - 不得将数据用于商业目的或二次销售
   - 不得恶意大量请求导致服务器负载过重

2. **🔐 隐私保护**：
   - 尊重目标网站的robots.txt协议
   - 不得获取用户隐私信息或敏感数据
   - 遵守相关数据保护法规

3. **⏱️ 访问频率**：
   - 合理控制请求频率，避免对目标服务器造成压力
   - 建议在请求间设置适当的延时
   - 遵守网站的访问限制和使用条款

### 风险提示 ⚡

1. **🔄 数据变化**：目标网站可能随时更改API接口或数据结构
2. **🛡️ 反爬措施**：网站可能采用反爬虫技术，导致爬虫失效
3. **📈 市场风险**：获取的金融数据仅供参考，不构成投资建议
4. **⚖️ 法律风险**：请确保使用行为符合当地法律法规

### 版权说明 📄

- 本项目代码采用 [Apache 2.0](../../LICENSE) 开源协议，详见 LICENSE 文件
- Coinglass 是 Coinglass 公司的注册商标，本项目与其无任何官方关联
- 所有数据版权归 Coinglass 平台所有
- 项目中使用的Logo、商标等知识产权归原权利人所有
- 更多详细信息请访问 [Coinglass 官网](https://www.coinglass.com/) 🌐
---

**⚡ Spider-Hub**

> 如果您觉得这个项目有用，请给我们一个 ⭐ Star，这是对我们最大的鼓励！

