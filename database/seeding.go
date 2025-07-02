package database

import (
	"database/sql"
	"log"

	"golang.org/x/crypto/bcrypt"
)

func Seed(db *sql.DB) error {
	var count int
	if err := db.QueryRow("SELECT COUNT(*) FROM users").Scan(&count); err != nil {
		return err
	}

	if count > 0 {
		log.Println("Users table already seeded")
		return nil
	}

	log.Println("Seeding users table...")

	passwordHash, err := bcrypt.GenerateFromPassword([]byte("admin"), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	result, err := db.Exec("INSERT INTO users (username, password) VALUES (?, ?)", "admin", passwordHash)
	if err != nil {
		return err
	}

	userID, err := result.LastInsertId()
	if err != nil {
		return err
	}

	log.Printf("User seeded successfully with ID %d", userID)
	return nil
}
