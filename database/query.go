package database

import (
	"fmt"
	"net/url"

	"github.com/TF2Stadium/Helen/models"
	"github.com/TF2Stadium/twitchbot/config"
)

func GetCurrentLobby(channel string) string {
	var playerid, lobbyid uint
	db.QueryRow("SELECT id FROM players WHERE twitch_name = $1", channel).Scan(&playerid)
	err := db.QueryRow("SELECT lobby_slots.lobby_id FROM lobbies INNER JOIN lobby_slots ON lobbies.id = lobby_slots.lobby_id WHERE lobbies.state = $1 AND lobby_slots.player_id = $2", models.LobbyStateWaiting, playerid).Scan(&lobbyid)
	if err != nil {
		return ""
	}

	lobbyURL, _ := url.Parse(config.Constants.FrontendURL)
	lobbyURL.Path = fmt.Sprintf("lobby/%d", lobbyid)
	return lobbyURL.String()
}
