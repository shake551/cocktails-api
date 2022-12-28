package datastore

import (
	"context"
	"database/sql"
	"github.com/shake551/cocktails-api/db"
	"github.com/shake551/cocktails-api/domain/model"
	"log"
)

type CocktailRepository struct{}

func NewCocktailRepository() *CocktailRepository {
	return &CocktailRepository{}
}

func (r CocktailRepository) GetLimit(ctx context.Context, limit int64, offset int64, keyword string) ([]model.Cocktail, error) {
	log.Printf("get cocktails with limit...")

	var rows *sql.Rows
	var err error

	if keyword != "" {
		query := `SELECT * FROM cocktails WHERE name LIKE CONCAT('%', ?, '%') LIMIT ? OFFSET ?`
		rows, err = db.DB.QueryContext(ctx, query, keyword, limit, offset)
	} else {
		query := `SELECT * FROM cocktails LIMIT ? OFFSET ?`
		rows, err = db.DB.QueryContext(ctx, query, limit, offset)
	}
	if err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	defer rows.Close()

	var cocktails []model.Cocktail
	for rows.Next() {
		nc := model.NullableCocktail{}
		if err := rows.Scan(&nc.ID, &nc.Name, &nc.ImageURL, &nc.CreatedAt, &nc.UpdatedAt); err != nil {
			return nil, err
		}

		c := model.Cocktail{
			ID:        nc.ID,
			Name:      nc.Name,
			ImageURL:  nc.ImageURL.String,
			CreatedAt: nc.CreatedAt,
			UpdatedAt: nc.CreatedAt,
		}
		cocktails = append(cocktails, c)
	}

	if len(cocktails) == 0 {
		return []model.Cocktail{}, nil
	}

	return cocktails, nil
}
