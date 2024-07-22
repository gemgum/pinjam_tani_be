package handler

type CartRequest struct {
	ProductID int `json:"products_id"`
	Quantity  int `json:"quantity"`
}
