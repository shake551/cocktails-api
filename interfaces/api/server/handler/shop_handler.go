package handler

import (
	"encoding/json"
	"github.com/go-chi/chi"
	"github.com/shake551/cocktails-api/application/usecase"
	"github.com/shake551/cocktails-api/domain/model"
	"log"
	"net/http"
	"strconv"
)

type ShopHandler interface {
	GetLimit(w http.ResponseWriter, r *http.Request)
	Create(w http.ResponseWriter, r *http.Request)
	GetByID(w http.ResponseWriter, r *http.Request)
	GetShopCocktailList(w http.ResponseWriter, r *http.Request)
	AddShopCocktail(w http.ResponseWriter, r *http.Request)
	GetShopCocktailDetail(w http.ResponseWriter, r *http.Request)
	GetUnprovidedOrderList(w http.ResponseWriter, r *http.Request)
	AddTable(w http.ResponseWriter, r *http.Request)
	GetTable(w http.ResponseWriter, r *http.Request)
	GetTableOrderList(w http.ResponseWriter, r *http.Request)
}

type shopHandler struct {
	u usecase.ShopUseCase
}

func NewShopHandler(u usecase.ShopUseCase) ShopHandler {
	return &shopHandler{u}
}

func (h *shopHandler) GetLimit(w http.ResponseWriter, r *http.Request) {
	v := r.URL.Query()
	if v == nil {
		return
	}

	var limit = int64(10)
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

	shops, err := h.u.GetLimit(r.Context(), limit, offset)
	if err != nil {
		log.Printf("failed to get cocktails. err: %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	b, err := json.Marshal(shops)
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

func (h *shopHandler) Create(w http.ResponseWriter, r *http.Request) {
	body := model.ShopParams{}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		log.Printf("bad request error. err: %v, body:%v", err, body)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	s, err := h.u.Create(r.Context(), body)
	if err != nil {
		log.Printf("failed to create shop. err: %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	b, err := json.Marshal(s)
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

func (h *shopHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "shopID"), 10, 64)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	cocktailsDetail, err := h.u.GetByID(r.Context(), id)
	if err != nil {
		log.Printf("failed to get shop with id. err: %v", err)
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

func (h *shopHandler) GetShopCocktailList(w http.ResponseWriter, r *http.Request) {
	shopID, err := strconv.ParseInt(chi.URLParam(r, "shopID"), 10, 64)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	v := r.URL.Query()
	if v == nil {
		return
	}

	var limit = int64(10)
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

	c, err := h.u.GetShopCocktailList(r.Context(), shopID, limit, offset)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	b, err := json.Marshal(c)
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

func (h *shopHandler) AddShopCocktail(w http.ResponseWriter, r *http.Request) {
	shopID, err := strconv.ParseInt(chi.URLParam(r, "shopID"), 10, 64)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	body := model.ShopCocktailParams{}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		log.Printf("bad request error. err: %v, body:%v", err, body)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	c, err := h.u.AddShopCocktail(r.Context(), shopID, body)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	b, err := json.Marshal(c)
	if err != nil {
		log.Printf("failed to parse json. err: %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Length", strconv.Itoa(len(b)))
	w.WriteHeader(http.StatusCreated)
	w.Write(b)
}

func (h *shopHandler) GetShopCocktailDetail(w http.ResponseWriter, r *http.Request) {
	shopID, err := strconv.ParseInt(chi.URLParam(r, "shopID"), 10, 64)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	cocktailID, err := strconv.ParseInt(chi.URLParam(r, "cocktailID"), 10, 64)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	d, err := h.u.GetShopCocktailDetail(r.Context(), shopID, cocktailID)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	b, err := json.Marshal(d)
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

func (h *shopHandler) GetUnprovidedOrderList(w http.ResponseWriter, r *http.Request) {
	shopID, err := strconv.ParseInt(chi.URLParam(r, "shopID"), 10, 64)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	v := r.URL.Query()
	if v == nil {
		return
	}

	var limit = int64(10)
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

	var unprovided bool
	if v.Get("unprovided") != "" {
		unprovided, err = strconv.ParseBool(v.Get("unprovided"))
		if err != nil {
			log.Printf("bad request error. err: %v, param:%v", err, v.Get("unprovided"))
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
	}

	if !unprovided {
		return
	}

	os, err := h.u.GetUnprovidedOrderList(r.Context(), shopID, limit, offset)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	b, err := json.Marshal(os)
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

func (h *shopHandler) AddTable(w http.ResponseWriter, r *http.Request) {
	shopID, err := strconv.ParseInt(chi.URLParam(r, "shopID"), 10, 64)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	t, err := h.u.AddTable(r.Context(), shopID)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	b, err := json.Marshal(t)
	if err != nil {
		log.Printf("failed to parse json. err: %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Length", strconv.Itoa(len(b)))
	w.WriteHeader(http.StatusCreated)
	w.Write(b)
}

func (h *shopHandler) GetTable(w http.ResponseWriter, r *http.Request) {
	shopID, err := strconv.ParseInt(chi.URLParam(r, "shopID"), 10, 64)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	tableID, err := strconv.ParseInt(chi.URLParam(r, "tableID"), 10, 64)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	t, err := h.u.GetTable(r.Context(), shopID, tableID)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	b, err := json.Marshal(t)
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

func (h *shopHandler) GetTableOrderList(w http.ResponseWriter, r *http.Request) {
	shopID, err := strconv.ParseInt(chi.URLParam(r, "shopID"), 10, 64)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	tableID, err := strconv.ParseInt(chi.URLParam(r, "tableID"), 10, 64)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	v := r.URL.Query()

	var unprovided bool
	if v.Get("unprovided") != "" {
		unprovided, err = strconv.ParseBool(v.Get("unprovided"))
		if err != nil {
			log.Printf("bad request error. err: %v, param:%v", err, v.Get("unprovided"))
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
	}

	os, err := h.u.GetTableOrderList(r.Context(), shopID, tableID, unprovided)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	b, err := json.Marshal(os)
	if err != nil {
		log.Printf("failed to parse json. err: %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Length", strconv.Itoa(len(b)))
	w.WriteHeader(http.StatusCreated)
	w.Write(b)
}
