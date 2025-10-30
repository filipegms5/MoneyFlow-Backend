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
	if err := r.db.Preload("FormaPagamento").Preload("Estabelecimento").Preload("Categoria").First(&t, id).Error; err != nil {
		return nil, err
	}
	return &t, nil
}

func (r *TransacaoRepository) GetAll() ([]models.Transacao, error) {
	var transacoes []models.Transacao
	if err := r.db.Preload("FormaPagamento").Preload("Estabelecimento").Preload("Categoria").Preload("Usuario", func(db *gorm.DB) *gorm.DB { return db.Select("id", "email") }).Find(&transacoes).Error; err != nil {
		return nil, err
	}
	return transacoes, nil
}

func (r *TransacaoRepository) GetByFormaPagamentoID(formaPagamentoID uint) ([]models.Transacao, error) {
	var transacoes []models.Transacao
	if err := r.db.Preload("FormaPagamento").Preload("Estabelecimento").Preload("Categoria").Where("forma_pagamento_id = ?", formaPagamentoID).Find(&transacoes).Error; err != nil {
		return nil, err
	}
	return transacoes, nil
}

func (r *TransacaoRepository) GetByEstabelecimentoID(estabelecimentoID uint) ([]models.Transacao, error) {
	var transacoes []models.Transacao
	if err := r.db.Preload("FormaPagamento").Preload("Estabelecimento").Preload("Categoria").Where("estabelecimento_id = ?", estabelecimentoID).Find(&transacoes).Error; err != nil {
		return nil, err
	}
	return transacoes, nil
}

func (r *TransacaoRepository) GetByTipo(tipo string) ([]models.Transacao, error) {
	var transacoes []models.Transacao
	if err := r.db.Preload("FormaPagamento").Preload("Estabelecimento").Preload("Categoria").Where("tipo = ?", tipo).Find(&transacoes).Error; err != nil {
		return nil, err
	}

	return transacoes, nil
}

func (r *TransacaoRepository) GetByUsuarioID(usuarioID uint) ([]models.Transacao, error) {
	var transacoes []models.Transacao
	if err := r.db.Preload("FormaPagamento").Preload("Estabelecimento").Preload("Categoria").Preload("Usuario", func(db *gorm.DB) *gorm.DB { return db.Select("id", "email") }).Where("usuario_id = ?", usuarioID).Find(&transacoes).Error; err != nil {
		return nil, err
	}
	return transacoes, nil
}

func (r *TransacaoRepository) GetByPeriodo(startISO, endISO string) ([]models.Transacao, error) {
	var transacoes []models.Transacao
	if err := r.db.Preload("FormaPagamento").Preload("Estabelecimento").Preload("Categoria").Preload("Usuario", func(db *gorm.DB) *gorm.DB { return db.Select("id", "email") }).
		Where("data BETWEEN ? AND ?", startISO, endISO).Find(&transacoes).Error; err != nil {
		return nil, err
	}
	return transacoes, nil
}

func (r *TransacaoRepository) GetByPeriodoAndUsuarioID(startISO, endISO string, usuarioID uint) ([]models.Transacao, error) {
	var transacoes []models.Transacao
	if err := r.db.Preload("FormaPagamento").Preload("Estabelecimento").Preload("Categoria").Preload("Usuario", func(db *gorm.DB) *gorm.DB { return db.Select("id", "email") }).
		Where("usuario_id = ? AND data BETWEEN ? AND ?", usuarioID, startISO, endISO).Find(&transacoes).Error; err != nil {
		return nil, err
	}
	return transacoes, nil
}

func (r *TransacaoRepository) GetByPeriodoAndUsuarioIDComRecorrentes(startISO, endISO string, usuarioID uint) ([]models.Transacao, error) {
	var transacoes []models.Transacao
	if err := r.db.Preload("FormaPagamento").Preload("Estabelecimento").Preload("Categoria").Preload("Usuario", func(db *gorm.DB) *gorm.DB { return db.Select("id", "email") }).
		Where("usuario_id = ? AND (data BETWEEN ? AND ? OR recorrente = true)", usuarioID, startISO, endISO).Find(&transacoes).Error; err != nil {
		return nil, err
	}
	return transacoes, nil
}

func (r *TransacaoRepository) GetRecentByUsuarioID(limit int, usuarioID uint) ([]models.Transacao, error) {
	var transacoes []models.Transacao
	if err := r.db.Preload("FormaPagamento").Preload("Estabelecimento").Preload("Categoria").Preload("Usuario", func(db *gorm.DB) *gorm.DB { return db.Select("id", "email") }).
		Where("usuario_id = ?", usuarioID).Order("data DESC").Limit(limit).Find(&transacoes).Error; err != nil {
		return nil, err
	}
	return transacoes, nil
}
