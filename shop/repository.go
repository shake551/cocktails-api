package shop

import (
	"context"
	"github.com/shake551/cocktails-api/db"
	"github.com/shake551/cocktails-api/domain/model"
	"log"
	"time"
)

type Repository interface {
	CreateTable(ctx context.Context, shopID int64) (*model.Table, error)
	GetTable(ctx context.Context, shopID int64, tableID int64) (*model.Table, error)
	Order(ctx context.Context, shopID int64, tableID int64, params OrderParams) ([]*model.Order, error)
	OrderProvide(ctx context.Context, shopID int64, tableID int64, orderID int64) error
	GetTableOrderList(ctx context.Context, shopID int64, tableID int64, unprovided bool) ([]*model.TableOrder, error)
	GetShopUnprovidedOrderList(ctx context.Context, shopID int64, limit int64, offset int64) ([]*model.TableOrder, error)
}

type OrderParams struct {
	CocktailIDs []int64 `json:"cocktail_ids"`
}

type ShopRepository struct{}

func NewShopRepository() Repository {
	return &ShopRepository{}
}

func (r ShopRepository) CreateTable(ctx context.Context, shopID int64) (*model.Table, error) {
	log.Println("create shop table ...")

	query := `INSERT INTO shop_tables (shop_id) VALUES (?)`
	res, err := db.DB.ExecContext(ctx, query, shopID)
	if err != nil {
		return nil, err
	}

	tableID, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}

	return &model.Table{ID: tableID, ShopID: shopID}, nil
}

func (r ShopRepository) GetTable(ctx context.Context, shopID int64, tableID int64) (*model.Table, error) {
	log.Printf("get table ... shopID: %d, tabelID: %d \n", shopID, tableID)

	q := `SELECT * FROM shop_tables WHERE id=? AND shop_id=?`
	rows, err := db.DB.QueryContext(ctx, q, tableID, shopID)
	if db.IsNoRows(err) {
		return &model.Table{}, nil
	}
	if err != nil {
		return &model.Table{}, err
	}

	defer rows.Close()

	var t model.Table
	for rows.Next() {
		if err := rows.Scan(&t.ID, &t.ShopID); err != nil {
			return &model.Table{}, err
		}
	}
	return &t, nil
}

func (r ShopRepository) Order(ctx context.Context, shopID int64, tableID int64, params OrderParams) ([]*model.Order, error) {
	log.Printf("receive order... shop_id: %d, table_id: %d \n", shopID, tableID)

	tx, err := db.DB.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}

	var orders []*model.Order

	now := time.Now().Unix()

	findCocktailQuery := `SELECT * FROM shop_cocktails WHERE shop_id=? AND cocktail_id=?`
	orderQuery := `INSERT INTO shop_orders (table_id, shop_cocktail_id, created_at, updated_at) VALUES (?, ?, ?, ?)`
	for _, cID := range params.CocktailIDs {
		_, err := db.DB.QueryContext(ctx, findCocktailQuery, shopID, cID)
		if db.IsNoRows(err) {
			tx.Rollback()
			log.Printf("does not exist shop_cocktails. shop_id: %d, cocktail_id: %d \n", shopID, cID)
			return nil, err
		}
		if err != nil {
			tx.Rollback()
			log.Printf("cannot find shop_cocktails. shop_id: %d, cocktail_id: %d\n", shopID, cID)
			return nil, err
		}

		res, err := db.DB.ExecContext(ctx, orderQuery, tableID, cID, now, now)
		if err != nil {
			tx.Rollback()
			log.Printf("fail create order. shop_id: %d, table_id: %d, cocktail_id: %d", shopID, tableID, cID)
			return nil, err
		}

		orderID, err := res.LastInsertId()
		if err != nil {
			return nil, err
		}

		orders = append(orders, &model.Order{ID: orderID, TableID: tableID, ShopCocktailID: cID, CreatedAt: now, UpdatedAt: now})
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return orders, nil
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

func (r ShopRepository) GetTableOrderList(ctx context.Context, shopID int64, tableID int64, unprovided bool) ([]*model.TableOrder, error) {
	log.Printf("get table order list ... shopID: %d, tableID: %d \n", shopID, tableID)

	q := `SELECT 
			cocktails.name,
			cocktails.image_url
		FROM shop_tables
			INNER JOIN shop_orders
			INNER JOIN cocktails
		WHERE shop_id=?
			AND shop_tables.id=? 
			AND shop_tables.id = shop_orders.table_id
			AND cocktails.id = shop_orders.shop_cocktail_id
			AND shop_orders.is_provided LIKE CONCAT('%', ?, '%')`

	var isProvided string
	if unprovided {
		isProvided = "1"
	}

	rows, err := db.DB.QueryContext(ctx, q, shopID, tableID, isProvided)
	if db.IsNoRows(err) {
		return []*model.TableOrder{}, err
	}
	if err != nil {
		return []*model.TableOrder{}, err
	}

	defer rows.Close()
	var orders []*model.TableOrder
	for rows.Next() {
		no := model.NullableTableOrder{}
		if err := rows.Scan(&no.Name, &no.ImageURL); err != nil {
			return nil, err
		}

		to := &model.TableOrder{
			Name:     no.Name,
			ImageURL: no.ImageURL.String,
		}
		orders = append(orders, to)
	}
	if len(orders) == 0 {
		return []*model.TableOrder{}, nil
	}

	return orders, nil
}

func (r ShopRepository) GetShopUnprovidedOrderList(ctx context.Context, shopID int64, limit int64, offset int64) ([]*model.TableOrder, error) {
	log.Printf("get shop unprovided prder list ... shopID: %d \n", shopID)

	q := `SELECT 
			cocktails.name,
			cocktails.image_url
		FROM shop_tables
			INNER JOIN shop_orders
			INNER JOIN cocktails
		WHERE shop_id=?
			AND cocktails.id = shop_orders.shop_cocktail_id
		LIMIT ? OFFSET ?`

	rows, err := db.DB.QueryContext(ctx, q, shopID, limit, offset)
	if err != nil {
		return []*model.TableOrder{}, err
	}

	defer rows.Close()
	var orders []*model.TableOrder
	for rows.Next() {
		no := model.NullableTableOrder{}
		if err := rows.Scan(&no.Name, &no.ImageURL); err != nil {
			return nil, err
		}

		to := &model.TableOrder{
			Name:     no.Name,
			ImageURL: no.ImageURL.String,
		}
		orders = append(orders, to)
	}
	if len(orders) == 0 {
		return []*model.TableOrder{}, nil
	}

	return orders, nil
}
