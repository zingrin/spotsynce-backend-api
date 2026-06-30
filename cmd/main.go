package main

import (
	"spot-sync/internal/config"
	"spot-sync/internal/server"
)

func main() {
	env := config.LoadEnv()

	db := config.ConnectToDB(env)

	server.Start(env, db)
}
