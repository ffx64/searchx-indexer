package handlers

import (
	"encoding/json"
	"log"
	"searchx-indexer/database"
	"searchx-indexer/repository/combolist"
)

func ProcessMessage(message []byte) {
	var baseMap map[string]interface{}
	if err := json.Unmarshal(message, &baseMap); err != nil {
		log.Println("[!] Error parsing JSON:", err)
		return
	}

	agent, ok := baseMap["agent"].(string)
	if !ok {
		log.Println("[!] Invalid or unknown agent")
		return
	}

	switch agent {
	case "combolist-agent":
		if db, err := database.GetManager().GetDB("combolist-db"); err != nil {
			log.Println("[!] Error GetManager() combolist-db")
		} else {
			HandleCombolistMessage(baseMap, combolist.NewCombolistRepository(db))
		}
	default:
		log.Println("[!] Unknown agent:", agent)
	}
}
