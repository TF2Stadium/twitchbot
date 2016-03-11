package database

import (
	"fmt"
	"github.com/TF2Stadium/Helen/models"
	"github.com/TF2Stadium/twitchbot/config"
)

func GetCurrentLobby(channel string) string {
	var playerid, lobbyid uint
	db.QueryRow("SELECT id FROM players WHERE twitch_name = $1", channel).Scan(&playerid)
	err := db.QueryRow("SELECT lobby_slots.lobby_id FROM lobbies INNER JOIN lobbies.id = lobby_slots.lobby_id WHERE lobbies.state = $1", models.LobbyStateWaiting).Scan(&lobbyid)
	if err != nil {
		return ""
	}

	return fmt.Sprintf("%s/lobby/%d", config.Constants.FrontendURL, lobbyid)
}
