package handler

import (
	"encoding/json"
	"github.com/go-chi/chi"
	"github.com/shake551/cocktails-api/application/usecase"
	"log"
	"net/http"
	"strconv"
)

type CocktailHandler interface {
	GetLimit(w http.ResponseWriter, r *http.Request)
	GetById(w http.ResponseWriter, r *http.Request)
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
