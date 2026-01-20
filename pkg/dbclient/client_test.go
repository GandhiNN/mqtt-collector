package dbclient

import (
	"context"
	"encoding/json"
	"io"
	"mqtt-collector/internal/models"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestClient_SendSample_Success(t *testing.T) {
	// Create test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request method
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}

		// Verify path
		if r.URL.Path != "/api/samples" {
			t.Errorf("expected /api/samples, got %s", r.URL.Path)
		}

		// Verify content type
		if r.Header.Get("Content-Type") != "application/json" {
			t.Errorf("expected application/json, got %s", r.Header.Get("Content-Type"))
		}

		// Read and verify body
		body, _ := io.ReadAll(r.Body)
		var sample models.Sample
		if err := json.Unmarshal(body, &sample); err != nil {
			t.Errorf("failed to unmarshal sample: %v", err)
		}

		if sample.BrokerID != "test-broker" {
			t.Errorf("expected broker_id 'test-broker', got '%s'", sample.BrokerID)
		}

		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	// Create client
	client := New(server.URL)

	// Create sample
	sample := models.Sample{
		BrokerID:    "test-broker",
		Topic:       "test/topic",
		PayloadType: models.PayloadJSON,
		Payload:     []byte(`{"test": true}`),
		Timestamp:   time.Now(),
	}

	// Send sample
	err := client.SendSample(context.Background(), sample)
	if err != nil {
		t.Errorf("SendSample() error = %v", err)
	}
}

func TestClient_SendSample_ServerError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("internal server error"))
	}))
	defer server.Close()

	client := New(server.URL)

	sample := models.Sample{
		BrokerID:    "test-broker",
		Topic:       "test/topic",
		PayloadType: models.PayloadJSON,
		Payload:     []byte(`{"test": true}`),
		Timestamp:   time.Now(),
	}

	err := client.SendSample(context.Background(), sample)
	if err == nil {
		t.Error("expected error for server error response, got nil")
	}
}

func TestClient_SendSample_ContextCanceled(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(100 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := New(server.URL)

	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	sample := models.Sample{
		BrokerID:    "test-broker",
		Topic:       "test/topic",
		PayloadType: models.PayloadJSON,
		Payload:     []byte(`{"test": true}`),
		Timestamp:   time.Now(),
	}

	err := client.SendSample(ctx, sample)
	if err == nil {
		t.Error("expected error for canceled context, got nil")
	}
}
