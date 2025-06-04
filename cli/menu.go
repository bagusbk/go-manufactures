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
			fmt.Println("3. Order")
			fmt.Println("4. Reports")
			fmt.Println("5. Logout")
			fmt.Print("Choose an option: ")

			input, _ := reader.ReadString('\n')
			input = strings.TrimSpace(input)

			if input == "1" {
				showProductMenu()
			} else if input == "2" {
				showUserMenu()
			} else if input == "3" {
				showOrderMenu()
			} else if input == "4" {
				showReportMenu()
			} else if input == "5" {
				handler.LoggedInStaff.Email = ""
				handler.LoggedInStaff.Position = ""
				fmt.Println("Logged out.")
			} else {
				fmt.Println("Invalid input")
			}
		}
	}
}

func showProductMenu() {
	fmt.Println("-- Manage Products --")
	fmt.Println("1. Insert Product")
	fmt.Println("2. View Products")
	fmt.Print("Choose an option: ")
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)
	if input == "1" {
		handler.InsertProduct()
	} else if input == "2" {
		handler.PrintProduct()
	} else {
		fmt.Println("Invalid input")
	}
}

func showUserMenu() {
	if handler.LoggedInStaff.Position != "admin" {
		fmt.Println("Unauthorized: Only admin can manage users.")
		return
	}

	fmt.Println("-- Manage Users --")
	fmt.Println("1. Insert User")
	fmt.Println("2. View Users")
	fmt.Print("Choose an option: ")
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)
	if input == "1" {
		handler.InsertUser()
	} else if input == "2" {
		handler.PrintUser()
	} else {
		fmt.Println("Invalid input")
	}
}

func showOrderMenu() {
	fmt.Println("-- Manage Orders --")
	fmt.Println("1. Insert Order")
	fmt.Println("2. View Orders")
	fmt.Println("3. Insert Order Items")
	fmt.Println("4. View Order Items")
	fmt.Println("5. Insert Payment")
	fmt.Println("6. View Payments")
	fmt.Print("Choose an option: ")
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)

	if input == "1" {
		handler.InsertOrder()
	} else if input == "2" {
		handler.PrintOrders()
	} else if input == "3" {
		handler.InsertOrderItem()
	} else if input == "4" {
		handler.PrintOrderItems()
	} else if input == "5" {
		handler.InsertPayment()
	} else if input == "6" {
		handler.PrintPayments()
	} else {
		fmt.Println("Invalid input")
	}
}

func showReportMenu() {
	pos := handler.LoggedInStaff.Position
	if pos != "admin" && pos != "manager" {
		fmt.Println("Unauthorized: Only admin or manager can view reports.")
		return
	}

	fmt.Println("-- Reports --")
	fmt.Println("[placeholder for future reports]")
}
