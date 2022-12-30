package shop

import (
	"encoding/json"
	"github.com/go-chi/chi"
	"log"
	"net/http"
	"strconv"
)

var sr = NewShopRepository()

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
