package api

import (
	"log"
	"mqtt-catalog/internal/repository"
	"net/http"
)

func NewRouter(repo *repository.TopicRepository) http.Handler {
	mux := http.NewServeMux()
	handler := NewHandler(repo)

	mux.HandleFunc("POST /api/samples", handler.CreateSample)
	mux.HandleFunc("GET /api/topics", handler.GetTopics)
	mux.HandleFunc("GET /api/topics/search", handler.GetTopic)
	mux.HandleFunc("GET /health", handler.HealthCheck)

	return corsMiddleware(loggingMiddleware(mux))
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}
