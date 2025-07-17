package entity

import (
	"time"

	"github.com/google/uuid"
)

type AgentEntity struct {
	ID                 uuid.UUID  `db:"uuid"`
	AuthKey            string     `db:"auth_key"`
	Platform           string     `db:"platform"`
	CollectionInterval int        `db:"collection_interval"`
	LastActivityAt     *time.Time `db:"last_activity_at"`
	DataProcessed      int        `db:"data_processed"`
	AgentStatus        string     `db:"agent_status"`
	LastIPAddress      *string    `db:"last_ip_address"`
}
