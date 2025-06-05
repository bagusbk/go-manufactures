package cli

import (
	"bufio"
	"fmt"
	"manufactures/handler"
	"os"
	"strings"
)

func showStaffMenu() {
	if handler.LoggedInStaff.Position != "admin" {
		fmt.Println("Unauthorized: Only admin can manage staff.")
		return
	}

	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Println("\n-- Manage Staff --")
		fmt.Println("1. Insert Staff")
		fmt.Println("2. View Staff")
		fmt.Println("3. Update Staff Role")
		fmt.Println("4. Delete Staff")
		fmt.Println("5. Back to Main Menu")
		fmt.Print("Choose an option: ")

		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		switch input {
		case "1":
			handler.InsertStaff()
		case "2":
			handler.PrintAllStaff()
		case "3":
			handler.UpdateStaffRole()
		case "4":
			handler.DeleteStaff()
		case "5":
			return
		default:
			fmt.Println("Invalid input. Please choose between 1–5.")
		}
	}
}
