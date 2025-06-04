package handler

import (
	"bufio"
	"fmt"
	"go-manufactures/config"
	"os"
	"strings"
)

func InsertUser() {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Enter full name: ")
	fullName, _ := reader.ReadString('\n')
	fullName = strings.TrimSpace(fullName)

	fmt.Print("Enter email: ")
	email, _ := reader.ReadString('\n')
	email = strings.TrimSpace(email)

	// Check for duplicate email
	var exists int
	checkErr := config.DB.QueryRow("SELECT COUNT(*) FROM user WHERE email = ?", email).Scan(&exists)
	if checkErr != nil {
		fmt.Println("Error checking email:", checkErr)
		return
	}
	if exists > 0 {
		fmt.Println("⚠️ Email already registered.")
		return
	}

	fmt.Print("Enter address: ")
	address, _ := reader.ReadString('\n')
	address = strings.TrimSpace(address)

	if fullName == "" || email == "" || address == "" {
		fmt.Println("All fields are required.")
		return
	}

	_, err := config.DB.Exec(`
		INSERT INTO user (full_name, email, address)
		VALUES (?, ?, ?)`,
		fullName, email, address,
	)
	if err != nil {
		fmt.Println("Error inserting user:", err)
		return
	}

	fmt.Println("✅ User inserted successfully.")
}

func PrintUser() {
	rows, err := config.DB.Query("SELECT user_id, full_name, email, address, created_at FROM user")
	if err != nil {
		fmt.Println("Error retrieving users:", err)
		return
	}
	defer rows.Close()

	fmt.Println("List of Users:")
	for rows.Next() {
		var id int
		var name, email, address, createdAt string
		if err := rows.Scan(&id, &name, &email, &address, &createdAt); err != nil {
			fmt.Println("Scan error:", err)
			return
		}
		fmt.Printf("ID: %d | Name: %s | Email: %s | Address: %s | Created At: %s\n",
			id, name, email, address, createdAt)
	}
}
