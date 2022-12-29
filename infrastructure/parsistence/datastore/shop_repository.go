package datastore

import (
	"context"
	"github.com/shake551/cocktails-api/db"
	"github.com/shake551/cocktails-api/domain/model"
	"log"
)

type ShopRepository struct{}

func NewShopRepository() *ShopRepository {
	return &ShopRepository{}
}

func (r ShopRepository) GetLimit(ctx context.Context, limit int64, offset int64) ([]model.Shop, error) {
	log.Println("get shops with limit ...")

	query := `SELECT * FROM shops LIMIT ? OFFSET ?`
	rows, err := db.DB.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var shops []model.Shop
	for rows.Next() {
		s := model.Shop{}
		if err := rows.Scan(&s.ID, &s.Name); err != nil {
			return nil, err
		}

		shops = append(shops, s)
	}

	if len(shops) == 0 {
		return []model.Shop{}, nil
	}

	return shops, nil
}

func (r ShopRepository) Create(ctx context.Context, params model.ShopParams) (*model.Shop, error) {
	log.Println("create shop...")

	query := `INSERT INTO shops (name) VALUES (?)`
	res, err := db.DB.ExecContext(ctx, query, params.Name)
	if err != nil {
		return nil, err
	}

	shopID, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}

	return &model.Shop{ID: shopID, Name: params.Name}, nil
}
