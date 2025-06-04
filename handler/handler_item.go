package handler

import (
	"bufio"
	"fmt"
	"manufactures/config"
	"os"
	"strconv"
	"strings"
)

type ItemReport struct {
	ItemID int     `json:"item_id"`
	Name   string  `json:"name"`
	Stock  int     `json:"stock"`
	Price  float64 `json:"price"`
}

func InsertProduct() {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Enter product name: ")
	name, _ := reader.ReadString('\n')
	name = strings.TrimSpace(name)

	fmt.Print("Enter price: ")
	priceStr, _ := reader.ReadString('\n')
	priceStr = strings.TrimSpace(priceStr)

	fmt.Print("Enter stock: ")
	stockStr, _ := reader.ReadString('\n')
	stockStr = strings.TrimSpace(stockStr)

	if name == "" || priceStr == "" || stockStr == "" {
		fmt.Println("All fields are required.")
		return
	}

	price, err := strconv.ParseFloat(priceStr, 64)
	if err != nil {
		fmt.Println("Invalid price format. Please enter a valid decimal number.")
		return
	}

	stock, err := strconv.Atoi(stockStr)
	if err != nil {
		fmt.Println("Invalid stock format. Please enter a valid integer number.")
		return
	}

	_, err = config.InitDB().Exec(`
		INSERT INTO item (name, price, stock)
		VALUES (?, ?, ?)`,
		name, price, stock,
	)
	if err != nil {
		fmt.Println("Error inserting product:", err)
		return
	}

	fmt.Println("✅ Product inserted successfully.")
}

func PrintProduct() {
	rows, err := config.InitDB().Query(`
		SELECT i.item_id, i.name, c.name AS category, i.price, i.stock
		FROM item i
		LEFT JOIN category c ON i.category_id = c.category_id
	`)
	if err != nil {
		fmt.Println("Error retrieving products:", err)
		return
	}
	defer rows.Close()

	fmt.Println("List of Products:")
	for rows.Next() {
		var id int
		var name, category string
		var price float64
		var stock int
		err := rows.Scan(&id, &name, &category, &price, &stock)
		if err != nil {
			fmt.Println("Scan error:", err)
			return
		}
		fmt.Printf("ID: %d | Name: %s | Category: %s | Price: %.2f | Stock: %d\n",
			id, name, category, price, stock)
	}
}
