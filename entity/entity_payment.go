package entity

type PaymentReport struct {
	UserID        int     `json:"user_id"`
	FullName      string  `json:"full_name"`
	Email         string  `json:"email"`
	PhoneNumber   string  `json:"phone_number"`
	Amount        float64 `json:"amount"`
	PaymentDate   string  `json:"payment_date"`
	PaymentMethod string  `json:"payment_method"`
	Status        string  `json:"status"`
}
