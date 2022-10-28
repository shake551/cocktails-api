package cocktail

import "database/sql"

type Cocktail struct {
	ID        int64          `json:"id"`
	Name      string         `json:"name"`
	ImageURL  sql.NullString `json:"image_url"`
	CreatedAt int64          `json:"created_at"`
	UpdatedAt int64          `json:"updated_at"`
	Materials []Material     `json:"materials"`
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
