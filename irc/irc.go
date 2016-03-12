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
	conn     *client.Conn
	mapMu    = new(sync.RWMutex)
	lastSaid = make(map[string]time.Time)
)

func Connect() {
	conn = client.SimpleClient(config.Constants.TwitchNick)
	if err := conn.ConnectTo("irc.twitch.tv", config.Constants.TwitchPassword); err != nil {
		log.Fatal(err)
	}

	disconnected := make(chan struct{})
	conn.HandleFunc("disconnected", func(c *client.Conn, l *client.Line) {
		close(disconnected)
	})

	conn.HandleFunc("PRIVMSG", func(conn *client.Conn, line *client.Line) {
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
					conn.Privmsg(line.Target(), line.Target()[1:]+" isn't in any lobby right now")
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
	if !hasJoined(channel) {
		return
	}

	mapMu.RLock()
	last, ok := lastSaid[channel]
	mapMu.RUnlock()

	if ok {
		duration := time.Since(last)
		if duration > 10*time.Second {
			time.Sleep(duration - 10*time.Second)
		}
	}

	mapMu.Lock()
	lastSaid[channel] = time.Now()
	mapMu.Unlock()

	conn.Privmsg("#"+channel, text)
}

var (
	joinedMu       = new(sync.RWMutex)
	joinedChannels []string
)

func hasJoined(channel string) bool {
	joinedMu.RLock()
	defer joinedMu.RUnlock()
	for _, c := range joinedChannels {
		if c == channel {
			return true
		}
	}

	return false
}

func Join(channel string) {
	if hasJoined(channel) {
		return
	}

	conn.Join("#" + channel)

	joinedMu.Lock()
	joinedChannels = append(joinedChannels, channel)
	joinedMu.Unlock()
}

func Leave(channel string) {
	if !hasJoined(channel) {
		return
	}

	conn.Part("#" + channel)

	joinedMu.Lock()
	for i, c := range joinedChannels {
		if c == channel {
			joinedChannels = append(joinedChannels[:i], joinedChannels[i+1:]...)
		}
	}
	joinedMu.Unlock()
}
