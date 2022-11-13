package shop

import (
	"context"
	"github.com/shake551/cocktails-api/db"
	"log"
)

type Repository interface {
	GetLimit(ctx context.Context, limit int64, offset int64) ([]Shop, error)
}

type ShopRepository struct{}

func NewShopRepository() Repository {
	return &ShopRepository{}
}

func (r ShopRepository) GetLimit(ctx context.Context, limit int64, offset int64) ([]Shop, error) {
	log.Println("get shops with limit ...")

	query := `SELECT * FROM shops LIMIT ? OFFSET ?`
	rows, err := db.DB.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}

	var shops []Shop
	for rows.Next() {
		s := Shop{}
		if err := rows.Scan(&s.ID, &s.Name); err != nil {
			return nil, err
		}

		shops = append(shops, s)
	}

	if len(shops) == 0 {
		return []Shop{}, nil
	}

	return shops, nil
}