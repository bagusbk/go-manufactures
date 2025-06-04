package cli

import (
	"bufio"
	"fmt"
	"go-manufactures/handler"
	"os"
	"strings"
)

func showOrderMenu() {
	fmt.Println("-- Manage Orders --")
	fmt.Println("1. Insert Order")
	fmt.Println("2. View Orders")
	fmt.Println("3. Insert Order Items")
	fmt.Println("4. View Order Items")
	fmt.Println("5. Insert Payment")
	fmt.Println("6. View Payments")
	fmt.Println("7. Back to Main Menu")
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
	} else if input == "7" {
		ShowMenu()
	} else {
		fmt.Println("Invalid input")
	}
}
