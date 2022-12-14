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
	GetUnprovidedOrderList(ctx context.Context, shopID int64, limit int64, offset int64) ([]*model.TableOrder, error)
	AddTable(ctx context.Context, shopID int64) (*model.Table, error)
	GetTable(ctx context.Context, shopID int64, tableID int64) (*model.Table, error)
	GetTableOrderList(ctx context.Context, shopID int64, tableID int64, unprovided bool) ([]*model.TableOrder, error)
	Order(ctx context.Context, shopID int64, tableID int64, params model.OrderParams) ([]*model.Order, error)
	OrderProvide(ctx context.Context, shopID int64, tableID int64, orderID int64) error
}
