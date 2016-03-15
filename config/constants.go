package config

import (
	"log"
	"net/url"

	"github.com/kelseyhightower/envconfig"
)

type config struct {
	TwitchNick     string `envconfig:"NICK" default:"TF2Stadium"`
	TwitchPassword string `envconfig:"PASSWORD"`

	RPCQueue    string `envconfig:"RPC_QUEUE" default:"twitchbot"`
	RabbitMQURL string `envconfig:"RABBITMQ_URL" default:"amqp://guest:guest@localhost:5672/"`

	DBAddr     string `envconfig:"DATABASE_ADDR" default:"127.0.0.1:5432"`
	DBDatabase string `envconfig:"DATABASE_NAME" default:"tf2stadium"`
	DBUsername string `envconfig:"DATABASE_USERNAME" default:"tf2stadium"`
	DBPassword string `envconfig:"DATABASE_PASSWORD" default:"dickbutt"`

	FrontendURL string `envconfig:"FRONTEND_URL"`
}

var Constants config

func SetupConstants() {
	envconfig.MustProcess("TWITCHBOT", &Constants)
	u, err := url.Parse(Constants.FrontendURL)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Frontend URL: %s", u.String())
}
