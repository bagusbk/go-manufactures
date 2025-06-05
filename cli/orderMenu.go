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
	fmt.Println("3. Back to Main Menu")
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
		ShowMenu()
	} else {
		fmt.Println("Invalid input")
	}
}
