package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
)

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits++
		next.ServeHTTP(w, r)
	})

	// log.Println("file hit value is", cfg.fileserverHits)
	// fs := http.StripPrefix("/app/", http.FileServer(http.Dir(".")))
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

func (cfg *apiConfig) adminMetricsHtmlHandler(
	next http.Handler,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		htmpTemplate := "<html><body><h1>Welcome, Chirpy Admin</h1><p>Chirpy has been visited %d times!</p></body></html>"
		res := fmt.Sprintf(htmpTemplate, cfg.fileserverHits)
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		if r.Method == "GET" {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(res))
			return
		}
		// next.ServeHTTP(w, r)
	}
}

func (cfg *apiConfig) validateChirp(next http.Handler) http.HandlerFunc {
	type response struct {
		Valid bool `json:"valid"`
	}

	res := response{Valid: true}

	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "appplication/json; charset=utf-8")
		decoder := json.NewDecoder(r.Body)
		chirpReq := chirpRequest{}
		err := decoder.Decode(&chirpReq)
		if err == nil && chirpReq.Body == "" {
			errMesg := "Something went wrong"
			respondWithError(w, r, 400, errMesg)
			return
		}

		if len(chirpReq.Body) > 140 {
			errMesg := "Chirp is too long"
			respondWithError(w, r, 400, errMesg)
			return
		}

		dat, _ := json.Marshal(res)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(dat))
		return
		// next.ServeHTTP(w, r)
	}
}

func respondWithError(
	w http.ResponseWriter,
	r *http.Request,
	code int,
	msg string,
) {
	errMesg := error{ErrorMesg: msg}
	dat, _ := json.Marshal(errMesg)
	w.WriteHeader(code)
	w.Write(dat)
	// return
}

func (cfg *apiConfig) metricsReset(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits = 0
		if r.Method == "GET" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	}
}

func middlewareLog(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s", r.URL, r.Method)
		next.ServeHTTP(w, r)
	})
}
