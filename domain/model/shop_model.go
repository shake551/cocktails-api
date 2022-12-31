package model

import "database/sql"

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
	CreatedAt      int64 `json:"created_at"`
	UpdatedAt      int64 `json:"updated_at"`
}

type TableOrder struct {
	Name     string `json:"name"`
	ImageURL string `json:"image_url"`
}

type NullableTableOrder struct {
	Name     string
	ImageURL sql.NullString
}

type ShopParams struct {
	Name string `json:"name"`
}

type ShopCocktailParams struct {
	CocktailIDs []int64 `json:"cocktail_ids"`
}

type OrderParams struct {
	CocktailIDs []int64 `json:"cocktail_ids"`
}
