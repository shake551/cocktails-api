package main

import (
	"context"
	"fmt"
	"github.com/shake551/cocktails-api/application/usecase"
	"github.com/shake551/cocktails-api/infrastructure/parsistence/datastore"
	"github.com/shake551/cocktails-api/interfaces/api/server/handler"
	"log"
	"mime"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/lestrrat-go/server-starter/listener"

	"github.com/shake551/cocktails-api/db"
)

const port = 8080
const logDir = "/var/log/app/cocktail"

func netListen(network, addr string) (net.Listener, error) {
	ls, err := listener.ListenAll()
	if err == nil {
		return ls[0], nil
	}
	return net.Listen(network, addr)
}

func getAccessLogFormatter() middleware.LogFormatter {
	err := os.MkdirAll(logDir, os.ModePerm|os.ModeDir)
	if err != nil {
		log.Fatalf("failed to prepare access log dir: %v", err)
	}

	logf, err := rotatelogs.New(
		filepath.Join(logDir, "access_log.%Y%m%d%H%M"),
		rotatelogs.WithLinkName(filepath.Join(logDir, "access_log")),
		rotatelogs.WithMaxAge(24*time.Hour),
		rotatelogs.WithRotationTime(time.Hour),
	)
	if err != nil {
		log.Fatalf("failed to open access log file: %v", err)
	}

	return &middleware.DefaultLogFormatter{Logger: log.New(logf, "", log.LstdFlags), NoColor: false}
}

func createRouter() chi.Router {
	mux := chi.NewRouter()
	mux.Use(cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "HEAD", "OPTIONS"},
		AllowedHeaders: []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
	}).Handler)
	mux.Use(middleware.RequestLogger(getAccessLogFormatter()))
	mux.Use(contentTypeRestrictionMiddleware("application/json"))

	cr := datastore.NewCocktailRepository()
	cu := usecase.NewCocktailUseCase(cr)
	ch := handler.NewCocktailHandler(cu)

	sr := datastore.NewShopRepository()
	su := usecase.NewShopUseCase(sr)
	sh := handler.NewShopHandler(su)

	// no auth
	mux.Group(func(mux chi.Router) {
		mux.MethodFunc("GET", "/health", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		})

		mux.MethodFunc("GET", "/cocktails", ch.GetLimit)
		mux.MethodFunc("POST", "/cocktails", ch.Create)
		mux.MethodFunc("GET", "/cocktails/{cocktailsID}", ch.GetById)
		mux.MethodFunc("GET", "/cocktails/list", ch.GetListByIDs)

		mux.MethodFunc("GET", "/shop", sh.GetLimit)
		mux.MethodFunc("POST", "/shop", sh.Create)
		mux.MethodFunc("GET", "/shop/{shopID}", sh.GetByID)
		mux.MethodFunc("GET", "/shop/{shopID}/cocktail", sh.GetShopCocktailList)
		mux.MethodFunc("POST", "/shop/{shopID}/cocktail", sh.AddShopCocktail)
		mux.MethodFunc("GET", "/shop/{shopID}/cocktail/{cocktailID}", sh.GetShopCocktailDetail)
		mux.MethodFunc("GET", "/shop/{shopID}/order", sh.GetUnprovidedOrderList)
		mux.MethodFunc("POST", "/shop/{shopID}/table", sh.AddTable)
		mux.MethodFunc("GET", "/shop/{shopID}/table/{tableID}", sh.GetTable)
		mux.MethodFunc("GET", "/shop/{shopID}/table/{tableID}/order", sh.GetTableOrderList)
		mux.MethodFunc("POST", "/shop/{shopID}/table/{tableID}/order", sh.Order)
		mux.MethodFunc("PUT", "/shop/{shopID}/table/{tableID}/order/{orderID}", sh.OrderProvide)
	})

	return mux
}

func main() {
	done, err := db.Initialize(os.Getenv("DSN"))
	if err != nil {
		log.Fatalf("failed to initialize db: %v", err)
	}
	defer done()

	mux := createRouter()
	server := http.Server{
		Handler: mux,
	}

	l, err := netListen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	go func() {
		log.Printf("starting server on %s", l.Addr())
		if err := server.Serve(l); err != nil {
			log.Fatalf("server closed with %v", err)
		}
	}()

	// graceful shutdown
	quit := make(chan os.Signal, 1)
	//signal.Notify(quit, syscall.SIGTERM, os.Interrupt)

	log.Printf("SIGNAL %v received, then shutting down...", <-quit)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("failed to graceful shutdown: %v", err)
	}
	log.Print("server shutdown")
}

func contentTypeRestrictionMiddleware(mediaType string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			switch r.Method {
			case "POST", "PUT", "PATCH":
				ct := r.Header.Get("Content-Type")
				if ct == "" {
					log.Print("Empty Content-Type")
					w.WriteHeader(http.StatusBadRequest)
					return
				}

				mt, _, err := mime.ParseMediaType(ct)
				if err != nil {
					log.Printf("Invalid Content-Type: %s", ct)
					w.WriteHeader(http.StatusBadRequest)
					return
				}

				if mt != mediaType {
					log.Printf("Unsupported Content-Type: %s", ct)
					w.WriteHeader(http.StatusUnsupportedMediaType)
					return
				}
			}

			next.ServeHTTP(w, r)
		})
	}
}
