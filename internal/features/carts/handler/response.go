package handler

type CartResponse struct {
	Id          uint    `json:"id"`
	UserID      uint    `json:"user_id,omitempty"`
	Images      string  `json:"images"`
	ProductName string  `json:"product_name"`
	Price       int     `json:"price"`
	Quantity    float64 `json:"quantity"`
	TotalPrice  float64 `json:"total_price"`
}

type CreateCartResponse struct {
	Id         uint `json:"id"`
	UserID     uint `json:"user_id,omitempty"`
	ProductID  uint `json:"product_id"`
	Quantity   uint `json:"quantity"`
	TotalPrice uint `json:"total_price"`
}
