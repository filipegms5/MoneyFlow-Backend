package repositories

import (
	"github.com/filipegms5/MoneyFlow-Backend/models"
	"gorm.io/gorm"
)

type FormaPagamentoRepository struct {
	db *gorm.DB
}

func NewFormaPagamentoRepository(db *gorm.DB) *FormaPagamentoRepository {
	return &FormaPagamentoRepository{db: db}
}

func (r *FormaPagamentoRepository) Create(formaPagamento *models.FormaPagamento) error {
	return r.db.Create(formaPagamento).Error
}

func (r *FormaPagamentoRepository) Update(formaPagamento *models.FormaPagamento) error {
	return r.db.Save(formaPagamento).Error
}

func (r *FormaPagamentoRepository) Delete(id uint) error {
	return r.db.Delete(&models.FormaPagamento{}, id).Error
}

func (r *FormaPagamentoRepository) GetAll() ([]models.FormaPagamento, error) {
	var formasPagamento []models.FormaPagamento
	err := r.db.Find(&formasPagamento).Error
	return formasPagamento, err
}

func (r *FormaPagamentoRepository) GetByID(id uint) (*models.FormaPagamento, error) {
	var formaPagamento models.FormaPagamento
	err := r.db.First(&formaPagamento, id).Error
	if err != nil {
		return nil, err
	}
	return &formaPagamento, nil
}

func (r *FormaPagamentoRepository) GetByNome(nome string) (*models.FormaPagamento, error) {
	var formaPagamento models.FormaPagamento
	err := r.db.Where("nome = ?", nome).First(&formaPagamento).Error
	if err != nil {
		return nil, err
	}
	return &formaPagamento, nil
}

func (r *FormaPagamentoRepository) GetFirst(limit int) ([]models.FormaPagamento, error) {
	var formasPagamento []models.FormaPagamento
	err := r.db.Order("id ASC").Limit(limit).Find(&formasPagamento).Error
	return formasPagamento, err
}
