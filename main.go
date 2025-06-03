package main

import (
	"fmt"
	"manufactures/config"
	"manufactures/handler"
)

func main() {
	db := config.InitDB()
	config.RunMigration(db)

	manufacturesHandler := &handler.ManufacturesHandler{DB: db}
	fmt.Println("Exiting application...")
}
