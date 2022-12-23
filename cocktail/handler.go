package cocktail

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi"
	"log"
	"net/http"
	"strconv"
)

type Handler interface {
	GetCocktailsHandler(w http.ResponseWriter, r *http.Request)
	FindCocktailsDetailByID(w http.ResponseWriter, r *http.Request)
	PostCocktailsHandler(w http.ResponseWriter, r *http.Request)
}

type CocktailsHandler struct {
	r Repository
}

func NewCocktailsHandler() *CocktailsHandler {
	return &CocktailsHandler{r: NewCocktailsRepository()}
}

func (h CocktailsHandler) GetCocktailsHandler(w http.ResponseWriter, r *http.Request) {
	v := r.URL.Query()
	if v == nil {
		return
	}

	var limit = int64(30)
	if v.Get("limit") != "" {
		l, err := strconv.ParseInt(v.Get("limit"), 10, 64)
		if err != nil {
			log.Printf("failed to get limit. err: %v", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		limit = l
	}

	var offset = int64(0)
	if v.Get("offset") != "" {
		o, err := strconv.ParseInt(v.Get("offset"), 10, 64)
		if err != nil {
			log.Printf("failed to get offset. err: %v", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		offset = o
	}

	var keyword = v.Get("keyword")

	cocktails, err := h.r.GetLimit(r.Context(), limit, offset, keyword)
	if err != nil {
		log.Printf("failed to get cocktails. err: %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	b, err := json.Marshal(cocktails)
	if err != nil {
		log.Printf("failed to parse json. err: %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Length", strconv.Itoa(len(b)))
	w.WriteHeader(http.StatusOK)
	w.Write(b)
}

func (h CocktailsHandler) FindCocktailsDetailByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "cocktailsID"), 10, 64)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	cocktailsDetail, err := h.r.FindCocktailsDetailByID(r.Context(), id)
	if err != nil {
		log.Printf("failed to get cocktails detail. err: %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	b, err := json.Marshal(cocktailsDetail)
	if err != nil {
		log.Printf("failed to parse json. err: %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Length", strconv.Itoa(len(b)))
	w.WriteHeader(http.StatusOK)
	w.Write(b)
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
		materials = append(materials, MaterialParams{Name: material.Name, Quantity: MaterialQuantity{
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
