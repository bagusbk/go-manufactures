package handler

import (
	"bufio"
	"database/sql"
	"fmt"
	"manufactures/config"
	"os"
	"regexp"
	"strconv"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
)

func hashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func isValidEmail(email string) bool {
	re := regexp.MustCompile(`^[a-z0-9._%+-]+@[a-z0-9.-]+\.[a-z]{2,}$`)
	return re.MatchString(email)
}

var LoggedInStaff struct {
	StaffID  int
	Email    string
	Position string
}

func LoginUser() string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter email: ")
	emailInput, _ := reader.ReadString('\n')
	emailInput = strings.TrimSpace(emailInput)

	fmt.Print("Enter password: ")
	passwordInput, _ := reader.ReadString('\n')
	passwordInput = strings.TrimSpace(passwordInput)

	var staffId int
	var email, passwordHash, position string
	err := config.InitDB().QueryRow("SELECT staff_id, email, password_hash, position FROM staff WHERE email = ?", emailInput).Scan(&staffId, &email, &passwordHash, &position)
	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Println("Email not found.")
		} else {
			fmt.Println("Error:", err)
		}
		return ""
	}

	if position == "admin" {
		if passwordInput != passwordHash {
			fmt.Println("Incorrect password.")
			return ""
		}
	} else {
		err = bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(passwordInput))
		if err != nil {
			fmt.Println("Incorrect password.")
			return ""
		}
	}

	LoggedInStaff.StaffID = staffId
	LoggedInStaff.Email = email
	LoggedInStaff.Position = position
	fmt.Printf("Login successful! Role: %s\n", position)
	return email
}

func InsertStaff() {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Enter full name: ")
	fullName, _ := reader.ReadString('\n')
	fullName = strings.TrimSpace(fullName)

	fmt.Print("Enter position (admin/manager/staff): ")
	position, _ := reader.ReadString('\n')
	position = strings.TrimSpace(position)

	validPositions := []string{"admin", "manager", "staff"}
	isValidPosition := false
	for _, valid := range validPositions {
		if position == valid {
			isValidPosition = true
			break
		}
	}

	if !isValidPosition {
		fmt.Println("Invalid position entered. Please enter one of the following: admin, manager, staff.")
		return
	}

	var email string
	for {
		fmt.Print("Enter email: ")
		email, _ = reader.ReadString('\n')
		email = strings.TrimSpace(email)

		if !isValidEmail(email) {
			fmt.Println("Invalid email format. Please enter a valid email.")
		} else {
			break
		}
	}

	var exists int
	checkErr := config.InitDB().QueryRow("SELECT COUNT(*) FROM staff WHERE email = ?", email).Scan(&exists)
	if checkErr != nil {
		fmt.Println("Error checking email:", checkErr)
		return
	}
	if exists > 0 {
		fmt.Println("⚠️ Email already registered.")
		return
	}

	fmt.Print("Enter password: ")
	password, _ := reader.ReadString('\n')
	password = strings.TrimSpace(password)

	hashedPassword, err := hashPassword(password)
	if err != nil {
		fmt.Println("Error hashing password:", err)
		return
	}

	if fullName == "" || email == "" || password == "" {
		fmt.Println("All fields are required.")
		return
	}

	_, err = config.InitDB().Exec(`INSERT INTO staff (full_name, position, email, password_hash) VALUES (?, ?, ?, ?)`, fullName, position, email, hashedPassword)
	if err != nil {
		fmt.Println("Error inserting staff:", err)
		return
	}

	fmt.Println("✅ Staff successfully inserted.")
}

func PrintAllStaff() {
	rows, err := config.InitDB().Query("SELECT staff_id, full_name, position, email, created_at FROM staff")
	if err != nil {
		fmt.Println("Error retrieving staff:", err)
		return
	}
	defer rows.Close()

	fmt.Println("List of Staff:")
	for rows.Next() {
		var id int
		var name, position, email string
		var createdAt string
		if err := rows.Scan(&id, &name, &position, &email, &createdAt); err != nil {
			fmt.Println("Scan error:", err)
			return
		}
		fmt.Printf("ID: %d | Name: %s | Position: %s | Email: %s | Created At: %s\n", id, name, position, email, createdAt)
	}
}

func DeleteStaff() {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Enter staff ID to delete: ")
	idStr, _ := reader.ReadString('\n')
	idStr = strings.TrimSpace(idStr)

	staffIDToDelete, err := strconv.Atoi(idStr)
	if err != nil {
		fmt.Println("Invalid staff ID format.")
		return
	}

	if staffIDToDelete == LoggedInStaff.StaffID {
		fmt.Println("❌ You cannot delete your own account!")
		return
	}

	var exists int
	err = config.InitDB().QueryRow("SELECT COUNT(*) FROM staff WHERE staff_id = ?", staffIDToDelete).Scan(&exists)
	if err != nil {
		fmt.Println("Error checking staff existence:", err)
		return
	}
	if exists == 0 {
		fmt.Println("⚠️ No staff found with that ID.")
		return
	}

	res, err := config.InitDB().Exec("DELETE FROM staff WHERE staff_id = ?", staffIDToDelete)
	if err != nil {
		fmt.Println("Error deleting staff:", err)
		return
	}

	fmt.Println("✅ Staff successfully deleted.")
}

func UpdateStaffRole() {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Enter Staff ID to update role: ")
	idStr, _ := reader.ReadString('\n')
	idStr = strings.TrimSpace(idStr)

	staffID, err := strconv.Atoi(idStr)
	if err != nil {
		fmt.Println("Invalid staff ID.")
		return
	}

	if staffID == LoggedInStaff.StaffID {
		fmt.Println("❌ You cannot change your own role.")
		return
	}

	var exists int
	err = config.InitDB().QueryRow("SELECT COUNT(*) FROM staff WHERE staff_id = ?", staffID).Scan(&exists)
	if err != nil {
		fmt.Println("Error checking staff existence:", err)
		return
	}
	if exists == 0 {
		fmt.Println("⚠️ No staff found with that ID.")
		return
	}

	fmt.Print("Enter new role (admin/manager/staff): ")
	newRole, _ := reader.ReadString('\n')
	newRole = strings.TrimSpace(newRole)

	validRoles := []string{"admin", "manager", "staff"}
	isValid := false
	for _, r := range validRoles {
		if newRole == r {
			isValid = true
			break
		}
	}
	if !isValid {
		fmt.Println("Invalid role. Allowed: admin, manager, staff.")
		return
	}

	var adminCount int
	err = config.InitDB().QueryRow("SELECT COUNT(*) FROM staff WHERE position = 'admin'").Scan(&adminCount)
	if err != nil {
		fmt.Println("Database error:", err)
		return
	}

	var currentRole string
	err = config.InitDB().QueryRow("SELECT position FROM staff WHERE staff_id = ?", staffID).Scan(&currentRole)
	if err != nil {
		fmt.Println("Staff not found.")
		return
	}

	if currentRole == "admin" && newRole != "admin" && adminCount == 1 {
		fmt.Println("⚠️ Cannot change role. This is the last admin.")
		return
	}

	_, err = config.InitDB().Exec("UPDATE staff SET position = ? WHERE staff_id = ?", newRole, staffID)
	if err != nil {
		fmt.Println("Error updating staff role:", err)
		return
	}

	fmt.Println("✅ Staff role updated successfully.")
}
