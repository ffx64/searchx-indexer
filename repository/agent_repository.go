package repository

import (
	"database/sql"
	"searchx-indexer/entity"
	"time"
)

type AgentRepository struct {
	DB *sql.DB
}

func NewAgentRepository(db *sql.DB) *AgentRepository {
	return &AgentRepository{DB: db}
}

func (r *AgentRepository) FindByAuthKey(authKey string) (*entity.AgentEntity, error) {
	var a entity.AgentEntity
	err := r.DB.QueryRow(`
		SELECT id, auth_key, platform, collection_interval, last_activity_at,
		       data_processed, agent_status, last_ip_address
		FROM agents
		WHERE auth_key = $1
	`, authKey).Scan(&a.ID, &a.AuthKey, &a.Platform, &a.CollectionInterval,
		&a.LastActivityAt, &a.DataProcessed, &a.AgentStatus, &a.LastIPAddress)

	if err != nil {
		return nil, err
	}
	return &a, nil
}

func (r *AgentRepository) UpdateActivity(authKey string, ip string) error {
	_, err := r.DB.Exec(`
		UPDATE agents
		SET last_activity_at = $1, last_ip_address = $2
		WHERE auth_key = $3
	`, time.Now(), ip, authKey)
	return err
}

func (r *AgentRepository) IncrementDataProcessed(authKey string, count int) error {
	_, err := r.DB.Exec(`
		UPDATE agents
		SET data_processed = data_processed + $1
		WHERE auth_key = $2
	`, count, authKey)
	return err
}
