package handler

import (
	"bufio"
	"fmt"
	"go-manufactures/config"
	"os"
	"strings"
)

func InsertProduct() {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Enter product name: ")
	name, _ := reader.ReadString('\n')
	name = strings.TrimSpace(name)

	fmt.Print("Enter price: ")
	price, _ := reader.ReadString('\n')
	price = strings.TrimSpace(price)

	fmt.Print("Enter stock: ")
	stock, _ := reader.ReadString('\n')
	stock = strings.TrimSpace(stock)

	if name == "" || price == "" || stock == "" {
		fmt.Println("All fields are required.")
		return
	}

	_, err := config.DB.Exec(`
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
	rows, err := config.DB.Query(`
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
