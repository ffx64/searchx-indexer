package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"searchx-indexer/models/discord"
	"searchx-indexer/repository/discord"
)

func HandleDiscordMessage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var message discord.Message
	if err := json.NewDecoder(r.Body).Decode(&message); err != nil {
		http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
		return
	}

	log.Printf("[+] Received Discord message: %+v\n", message)

	// Save the message to the database
	if err := repository.SaveDiscordMessage(message); err != nil {
		log.Printf("[!] Failed to save Discord message: %s\n", err)
		http.Error(w, "Failed to save message", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Message received and saved"))
}
