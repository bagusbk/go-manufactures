package cli

import (
	"bufio"
	"fmt"
	"go-manufactures/handler"
	"os"
	"strings"
)

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
