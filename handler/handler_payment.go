package handler

import (
	"bufio"
	"fmt"
	"log"
	"manufactures/config"
	"os"
)

type PaymentReport struct {
	UserID        int     `json:"user_id"`
	FullName      string  `json:"full_name"`
	Email         string  `json:"email"`
	PhoneNumber   string  `json:"phone_number"`
	Amount        float64 `json:"amount"`
	PaymentDate   string  `json:"payment_date"`
	PaymentMethod string  `json:"payment_method"`
	Status        string  `json:"status"`
}

// Handler untuk laporan pembayaran
func PrintPaymentReport() {
	rows, err := config.InitDB().Query(`
		SELECT u.user_id, u.full_name, u.email, u.phone_number, p.amount, p.payment_date, p.payment_method, p.status
		FROM payment p
		JOIN users u ON p.user_id = u.user_id
	`)
	if err != nil {
		fmt.Println("Error retrieving payments:", err)
		return
	}
	defer rows.Close()

	var reports []PaymentReport
	for rows.Next() {
		var report PaymentReport
		err := rows.Scan(&report.UserID, &report.FullName, &report.Email, &report.PhoneNumber, &report.Amount, &report.PaymentDate, &report.PaymentMethod, &report.Status)
		if err != nil {
			fmt.Println("Error scanning row:", err)
			return
		}
		reports = append(reports, report)
	}

	fmt.Println("Payment Report:")
	for _, report := range reports {
		fmt.Printf("UserID: %d | FullName: %s | Email: %s | Amount: %.2f | PaymentDate: %s | PaymentMethod: %s | Status: %s\n",
			report.UserID, report.FullName, report.Email, report.Amount, report.PaymentDate, report.PaymentMethod, report.Status)
	}
}

func CreatePayment() {
	// Gunakan koneksi yang sudah ada
	db := config.InitDB() // Hanya tangkap db tanpa error

	// Membaca input dari pengguna
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter order ID: ")
	orderID := getInputInt(reader)

	// Mengecek apakah pesanan ada di tabel orders
	var userID, totalAmount float64
	var orderStatus string
	err := db.QueryRow("SELECT user_id, total_amount, status FROM orders WHERE order_id = ?", orderID).Scan(&userID, &totalAmount, &orderStatus)
	if err != nil {
		fmt.Println("Order not found:", err)
		return
	}

	if orderStatus == "completed" {
		fmt.Println("This order has already been completed. Payment cannot be processed.")
		return
	}

	// Mengonfirmasi jumlah pembayaran
	var paymentAmount float64
	fmt.Printf("Total amount to pay for order %d is %.2f. Enter payment amount: ", orderID, totalAmount)
	paymentAmount = getInputFloat(reader)

	if paymentAmount != totalAmount {
		fmt.Println("Payment amount does not match the total amount.")
		return
	}

	// Meminta metode pembayaran
	fmt.Print("Enter payment method (e.g., cash, card): ")
	paymentMethod := getInputString(reader)

	// Memulai transaksi untuk memastikan atomik (semua langkah selesai atau tidak sama sekali)
	tx, err := db.Begin()
	if err != nil {
		log.Fatal("Transaction begin failed: ", err)
	}
	defer tx.Rollback()

	// Insert pembayaran ke tabel payment, menambahkan order_id
	_, err = tx.Exec("UPDATE payment SET amount = ?, payment_method = ?, status = 'paid', payment_date = CURRENT_TIMESTAMP WHERE order_id = ?", paymentAmount, paymentMethod, orderID)
	if err != nil {
		log.Fatal("Failed to update payment:", err)
	}

	// Update status di tabel orders menjadi 'completed'
	_, err = tx.Exec("UPDATE orders SET status = 'completed' WHERE order_id = ?", orderID)
	if err != nil {
		log.Fatal("Failed to update order status:", err)
	}

	rows, err := tx.Query("SELECT item_id, quantity FROM order_items WHERE order_id = ?", orderID)
	if err != nil {
		log.Fatal("Failed to retrieve order items:", err)
	}
	defer rows.Close()

	// Untuk setiap item, lakukan transaksi update secara terpisah
	for rows.Next() {
		var itemID, quantity int
		err := rows.Scan(&itemID, &quantity)
		if err != nil {
			log.Fatal("Failed to scan order item:", err)
		}

		// Mulai transaksi per item
		itemTx, err := db.Begin()
		if err != nil {
			log.Fatal("Transaction begin failed: ", err)
		}
		defer itemTx.Rollback()

		// Update stok untuk item tertentu
		_, err = itemTx.Exec("UPDATE item SET stock = stock - ? WHERE item_id = ?", quantity, itemID)
		if err != nil {
			log.Fatal("Failed to update item stock:", err)
		}

		if err := itemTx.Commit(); err != nil {
			log.Fatal("Transaction commit failed:", err)
		}
	}

	// Commit transaksi jika semua langkah berhasil
	if err := tx.Commit(); err != nil {
		log.Fatal("Transaction commit failed:", err)
	}

	fmt.Println("Payment successful and order completed!")
}
