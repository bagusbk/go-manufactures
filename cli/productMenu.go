package cli

import (
	"bufio"
	"fmt"
	"manufactures/handler"
	"os"
	"strings"
)

func showProductMenu() {
	fmt.Println("-- Manage Products --")
	fmt.Println("1. Insert Product")
	fmt.Println("2. View Products")
	fmt.Println("3. Back to Main Menu")
	fmt.Print("Choose an option: ")
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)
	if input == "1" {
		handler.InsertProduct()
	} else if input == "2" {
		handler.PrintProduct()
	} else if input == "3" {
		ShowMenu()
	} else {
		fmt.Println("Invalid input")
	}
}
