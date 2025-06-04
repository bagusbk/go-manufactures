package cli

import (
	"fmt"
	"go-manufactures/handler"
)

func showReportMenu() {
	pos := handler.LoggedInStaff.Position
	if pos != "admin" && pos != "manager" {
		fmt.Println("Unauthorized: Only admin or manager can view reports.")
		return
	}

	fmt.Println("-- Reports --")
	fmt.Println("[placeholder for future reports]")
}
