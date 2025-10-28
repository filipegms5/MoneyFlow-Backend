package repositories

import (
	"github.com/filipegms5/MoneyFlow-Backend/models"
	"gorm.io/gorm"
)

type TransacaoRepository struct {
	db *gorm.DB
}

func NewTransacaoRepository(db *gorm.DB) *TransacaoRepository {
	return &TransacaoRepository{db: db}
}

func (r *TransacaoRepository) Create(t *models.Transacao) error {
	return r.db.Create(t).Error
}

func (r *TransacaoRepository) Update(t *models.Transacao) error {
	return r.db.Save(t).Error
}

func (r *TransacaoRepository) Delete(id uint) error {
	return r.db.Delete(&models.Transacao{}, id).Error
}
func (r *TransacaoRepository) GetByID(id uint) (*models.Transacao, error) {
	var t models.Transacao
	if err := r.db.Preload("FormaPagamento").Preload("Estabelecimento").First(&t, id).Error; err != nil {
		return nil, err
	}
	return &t, nil
}

func (r *TransacaoRepository) GetAll() ([]models.Transacao, error) {
	var transacoes []models.Transacao
	if err := r.db.Preload("FormaPagamento").Preload("Estabelecimento").Find(&transacoes).Error; err != nil {
		return nil, err
	}
	return transacoes, nil
}

func (r *TransacaoRepository) GetByFormaPagamentoID(formaPagamentoID uint) ([]models.Transacao, error) {
	var transacoes []models.Transacao
	if err := r.db.Preload("FormaPagamento").Preload("Estabelecimento").Where("forma_pagamento_id = ?", formaPagamentoID).Find(&transacoes).Error; err != nil {
		return nil, err
	}
	return transacoes, nil
}

func (r *TransacaoRepository) GetByEstabelecimentoID(estabelecimentoID uint) ([]models.Transacao, error) {
	var transacoes []models.Transacao
	if err := r.db.Preload("FormaPagamento").Preload("Estabelecimento").Where("estabelecimento_id = ?", estabelecimentoID).Find(&transacoes).Error; err != nil {
		return nil, err
	}
	return transacoes, nil
}

func (r *TransacaoRepository) GetByTipo(tipo string) ([]models.Transacao, error) {
	var transacoes []models.Transacao
	if err := r.db.Preload("FormaPagamento").Preload("Estabelecimento").Where("tipo = ?", tipo).Find(&transacoes).Error; err != nil {
		return nil, err
	}

	return transacoes, nil
}
