package repository

import (
	"context"
	"github.com/shake551/cocktails-api/domain/model"
)

type CocktailRepository interface {
	GetLimit(ctx context.Context, limit int64, offset int64, keyword string) ([]model.Cocktail, error)
	GetByID(ctx context.Context, id int64) (model.CocktailsDetail, error)
	Create(ctx context.Context, params model.CocktailsParams) (*model.CocktailsDetail, error)
}