package db

import (
	"database/sql"
	"encoding/json"
	"log"
	"time"

	"github.com/harshadmanglani/polaris"
	_ "github.com/lib/pq"
)

type PostgresDataStore struct {
	db *sql.DB
}

func Init() {
	mockStorage := &MockStorage{
		store: make(map[string]interface{}),
	}
	polaris.InitRegistry(mockStorage)
}

type Incident struct {
	ID        string    `json:"id"`
	Timestamp time.Time `json:"timestamp"`
	Service   string    `json:"service"`
	Reasoning string    `json:"reasoning"`
	Summary   string    `json:"summary"`
	Type      string    `json:"type"`
}

func NewPostgresDataStore(connectionString string) (*PostgresDataStore, error) {
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	// Create workflows table if not exists
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS workflows (
			key VARCHAR PRIMARY KEY,
			value JSONB
		)
	`)
	if err != nil {
		return nil, err
	}

	// Create incidents table if not exists
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS incidents (
			id VARCHAR PRIMARY KEY,
			timestamp TIMESTAMP,
			service VARCHAR,
			reasoning TEXT,
			summary TEXT,
			type VARCHAR
		)
	`)
	if err != nil {
		return nil, err
	}

	return &PostgresDataStore{db: db}, nil
}

// Write stores a workflow
func (p *PostgresDataStore) Write(key string, value interface{}) {
	jsonData, err := json.Marshal(value)
	if err != nil {
		return
	}

	_, err = p.db.Exec(`
		INSERT INTO workflows (key, value)
		VALUES ($1, $2)
		ON CONFLICT (key) DO UPDATE SET value = $2
	`, key, jsonData)
	log.Printf("[Error] writing workflow to database: %v", err)
}

// Read retrieves a workflow
func (p *PostgresDataStore) Read(key string) (interface{}, bool) {
	var jsonData []byte
	err := p.db.QueryRow("SELECT value FROM workflows WHERE key = $1", key).Scan(&jsonData)
	if err != nil {
		return nil, false
	}

	return jsonData, true
}

// WriteIncident creates a new incident
func (p *PostgresDataStore) WriteIncident(incident *Incident) error {
	_, err := p.db.Exec(`
		INSERT INTO incidents (id, timestamp, service, reasoning, summary, type)
		VALUES ($1, $2, $3, $4, $5, $6)
	`, incident.ID, incident.Timestamp, incident.Service, incident.Reasoning, incident.Summary, incident.Type)
	return err
}

// ReadIncident retrieves an incident by ID
func (p *PostgresDataStore) ReadIncident(id string) (*Incident, bool) {
	incident := &Incident{}
	err := p.db.QueryRow(`
		SELECT id, timestamp, service, reasoning, summary, type 
		FROM incidents WHERE id = $1
	`, id).Scan(&incident.ID, &incident.Timestamp, &incident.Service, &incident.Reasoning, &incident.Summary, &incident.Type)

	if err != nil {
		return nil, false
	}
	return incident, true
}

// UpdateIncident updates an existing incident
func (p *PostgresDataStore) UpdateIncident(incident *Incident) error {
	result, err := p.db.Exec(`
		UPDATE incidents 
		SET timestamp = $2, service = $3, reasoning = $4, summary = $5, type = $6
		WHERE id = $1
	`, incident.ID, incident.Timestamp, incident.Service, incident.Reasoning, incident.Summary, incident.Type)

	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return sql.ErrNoRows
	}
	return nil
}

type MockStorage struct {
	store map[string]interface{}
}

func (ms *MockStorage) Read(key string) (interface{}, bool) {
	val, ok := ms.store[key]
	return val, ok
}

func (ms *MockStorage) Write(key string, val interface{}) {
	ms.store[key] = val
}
