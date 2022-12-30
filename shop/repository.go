package shop

import (
	"context"
	"github.com/shake551/cocktails-api/db"
	"log"
)

type Repository interface {
	OrderProvide(ctx context.Context, shopID int64, tableID int64, orderID int64) error
}

type ShopRepository struct{}

func NewShopRepository() Repository {
	return &ShopRepository{}
}

func (r ShopRepository) OrderProvide(ctx context.Context, shopID int64, tableID int64, orderID int64) error {
	log.Printf("order provide ... shopID: %d, tableID: %d, orderID: %d \n", shopID, tableID, orderID)

	selectQuery := `SELECT 
			shop_orders.*
		FROM shop_tables
			INNER JOIN shop_orders
		WHERE shop_id=?
			AND shop_tables.id=? 
			AND shop_tables.id = shop_orders.table_id
			AND shop_orders.id = ?`

	_, err := db.DB.QueryContext(ctx, selectQuery, shopID, tableID, orderID)
	if err != nil {
		return err
	}

	q := `UPDATE shop_orders SET is_provided=true WHERE id=?`
	_, err = db.DB.QueryContext(ctx, q, orderID)
	return err
}
