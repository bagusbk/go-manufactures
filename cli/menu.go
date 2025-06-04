package cli

import (
	"bufio"
	"fmt"
	"go-manufactures/handler"
	"os"
	"strings"
)

func ShowMenu() {
	reader := bufio.NewReader(os.Stdin)
	for {
		if handler.LoggedInStaff.Email == "" {
			fmt.Println("============================")
			fmt.Println("Welcome to ManufactureStore CLI")
			fmt.Println("============================")
			fmt.Println("1. Login")
			fmt.Println("2. Exit")
			fmt.Print("Choose an option: ")

			input, _ := reader.ReadString('\n')
			input = strings.TrimSpace(input)

			if input == "1" {
				email := handler.LoginUser()
				if email != "" {
					handler.LoggedInStaff.Email = email
				}
			} else if input == "2" {
				fmt.Println("Goodbye!")
				return
			} else {
				fmt.Println("Invalid input")
			}
		} else {
			fmt.Println("============================")
			fmt.Printf("Welcome, %s (Role: %s)\n", handler.LoggedInStaff.Email, handler.LoggedInStaff.Position)
			fmt.Println("============================")
			fmt.Println("1. Manage Products")
			fmt.Println("2. Manage User")
			fmt.Println("3. Manage Staff")
			fmt.Println("4. Order")
			fmt.Println("5. Reports")
			fmt.Println("6. Logout")
			fmt.Print("Choose an option: ")

			input, _ := reader.ReadString('\n')
			input = strings.TrimSpace(input)

			switch input {
			case "1":
				showProductMenu()
			case "2":
				showUserMenu()
			case "3":
				showStaffMenu()
			case "4":
				showOrderMenu()
			case "5":
				showReportMenu()
			case "6":
				handler.LoggedInStaff.Email = ""
				handler.LoggedInStaff.Position = ""
				fmt.Println("Logged out.")
			default:
				fmt.Println("Invalid input")
			}
		}
	}
}
