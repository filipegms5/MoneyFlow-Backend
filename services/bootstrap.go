package services

import (
	"context"
	"time"

	"strings"

	"github.com/filipegms5/MoneyFlow-Backend/models"
	"gorm.io/gorm"
)

// BackfillCategoriasMissing encontra estabelecimentos sem categoria e tenta atribuir automaticamente
func BackfillCategoriasMissing(db *gorm.DB) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var estabelecimentos []models.Estabelecimento
	if err := db.Where("categoria_id = 0 OR categoria_id IS NULL").Find(&estabelecimentos).Error; err != nil {
		return
	}

	for _, est := range estabelecimentos {
		if strings.TrimSpace(est.CNPJ) == "" {
			if cat, err := EnsureCategoria(db, int(9999999), "Outros"); err == nil {
				_ = db.Model(&models.Estabelecimento{}).Where("id = ?", est.ID).Update("categoria_id", cat.ID).Error
			}
			continue
		}
		cnae, err := FetchCNAEFiscalByCNPJ(ctx, est.CNPJ)
		if err != nil || cnae == 0 {
			continue
		}
		nome := MapCNAEToCategory(cnae)
		cat, err := EnsureCategoria(db, cnae, nome)
		if err != nil {
			continue
		}
		// atualiza estabelecimento
		_ = db.Model(&models.Estabelecimento{}).Where("id = ?", est.ID).Update("categoria_id", cat.ID).Error
		// pequenas pausas para evitar rate limit
		time.Sleep(100 * time.Millisecond)
	}
}
