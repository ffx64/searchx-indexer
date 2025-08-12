package service

import (
	"log"
	"searchx-indexer/entity"
	"searchx-indexer/repository"
	"time"

	"github.com/google/uuid"
)

type CombolistService struct {
	CombolistRepository *repository.ComboListRepository
	AgentRepository     *repository.AgentRepository
}

func NewComboListService(combo *repository.ComboListRepository, agent *repository.AgentRepository) *CombolistService {
	return &CombolistService{CombolistRepository: combo, AgentRepository: agent}
}

func (s *CombolistService) BulkInsert(hash string, meta entity.CombolistMetadataEntity, data []entity.CombolistDataEntity) error {
	existing, err := s.CombolistRepository.FindMetadataByHash(hash)
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
		metadataID, err = s.CombolistRepository.InsertMetadata(&meta)
		if err != nil {
			return err
		}
	}

	for i, d := range data {
		d.MetadataID = metadataID

		err = s.CombolistRepository.InsertData(d)

		if err != nil {
			return err
		}

		log.Println("[", i, "/", len(data), "]", "combolist data insert:", d.Username, d.Password)
	}

	return nil
}
