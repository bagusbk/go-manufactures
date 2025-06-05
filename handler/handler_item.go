package handler

import (
	"bufio"
	"fmt"
	"log"
	"manufactures/config"
	"manufactures/entity"
	"os"
	"strconv"
	"strings"
)

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
		SELECT item_id, name, stock, price FROM item
	`)
	if err != nil {
		fmt.Println("Error retrieving items:", err)
		return
	}
	defer rows.Close()

	var listProduct []entity.Item
	for rows.Next() {
		var report entity.Item
		err := rows.Scan(&report.ItemID, &report.Name, &report.Stock, &report.Price)
		if err != nil {
			fmt.Println("Error scanning row:", err)
			return
		}
		listProduct = append(listProduct, report)
	}

	fmt.Println("List of Products:")
	for _, report := range listProduct {
		fmt.Printf("ProductID: %d | Product Name: %s | Stock: %d | Price: %.2f\n",
			report.ItemID, report.Name, report.Stock, report.Price)
	}
}

func PrintMostSoldItemsReport() {
	rows, err := config.InitDB().Query(`
        SELECT i.item_id, i.name, SUM(oi.quantity) AS total_sold
        FROM item i
        JOIN order_items oi ON i.item_id = oi.item_id
        JOIN orders o ON oi.order_id = o.order_id
        JOIN payment p ON o.order_id = p.order_id
        WHERE p.status = 'paid'  -- Hanya hitung pesanan yang sudah dibayar
        GROUP BY i.item_id
        ORDER BY total_sold DESC;
    `)
	if err != nil {
		log.Fatal("Error retrieving most sold items:", err)
	}
	defer rows.Close()

	fmt.Println("Most Sold Items Report:")
	for rows.Next() {
		var itemID int
		var name string
		var totalSold int
		err := rows.Scan(&itemID, &name, &totalSold)
		if err != nil {
			log.Fatal("Error scanning row:", err)
		}
		fmt.Printf("Item ID: %d | Name: %s | Total Sold: %d\n", itemID, name, totalSold)
	}
}
