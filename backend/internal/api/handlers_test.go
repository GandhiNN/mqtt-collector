package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"mqtt-catalog/internal/database"
	"mqtt-catalog/internal/repository"
	"mqtt-catalog/pkg/models"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

func setupTestDB(t *testing.T) *sql.DB {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("Failed to open test database: %v", err)
	}

	if err := database.RunMigrations(db); err != nil {
		t.Fatalf("Failed to run migrations: %v", err)
	}

	return db
}

func TestHandler_Createsample(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := repository.NewTopicRepository(db)
	handler := NewHandler(repo)

	sample := models.Sample{
		BrokerID:    "test-broker",
		Topic:       "test/topic",
		PayloadType: models.PayloadJSON,
		Payload:     []byte(`{"temp": 22.5}`),
		Timestamp:   time.Now(),
	}

	body, _ := json.Marshal(sample)
	req := httptest.NewRequest(http.MethodPost, "/api/samples", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.CreateSample(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("expected status %d, got %d", http.StatusCreated, w.Code)
	}
}

func TestHandler_CreateSample_InvalidJSON(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := repository.NewTopicRepository(db)
	handler := NewHandler(repo)

	req := httptest.NewRequest(
		http.MethodPost,
		"/api/samples",
		bytes.NewBufferString("invalid json"),
	)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.CreateSample(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestHandler_CreateSample_MissingFields(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := repository.NewTopicRepository(db)
	handler := NewHandler(repo)

	sample := models.Sample{
		PayloadType: models.PayloadJSON,
		Payload:     []byte(`{"temp": 22.5}`),
		Timestamp:   time.Now(),
	}

	body, _ := json.Marshal(sample)
	req := httptest.NewRequest(http.MethodPost, "/api/samples", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.CreateSample(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestHandler_GetTopics(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := repository.NewTopicRepository(db)
	handler := NewHandler(repo)

	sample := models.Sample{
		BrokerID:    "test-broker",
		Topic:       "test/topic",
		PayloadType: models.PayloadJSON,
		Payload:     []byte(`{"temp": 22.5}`),
		Timestamp:   time.Now(),
	}
	repo.Upsert(sample)

	req := httptest.NewRequest(http.MethodGet, "/api/topics", nil)
	w := httptest.NewRecorder()

	handler.GetTopics(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}

	var response models.TopicListResponse
	if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if response.Total != 1 {
		t.Errorf("expected total 1, got %d", response.Total)
	}

	if len(response.Topics) != 1 {
		t.Errorf("expected 1 topic, got %d", len(response.Topics))
	}
}

func TestHandler_GetTopics_WithBrokerFilter(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := repository.NewTopicRepository(db)
	handler := NewHandler(repo)

	repo.Upsert(models.Sample{
		BrokerID:    "broker1",
		Topic:       "topic1",
		PayloadType: models.PayloadJSON,
		Payload:     []byte(`{}`),
		Timestamp:   time.Now(),
	})
	repo.Upsert(models.Sample{
		BrokerID:    "broker2",
		Topic:       "topic2",
		PayloadType: models.PayloadJSON,
		Payload:     []byte(`{}`),
		Timestamp:   time.Now(),
	})

	req := httptest.NewRequest(http.MethodGet, "/api/topics?broker_id=broker1", nil)
	w := httptest.NewRecorder()

	handler.GetTopics(w, req)

	var response models.TopicListResponse
	json.NewDecoder(w.Body).Decode(&response)

	if response.Total != 1 {
		t.Errorf("expected total 1, got %d", response.Total)
	}

	if response.Topics[0].BrokerID != "broker1" {
		t.Errorf("expected broker_id 'broker1' got '%s'", response.Topics[0].BrokerID)
	}
}
