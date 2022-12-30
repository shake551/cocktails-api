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

func (r ShopRepository) GetByID(ctx context.Context, id int64) (model.Shop, error) {
	log.Println("find shop with shop id ...")

	query := `SELECT * FROM shops WHERE id = ?`
	rows, err := db.DB.QueryContext(ctx, query, id)
	if db.IsNoRows(err) {
		return model.Shop{}, err
	}
	if err != nil {
		return model.Shop{}, err
	}

	defer rows.Close()

	s := model.Shop{}
	for rows.Next() {
		if err := rows.Scan(&s.ID, &s.Name); err != nil {
			return model.Shop{}, err
		}
	}

	return s, nil
}

func (r ShopRepository) GetShopCocktailList(ctx context.Context, shopID int64, limit int64, offset int64) ([]model.Cocktail, error) {
	log.Printf("get shop cocktail list ... %d \n", shopID)

	q := `SELECT
    		cocktails.*
		FROM 
		    cocktails
		    INNER JOIN shop_cocktails
		WHERE shop_cocktails.shop_id = ? 
		    AND shop_cocktails.cocktail_id = cocktails.id
		LIMIT ? OFFSET ?`

	rows, err := db.DB.QueryContext(ctx, q, shopID, limit, offset)
	if err != nil {
		log.Println(err)
		return []model.Cocktail{}, err
	}

	defer rows.Close()
	var cocktails []model.Cocktail
	for rows.Next() {
		nc := model.NullableCocktail{}
		if err := rows.Scan(&nc.ID, &nc.Name, &nc.ImageURL, &nc.CreatedAt, &nc.UpdatedAt); err != nil {
			log.Println(err)
			return []model.Cocktail{}, err
		}

		c := model.Cocktail{
			ID:        nc.ID,
			Name:      nc.Name,
			ImageURL:  nc.ImageURL.String,
			CreatedAt: nc.CreatedAt,
			UpdatedAt: nc.UpdatedAt,
		}
		cocktails = append(cocktails, c)
	}

	if len(cocktails) == 0 {
		return []model.Cocktail{}, nil
	}
	return cocktails, nil
}

func (r ShopRepository) AddShopCocktail(ctx context.Context, shopID int64, params model.ShopCocktailParams) ([]*model.ShopCocktail, error) {
	log.Printf("add shop cocktails... shop_id: %d\n", shopID)

	tx, err := db.DB.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}

	var cocktails []*model.ShopCocktail

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

		cocktails = append(cocktails, &model.ShopCocktail{ShopID: shopID, CocktailID: cID})
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return cocktails, nil
}

func (r ShopRepository) GetShopCocktailDetail(ctx context.Context, shopID int64, cocktailID int64) (model.CocktailDetail, error) {
	log.Printf("get shop cocktail detail ... shopID: %d, cocktailID: %d \n", shopID, cocktailID)

	q := `
		SELECT
		    cocktails.id,
			cocktails.name,
			cocktails.image_url,
			materials.id,
			materials.name
		FROM cocktails
		INNER JOIN shop_cocktails
			ON shop_cocktails.cocktail_id = cocktails.id
		INNER JOIN cocktail_materials
			ON cocktails.id = cocktail_materials.cocktail_id
			INNER JOIN materials
				ON cocktail_materials.material_id = materials.id
		WHERE shop_cocktails.shop_id= ?
			AND cocktails.id = ?
	`

	rows, err := db.DB.QueryContext(ctx, q, shopID, cocktailID)
	if db.IsNoRows(err) {
		return model.CocktailDetail{}, nil
	}
	if err != nil {
		return model.CocktailDetail{}, err
	}

	defer rows.Close()

	var ncd model.NullableCocktailDetailRow
	var materials []model.Material
	for rows.Next() {
		if err := rows.Scan(&ncd.ID, &ncd.Name, &ncd.ImageURL, &ncd.MaterialID, &ncd.MaterialName); err != nil {
			return model.CocktailDetail{}, err
		}

		materials = append(materials, model.Material{
			ID:   ncd.MaterialID,
			Name: ncd.MaterialName,
		})
	}

	d := model.CocktailDetail{
		ID:        ncd.ID,
		Name:      ncd.Name,
		ImageURL:  ncd.ImageURL.String,
		Materials: materials,
	}

	return d, nil
}

func (r ShopRepository) GetUnprovidedOrderList(ctx context.Context, shopID int64, limit int64, offset int64) ([]*model.TableOrder, error) {
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

func (r ShopRepository) AddTable(ctx context.Context, shopID int64) (*model.Table, error) {
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
