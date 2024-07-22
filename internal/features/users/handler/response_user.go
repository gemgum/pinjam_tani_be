package handler

type UserResponse struct {
	ID          uint   `json:"user_id"`
	Images      string `json:"images"`
	UserName    string `json:"username"`
	Email       string `json:"email"`
	PhoneNumber string `json:"phone_number"`
	Address     string `json:"address"`
}
