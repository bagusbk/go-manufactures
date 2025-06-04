package main

import (
	"fmt"
	"manufactures/cli"
	"manufactures/config"
)

func main() {
	db := config.InitDB()
	config.RunMigration(db)

	// manufacturesHandler := &handler.ManufacturesHandler{DB: db}

	cli.ShowMenu()

	fmt.Println("Exiting application...")
}
