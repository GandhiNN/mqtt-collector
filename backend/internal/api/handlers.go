package api

import (
	"encoding/json"
	"fmt"
	"log"
	"mqtt-catalog/internal/repository"
	"mqtt-catalog/pkg/models"
	"net/http"
	"strconv"
)

type Handler struct {
	repo *repository.TopicRepository
}

func NewHandler(repo *repository.TopicRepository) *Handler {
	return &Handler{repo: repo}
}

func (h *Handler) CreateSample(w http.ResponseWriter, r *http.Request) {
	var sample models.Sample
	if err := json.NewDecoder(r.Body).Decode(&sample); err != nil {
		http.Error(w, fmt.Sprintf("invalid request body: %v", err), http.StatusBadRequest)
		return
	}

	if sample.BrokerID == "" || sample.Topic == "" {
		http.Error(w, "broker_id and topic are required", http.StatusBadRequest)
		return
	}

	if err := h.repo.Upsert(sample); err != nil {
		log.Printf("Error upserting sample: %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func (h *Handler) GetTopics(w http.ResponseWriter, r *http.Request) {
	limit := parseIntQuery(r, "limit", 100)
	offset := parseIntQuery(r, "offset", 0)
	brokerID := r.URL.Query().Get("broker_id")

	var topics []models.Topic
	var total int
	var err error

	if brokerID != "" {
		topics, total, err = h.repo.GetByBroker(brokerID, limit, offset)
	} else {
		topics, total, err = h.repo.GetAll(limit, offset)
	}

	if err != nil {
		log.Printf("Error getting topics: %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	response := models.TopicListResponse{
		Topics: topics,
		Total:  total,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *Handler) GetTopic(w http.ResponseWriter, r *http.Request) {
	brokerID := r.URL.Query().Get("broker_id")
	topic := r.URL.Query().Get("topic")

	if brokerID == "" || topic == "" {
		http.Error(w, "broker_id and topic are required", http.StatusBadRequest)
		return
	}

	t, err := h.repo.GetByBrokerAndTopic(brokerID, topic)
	if err != nil {
		log.Printf("Error getting topic: %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	if t == nil {
		http.Error(w, "topic not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(t)
}

func (h *Handler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "healthy"})
}

func parseIntQuery(r *http.Request, key string, defaultValue int) int {
	val := r.URL.Query().Get(key)
	if val == "" {
		return defaultValue
	}

	parsed, err := strconv.Atoi(val)
	if err != nil {
		return defaultValue
	}

	return parsed
}
