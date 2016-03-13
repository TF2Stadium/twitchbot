package main

import (
	"github.com/TF2Stadium/twitchbot/config"
	"github.com/TF2Stadium/twitchbot/database"
	"github.com/TF2Stadium/twitchbot/irc"
	"github.com/TF2Stadium/twitchbot/rpc"
	"log"
)

func main() {
	log.SetFlags(log.Lshortfile)
	config.SetupConstants()

	database.Connect()
	irc.Connect()
	rpc.StartRPC()
}
