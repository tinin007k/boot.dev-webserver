package main

import (
	"log"
	"net/http"
	"strconv"
)

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits++
		next.ServeHTTP(w, r)
	})

	//log.Println("file hit value is", cfg.fileserverHits)
	//fs := http.StripPrefix("/app/", http.FileServer(http.Dir(".")))

}

func addHeaders(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Cache-Control", "no-cache")
		next.ServeHTTP(w, r)
	}
}

func (cfg *apiConfig) metricsHandler(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		res := "Hits: " + strconv.Itoa(cfg.fileserverHits)
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		if r.Method == "GET" {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(res))
			return
		}
		next.ServeHTTP(w, r)
	}
}

func (cfg *apiConfig) metricsReset(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits = 0
		if r.Method == "GET" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func middlewareLog(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s", r.URL, r.Method)
		next.ServeHTTP(w, r)
	})
}
