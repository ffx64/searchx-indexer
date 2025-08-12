package security

import (
	"errors"
	"fmt"
	"searchx-indexer/entity"
	"searchx-indexer/repository"
)

func AuthAgent(repo *repository.AgentRepository, authkey string, platform string) (*entity.AgentEntity, error) {
	if authkey == "" {
		return nil, errors.New("missing authorization header")
	}

	agent, err := repo.FindByAuthKey(authkey)
	if err != nil {
		return nil, fmt.Errorf("agent not found: %w", err)
	}

	if agent.AgentStatus != "active" {
		return nil, errors.New("agent is not active")
	}

	if agent.Platform != platform {
		repo.UpdateStatus(authkey, "compromised")
		return nil, errors.New("agent platform mismatch")
	}

	return agent, nil
}
