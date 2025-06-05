package cli

import (
	"bufio"
	"fmt"
	"manufactures/handler"
	"os"
	"strings"
)

func showOrderMenu() {
	fmt.Println("-- Manage Orders --")
	fmt.Println("1. Insert Order")
	fmt.Println("2. Insert Payment")
	fmt.Println("3. View Order")
	fmt.Println("4. View Payment")
	fmt.Println("5. Back to Main Menu")
	fmt.Print("Choose an option: ")
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)

	if input == "1" {
		handler.CreateOrder()
	} else if input == "2" {
		// handler.PrintOrders()
		handler.CreatePayment()
	} else if input == "3" {
		// handler.PrintOrders()
		handler.PrintOrderReport()
	} else if input == "4" {
		// handler.PrintOrders()
		handler.PrintPaymentReport()
	} else if input == "5" {
		ShowMenu()
	} else {
		fmt.Println("Invalid input")
	}
}
