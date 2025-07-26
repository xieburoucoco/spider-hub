package coinglass

import (
	"bytes"
	"compress/gzip"
	"crypto/aes"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"time"

	"github.com/go-resty/resty/v2"
)

// CoinglassResponse API响应结构
// 🔐 所有Coinglass API都返回加密的响应数据
// 📦 需要通过特定的解密流程才能获取原始JSON数据
type CoinglassResponse struct {
	Code    string `json:"code"`    // 📋 响应状态码
	Msg     string `json:"msg"`     // 💬 响应消息
	Data    string `json:"data"`    // 🔒 加密的数据内容
	Success bool   `json:"success"` // ✅ 请求是否成功
}

// Spider Coinglass爬虫
// 🕷️ 专门用于抓取Coinglass平台加密货币数据的爬虫实例
// 🔧 内置HTTP客户端和加密解密功能
type Spider struct {
	client *resty.Client // 🌐 HTTP请求客户端
}

// NewSpider 创建新的Coinglass爬虫实例
// 🏗️ 初始化爬虫配置，设置超时时间和请求参数
// 📝 返回可用于数据抓取的Spider实例
//
// 使用示例:
//
//	spider := NewSpider()
//	data, err := spider.GetData("https://capi.coinglass.com/api/...")
func NewSpider() *Spider {
	client := resty.New()
	client.SetTimeout(30 * time.Second) // ⏱️ 设置30秒超时

	return &Spider{
		client: client,
	}
}

// GetData 获取任意API数据（通用方法）
// 🌟 这是Spider的核心方法，用于获取任意Coinglass API的数据
// 🔄 自动处理加密请求、数据解密、gzip解压等复杂流程
// 📊 支持所有Coinglass公开API接口
//
// 参数:
//
//	apiURL - 🌐 完整的API请求URL，包含查询参数
//
// 返回:
//
//	string - 📄 解密后的JSON字符串数据
//	error  - ❌ 如果请求或解密过程中出现错误
//
// 支持的API接口:
//
//	📈 /api/openInterest/v3/chart - 持仓量图表数据
//	🏠 /api/futures/home/statistics - 期货首页统计
//	🏢 /api/derivative/exchange/list - 衍生品交易所列表
//	💰 /api/home/v2/coinMarkets - 币种市场数据
//	🔗 /api/exchange/futures/pairInfo - 期货交易对信息
//	🪙 /api/spot/support/coin - 现货支持币种
func (s *Spider) GetData(apiURL string) (string, error) {
	// ⏰ 生成时间戳作为加密密钥的一部分
	cacheTsV2 := fmt.Sprintf("%d", time.Now().UnixMilli())

	// 🔐 第一步：获取加密的响应数据和动态密钥
	response, userHeader, err := s.getEncryptedData(apiURL, cacheTsV2)
	if err != nil {
		return "", fmt.Errorf("获取加密数据失败: %v", err)
	}

	// 🔓 第二步：解密数据并解压gzip
	decryptedData, err := s.decryptData(response, cacheTsV2, userHeader)
	if err != nil {
		return "", fmt.Errorf("解密数据失败: %v", err)
	}

	return decryptedData, nil
}

// getEncryptedData 获取加密响应数据
// 🌐 向Coinglass API发送请求，获取加密的响应数据
// 🔑 同时获取用于解密的动态密钥（通过response header传递）
// 🎭 模拟真实浏览器的请求头，避免被反爬虫检测
//
// 参数:
//
//	apiURL    - 🌐 API请求地址
//	cacheTsV2 - ⏰ 时间戳，用于生成解密密钥
//
// 返回:
//
//	*CoinglassResponse - 📦 加密的API响应结构
//	string            - 🔑 用户动态密钥（从response header获取）
//	error             - ❌ 请求错误信息
func (s *Spider) getEncryptedData(apiURL, cacheTsV2 string) (*CoinglassResponse, string, error) {
	// 🎭 设置完整的浏览器请求头，模拟真实用户访问
	headers := map[string]string{
		"accept":             "application/json",                                                                                                // 📋 接受JSON响应
		"accept-language":    "zh-CN,zh;q=0.9",                                                                                                  // 🌏 语言偏好
		"cache-ts-v2":        cacheTsV2,                                                                                                         // ⏰ 时间戳密钥
		"encryption":         "true",                                                                                                            // 🔐 启用加密
		"language":           "zh",                                                                                                              // 🈳 界面语言
		"origin":             "https://www.coinglass.com",                                                                                       // 🏠 请求来源
		"priority":           "u=1, i",                                                                                                          // 🚦 请求优先级
		"referer":            "https://www.coinglass.com/",                                                                                      // 🔗 引用页面
		"sec-ch-ua":          `"Not)A;Brand";v="8", "Chromium";v="138", "Google Chrome";v="138"`,                                                // 🌐 浏览器标识
		"sec-ch-ua-mobile":   "?0",                                                                                                              // 📱 非移动设备
		"sec-ch-ua-platform": `"Windows"`,                                                                                                       // 💻 操作系统
		"sec-fetch-dest":     "empty",                                                                                                           // 🎯 请求目标
		"sec-fetch-mode":     "cors",                                                                                                            // 🔄 CORS模式
		"sec-fetch-site":     "same-site",                                                                                                       // 🏠 同站请求
		"user-agent":         "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/138.0.0.0 Safari/537.36", // 🕷️ 用户代理
	}

	// 🚀 发送HTTP GET请求
	resp, err := s.client.R().SetHeaders(headers).Get(apiURL)
	if err != nil {
		return nil, "", fmt.Errorf("请求失败: %v", err)
	}

	// 🔍 检查HTTP状态码
	if resp.StatusCode() != 200 {
		return nil, "", fmt.Errorf("HTTP状态码错误: %d", resp.StatusCode())
	}

	// 🔑 从响应头获取动态解密密钥
	userHeader := resp.Header().Get("user")
	if userHeader == "" {
		return nil, "", fmt.Errorf("未找到user header")
	}

	// 📦 解析JSON响应体
	var response CoinglassResponse
	if err := json.Unmarshal(resp.Body(), &response); err != nil {
		return nil, "", fmt.Errorf("解析响应失败: %v", err)
	}

	return &response, userHeader, nil
}

// decryptData 解密数据
// 🔓 执行完整的数据解密流程，包括多层加密解密和数据解压
// 🔐 使用时间戳密钥解密用户动态密钥
// 🗜️ 使用动态密钥解密响应数据并解压gzip
//
// 解密流程:
//  1. 📅 使用时间戳生成第一层密钥
//  2. 🔑 解密user header获取动态密钥
//  3. 📦 使用动态密钥解密响应数据
//  4. 🗜️ 解压gzip获取最终JSON数据
//
// 参数:
//
//	response   - 📦 加密的API响应
//	cacheTsV2  - ⏰ 时间戳密钥
//	userHeader - 🔑 加密的动态密钥
//
// 返回:
//
//	string - 📄 解密后的JSON字符串
//	error  - ❌ 解密过程中的错误
func (s *Spider) decryptData(response *CoinglassResponse, cacheTsV2, userHeader string) (string, error) {
	// 🔑 第一步：生成时间戳密钥（Base64编码后取前16位）
	timestampEncoded := base64.StdEncoding.EncodeToString([]byte(cacheTsV2))
	timestampKey := []byte(timestampEncoded[:16])

	// 🔓 第二步：解密user header获取动态密钥
	userCiphertext, err := base64.StdEncoding.DecodeString(userHeader)
	if err != nil {
		return "", fmt.Errorf("解码user header失败: %v", err)
	}

	// 🔐 使用时间戳密钥解密用户动态密钥
	dynamicKeyBytes, err := s.decryptAES(userCiphertext, timestampKey)
	if err != nil {
		return "", fmt.Errorf("解密user header失败: %v", err)
	}

	// 🔄 转换动态密钥格式并解析
	dynamicZipKey := hex.EncodeToString(dynamicKeyBytes)
	dynamicKey, err := s.parseDecryptedHex(dynamicZipKey)
	if err != nil {
		return "", fmt.Errorf("解析动态密钥失败: %v", err)
	}

	// 🔓 第三步：解密响应数据
	ciphertext, err := base64.StdEncoding.DecodeString(response.Data)
	if err != nil {
		return "", fmt.Errorf("解码响应数据失败: %v", err)
	}

	// 🔐 使用动态密钥解密响应数据
	dynamicKeyBytes = []byte(dynamicKey)
	decryptedData, err := s.decryptAES(ciphertext, dynamicKeyBytes)
	if err != nil {
		return "", fmt.Errorf("解密响应数据失败: %v", err)
	}

	// 🗜️ 第四步：解压gzip获取最终数据
	hexString := hex.EncodeToString(decryptedData)
	result, err := s.parseDecryptedHex(hexString)
	if err != nil {
		return "", fmt.Errorf("解压最终数据失败: %v", err)
	}

	return result, nil
}

// decryptAES 执行AES-ECB解密
// 🔐 使用AES算法的ECB模式进行数据解密
// 🔑 支持128位密钥长度（16字节）
// 📦 自动处理PKCS7填充的去除
//
// AES-ECB特点:
//
//	🔒 每个数据块独立加密，相同明文产生相同密文
//	⚡ 加密速度快，适合大量数据处理
//	🚫 安全性相对较低，但足够用于API数据传输
//
// 参数:
//
//	ciphertext - 🔒 待解密的密文数据
//	key        - 🔑 16字节的AES密钥
//
// 返回:
//
//	[]byte - 📄 解密后的原始数据
//	error  - ❌ 解密过程中的错误
func (s *Spider) decryptAES(ciphertext, key []byte) ([]byte, error) {
	// 🔍 验证密文长度必须是AES块大小的倍数
	if len(ciphertext)%aes.BlockSize != 0 {
		return nil, fmt.Errorf("密文长度不是块大小的倍数")
	}

	// 🔑 验证密钥长度必须为16字节（AES-128）
	if len(key) != 16 {
		return nil, fmt.Errorf("密钥长度必须为16字节")
	}

	// 🔐 初始化AES解密器
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("初始化AES失败: %v", err)
	}

	// 🔄 逐块解密数据（ECB模式）
	decrypted := make([]byte, len(ciphertext))
	for i := 0; i < len(ciphertext); i += aes.BlockSize {
		block.Decrypt(decrypted[i:i+aes.BlockSize], ciphertext[i:i+aes.BlockSize])
	}

	// 📦 去除PKCS7填充
	return s.pkcs7Unpad(decrypted)
}

// pkcs7Unpad 去除PKCS7填充
// �� PKCS7是一种标准的数据填充方式，用于确保数据长度符合加密算法要求
// 🔍 填充值等于填充字节的数量，例如填充3个字节则每个字节的值都是3
// ✂️ 解密后需要去除这些填充字节以获取原始数据
//
// PKCS7填充规则:
//
//	📏 如果数据长度刚好是块大小的倍数，仍需添加一个完整的填充块
//	🔢 填充字节的值等于填充的字节数
//	✅ 例如：...ABC + [5,5,5,5,5] 表示填充了5个字节
//
// 参数:
//
//	data - 📄 包含PKCS7填充的解密数据
//
// 返回:
//
//	[]byte - 📄 去除填充后的原始数据
//	error  - ❌ 填充格式错误
func (s *Spider) pkcs7Unpad(data []byte) ([]byte, error) {
	length := len(data)
	if length == 0 {
		return nil, fmt.Errorf("无效的解密数据长度")
	}

	// 🔢 获取填充长度（最后一个字节的值）
	paddingLen := int(data[length-1])

	// 🔍 验证填充长度的合理性
	if paddingLen > length || paddingLen > aes.BlockSize {
		return nil, fmt.Errorf("无效的PKCS7填充")
	}

	// ✅ 验证填充字节的正确性
	for i := length - paddingLen; i < length; i++ {
		if data[i] != byte(paddingLen) {
			return nil, fmt.Errorf("无效的PKCS7填充字节")
		}
	}

	// ✂️ 返回去除填充后的数据
	return data[:length-paddingLen], nil
}

// parseDecryptedHex 解析十六进制字符串并解压gzip
// 🔄 将十六进制字符串转换为字节数组
// 🗜️ 使用gzip算法解压数据获取最终的JSON字符串
// 📊 Coinglass使用gzip压缩来减少数据传输量
//
// 处理流程:
//  1. 🔤 十六进制字符串 -> 字节数组
//  2. 🗜️ gzip解压 -> 原始JSON字符串
//  3. 📄 返回可解析的JSON数据
//
// 参数:
//
//	hexString - 🔤 十六进制编码的gzip压缩数据
//
// 返回:
//
//	string - 📄 解压后的JSON字符串
//	error  - ❌ 解析或解压过程中的错误
func (s *Spider) parseDecryptedHex(hexString string) (string, error) {
	// 🔤 将十六进制字符串转换为字节数组
	data, err := hex.DecodeString(hexString)
	if err != nil {
		return "", fmt.Errorf("解码十六进制失败: %v", err)
	}

	// 📖 创建字节读取器
	reader := bytes.NewReader(data)

	// 🗜️ 创建gzip解压器
	gzipReader, err := gzip.NewReader(reader)
	if err != nil {
		return "", fmt.Errorf("创建gzip reader失败: %v", err)
	}
	defer gzipReader.Close() // 🔒 确保资源释放

	// 📄 解压数据到缓冲区
	var decompressed bytes.Buffer
	if _, err := io.Copy(&decompressed, gzipReader); err != nil {
		return "", fmt.Errorf("解压gzip失败: %v", err)
	}

	// 📊 返回解压后的JSON字符串
	return decompressed.String(), nil
}
