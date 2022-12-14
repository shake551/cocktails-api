package datastore

import (
	"context"
	"database/sql"
	"github.com/shake551/cocktails-api/db"
	"github.com/shake551/cocktails-api/domain/model"
	"log"
	"strings"
	"time"
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

func (r CocktailRepository) GetByID(ctx context.Context, id int64) (model.CocktailDetail, error) {
	log.Printf("get cocktails with cocktail id...")

	query := `
		SELECT
		    cocktails.id,
			cocktails.name,
			cocktails.image_url,
			materials.id,
			materials.name,
			cocktail_materials.quantity,
			cocktail_materials.unit
		FROM cocktails
		INNER JOIN cocktail_materials
			ON cocktails.id = cocktail_materials.cocktail_id
			INNER JOIN materials
				ON cocktail_materials.material_id = materials.id
		WHERE cocktails.id = ?
	`

	rows, err := db.DB.QueryContext(ctx, query, id)
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

		if err := rows.Scan(&ncd.ID, &ncd.Name, &ncd.ImageURL, &ncd.MaterialID, &ncd.MaterialName, &ncd.Quantity, &ncd.Unit); err != nil {
			return model.CocktailDetail{}, err
		}

		materials = append(materials, model.Material{
			ID:   ncd.MaterialID,
			Name: ncd.MaterialName,
			Quantity: model.MaterialQuantity{
				Quantity: ncd.Quantity,
				Unit:     ncd.Unit,
			},
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

func (r CocktailRepository) Create(ctx context.Context, params model.CocktailParams) (*model.CocktailDetail, error) {
	log.Printf("create cocktails...")

	tx, err := db.DB.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}

	now := time.Now().Unix()
	cocktailsQuery := `INSERT INTO cocktails (name,created_at,updated_at) VALUES (?,?,?)`
	res, err := db.DB.ExecContext(ctx, cocktailsQuery, params.Name, now, now)
	if err != nil {
		tx.Rollback()
		log.Printf("failed to create message. err: %v", err)
		return nil, err
	}
	cocktailID, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}

	materials := []model.Material{}

	materialSelectQuery := `SELECT EXISTS (SELECT * FROM materials WHERE materials.name = ?)`
	materialInsertQuery := `INSERT INTO materials (name, created_at, updated_at) VALUES (?, ?, ?)`
	cocktailMaterialQuery := `INSERT INTO cocktail_materials (cocktail_id, material_id, quantity, unit) VALUES (?, ?, ?, ?)`
	for _, m := range params.Materials {
		rows, err := db.DB.QueryContext(ctx, materialSelectQuery, m.Name)
		if err != nil {
			log.Println(err)
		}

		defer rows.Close()

		var recordCount int64
		var materialID int64
		for rows.Next() {
			err := rows.Scan(&recordCount)
			if recordCount == 0 {
				res, err = db.DB.ExecContext(ctx, materialInsertQuery, m.Name, now, now)
				if err != nil {
					tx.Rollback()
					log.Printf("failed to create message. err: %v", err)
					return nil, err
				}

				materialID, err = res.LastInsertId()
				if err != nil {
					return nil, err
				}
			}
		}

		_, err = db.DB.ExecContext(ctx, cocktailMaterialQuery, cocktailID, materialID, m.Quantity.Quantity, m.Quantity.Unit)
		if err != nil {
			tx.Rollback()
			log.Printf("failed to create message. err: %v", err)
			return nil, err
		}

		material := model.Material{
			ID:       materialID,
			Name:     m.Name,
			Quantity: m.Quantity,
		}

		materials = append(materials, material)
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return &model.CocktailDetail{
		ID:        cocktailID,
		Name:      params.Name,
		Materials: materials,
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}

func (r CocktailRepository) GetListByIDs(ctx context.Context, ids []int64) ([]model.Cocktail, error) {
	log.Println("get cocktails with id list ...")

	var rows *sql.Rows
	var err error

	repeat := strings.Repeat("?,", len(ids)-1) + "?"
	var cocktailIds []interface{}
	for _, id := range ids {
		cocktailIds = append(cocktailIds, id)
	}

	query := `SELECT * FROM cocktails where id IN ( ` + repeat + ` ) `
	rows, err = db.DB.QueryContext(ctx, query, cocktailIds...)
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
