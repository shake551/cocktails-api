package shop

type Shop struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type Table struct {
	ID     int64 `json:"id"`
	ShopID int64 `json:"shop_id"`
}

type ShopCocktail struct {
	ShopID     int64 `json:"shop_id"`
	CocktailID int64 `json:"cocktail_id"`
}

type Order struct {
	ID             int64 `json:"id"`
	TableID        int64 `json:"table_id"`
	ShopCocktailID int64 `json:"shop_cocktail_id"`
}
