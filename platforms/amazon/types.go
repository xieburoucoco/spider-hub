package amazon

// ProductDetail 产品详细信息结构
type ProductDetail struct {
	Title             string  `json:"title"`
	ByDesc            string  `json:"by_desc"`
	Feature           string  `json:"feature"`
	BucketdividerDesc string  `json:"bucketdivider_desc"`
	ProductParameters string  `json:"product_parameters"`
	ProductDesc       string  `json:"product_desc"`
	ProductDiscount   *string `json:"product_discount"`
	ProductPrice      *string `json:"product_price"`
	Language          string  `json:"language"`
}

// ProductResult 产品结果结构
type ProductResult struct {
	LinkURL  string   `json:"link_url"`
	Title    string   `json:"title"`
	Desc     string   `json:"desc"`
	Language string   `json:"language"`
	Images   []string `json:"images"`
	Videos   []string `json:"videos"`
	Price    *string  `json:"price"`
	Discount *string  `json:"discount"`
}

// ImageSearchProduct 图片搜索商品结构
type ImageSearchProduct struct {
	GLProductGroup               string             `json:"glProductGroup"`
	ByLine                       string             `json:"byLine"`
	Price                        string             `json:"price"`
	ListPrice                    *string            `json:"listPrice"`
	CurrencyPriceRange           *string            `json:"currencyPriceRange"`
	VariationalSomePrimeEligible *string            `json:"variationalSomePrimeEligible"`
	ImageURL                     string             `json:"imageUrl"`
	ASIN                         string             `json:"asin"`
	Availability                 string             `json:"availability"`
	Title                        string             `json:"title"`
	IsAdultProduct               string             `json:"isAdultProduct"`
	IsEligibleForPrimeShipping   *string            `json:"isEligibleForPrimeShipping"`
	AverageOverallRating         float64            `json:"averageOverallRating"`
	TotalReviewCount             string             `json:"totalReviewCount"`
	ColorSwatches                []interface{}      `json:"colorSwatches"`
	TwisterVariations            []TwisterVariation `json:"twisterVariations"`
	LinkURL                      string             `json:"link_url"`
}

// TwisterVariation 变体信息
type TwisterVariation struct {
	ASIN     string `json:"asin"`
	ImageURL string `json:"imageUrl"`
}

// StyleSnapResponse 图片搜索API响应
type StyleSnapResponse struct {
	SearchResults []SearchResult `json:"searchResults"`
}

// SearchResult 搜索结果
type SearchResult struct {
	BBXASINMetadataList []ImageSearchProduct `json:"bbxAsinMetadataList"`
}
