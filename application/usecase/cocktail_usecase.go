package usecase

import (
	"context"
	"github.com/shake551/cocktails-api/domain/model"
	"github.com/shake551/cocktails-api/domain/repository"
)

type CocktailUseCase interface {
	GetLimit(ctx context.Context, limit int64, offset int64, keyword string) ([]model.Cocktail, error)
	GetById(ctx context.Context, id int64) (model.CocktailDetail, error)
	Create(ctx context.Context, params model.CocktailParams) (*model.CocktailDetail, error)
	GetListByIDs(ctx context.Context, ids []int64) ([]model.Cocktail, error)
}

type cocktailUseCase struct {
	repository.CocktailRepository
}

func NewCocktailUseCase(r repository.CocktailRepository) CocktailUseCase {
	return &cocktailUseCase{r}
}

func (u *cocktailUseCase) GetLimit(ctx context.Context, limit int64, offset int64, keyword string) ([]model.Cocktail, error) {
	return u.CocktailRepository.GetLimit(ctx, limit, offset, keyword)
}

func (u *cocktailUseCase) GetById(ctx context.Context, id int64) (model.CocktailDetail, error) {
	return u.CocktailRepository.GetByID(ctx, id)
}

func (u *cocktailUseCase) Create(ctx context.Context, params model.CocktailParams) (*model.CocktailDetail, error) {
	return u.CocktailRepository.Create(ctx, params)
}

func (u *cocktailUseCase) GetListByIDs(ctx context.Context, ids []int64) ([]model.Cocktail, error) {
	return u.CocktailRepository.GetListByIDs(ctx, ids)
}
