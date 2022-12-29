package repository

import (
	"context"
	"github.com/shake551/cocktails-api/domain/model"
)

type ShopRepository interface {
	GetLimit(ctx context.Context, limit int64, offset int64) ([]model.Shop, error)
	Create(ctx context.Context, params model.ShopParams) (*model.Shop, error)
}
