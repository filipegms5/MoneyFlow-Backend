package database

import (
	"github.com/filipegms5/MoneyFlow-Backend/models"

	"gorm.io/gorm"
)

func SetupDatabase(db *gorm.DB) {
	db.AutoMigrate(&models.Estabelecimento{}, &models.FormaPagamento{}, &models.Transacao{}, &models.Usuario{})
}
