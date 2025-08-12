package entity

import (
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

type CombolistMetadataEntity struct {
	ID          uuid.UUID      `db:"uuid"`
	Source      string         `db:"source"`
	CollectedAt time.Time      `db:"collected_at"`
	Tags        pq.StringArray `db:"tags"`
	Notes       string         `db:"notes"`
	Hash        string         `db:"hash"`
}
