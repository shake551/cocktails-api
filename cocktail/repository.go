package cocktail

import (
	"context"
	"github.com/shake551/cocktails-api/db"
	"github.com/shake551/cocktails-api/domain/model"
	"log"
	"time"
)

type Repository interface {
	Create(ctx context.Context, params CocktailsParams) (*model.CocktailsDetail, error)
}

type CocktailsParams struct {
	Name      string
	Materials []MaterialParams
}

type MaterialParams struct {
	Name     string                 `json:"name"`
	Quantity model.MaterialQuantity `json:"quantity"`
}

type CocktailsRepository struct{}

func NewCocktailsRepository() Repository {
	return &CocktailsRepository{}
}

func (r CocktailsRepository) Create(ctx context.Context, params CocktailsParams) (*model.CocktailsDetail, error) {
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

	return &model.CocktailsDetail{
		ID:        cocktailID,
		Name:      params.Name,
		Materials: materials,
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}
