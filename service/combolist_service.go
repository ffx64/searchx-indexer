package service

import (
	"log"
	"searchx-indexer/entity"
	"searchx-indexer/repository"
	"time"

	"github.com/google/uuid"
)

type ComboListService struct {
	ComboRepo *repository.ComboListRepository
	AgentRepo *repository.AgentRepository
}

func NewComboListService(combo *repository.ComboListRepository, agent *repository.AgentRepository) *ComboListService {
	return &ComboListService{ComboRepo: combo, AgentRepo: agent}
}

func (s *ComboListService) BulkInsert(hash string, meta entity.ComboListMetadataEntity, data []entity.ComboListDataEntity) error {
	existing, err := s.ComboRepo.FindMetadataByHash(hash)
	if err != nil {
		return err
	}

	var metadataID uuid.UUID
	if existing != nil {
		metadataID = existing.ID
	} else {
		if meta.CollectedAt.IsZero() {
			meta.CollectedAt = time.Now()
		}
		meta.Hash = hash
		metadataID, err = s.ComboRepo.InsertMetadata(&meta)
		if err != nil {
			return err
		}
	}

	for i, d := range data {
		d.MetadataID = metadataID

		err = s.ComboRepo.InsertData(d)

		if err != nil {
			return err
		}

		log.Println("[", i, "/", len(data), "]", "combolist data insert:", d.Username, d.Password)
	}

	return nil
}
