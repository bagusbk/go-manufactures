package cli

import (
	"bufio"
	"fmt"
	"manufactures/entity"
	"manufactures/handler"
	"os"
	"strings"
)

func showUserMenu() {
	if entity.LoggedInStaff.Position != "admin" {
		fmt.Println("Unauthorized: Only admin can manage users.")
		return
	}

	fmt.Println("-- Manage Users --")
	fmt.Println("1. Insert User")
	fmt.Println("2. View Users")
	fmt.Println("3. Back to Main Menu")
	fmt.Print("Choose an option: ")
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)
	if input == "1" {
		handler.InsertUser()
	} else if input == "2" {
		handler.PrintUser()
	} else if input == "3" {
		ShowMenu()
	} else {
		fmt.Println("Invalid input")
	}
}
