package database

import (
	"database/sql"

	"github.com/TF2Stadium/twitchbot/config"
	_ "github.com/lib/pq"
	"log"
	"net/url"
)

var db *sql.DB

func Connect() {
	dburl := url.URL{
		Scheme:   "postgres",
		Host:     config.Constants.DBAddr,
		Path:     config.Constants.DBDatabase,
		RawQuery: "sslmode=disable",
	}

	log.Printf("Connecting to DB on %s", dburl.String())
	dburl.User = url.UserPassword(config.Constants.DBUsername, config.Constants.DBPassword)
	var err error
	db, err = sql.Open("postgres", dburl.String())
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Connected")
}
