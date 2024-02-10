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

	r.Get("/healthz", customHandler(r))
	r.Get("/metrics", apiCfg.metricsHandler(r))
	r.Handle("/reset", apiCfg.metricsReset(r))
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
		w.Header().Set("Access-Control-Allow-Methods", "GET")
		w.Header().Set("Access-Control-Allow-Headers", "*")
		w.Header().Set("Cache-Control", "no-cache")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}
