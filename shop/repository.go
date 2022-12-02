package shop

import (
	"context"
	"github.com/shake551/cocktails-api/db"
	"log"
)

type Repository interface {
	GetLimit(ctx context.Context, limit int64, offset int64) ([]Shop, error)
	FindByID(ctx context.Context, id int64) (Shop, error)
	Create(ctx context.Context, params ShopParams) (*Shop, error)
	CreateTable(ctx context.Context, shopID int64) (*Table, error)
	AddShopCocktails(ctx context.Context, shopID int64, params ShopCocktailParams) ([]*ShopCocktail, error)
}

type ShopParams struct {
	Name string `json:"name"`
}

type ShopCocktailParams struct {
	CocktailIDs []int64 `json:"cocktail_ids"`
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

func (r ShopRepository) FindByID(ctx context.Context, id int64) (Shop, error) {
	log.Println("find shop with shop id ...")

	query := `SELECT * FROM shops WHERE id = ?`
	rows, err := db.DB.QueryContext(ctx, query, id)
	if db.IsNoRows(err) {
		return Shop{}, err
	}
	if err != nil {
		return Shop{}, err
	}

	defer rows.Close()

	s := Shop{}
	for rows.Next() {
		if err := rows.Scan(&s.ID, &s.Name); err != nil {
			return Shop{}, err
		}
	}

	return s, nil
}

func (r ShopRepository) Create(ctx context.Context, params ShopParams) (*Shop, error) {
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

	return &Shop{ID: shopID, Name: params.Name}, nil
}

func (r ShopRepository) CreateTable(ctx context.Context, shopID int64) (*Table, error) {
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

	return &Table{ID: tableID, ShopID: shopID}, nil
}

func (r ShopRepository) AddShopCocktails(ctx context.Context, shopID int64, params ShopCocktailParams) ([]*ShopCocktail, error) {
	log.Printf("add shop cocktails... shop_id: %d\n", shopID)

	tx, err := db.DB.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}

	var cocktails []*ShopCocktail

	findCocktailQuery := `SELECT id FROM cocktails WHERE id=?`
	createShopCocktailQuery := `INSERT INTO shop_cocktails (shop_id, cocktail_id) VALUES (?, ?)`
	for _, cID := range params.CocktailIDs {
		_, err := db.DB.QueryContext(ctx, findCocktailQuery, cID)
		if db.IsNoRows(err) {
			tx.Rollback()
			log.Printf("does not exist cocktails. cokctail_id: %d \n", cID)
			return nil, err
		}
		if err != nil {
			tx.Rollback()
			log.Printf("cannot find shop_cocktails. shop_id: %d, cokctail_id: %d\n", shopID, cID)
			return nil, err
		}

		_, err = db.DB.ExecContext(ctx, createShopCocktailQuery, shopID, cID)
		if err != nil {
			tx.Rollback()
			log.Printf("fail create shop_cocktail. shop_id: %d, cocktail_id: %d", shopID, cID)
			return nil, err
		}

		cocktails = append(cocktails, &ShopCocktail{ShopID: shopID, CocktailID: cID})
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return cocktails, nil
}