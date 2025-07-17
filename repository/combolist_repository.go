package repository

import (
	"database/sql"
	"errors"
	"searchx-indexer/entity"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

type ComboListRepository struct {
	DB *sql.DB
}

func NewComboListRepository(db *sql.DB) *ComboListRepository {
	return &ComboListRepository{DB: db}
}

func (r *ComboListRepository) FindMetadataByHash(hash string) (*entity.ComboListMetadataEntity, error) {
	var m entity.ComboListMetadataEntity

	err := r.DB.QueryRow(`
		SELECT id, source, collected_at, tags, notes, hash FROM combolist_metadata WHERE hash = $1
	`, hash).Scan(&m.ID, &m.Source, &m.CollectedAt, &m.Tags, &m.Notes, &m.Hash)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}

	return &m, err
}

func (r *ComboListRepository) InsertMetadata(meta *entity.ComboListMetadataEntity) (uuid.UUID, error) {
	tx, err := r.DB.Begin()
	if err != nil {
		return uuid.UUID{}, err
	}

	var id uuid.UUID
	err = tx.QueryRow(`
		INSERT INTO combolist_metadata (source, collected_at, tags, notes, hash) VALUES ($1, $2, $3, $4, $5) ON CONFLICT (hash) DO NOTHING RETURNING id
	`, meta.Source, meta.CollectedAt, pq.StringArray(meta.Tags), meta.Notes, meta.Hash).Scan(&id)

	if err == sql.ErrNoRows {
		err = tx.QueryRow(`SELECT id FROM combolist_metadata WHERE hash = $1`, meta.Hash).Scan(&id)

		if err != nil {
			tx.Rollback()
			return uuid.UUID{}, err
		}
	} else if err != nil {
		tx.Rollback()
		return uuid.UUID{}, err
	}

	if err := tx.Commit(); err != nil {
		return uuid.UUID{}, err
	}

	return id, nil
}

func (r *ComboListRepository) InsertData(data entity.ComboListDataEntity) error {
	tx, err := r.DB.Begin()

	if err != nil {
		return err
	}

	stmt, err := tx.Prepare(`
		INSERT INTO combolist_data (email, password, username, domain, metadata_id)
		VALUES ($1, $2, $3, $4, $5) ON CONFLICT DO NOTHING
	`)

	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.Exec(data.Email, data.Password, data.Username, data.Domain, data.MetadataID)

	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}
