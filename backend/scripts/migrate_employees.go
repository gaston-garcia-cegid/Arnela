package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

// Employee represents an existing employee record
type Employee struct {
	ID        uuid.UUID
	Email     string
	FirstName string
	LastName  string
	DNI       string
	IsActive  bool
	UserID    *uuid.UUID
}

func main() {
	// Database connection string
	// You can override with DATABASE_URL environment variable
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		// Default values from .env.example
		dbHost := getEnv("DB_HOST", "localhost")
		dbPort := getEnv("DB_PORT", "5432")
		dbUser := getEnv("DB_USER", "arnela_user")
		dbPassword := getEnv("DB_PASSWORD", "arnela_secure_pass_2024")
		dbName := getEnv("DB_NAME", "arnela_db")
		dbSSLMode := getEnv("DB_SSLMODE", "disable")

		dsn = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
			dbHost, dbPort, dbUser, dbPassword, dbName, dbSSLMode)
	}

	fmt.Printf("Connecting to database...\n")

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	ctx := context.Background()

	// Start transaction
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		log.Fatalf("Failed to start transaction: %v", err)
	}
	defer tx.Rollback()

	// Get all employees without user_id
	query := `
		SELECT id, email, first_name, last_name, dni, is_active, user_id
		FROM employees
		WHERE deleted_at IS NULL
		  AND user_id IS NULL
		  AND email IS NOT NULL
		  AND dni IS NOT NULL
	`

	rows, err := tx.QueryContext(ctx, query)
	if err != nil {
		log.Fatalf("Failed to query employees: %v", err)
	}
	defer rows.Close()

	var employees []Employee
	for rows.Next() {
		var emp Employee
		err := rows.Scan(&emp.ID, &emp.Email, &emp.FirstName, &emp.LastName, &emp.DNI, &emp.IsActive, &emp.UserID)
		if err != nil {
			log.Fatalf("Failed to scan employee: %v", err)
		}
		employees = append(employees, emp)
	}

	if len(employees) == 0 {
		fmt.Println("No employees found without user_id")
		return
	}

	fmt.Printf("Found %d employees without user accounts\n", len(employees))

	// Create user for each employee
	createUserQuery := `
		INSERT INTO users (id, email, password_hash, first_name, last_name, role, is_active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, NOW(), NOW())
		ON CONFLICT (email) DO NOTHING
		RETURNING id
	`

	updateEmployeeQuery := `
		UPDATE employees
		SET user_id = $1, updated_at = NOW()
		WHERE id = $2
	`

	successCount := 0
	skipCount := 0

	for _, emp := range employees {
		// Hash the DNI as password
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(emp.DNI), bcrypt.DefaultCost)
		if err != nil {
			log.Printf("Failed to hash password for %s: %v", emp.Email, err)
			continue
		}

		// Create user
		userID := uuid.New()
		var returnedID uuid.UUID
		err = tx.QueryRowContext(
			ctx,
			createUserQuery,
			userID,
			emp.Email,
			string(hashedPassword),
			emp.FirstName,
			emp.LastName,
			"employee",
			emp.IsActive,
		).Scan(&returnedID)

		if err != nil {
			if err == sql.ErrNoRows {
				// User already exists (ON CONFLICT DO NOTHING)
				// Try to find existing user
				var existingUserID uuid.UUID
				err = tx.QueryRowContext(ctx, "SELECT id FROM users WHERE email = $1", emp.Email).Scan(&existingUserID)
				if err != nil {
					log.Printf("Skipping %s: email already exists but couldn't fetch user: %v", emp.Email, err)
					skipCount++
					continue
				}
				userID = existingUserID
			} else {
				log.Printf("Failed to create user for %s: %v", emp.Email, err)
				continue
			}
		}

		// Update employee with user_id
		_, err = tx.ExecContext(ctx, updateEmployeeQuery, userID, emp.ID)
		if err != nil {
			log.Printf("Failed to update employee %s with user_id: %v", emp.Email, err)
			continue
		}

		fmt.Printf("✓ Created user for %s %s (%s) - Login: %s / %s\n",
			emp.FirstName, emp.LastName, emp.Email, emp.Email, emp.DNI)
		successCount++
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		log.Fatalf("Failed to commit transaction: %v", err)
	}

	fmt.Printf("\n=== Migration Summary ===\n")
	fmt.Printf("Total employees processed: %d\n", len(employees))
	fmt.Printf("Successfully migrated: %d\n", successCount)
	fmt.Printf("Skipped (already existed): %d\n", skipCount)
	fmt.Printf("Failed: %d\n", len(employees)-successCount-skipCount)

	// Verify migration
	var total, withUser, withoutUser int
	err = db.QueryRowContext(ctx, `
		SELECT 
			COUNT(*) as total,
			COUNT(user_id) as with_user,
			COUNT(*) - COUNT(user_id) as without_user
		FROM employees
		WHERE deleted_at IS NULL
	`).Scan(&total, &withUser, &withoutUser)

	if err != nil {
		log.Printf("Failed to verify migration: %v", err)
	} else {
		fmt.Printf("\n=== Verification ===\n")
		fmt.Printf("Total employees: %d\n", total)
		fmt.Printf("Employees with user: %d\n", withUser)
		fmt.Printf("Employees without user: %d\n", withoutUser)
	}
}

// getEnv gets an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
