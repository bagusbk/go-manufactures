package handler

import (
	"bufio"
	"fmt"
	"manufactures/config"
	"os"
	"regexp"
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
	checkErr := config.InitDB().QueryRow("SELECT COUNT(*) FROM users WHERE email = ?", email).Scan(&exists)
	if checkErr != nil {
		fmt.Println("Error checking email:", checkErr)
		return
	}
	if exists > 0 {
		fmt.Println("⚠️ Email already registered.")
		return
	}

	fmt.Print("Enter phone number: ")
	phoneStr, _ := reader.ReadString('\n')
	phoneStr = strings.TrimSpace(phoneStr)

	fmt.Print("Enter address: ")
	address, _ := reader.ReadString('\n')
	address = strings.TrimSpace(address)

	if fullName == "" || email == "" || address == "" || phoneStr == "" {
		fmt.Println("All fields are required.")
		return
	}
	phoneRegex := `^\+?[0-9]{10,15}$` // regex for phone numbers, optional + and 10-15 digits
	matched, err := regexp.MatchString(phoneRegex, phoneStr)
	if err != nil {
		fmt.Println("Error validating phone number format:", err)
		return
	}
	if !matched {
		fmt.Println("Invalid phone number format. Please enter a valid phone number (e.g., +1234567890).")
		return
	}

	_, err = config.InitDB().Exec(`
		INSERT INTO users (full_name, email, address, phone_number)
		VALUES (?, ?, ?, ?)`,
		fullName, email, address,
	)
	if err != nil {
		fmt.Println("Error inserting user:", err)
		return
	}

	fmt.Println("✅ User inserted successfully.")
}

func PrintUser() {
	rows, err := config.InitDB().Query("SELECT user_id, full_name, email, address, created_at FROM user")
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
