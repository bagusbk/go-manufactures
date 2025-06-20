package handler

import (
	"bufio"
	"fmt"
	"log"
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
	rows, err := config.InitDB().Query("SELECT user_id, full_name, email, address, phone_number FROM users")
	if err != nil {
		fmt.Println("Error retrieving users:", err)
		return
	}
	defer rows.Close()

	fmt.Println("List of Users:")
	for rows.Next() {
		var id int
		var name, email, address, phoneNumber string
		if err := rows.Scan(&id, &name, &email, &address, &phoneNumber); err != nil {
			fmt.Println("Scan error:", err)
			return
		}
		fmt.Printf("ID: %d | Name: %s | Email: %s | Address: %s | Phone Number: %s\n",
			id, name, email, address, phoneNumber)
	}
}

func PrintFrequentBuyersReport() {
	rows, err := config.InitDB().Query(`
        SELECT u.user_id, u.full_name, COUNT(o.order_id) AS order_count, 
               i.name AS item_name, SUM(oi.quantity) AS total_quantity
        FROM users u
        JOIN orders o ON u.user_id = o.user_id
        JOIN order_items oi ON o.order_id = oi.order_id
        JOIN item i ON oi.item_id = i.item_id
        JOIN payment p ON o.order_id = p.order_id
        WHERE p.status = 'paid' AND o.status = 'completed'
        GROUP BY u.user_id, i.item_id
        ORDER BY order_count DESC;
    `)
	if err != nil {
		log.Fatal("Error retrieving frequent buyers:", err)
	}
	defer rows.Close()

	fmt.Println("Frequent Buyers Report:")
	var previousUserID int
	var previousFullName string
	var orderCount int
	var itemName string
	var totalQuantity int

	for rows.Next() {
		var userID int
		err := rows.Scan(&userID, &previousFullName, &orderCount, &itemName, &totalQuantity)
		if err != nil {
			log.Fatal("Error scanning row:", err)
		}

		// Print the user information for the first item or when user changes
		if userID != previousUserID {
			if previousUserID != 0 {
				fmt.Println() // Empty line between different users
			}
			fmt.Printf("User ID: %d | Name: %s | Order Count: %d\n", userID, previousFullName, orderCount)
		}

		// Print the item name and total quantity for the user
		fmt.Printf("   - Item: %s | Quantity: %d\n", itemName, totalQuantity)

		// Update previousUserID to current userID for the next iteration
		previousUserID = userID
	}
}
