package shop

import (
	"encoding/json"
	"github.com/go-chi/chi"
	"log"
	"net/http"
	"strconv"
)

var sr = NewShopRepository()

func CreateShopTableHandler(w http.ResponseWriter, r *http.Request) {
	shopID, err := strconv.ParseInt(chi.URLParam(r, "shopID"), 10, 64)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	t, err := sr.CreateTable(r.Context(), shopID)
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

func GetTableHandler(w http.ResponseWriter, r *http.Request) {
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

	t, err := sr.GetTable(r.Context(), shopID, tableID)
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

func AddShopCocktailHandler(w http.ResponseWriter, r *http.Request) {
	shopID, err := strconv.ParseInt(chi.URLParam(r, "shopID"), 10, 64)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	body := ShopCocktailParams{}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		log.Printf("bad request error. err: %v, body:%v", err, body)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	c, err := sr.AddShopCocktails(r.Context(), shopID, body)
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

func OrderHandler(w http.ResponseWriter, r *http.Request) {
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

	body := OrderParams{}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		log.Printf("bad request error. err: %v, body:%v", err, body)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	o, err := sr.Order(r.Context(), shopID, tableID, body)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	b, err := json.Marshal(o)
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

func OrderProvide(w http.ResponseWriter, r *http.Request) {
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

	orderID, err := strconv.ParseInt(chi.URLParam(r, "orderID"), 10, 64)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	err = sr.OrderProvide(r.Context(), shopID, tableID, orderID)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
}

func GetTableOrderListHandler(w http.ResponseWriter, r *http.Request) {
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

	os, err := sr.GetTableOrderList(r.Context(), shopID, tableID, unprovided)
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

func GetShopUnprovidedOrderList(w http.ResponseWriter, r *http.Request) {
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

	os, err := sr.GetShopUnprovidedOrderList(r.Context(), shopID, limit, offset)
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

func GetShopCocktailsList(w http.ResponseWriter, r *http.Request) {
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

	c, err := sr.GetShopCocktailsList(r.Context(), shopID, limit, offset)
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

func GetShopCocktailDetailHandler(w http.ResponseWriter, r *http.Request) {
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

	d, err := sr.GetShopCocktailDetail(r.Context(), shopID, cocktailID)
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
