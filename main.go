package main

import (
	"fmt"
	"manufactures/config"
)

func main() {
	db := config.InitDB()
	config.RunMigration(db)

	fmt.Println("Exiting application...")
}
