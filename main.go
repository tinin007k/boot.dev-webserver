package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type apiConfig struct {
	fileserverHits int
}

type chirpRequest struct {
	Body string `json:"body"`
}

type error struct {
	ErrorMesg string `json:"error"`
}

func main() {
	fmt.Println("Hello from web-server")
	r := chi.NewRouter()
	//mux := http.NewServeMux()
	apiCfg := apiConfig{fileserverHits: 0}
	fsHandler := apiCfg.middlewareMetricsInc(
		http.StripPrefix("/app", http.FileServer(http.Dir("."))),
	)
	r.Handle("/app", fsHandler)
	r.Handle("/app/*", fsHandler)

	apiRoute := chi.NewRouter()

	apiRoute.Get("/healthz", customHandler(apiRoute))
	apiRoute.Get("/metrics", apiCfg.metricsHandler(apiRoute))
	apiRoute.Get("/reset", apiCfg.metricsReset(apiRoute))
	apiRoute.Post("/validate_chirp", apiCfg.validateChirp(apiRoute))
	r.Mount("/api", apiRoute)

	adminRoute := chi.NewRouter()

	adminRoute.Get("/metrics", apiCfg.adminMetricsHtmlHandler(adminRoute))
	r.Mount("/admin", adminRoute)

	corsMux := middlewareLog(middlewareCors(r))
	srvErr := http.ListenAndServe(":8080", corsMux)
	if srvErr != nil {
		log.Fatal(srvErr)
	}
}

func customHandler(handler http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		response := "OK"
		if r.Method == "GET" {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(response))
			return
		}
		handler.ServeHTTP(w, r)
	}
}

/*
*
First assignment - Using and exploring the CORS
*/
func middlewareCors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*") //https://www.boot.dev
		w.Header().Set("Access-Control-Allow-Methods", "*")
		w.Header().Set("Access-Control-Allow-Headers", "*")
		w.Header().Set("Cache-Control", "no-cache")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}
