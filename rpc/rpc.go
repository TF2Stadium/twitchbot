package rpc

import (
	"fmt"
	"log"
	"net/rpc"
	"net/url"

	"github.com/TF2Stadium/twitchbot/config"
	"github.com/TF2Stadium/twitchbot/irc"
	"github.com/streadway/amqp"
	"github.com/vibhavp/amqp-rpc"
)

type TwitchBot struct{}

func StartRPC() {
	conn, err := amqp.Dial(config.Constants.RabbitMQURL)
	if err != nil {
		log.Fatal(err)
	}

	serverCodec, err := amqprpc.NewServerCodec(conn, config.Constants.RPCQueue, amqprpc.JSONCodec{})
	if err != nil {
		log.Fatal(err)
	}

	rpc.Register(new(TwitchBot))
	rpc.ServeCodec(serverCodec)
}

func (TwitchBot) Join(channel string, _ *struct{}) error {
	irc.Join(channel)
	return nil
}

func (TwitchBot) Leave(channel string, _ *struct{}) error {
	irc.Leave(channel)
	return nil
}

func (TwitchBot) Announce(action struct {
	Channel string
	LobbyID uint
}, _ *struct{}) error {

	lobbyURL, _ := url.Parse(config.Constants.FrontendURL)
	lobbyURL.Path = fmt.Sprintf("lobby/%d", action.LobbyID)
	irc.Say(action.Channel+" just joined "+lobbyURL.String(), action.Channel)
	return nil
}
