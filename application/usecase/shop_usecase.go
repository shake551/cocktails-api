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
	GetShopCocktailList(ctx context.Context, shopID int64, limit int64, offset int64) ([]model.Cocktail, error)
	AddShopCocktail(ctx context.Context, shopID int64, params model.ShopCocktailParams) ([]*model.ShopCocktail, error)
	GetShopCocktailDetail(ctx context.Context, shopID int64, cocktailID int64) (model.CocktailDetail, error)
	GetUnprovidedOrderList(ctx context.Context, shopID int64, limit int64, offset int64) ([]*model.TableOrder, error)
	AddTable(ctx context.Context, shopID int64) (*model.Table, error)
	GetTable(ctx context.Context, shopID int64, tableID int64) (*model.Table, error)
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

func (u *shopUseCase) GetShopCocktailList(ctx context.Context, shopID int64, limit int64, offset int64) ([]model.Cocktail, error) {
	return u.ShopRepository.GetShopCocktailList(ctx, shopID, limit, offset)
}

func (u *shopUseCase) AddShopCocktail(ctx context.Context, shopID int64, params model.ShopCocktailParams) ([]*model.ShopCocktail, error) {
	return u.ShopRepository.AddShopCocktail(ctx, shopID, params)
}

func (u *shopUseCase) GetShopCocktailDetail(ctx context.Context, shopID int64, cocktailID int64) (model.CocktailDetail, error) {
	return u.ShopRepository.GetShopCocktailDetail(ctx, shopID, cocktailID)
}

func (u *shopUseCase) GetUnprovidedOrderList(ctx context.Context, shopID int64, limit int64, offset int64) ([]*model.TableOrder, error) {
	return u.ShopRepository.GetUnprovidedOrderList(ctx, shopID, limit, offset)
}

func (u *shopUseCase) AddTable(ctx context.Context, shopID int64) (*model.Table, error) {
	return u.ShopRepository.AddTable(ctx, shopID)
}

func (u *shopUseCase) GetTable(ctx context.Context, shopID int64, tableID int64) (*model.Table, error) {
	return u.ShopRepository.GetTable(ctx, shopID, tableID)
}
