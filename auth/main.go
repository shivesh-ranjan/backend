package main

import (
	"context"
	"database/sql"
	"log"

	"shivesh-ranjan.github.io/m/api"
	db "shivesh-ranjan.github.io/m/db/sqlc"
	"shivesh-ranjan.github.io/m/utils"
)

func main() {
	config, err := utils.LoadConfig(".")
	if err != nil {
		log.Fatal("Can't load config: ", err)
	}

	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("Can't connect to Database: ", err)
	}

	args := db.CreateUserParams{
		Username: "shaw",
		Name:     "Derek Shaw",
		About:    "Devops Engineer | Calisthenics Athlete",
		Photo:    "https://www.behindthevoiceactors.com/_img/chars/thumbs/agent-47-hitman-78.7_thumb.jpg",
		Role:     "admin",
	}
	args.Password, err = utils.HashPassword(config.AdminPassword)
	if err != nil {
		log.Print("While hashing Admin Password: ", err)
	}
	store := db.NewStore(conn)
	user, err := db.Store.CreateUser(store, context.Background(), args)
	if err != nil {
		log.Print("This happened during Admin Creation: ", err)
	} else {
		log.Print("Admin User Added Successfully: ", user.Username)
	}

	server, err := api.NewServer(config, store)
	if err != nil {
		log.Fatal("cannot create server:", err)
	}
	err = server.Start(config.ServerAddress)
	if err != nil {
		log.Fatal("cannot start server:", err)
	}
}
