package cocktail

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
)

var cr = NewCocktailsRepository()

func GetCocktailsHandler(w http.ResponseWriter, r *http.Request) {
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

	channels, err := cr.GetLimit(r.Context(), limit, offset, keyword)
	if err != nil {
		log.Printf("failed to get channels. err: %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	b, err := json.Marshal(channels)
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
	Name      string           `json:"name"`
	Materials []MaterialParams `json:"materials"`
}

func PostCocktailsHandler(w http.ResponseWriter, r *http.Request) {
	body := &PostCocktailsBody{}
	if err := json.NewDecoder(r.Body).Decode(body); err != nil {
		log.Printf("bad request error. err: %v, body: %v", err, body)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	params := CocktailsParams{
		Name:      body.Name,
		Materials: body.Materials,
	}

	coc, err := cr.Create(r.Context(), params)
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
