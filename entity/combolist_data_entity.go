package entity

import "github.com/google/uuid"

type CombolistDataEntity struct {
	ID         uuid.UUID `db:"uuid"`
	Email      string    `db:"email"`
	Password   string    `db:"password"`
	Username   string    `db:"username"`
	Domain     string    `db:"domain"`
	MetadataID uuid.UUID `db:"metadata_id"`
}
