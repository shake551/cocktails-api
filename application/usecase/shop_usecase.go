package usecase

import (
	"context"
	"github.com/shake551/cocktails-api/domain/model"
	"github.com/shake551/cocktails-api/domain/repository"
)

type ShopUseCase interface {
	GetLimit(ctx context.Context, limit int64, offset int64) ([]model.Shop, error)
	Create(ctx context.Context, params model.ShopParams) (*model.Shop, error)
	GetByID(ctx context.Context, id int64) (model.Shop, error)
}

type shopUseCase struct {
	repository.ShopRepository
}

func NewShopUseCase(r repository.ShopRepository) ShopUseCase {
	return &shopUseCase{r}
}

func (u *shopUseCase) GetLimit(ctx context.Context, limit int64, offset int64) ([]model.Shop, error) {
	return u.GetLimit(ctx, limit, offset)
}

func (u *shopUseCase) Create(ctx context.Context, params model.ShopParams) (*model.Shop, error) {
	return u.ShopRepository.Create(ctx, params)
}

func (u *shopUseCase) GetByID(ctx context.Context, id int64) (model.Shop, error) {
	return u.ShopRepository.GetByID(ctx, id)
}
