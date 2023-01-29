package application_test

import (
	"context"
	"testing"

	"github.com/shake551/cocktails-api/application/usecase"
	model "github.com/shake551/cocktails-api/domain/model"
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
			uc := usecase.NewCocktailUseCase(r)
			res, err := uc.GetLimit(context.Background(), tt.limit, tt.offset, tt.keyword)
			assert.Equal(t, res, tt.want)
			assert.Nil(t, err)
		})
	}
}
