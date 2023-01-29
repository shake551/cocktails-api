package repository

import (
	"context"
	"github.com/shake551/cocktails-api/domain/model"
)

//go:generate mockery --dir . --name CocktailRepository --outpkg repository_mock --output ../repository_mock --case underscore
type CocktailRepository interface {
	GetLimit(ctx context.Context, limit int64, offset int64, keyword string) ([]model.Cocktail, error)
	GetByID(ctx context.Context, id int64) (model.CocktailDetail, error)
	Create(ctx context.Context, params model.CocktailParams) (*model.CocktailDetail, error)
	GetListByIDs(ctx context.Context, ids []int64) ([]model.Cocktail, error)
}

