package handler

type ProductResponse struct {
	SellerImages string `json:"seller_pictures,omitempty"`
	UserName     string `json:"sellers,omitempty"`
	ProductID    uint   `json:"product_id,omitempty"`
	Images       string `json:"images"`
	ProductName  string `json:"product_name"`
	Price        int    `json:"price"`
	Stock        int    `json:"stock"`
	City         string `json:"city,omitempty"`
	Description  string `json:"description,omitempty"`
}
