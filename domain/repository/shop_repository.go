package repository

import (
	"context"
	"github.com/shake551/cocktails-api/domain/model"
)

type ShopRepository interface {
	GetLimit(ctx context.Context, limit int64, offset int64) ([]model.Shop, error)
	Create(ctx context.Context, params model.ShopParams) (*model.Shop, error)
	GetByID(ctx context.Context, id int64) (model.Shop, error)
	GetShopCocktailList(ctx context.Context, shopID int64, limit int64, offset int64) ([]model.Cocktail, error)
	AddShopCocktail(ctx context.Context, shopID int64, params model.ShopCocktailParams) ([]*model.ShopCocktail, error)
	GetShopCocktailDetail(ctx context.Context, shopID int64, cocktailID int64) (model.CocktailDetail, error)
}
