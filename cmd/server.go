package main

import (
	"log"
	"shifty-backend/configs"
	"shifty-backend/pkg/database"
)

func main() {
	cfg, err := configs.LoadConfig()
	if err != nil {
		log.Fatal("Can not load config: ", err)
	}

	db := database.ConnectPostgres(cfg)
	if err:= db.AutoMigrate(),
	
}
