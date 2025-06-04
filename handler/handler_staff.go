package handler

import (
	"bufio"
	"database/sql"
	"fmt"
	"os"
	"strings"

	"go-manufactures/config"

	_ "github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
)

func hashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

var LoggedInStaff struct {
	Email    string
	Position string
}

func LoginUser() string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter email: ")
	emailInput, _ := reader.ReadString('\n')
	emailInput = strings.TrimSpace(emailInput)

	fmt.Print("Enter password: ")
	passwordInput, _ := reader.ReadString('\n')
	passwordInput = strings.TrimSpace(passwordInput)

	var email, passwordHash, position string
	err := config.DB.QueryRow("SELECT email, password_hash, position FROM staff WHERE email = ?", emailInput).Scan(&email, &passwordHash, &position)
	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Println("Email not found.")
		} else {
			fmt.Println("Error:", err)
		}
		return ""
	}

	if position == "admin" {
		if passwordInput != passwordHash {
			fmt.Println("Incorrect password.")
			return ""
		}
	} else {
		err = bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(passwordInput))
		if err != nil {
			fmt.Println("Incorrect password.")
			return ""
		}
	}

	// if passwordInput != passwordHash {
	// 	fmt.Println("Incorrect password.")
	// 	return ""
	// }

	LoggedInStaff.Email = email
	LoggedInStaff.Position = position
	fmt.Printf("Login successful! Role: %s\n", position)
	return email
}

func InsertStaff() {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Enter full name: ")
	fullName, _ := reader.ReadString('\n')
	fullName = strings.TrimSpace(fullName)

	fmt.Print("Enter position (admin/manager/staff): ")
	position, _ := reader.ReadString('\n')
	position = strings.TrimSpace(position)

	fmt.Print("Enter email: ")
	email, _ := reader.ReadString('\n')
	email = strings.TrimSpace(email)

	// Check for duplicate email
	var exists int
	checkErr := config.DB.QueryRow("SELECT COUNT(*) FROM staff WHERE email = ?", email).Scan(&exists)
	if checkErr != nil {
		fmt.Println("Error checking email:", checkErr)
		return
	}
	if exists > 0 {
		fmt.Println("⚠️ Email already registered.")
		return
	}

	fmt.Print("Enter password: ")
	password, _ := reader.ReadString('\n')
	password = strings.TrimSpace(password)

	hashedPassword, err := hashPassword(password)
	if err != nil {
		fmt.Println("Error hashing password:", err)
		return
	}

	if fullName == "" || email == "" || password == "" {
		fmt.Println("All fields are required.")
		return
	}

	_, err = config.DB.Exec(`
		INSERT INTO staff (full_name, position, email, password_hash)
		VALUES (?, ?, ?, ?)`,
		fullName, position, email, hashedPassword,
	)
	if err != nil {
		fmt.Println("Error inserting staff:", err)
		return
	}

	fmt.Println("✅ Staff successfully inserted.")
}

func PrintAllStaff() {
	rows, err := config.DB.Query("SELECT staff_id, full_name, position, email, created_at FROM staff")
	if err != nil {
		fmt.Println("Error retrieving staff:", err)
		return
	}
	defer rows.Close()

	fmt.Println("List of Staff:")
	for rows.Next() {
		var id int
		var name, position, email string
		var createdAt string
		if err := rows.Scan(&id, &name, &position, &email, &createdAt); err != nil {
			fmt.Println("Scan error:", err)
			return
		}
		fmt.Printf("ID: %d | Name: %s | Position: %s | Email: %s | Created At: %s\n", id, name, position, email, createdAt)
	}
}

func DeleteStaff() {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Enter staff ID to delete: ")
	idStr, _ := reader.ReadString('\n')
	idStr = strings.TrimSpace(idStr)

	res, err := config.DB.Exec("DELETE FROM staff WHERE staff_id = ?", idStr)
	if err != nil {
		fmt.Println("Error deleting staff:", err)
		return
	}

	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		fmt.Println("⚠️ No staff found with that ID.")
		return
	}

	fmt.Println("✅ Staff successfully deleted.")
}
