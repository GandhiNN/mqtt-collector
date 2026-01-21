package repository

import (
	"database/sql"
	"mqtt-catalog/internal/database"
	"mqtt-catalog/pkg/models"
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

func TestTopicRepository_Upsert_Insert(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewTopicRepository(db)

	sample := models.Sample{
		BrokerID:    "test-broker",
		Topic:       "test/topic",
		PayloadType: models.PayloadJSON,
		Payload:     []byte(`{"temp": 22.5}`),
		Timestamp:   time.Now(),
	}

	err := repo.Upsert(sample)
	if err != nil {
		t.Fatalf("Upsert() error = %v", err)
	}

	topic, err := repo.GetByBrokerAndTopic("test-broker", "test/topic")
	if err != nil {
		t.Fatalf("GetByBrokerAndTopic() error = %v", err)
	}

	if topic == nil {
		t.Fatal("expected topic to be found")
	}

	if topic.BrokerID != "test-broker" {
		t.Errorf("BrokerID = %v, want test-broker", topic.BrokerID)
	}
}

func TestTopicRepository_Upsert_Update(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewTopicRepository(db)

	sample1 := models.Sample{
		BrokerID:    "test-broker",
		Topic:       "test/topic",
		PayloadType: models.PayloadJSON,
		Payload:     []byte(`{"temp": 22.5}`),
		Timestamp:   time.Now(),
	}
	repo.Upsert(sample1)

	sample2 := models.Sample{
		BrokerID:    "test-broker",
		Topic:       "test/topic",
		PayloadType: models.PayloadXML,
		Payload:     []byte(`<temp>25.0</temp>`),
		Timestamp:   time.Now(),
	}
	repo.Upsert(sample2)

	topic, _ := repo.GetByBrokerAndTopic("test-broker", "test/topic")

	if topic.PayloadType != models.PayloadXML {
		t.Errorf("PayloadType = %v, want xml", topic.PayloadType)
	}

	if string(topic.SamplePayload) != `<temp>25.0</temp>` {
		t.Errorf("SamplePayload = %v, want <temp>25.0</temp>", string(topic.SamplePayload))
	}
}

func TestTopicRepository_GetAll(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewTopicRepository(db)

	for i := 0; i < 5; i++ {
		repo.Upsert(models.Sample{
			BrokerID:    "test-broker",
			Topic:       "test/topic" + string(rune(i)),
			PayloadType: models.PayloadJSON,
			Payload:     []byte(`{}`),
			Timestamp:   time.Now(),
		})
	}

	topics, total, err := repo.GetAll(10, 0)
	if err != nil {
		t.Fatalf("GetAll() error = %v", err)
	}

	if total != 5 {
		t.Errorf("total = %v, want 5", total)
	}

	if len(topics) != 5 {
		t.Errorf("len(topics) = %v, want 5", len(topics))
	}
}

func TestTopicRepository_GetAll_Pagination(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewTopicRepository(db)

	for i := 0; i < 10; i++ {
		repo.Upsert(models.Sample{
			BrokerID:    "test-broker",
			Topic:       "test/topic" + string(rune(i)),
			PayloadType: models.PayloadJSON,
			Payload:     []byte(`{}`),
			Timestamp:   time.Now(),
		})
	}

	topics, total, err := repo.GetAll(5, 0)
	if err != nil {
		t.Fatalf("GetAll() error = %v", err)
	}

	if total != 10 {
		t.Errorf("total = %v, want 10", total)
	}

	if len(topics) != 5 {
		t.Errorf("len(topics) = %v, want 5", len(topics))
	}

	topics2, _, _ := repo.GetAll(5, 5)

	if len(topics2) != 5 {
		t.Errorf("len(topics2) = %v, want 5", len(topics2))
	}
}

func TestTopicRepository_GetByBroker(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewTopicRepository(db)

	for i := 0; i < 3; i++ {
		repo.Upsert(models.Sample{
			BrokerID:    "broker1",
			Topic:       "topic" + string(rune(i)),
			PayloadType: models.PayloadJSON,
			Payload:     []byte(`{}`),
			Timestamp:   time.Now(),
		})
	}

	for i := 0; i < 2; i++ {
		repo.Upsert(models.Sample{
			BrokerID:    "broker2",
			Topic:       "topic" + string(rune(i)),
			PayloadType: models.PayloadJSON,
			Payload:     []byte(`{}`),
			Timestamp:   time.Now(),
		})
	}

	topics, total, err := repo.GetByBroker("broker1", 10, 0)
	if err != nil {
		t.Fatalf("GetByBroker() error = %v", err)
	}

	if total != 3 {
		t.Errorf("len(topics) = %v, want 3", len(topics))
	}

	for _, topic := range topics {
		if topic.BrokerID != "broker1" {
			t.Errorf("BrokerID = %v, want broker1", topic.BrokerID)
		}
	}
}

func TestTopicRepository_GetByBrokerAndTopic_NotFound(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewTopicRepository(db)

	topic, err := repo.GetByBrokerAndTopic("nonexistent", "nonexistent")
	if err != nil {
		t.Fatalf("GetByBrokerAndTopic() error = %v", err)
	}

	if topic != nil {
		t.Error("expected nil, got topic")
	}
}
