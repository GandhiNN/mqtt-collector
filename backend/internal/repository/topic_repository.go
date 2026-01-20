package repository

import (
	"database/sql"
	"fmt"
	"mqtt-catalog/pkg/models"
	"time"
)

type TopicRepository struct {
	db *sql.DB
}

func NewTopicRepository(db *sql.DB) *TopicRepository {
	return &TopicRepository{db: db}
}

func (r *TopicRepository) Upsert(sample models.Sample) error {
	query := `
		INSERT INTO topics (broker_id, topic, payload_type, sample_payload, last_seen, created_at)
		VALUES (?, ?, ?, ?, ?, ?)
		ON CONFLICT(broker_id, topic)
		DO UPDATE SET
			payload_type = excluded.payload_type,
			sample_payload = excluded.sample_payload,
			last_seen = excluded.last_seen
	`
	// For Postgres, use $1, $2 syntax instead of ?
	var err error
	var version string
	checkErr := r.db.QueryRow("SELECT sqlite_version()").Scan(&version)
	isSQLite := checkErr == nil

	if !isSQLite {
		// Postgres query
		query = `
			INSERT INTO topics (broker_id, topic, payload_type, sample_payload, last_seen, created_at)
			VALUES ($1, $2, $3, $4, $5, $6)
			ON CONFLICT(broker_id, topic)
			DO UPDATE SET
				payload_type = EXCLUDED.payload_type,
				sample_payload = EXCLUDED.sample_payload,
				last_seen = EXCLUDED.last_seen
		`
	}

	now := time.Now()
	_, err = r.db.Exec(query,
		sample.BrokerID,
		sample.Topic,
		sample.PayloadType,
		sample.Payload,
		sample.Timestamp,
		now,
	)

	if err != nil {
		return fmt.Errorf("upsert topic: %w", err)
	}

	return nil
}

func (r *TopicRepository) GetAll(limit, offset int) ([]models.Topic, int, error) {
	// Get total count
	var total int
	err := r.db.QueryRow("SELECT COUNT(*) FROM topics").Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("count topics: %w", err)
	}

	// Get paginated results
	query := `
		SELECT id, broker_id, topic, payload_type, sample_payload, last_seen, created_at
		FROM topics
		ORDER BY last_seen DESC
		LIMIT ? OFFSET ?
	`

	// Check if Postgres
	var version string
	checkErr := r.db.QueryRow("SELECT sqlite_version()").Scan(&version)
	if checkErr != nil {
		// Postgres
		query = `
			SELECT id, broker_id, topic, payload_type, sample_payload, last_seen, created_at
			FROM topics
			ORDER BY last_seen DESC
			LIMIT $1 OFFSET $2
		`
	}

	rows, err := r.db.Query(query, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("query topics: %w", err)
	}
	defer rows.Close()

	var topics []models.Topic
	for rows.Next() {
		var t models.Topic
		err := rows.Scan(
			&t.ID,
			&t.BrokerID,
			&t.Topic,
			&t.PayloadType,
			&t.SamplePayload,
			&t.LastSeen,
			&t.CreatedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("scan topic: %w", err)
		}
		topics = append(topics, t)
	}
	return topics, total, nil
}
