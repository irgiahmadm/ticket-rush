package domain

type Order struct {
	ID        int    `json:"id"`
	UserID    string `json:"user_id"`
	ProductID int    `json:"product_id"`
	Status    string `json:"status"`
}
