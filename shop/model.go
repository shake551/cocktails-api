package shop

type Shop struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type Table struct {
	ID     int64 `json:"id"`
	ShopID int64 `json:"shop_id"`
}
