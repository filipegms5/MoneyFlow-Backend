package repositories

import (
	"github.com/filipegms5/MoneyFlow-Backend/models"
	"gorm.io/gorm"
)

type EstabelecimentoRepository struct {
	db *gorm.DB
}

func NewEstabelecimentoRepository(db *gorm.DB) *EstabelecimentoRepository {
	return &EstabelecimentoRepository{db: db}
}

func (r *EstabelecimentoRepository) Create(estabelecimento *models.Estabelecimento) error {
	return r.db.Create(estabelecimento).Error
}

func (r *EstabelecimentoRepository) Update(estabelecimento *models.Estabelecimento) error {
	return r.db.Save(estabelecimento).Error
}

func (r *EstabelecimentoRepository) Delete(id uint) error {
	return r.db.Delete(&models.Estabelecimento{}, id).Error
}

func (r *EstabelecimentoRepository) GetAll() ([]models.Estabelecimento, error) {
	var estabelecimentos []models.Estabelecimento
	err := r.db.Find(&estabelecimentos).Error
	return estabelecimentos, err
}

func (r *EstabelecimentoRepository) GetByID(id uint) (*models.Estabelecimento, error) {
	var estabelecimento models.Estabelecimento
	err := r.db.First(&estabelecimento, id).Error
	if err != nil {
		return nil, err
	}
	return &estabelecimento, nil
}

func (r *EstabelecimentoRepository) GetByCnpj(cnpj string) (*models.Estabelecimento, error) {
	var estabelecimento models.Estabelecimento
	err := r.db.Where("cnpj = ?", cnpj).First(&estabelecimento).Error
	if err != nil {
		return nil, err
	}
	return &estabelecimento, nil
}
