package repository

import (
	"errors"
	"searchx-indexer/database"
	"searchx-indexer/models/discord"
)

func SaveDiscordMessage(message discord.Message) error {
	// Filtrar mensagens de bots
	if message.IsBot {
		return errors.New("message is from a bot, ignoring")
	}

	db, err := database.GetManager().GetDB("discord-db")
	if err != nil {
		return err
	}

	query := `
        INSERT INTO discord_messages (
            id, content, timestamp, edited_at, user_id, username, avatar,
            channel_id, channel_name, guild_id, guild_name, is_pinned, is_tts
        ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
    `

	_, err = db.Exec(query,
		message.ID, message.Content, message.Timestamp, message.EditedAt,
		message.UserID, message.Username, message.Avatar,
		message.ChannelID, message.ChannelName, message.GuildID, message.GuildName,
		message.IsPinned, message.IsTTS,
	)

	return err
}
