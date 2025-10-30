package database

import (
	"github.com/filipegms5/MoneyFlow-Backend/models"

	"gorm.io/gorm"
)

func SetupDatabase(db *gorm.DB) {
	db.AutoMigrate(&models.Categoria{}, &models.Estabelecimento{}, &models.FormaPagamento{}, &models.Transacao{}, &models.Usuario{}, &models.MetaFinanceira{})
	SeedCategoriasIfEmpty(db)
}

func SeedCategoriasIfEmpty(db *gorm.DB) {
	categorias := []models.Categoria{
		{ID: 9999999, Nome: "Outros", CnaeId: "9999999"},
		{ID: 4711302, Nome: "Supermercados", CnaeId: "4711302"},
		{ID: 4721102, Nome: "Padarias", CnaeId: "4721102"},
		{ID: 4771701, Nome: "Farmácias", CnaeId: "4771701"},
		{ID: 4722901, Nome: "Açougues", CnaeId: "4722901"},
		{ID: 5611201, Nome: "Restaurantes", CnaeId: "5611201"},
		{ID: 5611203, Nome: "Lanchonetes", CnaeId: "5611203"},
		{ID: 4781400, Nome: "Lojas de roupas", CnaeId: "4781400"},
		{ID: 4761003, Nome: "Papelarias", CnaeId: "4761003"},
		{ID: 4761001, Nome: "Livrarias", CnaeId: "4761001"},
		{ID: 4774100, Nome: "Óticas", CnaeId: "4774100"},
		{ID: 4783102, Nome: "Relojoarias", CnaeId: "4783102"},
		{ID: 4763601, Nome: "Brinquedos", CnaeId: "4763601"},
		{ID: 4754701, Nome: "Lojas de móveis", CnaeId: "4754701"},
		{ID: 4751201, Nome: "Lojas de informática", CnaeId: "4751201"},
		{ID: 4752100, Nome: "Lojas de telefonia", CnaeId: "4752100"},
		{ID: 4731800, Nome: "Postos de combustíveis", CnaeId: "4731800"},
		{ID: 9602501, Nome: "Cabeleireiros", CnaeId: "9602501"},
		{ID: 4763602, Nome: "Artigos esportivos", CnaeId: "4763602"},
		{ID: 4782201, Nome: "Lojas de calçados", CnaeId: "4782201"},
		{ID: 4783101, Nome: "Joalherias", CnaeId: "4783101"},
		{ID: 4729602, Nome: "Lojas de conveniência", CnaeId: "4729602"},
	}

	for _, cat := range categorias {
		var existing models.Categoria
		if err := db.First(&existing, cat.ID).Error; err != nil {
			_ = db.Create(&cat).Error
		}
	}
}
