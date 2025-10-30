package main

import (
	"os"

	"github.com/filipegms5/MoneyFlow-Backend/database"
	"github.com/filipegms5/MoneyFlow-Backend/router"
	"github.com/filipegms5/MoneyFlow-Backend/services"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// Instruções de como executar o projeto estão no README
func main() {

	port := os.Getenv("PORT")
	if port == "" {
		port = ":8000" // Porta padrão caso não seja especificada
	}

	// Conecta ao banco de dados
	db, err := gorm.Open(sqlite.Open("MoneyFlow.db"), &gorm.Config{})
	if err != nil {
		panic(err.Error())
	}

	database.SetupDatabase(db)
	router := router.SetupRouter(db)

	// Backfill categorias para estabelecimentos sem categoria na inicialização
	services.BackfillCategoriasMissing(db)

	router.Run(port)
}
