package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	fmt.Println("Hello from web-server")
	mux := http.NewServeMux()
	serveImageFromAssets(mux)
	mux.Handle("/healthz", customHandler(mux))
	corsMux := middlewareCors(mux)
	srvErr := http.ListenAndServe(":8080", corsMux)
	if srvErr != nil {
		log.Fatal(srvErr)
	}
}

func serveImageFromAssets(mux *http.ServeMux) {
	mux.Handle("/app/", http.StripPrefix("/app/", http.FileServer(http.Dir("."))))
	//mux.Handle("/assets/", http.FileServer(http.Dir("./logo.png")))
}

func customHandler(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		response := "OK"
		if r.Method == "GET" {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(response))
			return
		}
		handler.ServeHTTP(w, r)
	})
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
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}
