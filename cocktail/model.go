package cocktail

import "database/sql"

type Cocktail struct {
	ID        int64  `json:"id"`
	Name      string `json:"name"`
	ImageURL  string `json:"image_url"`
	CreatedAt int64  `json:"created_at"`
	UpdatedAt int64  `json:"updated_at"`
}

type NullableCocktail struct {
	ID        int64
	Name      string
	ImageURL  sql.NullString
	CreatedAt int64
	UpdatedAt int64
}

type CocktailsDetail struct {
	ID        int64      `json:"id"`
	Name      string     `json:"name"`
	ImageURL  string     `json:"image_url"`
	Materials []Material `json:"materials"`
	CreatedAt int64      `json:"created_at"`
	UpdatedAt int64      `json:"updated_at"`
}

type NullableCocktailDetailRow struct {
	ID           int64
	Name         string
	ImageURL     sql.NullString
	MaterialID   int64
	MaterialName string
	Quantity     int64
	Unit         string
}

type Material struct {
	ID       int64            `json:"id"`
	Name     string           `json:"name"`
	Quantity MaterialQuantity `json:"quantity"`
}

type MaterialQuantity struct {
	Quantity int64  `json:"quantity"`
	Unit     string `json:"unit"`
}
