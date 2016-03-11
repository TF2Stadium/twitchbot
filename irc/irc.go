package irc

import (
	"log"

	"github.com/TF2Stadium/twitchbot/config"
	"github.com/TF2Stadium/twitchbot/database"
	"github.com/fluffle/goirc/client"
	"sync"
	"time"
)

var (
	Conn     *client.Conn
	mapMu    = new(sync.RWMutex)
	lastSaid = make(map[string]time.Time)
)

func Connect() {
	config.SetupConstants()

	Conn = client.SimpleClient(config.Constants.TwitchNick)
	if err := Conn.ConnectTo("irc.twitch.tv", config.Constants.TwitchPassword); err != nil {
		log.Fatal(err)
	}

	disconnected := make(chan struct{})
	Conn.HandleFunc("disconnected", func(c *client.Conn, l *client.Line) {
		close(disconnected)
	})

	Conn.HandleFunc("PRIVMSG", func(conn *client.Conn, line *client.Line) {
		if !line.Public() {
			return
		}

		mapMu.RLock()
		if time.Since(lastSaid[line.Target()]) < 10*time.Second {
			mapMu.RUnlock()
			return
		}
		mapMu.RUnlock()

		mapMu.Lock()
		lastSaid[line.Target()] = time.Now()
		mapMu.Unlock()

		if line.Text() == "!lobby" {
			if line.Public() && line.Target()[0] == '#' {
				lobbyURL := database.GetCurrentLobby(line.Target()[1:])
				if lobbyURL == "" {
					return
				}
				conn.Privmsg(line.Target(), lobbyURL)
			}
		}
	})

	log.Println("Connected to Twitch IRC")
	go func() {
		<-disconnected
		log.Fatal("Disconnected")
	}()
}

func Say(text, channel string) {
	mapMu.RLock()
	duration := time.Since(lastSaid[channel])
	mapMu.RUnlock()

	time.Sleep(duration)

	mapMu.Lock()
	lastSaid[channel] = time.Now()
	mapMu.Unlock()

	Conn.Privmsg(channel, text)
}
