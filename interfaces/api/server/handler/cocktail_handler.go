package handler

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi"
	"github.com/shake551/cocktails-api/application/usecase"
	"github.com/shake551/cocktails-api/domain/model"
	"log"
	"net/http"
	"strconv"
)

type CocktailHandler interface {
	GetLimit(w http.ResponseWriter, r *http.Request)
	GetById(w http.ResponseWriter, r *http.Request)
	Create(w http.ResponseWriter, r *http.Request)
}

type cocktailHandler struct {
	u usecase.CocktailUseCase
}

func NewCocktailHandler(u usecase.CocktailUseCase) CocktailHandler {
	return &cocktailHandler{u}
}

func (h *cocktailHandler) GetLimit(w http.ResponseWriter, r *http.Request) {
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

	cocktails, err := h.u.GetLimit(r.Context(), limit, offset, keyword)
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

func (h *cocktailHandler) GetById(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "cocktailsID"), 10, 64)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	cocktailsDetail, err := h.u.GetById(r.Context(), id)
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

func (h *cocktailHandler) Create(w http.ResponseWriter, r *http.Request) {
	body := &PostCocktailsBody{}
	if err := json.NewDecoder(r.Body).Decode(body); err != nil {
		log.Printf("bad request error. err: %v, body: %v", err, body)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	var materials []model.MaterialParams
	for _, material := range body.Materials {
		quantity, err := material.Quantity.Quantity.Int64()
		if err != nil {
			fmt.Println(err)
		}
		materials = append(materials, model.MaterialParams{Name: material.Name, Quantity: model.MaterialQuantity{
			Quantity: quantity,
			Unit:     material.Quantity.Unit,
		}})
	}

	params := model.CocktailParams{
		Name:      body.Name,
		Materials: materials,
	}

	coc, err := h.u.Create(r.Context(), params)
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
