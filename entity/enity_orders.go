package entity

type OrderReport struct {
	OrderID     int     `json:"order_id"`
	UserID      int     `json:"user_id"`
	FullName    string  `json:"full_name"`
	OrderDate   string  `json:"order_date"`
	Status      string  `json:"status"`
	TotalAmount float64 `json:"total_amount"`
	StaffID     int     `json:"staff_id"`
	StaffName   string  `json:"staff_name"`
}

type OrderItem struct {
	ItemID   int
	Quantity int
	Price    float64
}

type Order struct {
	UserID      int
	TotalAmount float64
	OrderItems  []OrderItem
}

type Payment struct {
	UserID        int
	Amount        float64
	PaymentMethod string
}
