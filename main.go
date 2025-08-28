package main

import (
	"os"

	"github.com/filipegms5/MoneyFlow-Backend/database"
	"github.com/filipegms5/MoneyFlow-Backend/router"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// Instruction os how to run the project on README
func main() {

	port := os.Getenv("PORT")
	if port == "" {
		port = ":8000" // Default port if not specified
	}

	// Connect to the database
	db, err := gorm.Open(sqlite.Open("MoneyFlow.db"), &gorm.Config{})
	if err != nil {
		panic(err.Error())
	}

	database.SetupDatabase(db)
	router := router.SetupRouter(db)

	router.Run(port)
}
