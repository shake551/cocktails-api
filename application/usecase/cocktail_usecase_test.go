package usecase

import (
	"context"
	"testing"

	"github.com/shake551/cocktails-api/domain/model"
	"github.com/shake551/cocktails-api/domain/repository_mock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetLimit(t *testing.T) {
	type testcase struct {
		name    string
		limit   int64
		offset  int64
		keyword string
		want    []model.Cocktail
	}

	tests := []testcase{
		{
			name:    "get two cocktail",
			limit:   int64(2),
			offset:  int64(0),
			keyword: "",
			want: []model.Cocktail{
				{
					ID:        1,
					Name:      "hoge",
					ImageURL:  "",
					CreatedAt: 1000000000,
					UpdatedAt: 1000000000,
				},
				{
					ID:        2,
					Name:      "fuga",
					ImageURL:  "",
					CreatedAt: 1000000000,
					UpdatedAt: 1000000000,
				},
			},
		},
	}

	r := new(repository_mock.CocktailRepository)

	cocktails := []model.Cocktail{
		{
			ID:        1,
			Name:      "hoge",
			ImageURL:  "",
			CreatedAt: 1000000000,
			UpdatedAt: 1000000000,
		},
		{
			ID:        2,
			Name:      "fuga",
			ImageURL:  "",
			CreatedAt: 1000000000,
			UpdatedAt: 1000000000,
		},
	}

	r.On("GetLimit", mock.Anything, int64(2), int64(0), "").Return(cocktails, nil)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uc := &cocktailUseCase{r}
			res, err := uc.GetLimit(context.Background(), tt.limit, tt.offset, tt.keyword)
			assert.Equal(t, res, tt.want)
			assert.Nil(t, err)
		})
	}
}

func TestGetById(t *testing.T) {
	type testcase struct {
		Name string
		ID   int64
		Want model.CocktailDetail
	}

	tests := []testcase{
		{
			Name: "success",
			ID:   1,
			Want: model.CocktailDetail{
				ID:       1,
				Name:     "ゴットファーザー",
				ImageURL: "",
				Materials: []model.Material{
					{
						ID:   1,
						Name: "ウイスキー",
						Quantity: model.MaterialQuantity{
							Quantity: 30,
							Unit:     "ml",
						},
					},
					{
						ID:   2,
						Name: "アマレット",
						Quantity: model.MaterialQuantity{
							Quantity: 10,
							Unit:     "ml",
						},
					},
				},
			},
		},
	}

	r := new(repository_mock.CocktailRepository)

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			r.On("GetByID", mock.Anything, int64(1)).Return(tc.Want, nil)
			uc := &cocktailUseCase{r}
			res, err := uc.GetById(context.Background(), tc.ID)

			assert.Equal(t, res, tc.Want)
			assert.Nil(t, err)
		})
	}
}

func TestCreate(t *testing.T) {
	type testcase struct {
		Name  string
		Input model.CocktailParams
		Want  *model.CocktailDetail
	}

	tests := []testcase{
		{
			Name: "success",
			Input: model.CocktailParams{
				Name: "ゴットファーザー",
				Materials: []model.MaterialParams{
					{
						Name: "ウイスキー",
						Quantity: model.MaterialQuantity{
							Quantity: 30,
							Unit:     "ml",
						},
					},
					{
						Name: "アマレット",
						Quantity: model.MaterialQuantity{
							Quantity: 10,
							Unit:     "ml",
						},
					},
				},
			},
			Want: &model.CocktailDetail{
				ID:       int64(1),
				Name:     "ゴットファーザー",
				ImageURL: "",
				Materials: []model.Material{
					{
						ID:   int64(1),
						Name: "ウイスキー",
						Quantity: model.MaterialQuantity{
							Quantity: 30,
							Unit:     "ml",
						},
					},
					{
						ID:   int64(2),
						Name: "アマレット",
						Quantity: model.MaterialQuantity{
							Quantity: 10,
							Unit:     "ml",
						},
					},
				},
			},
		},
	}

	r := new(repository_mock.CocktailRepository)

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			r.On("Create", mock.Anything, tc.Input).Return(tc.Want, nil)
			uc := &cocktailUseCase{r}

			res, err := uc.Create(context.Background(), tc.Input)

			assert.Equal(t, res, tc.Want)
			assert.Nil(t, err)
		})
	}
}
