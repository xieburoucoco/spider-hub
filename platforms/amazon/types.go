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
