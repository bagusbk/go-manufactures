package handler

import (
	"bufio"
	"fmt"
	"log"
	"manufactures/config"
	"manufactures/entity"
	"os"
	"strings"
	"time"
)

// Fungsi untuk mendapatkan input integer dari user
func getInputInt(reader *bufio.Reader) int {
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)
	var result int
	fmt.Sscanf(input, "%d", &result)
	return result
}

// Fungsi untuk mendapatkan input string dari user
func getInputString(reader *bufio.Reader) string {
	input, _ := reader.ReadString('\n')
	return strings.TrimSpace(input)
}

// Fungsi untuk mendapatkan input float dari user
func getInputFloat(reader *bufio.Reader) float64 {
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)
	var result float64
	fmt.Sscanf(input, "%f", &result)
	return result
}

func readInput() string {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()        // Membaca input
	return scanner.Text() // Mengembalikan input sebagai string
}

// Handler untuk laporan pesanan
func PrintOrderReport() {
	rows, err := config.InitDB().Query(`
		SELECT o.order_id, o.user_id, u.full_name, o.order_date, o.status, o.total_amount, o.staff_id, s.full_name
		FROM orders o
		JOIN users u ON o.user_id = u.user_id
		JOIN staff s ON o.staff_id = s.staff_id
	`)
	if err != nil {
		fmt.Println("Error retrieving orders:", err)
		return
	}
	defer rows.Close()

	var reports []entity.OrderReport
	for rows.Next() {
		var report entity.OrderReport
		err := rows.Scan(&report.OrderID, &report.UserID, &report.FullName, &report.OrderDate, &report.Status, &report.TotalAmount, &report.StaffID, &report.StaffName)
		if err != nil {
			fmt.Println("Error scanning row:", err)
			return
		}
		reports = append(reports, report)
	}

	fmt.Println("Order Report:")
	for _, report := range reports {
		fmt.Printf("OrderID: %d | UserID: %d | Full Name User: %s | Order Date: %s | Status: %s | Total Amount: %.2f | Staff Name: %s\n",
			report.OrderID, report.UserID, report.FullName, report.OrderDate, report.Status, report.TotalAmount, report.StaffName)
	}
}

func CreateOrder() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter user ID: ")
	userID := getInputInt(reader)

	var orderItems []entity.OrderItem
	var totalAmount float64

	staffID := entity.LoggedInStaff.StaffID
	var count int
	err := config.InitDB().QueryRow("SELECT COUNT(*) FROM staff WHERE staff_id = ?", staffID).Scan(&count)
	if err != nil || count == 0 {
		fmt.Println("Staff ID not found in the database!")
		return
	}

	for {
		fmt.Print("Enter item ID (or 0 to finish): ")
		itemID := getInputInt(reader)
		if itemID == 0 {
			break
		}

		fmt.Print("Enter quantity: ")
		quantity := getInputInt(reader)

		var price float64
		var stock int
		// Mendapatkan harga dari item berdasarkan item ID
		err := config.InitDB().QueryRow("SELECT price, stock FROM item WHERE item_id = ?", itemID).Scan(&price, &stock)
		if err != nil {
			fmt.Println("Error retrieving item:", err)
			return
		}

		if stock <= 0 {
			fmt.Printf("Item ID %d is out of stock and cannot be ordered.\n", itemID)
			continue
		}

		if quantity > stock {
			fmt.Printf("Only %d items in stock for Item ID %d. Please enter a lower quantity.\n", stock, itemID)
			continue
		}

		// Hitung total harga untuk item ini
		orderItem := entity.OrderItem{
			ItemID:   itemID,
			Quantity: quantity,
			Price:    price,
		}

		totalAmount += price * float64(quantity)
		orderItems = append(orderItems, orderItem)
	}

	// Insert order ke tabel 'orders'
	tx, err := config.InitDB().Begin()
	if err != nil {
		log.Fatal("Transaction begin failed: ", err)
	}
	defer tx.Rollback()

	// Insert order
	orderStmt, err := tx.Prepare("INSERT INTO orders (user_id, total_amount, staff_id) VALUES (?, ?, ?)")
	if err != nil {
		log.Fatal("Prepare statement failed: ", err)
	}
	defer orderStmt.Close()

	res, err := orderStmt.Exec(userID, totalAmount, staffID)
	if err != nil {
		log.Fatal("Order insertion failed: ", err)
	}

	// Mendapatkan order ID yang baru saja dimasukkan
	orderID, err := res.LastInsertId()
	if err != nil {
		log.Fatal("Failed to get last insert ID: ", err)
	}

	// Insert order items
	orderItemStmt, err := tx.Prepare("INSERT INTO order_items (order_id, item_id, quantity, price) VALUES (?, ?, ?, ?)")
	if err != nil {
		log.Fatal("Prepare statement for order items failed: ", err)
	}
	defer orderItemStmt.Close()

	for _, item := range orderItems {
		_, err := orderItemStmt.Exec(orderID, item.ItemID, item.Quantity, item.Price)
		if err != nil {
			log.Fatal("Failed to insert order item:", err)
		}
	}

	paymentStmt, err := tx.Prepare("INSERT INTO payment (user_id, order_id, amount, payment_method, status, payment_date) VALUES (?, ?, ?, ?, ?, NULL)")
	if err != nil {
		log.Fatal("Failed to prepare payment statement:", err)
	}
	defer paymentStmt.Close()

	// Pada awal, pembayaran dianggap "pending"
	_, err = paymentStmt.Exec(userID, orderID, totalAmount, "cash", "pending")
	if err != nil {
		log.Fatal("Failed to insert payment:", err)
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		log.Fatal("Transaction commit failed:", err)
	}

	fmt.Println("Order created successfully! Order ID:", orderID)
}

func PrintOrderDetailsByDateRange() {
	// Input tanggal dari pengguna
	fmt.Println("Masukkan tanggal mulai (YYYY-MM-DD):")
	startDate := readInput()

	fmt.Println("Masukkan tanggal akhir (YYYY-MM-DD):")
	endDate := readInput()

	// Parse tanggal yang diterima untuk query
	start, err := time.Parse("2006-01-02", startDate)
	if err != nil {
		log.Println("Tanggal mulai tidak valid:", err)
		return
	}
	end, err := time.Parse("2006-01-02", endDate)
	if err != nil {
		log.Println("Tanggal akhir tidak valid:", err)
		return
	}

	// Inisialisasi koneksi ke database dan eksekusi query
	rows, err := config.InitDB().Query(`
		SELECT o.order_id, o.order_date, u.full_name, o.total_amount, o.status
		FROM orders o
		JOIN users u ON o.user_id = u.user_id
		WHERE o.order_date BETWEEN ? AND ?
		ORDER BY o.order_date DESC
	`, start, end)
	if err != nil {
		fmt.Println("Error retrieving orders:", err)
		return
	}
	defer rows.Close()

	// Menampung hasil query dalam slice
	var reports []entity.OrderReport
	for rows.Next() {
		var report entity.OrderReport
		err := rows.Scan(&report.OrderID, &report.OrderDate, &report.FullName, &report.TotalAmount, &report.Status)
		if err != nil {
			fmt.Println("Error scanning row:", err)
			return
		}
		reports = append(reports, report)
	}

	// Menampilkan hasil laporan
	fmt.Printf("Laporan Pesanan dari %s sampai %s:\n", startDate, endDate)
	for _, report := range reports {
		fmt.Printf("OrderID: %d | User: %s | Order Date: %s | Total Amount: %.2f | Status: %s\n",
			report.OrderID, report.FullName, report.OrderDate, report.TotalAmount, report.Status)
	}
}
