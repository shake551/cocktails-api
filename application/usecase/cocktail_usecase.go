package usecase

import (
	"context"
	"github.com/shake551/cocktails-api/domain/model"
	"github.com/shake551/cocktails-api/domain/repository"
)

type CocktailUseCase interface {
	GetLimit(ctx context.Context, limit int64, offset int64, keyword string) ([]model.Cocktail, error)
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
