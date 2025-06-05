package cli

import (
	"bufio"
	"fmt"
	"manufactures/entity"
	"manufactures/handler"
	"os"
	"strings"
)

func showReportMenu() {
	pos := entity.LoggedInStaff.Position
	if pos != "admin" && pos != "manager" {
		fmt.Println("Unauthorized: Only admin or manager can view reports.")
		return
	}

	fmt.Println("-- Reports --")
	fmt.Println("1. Report User Pembelian Terbanyak")
	fmt.Println("2. Report Pesanan Range Date")
	fmt.Println("3. Report Barang Terjual")
	fmt.Println("7. Back to Main Menu")
	fmt.Print("Choose an option: ")
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)

	if input == "1" {
		// handler.PrintItemReport()
		handler.PrintFrequentBuyersReport()
	} else if input == "2" {
		handler.PrintOrderDetailsByDateRange()
	} else if input == "3" {
		handler.PrintMostSoldItemsReport()
	} else if input == "7" {
		ShowMenu()
	} else {
		fmt.Println("Invalid input")
	}
}
