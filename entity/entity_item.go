package entity

type Item struct {
	ItemID int     `json:"item_id"`
	Name   string  `json:"name"`
	Stock  int     `json:"stock"`
	Price  float64 `json:"price"`
}
