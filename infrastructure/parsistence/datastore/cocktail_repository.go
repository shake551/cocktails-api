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

func (r CocktailRepository) GetByID(ctx context.Context, id int64) (model.CocktailsDetail, error) {
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
		return model.CocktailsDetail{}, nil
	}
	if err != nil {
		return model.CocktailsDetail{}, err
	}

	defer rows.Close()

	var ncd model.NullableCocktailDetailRow
	var materials []model.Material
	for rows.Next() {

		if err := rows.Scan(&ncd.ID, &ncd.Name, &ncd.ImageURL, &ncd.MaterialID, &ncd.MaterialName, &ncd.Quantity, &ncd.Unit); err != nil {
			return model.CocktailsDetail{}, err
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

	d := model.CocktailsDetail{
		ID:        ncd.ID,
		Name:      ncd.Name,
		ImageURL:  ncd.ImageURL.String,
		Materials: materials,
	}

	return d, nil
}
