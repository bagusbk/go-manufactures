package cli

import (
	"bufio"
	"fmt"
	"manufactures/handler"
	"os"
	"strings"
)

func showReportMenu() {
	pos := handler.LoggedInStaff.Position
	if pos != "admin" && pos != "manager" {
		fmt.Println("Unauthorized: Only admin or manager can view reports.")
		return
	}

	fmt.Println("-- Reports --")
	fmt.Println("1. Report Stock Item")
	fmt.Println("2. Report Pesanan")
	fmt.Println("3. Report Pembayaran")
	fmt.Println("7. Back to Main Menu")
	fmt.Print("Choose an option: ")
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)

	if input == "1" {
		handler.PrintItemReport()
	} else if input == "2" {
		handler.PrintOrderReport()
	} else if input == "3" {
		handler.PrintPaymentReport()
	} else if input == "7" {
		ShowMenu()
	} else {
		fmt.Println("Invalid input")
	}
}
