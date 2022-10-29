package cocktail

import (
	"context"
	"fmt"
	"github.com/shake551/cocktails-api/db"
	"log"
	"time"
)

type Repository interface {
	GetLimit(ctx context.Context, limit int64, offset int64, keyword string) ([]Cocktail, error)
	Create(ctx context.Context, params CocktailsParams) (*CocktailsDetail, error)
}

type CocktailsParams struct {
	Name      string
	Materials []MaterialParams
}

type MaterialParams struct {
	Name     string           `json:"name"`
	Quantity MaterialQuantity `json:"quantity"`
}

type CocktailsRepository struct{}

func NewCocktailsRepository() Repository {
	return &CocktailsRepository{}
}

func (r CocktailsRepository) GetLimit(ctx context.Context, limit int64, offset int64, keyword string) ([]Cocktail, error) {
	log.Printf("get cocktails with limit...")

	query := `SELECT * FROM cocktails LIMIT ? OFFSET ?`
	rows, err := db.DB.QueryContext(ctx, query, limit, offset)

	if keyword != "" {
		query = `SELECT * FROM cocktails WHERE name LIKE CONCAT('%', ?, '%') LIMIT ? OFFSET ?`
		rows, err = db.DB.QueryContext(ctx, query, keyword, limit, offset)
	}
	if err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	defer rows.Close()

	var cocktails []Cocktail
	for rows.Next() {
		nc := NullableCocktail{}
		if err := rows.Scan(&nc.ID, &nc.Name, &nc.ImageURL, &nc.CreatedAt, &nc.UpdatedAt); err != nil {
			return nil, err
		}

		c := Cocktail{
			ID:        nc.ID,
			Name:      nc.Name,
			ImageURL:  nc.ImageURL.String,
			CreatedAt: nc.CreatedAt,
			UpdatedAt: nc.CreatedAt,
		}
		cocktails = append(cocktails, c)
	}

	if len(cocktails) == 0 {
		return []Cocktail{}, nil
	}

	return cocktails, nil
}

func (r CocktailsRepository) Create(ctx context.Context, params CocktailsParams) (*CocktailsDetail, error) {
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

	materials := []Material{}

	materialSelectQuery := `SELECT id FROM materials WHERE name=?`
	materialInsertQuery := `INSERT INTO materials (name, created_at, updated_at) VALUES (?, ?, ?)`
	cocktailMaterialQuery := `INSERT INTO cocktail_materials (cocktail_id, material_id, quantity, unit) VALUES (?, ?, ?, ?)`
	for _, m := range params.Materials {
		rows, err := db.DB.QueryContext(ctx, materialSelectQuery, m.Name)

		if db.IsNoRows(err) {
			rows, err = db.DB.QueryContext(ctx, materialInsertQuery, m.Name, now, now)
			if err != nil {
				fmt.Println(err)
			}
		}

		var materialID int64

		defer rows.Close()

		for rows.Next() {
			err := rows.Scan(&materialID)
			if err != nil {
				log.Fatal(err)
			}
		}

		if err != nil {
			fmt.Println(err)
		}

		_, err = db.DB.ExecContext(ctx, cocktailMaterialQuery, cocktailID, materialID, m.Quantity.Quantity, m.Quantity.Unit)
		if err != nil {
			tx.Rollback()
			log.Printf("failed to create message. err: %v", err)
			return nil, err
		}

		material := Material{
			ID:       materialID,
			Name:     m.Name,
			Quantity: m.Quantity,
		}

		materials = append(materials, material)
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return &CocktailsDetail{
		ID:        cocktailID,
		Name:      params.Name,
		Materials: materials,
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}
