package handler

import (
	"database/sql"
	"fmt"
	"strconv"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// InsertProductTestHelper mensimulasikan fungsi InsertProduct untuk menguji penambahan produk
func InsertProductTestHelper(db *sql.DB, name, priceStr, stockStr string) error {
	// Jika ada input yang tidak valid (kosong), kembalikan error
	if name == "" || priceStr == "" || stockStr == "" {
		return fmt.Errorf("All fields are required.") // Semua field harus diisi
	}

	// Mengonversi price (harga) ke tipe float64
	price, err := strconv.ParseFloat(priceStr, 64)
	if err != nil {
		return fmt.Errorf("Invalid price format.") // Format harga tidak valid
	}

	// Mengonversi stock (stok) ke tipe int
	stock, err := strconv.Atoi(stockStr)
	if err != nil {
		return fmt.Errorf("Invalid stock format.") // Format stok tidak valid
	}

	// Mensimulasikan eksekusi query untuk memasukkan data produk ke dalam database
	_, err = db.Exec("INSERT INTO item (name, price, stock) VALUES (?, ?, ?)", name, price, stock)
	return err
}

// Unit Test untuk InsertProductTestHelper
func TestInsertProduct(t *testing.T) {
	// Membuat koneksi mock ke database menggunakan sqlmock
	db, mock, err := sqlmock.New() // Menginisialisasi koneksi mock
	require.NoError(t, err)        // Memastikan tidak ada error dalam setup mock database
	defer db.Close()               // Pastikan koneksi ditutup setelah selesai

	// Setup ekspektasi mock untuk eksekusi query INSERT produk
	mock.ExpectExec("INSERT INTO item").WithArgs("Product1", 100.0, 10).WillReturnResult(sqlmock.NewResult(1, 1))

	// Uji dengan input valid
	err = InsertProductTestHelper(db, "Product1", "100.0", "10")
	require.NoError(t, err, "InsertProductTestHelper should succeed with valid inputs") // Harus berhasil dengan input valid

	// Uji input invalid: nama kosong
	err = InsertProductTestHelper(db, "", "100.0", "10")
	require.Error(t, err, "InsertProductTestHelper should fail when name is empty") // Harus gagal karena nama kosong

	// Uji input invalid: format harga tidak valid
	err = InsertProductTestHelper(db, "Product1", "abc", "10")
	require.Error(t, err, "InsertProductTestHelper should fail when price is invalid") // Harus gagal karena format harga salah

	// Uji input invalid: format stok tidak valid
	err = InsertProductTestHelper(db, "Product1", "100.0", "xyz")
	require.Error(t, err, "InsertProductTestHelper should fail when stock is invalid") // Harus gagal karena format stok salah
}

// Unit Test untuk InsertProductTestHelper - Kasus Gagal
func TestInsertProductFail(t *testing.T) {
	// Membuat koneksi mock ke database menggunakan sqlmock
	db, _, err := sqlmock.New() // Menginisialisasi koneksi mock
	require.NoError(t, err)     // Memastikan tidak ada error dalam setup mock database
	defer db.Close()            // Pastikan koneksi ditutup setelah selesai

	// Test input invalid: nama kosong
	err = InsertProductTestHelper(db, "", "100.0", "10")
	require.Error(t, err, "InsertProductTestHelper should fail when name is empty") // Harus gagal karena nama kosong

	// Test input invalid: harga tidak valid
	err = InsertProductTestHelper(db, "Product1", "abc", "10")
	require.Error(t, err, "InsertProductTestHelper should fail when price is invalid") // Harus gagal karena format harga salah

	// Test input invalid: stok tidak valid
	err = InsertProductTestHelper(db, "Product1", "100.0", "xyz")
	require.Error(t, err, "InsertProductTestHelper should fail when stock is invalid") // Harus gagal karena format stok salah

	// Test input invalid: semua field kosong
	err = InsertProductTestHelper(db, "", "", "")
	require.Error(t, err, "InsertProductTestHelper should fail when all fields are empty") // Harus gagal karena semua field kosong
}

// PrintMostSoldItemsReportTest menghasilkan laporan untuk item yang paling laris (mensimulasikan query database)
func PrintMostSoldItemsReportTest(db *sql.DB) string {
	// Menjalankan query untuk mendapatkan item yang paling laris
	rows, err := db.Query("SELECT item_id, name, total_sold FROM items")
	if err != nil {
		return fmt.Sprintf("Error executing query: %v", err) // Menangani error jika query gagal
	}
	defer rows.Close() // Pastikan koneksi rows ditutup setelah selesai

	var result string
	// Iterasi untuk membaca hasil query
	for rows.Next() {
		var itemID int
		var name string
		var totalSold int
		err := rows.Scan(&itemID, &name, &totalSold)
		if err != nil {
			return fmt.Sprintf("Error scanning row: %v", err) // Menangani error saat pemindaian data
		}
		// Menambahkan hasil item ke dalam string laporan
		result += fmt.Sprintf("Item ID: %d | Name: %s | Total Sold: %d\n", itemID, name, totalSold)
	}

	return result
}

// Unit Test untuk PrintMostSoldItemsReportTest
func TestPrintMostSoldItemsReport(t *testing.T) {
	// Membuat koneksi mock ke database menggunakan sqlmock
	db, mock, err := sqlmock.New() // Menginisialisasi koneksi mock
	require.NoError(t, err)        // Memastikan tidak ada error dalam setup mock database
	defer db.Close()               // Pastikan koneksi ditutup setelah selesai

	// Setup mock rows untuk hasil query
	rows := sqlmock.NewRows([]string{"item_id", "name", "total_sold"}).
		AddRow(1, "Magnum", 150).
		AddRow(2, "Dove", 120).
		AddRow(3, "Vaseline", 100)

	// Setup ekspektasi mock untuk query SELECT
	mock.ExpectQuery("SELECT item_id, name, total_sold FROM items").
		WillReturnRows(rows) // Mengembalikan mock rows sebagai hasil query

	// Panggil fungsi dan tangkap hasil output
	report := PrintMostSoldItemsReportTest(db)

	// Output yang diharapkan
	expectedOutput := "Item ID: 1 | Name: Magnum | Total Sold: 150\n" +
		"Item ID: 2 | Name: Dove | Total Sold: 120\n" +
		"Item ID: 3 | Name: Vaseline | Total Sold: 100\n"

	// Validasi hasil
	assert.Equal(t, expectedOutput, report, "The report output should match the expected format.") // Memastikan hasil output sesuai dengan yang diharapkan
}
