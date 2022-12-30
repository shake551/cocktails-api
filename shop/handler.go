package shop

import (
	"github.com/go-chi/chi"
	"net/http"
	"strconv"
)

var sr = NewShopRepository()

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
