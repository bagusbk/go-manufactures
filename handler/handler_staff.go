package handler

import (
	"bufio"
	"database/sql"
	"fmt"
	"os"
	"strings"

	"go-manufactures/config"

	_ "github.com/go-sql-driver/mysql"
)

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

	if passwordInput != passwordHash {
		fmt.Println("Incorrect password.")
		return ""
	}

	LoggedInStaff.Email = email
	LoggedInStaff.Position = position
	fmt.Printf("Login successful! Role: %s\n", position)
	return email
}
