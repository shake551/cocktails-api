package cocktail

import (
	"encoding/json"
	"fmt"
	"github.com/shake551/cocktails-api/domain/model"
	"log"
	"net/http"
	"strconv"
)

type Handler interface {
	PostCocktailsHandler(w http.ResponseWriter, r *http.Request)
}

type CocktailsHandler struct {
	r Repository
}

func NewCocktailsHandler(r Repository) *CocktailsHandler {
	return &CocktailsHandler{r: r}
}

type PostCocktailsBody struct {
	Name      string                  `json:"name"`
	Materials []PostCocktailsMaterial `json:"materials"`
}

type PostCocktailsMaterial struct {
	Name     string                `json:"name"`
	Quantity PostCocktailsQuantity `json:"quantity"`
}

type PostCocktailsQuantity struct {
	Quantity json.Number `json:"quantity"`
	Unit     string      `json:"unit"`
}

func (h CocktailsHandler) PostCocktailsHandler(w http.ResponseWriter, r *http.Request) {
	body := &PostCocktailsBody{}
	if err := json.NewDecoder(r.Body).Decode(body); err != nil {
		log.Printf("bad request error. err: %v, body: %v", err, body)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	var materials []MaterialParams
	for _, material := range body.Materials {
		quantity, err := material.Quantity.Quantity.Int64()
		if err != nil {
			fmt.Println(err)
		}
		materials = append(materials, MaterialParams{Name: material.Name, Quantity: model.MaterialQuantity{
			Quantity: quantity,
			Unit:     material.Quantity.Unit,
		}})
	}

	params := CocktailsParams{
		Name:      body.Name,
		Materials: materials,
	}

	coc, err := h.r.Create(r.Context(), params)
	if err != nil {
		log.Printf("failed to create cocktail. err: %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	b, err := json.Marshal(coc)
	if err != nil {
		log.Printf("failed to parse json")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Length", strconv.Itoa(len(b)))
	w.WriteHeader(http.StatusOK)
	w.Write(b)
}
