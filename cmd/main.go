package main

import (
	"database/sql"
	"github/com/vybraan/snare"
	"github/com/vybraan/snare/database"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {

	db, err := sql.Open("sqlite3", "./app.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if err := database.Migrate(db); err != nil {
		log.Fatal(err)
	}
	if err := database.Seed(db); err != nil {
		log.Fatal(err)
	}

	auth := snare.New(db)

	r := gin.Default()
	auth.RegisterRoutes(r)

	r.Run(":8080")
}
