package main

import (
	"log"

	"go-challenge/config"
	db "go-challenge/database"
)

func main() {
	config.Init()
	db.Init()
	defer db.Close()

	// CSVからデータをシード
	if err := db.SeedFromCSV(db.GetDB(), "sample/locations.csv", "sample/evses.csv"); err != nil {
		log.Printf("Warning: Failed to seed data: %v", err)
	}

	serverInit()
}
