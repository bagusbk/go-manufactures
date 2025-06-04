package cli

import (
	"bufio"
	"fmt"
	"go-manufactures/handler"
	"os"
	"strings"
)

func showStaffMenu() {
	if handler.LoggedInStaff.Position != "admin" {
		fmt.Println("Unauthorized: Only admin can manage staff.")
		return
	}

	reader := bufio.NewReader(os.Stdin)
	fmt.Println("-- Manage Staff --")
	fmt.Println("1. Insert Staff")
	fmt.Println("2. View Staff")
	fmt.Println("3. Delete Staff")
	fmt.Print("Choose an option: ")

	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)

	if input == "1" {
		handler.InsertStaff()
	} else if input == "2" {
		handler.PrintAllStaff()
	} else if input == "3" {
		handler.DeleteStaff()
	} else {
		fmt.Println("Invalid input")
	}
}
