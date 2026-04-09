package repositories

import (
	"time"

	"github.com/filipegms5/MoneyFlow-Backend/models"
	"gorm.io/gorm"
)

type CategoriaGasto struct {
	CategoriaID uint    `json:"categoria_id"`
	Nome        string  `json:"nome"`
	Total       float64 `json:"total"`
}

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

// Agregado: gastos por categoria dos últimos 30 dias incluindo hoje, considerando transação.categoria_id ou estabelecimento.categoria_id
func (r *TransacaoRepository) GetGastosPorCategoriaUltimoMes(usuarioID uint) ([]CategoriaGasto, error) {
	// janela dos últimos 30 dias até hoje (inclusive)
	now := time.Now().UTC()
	startTime := now.AddDate(0, 0, -30)
	start := startTime.Format("2006-01-02")
	end := now.Format("2006-01-02")

	result := make([]CategoriaGasto, 0)
	query := `
	SELECT COALESCE(c.id, 9999999) AS categoria_id, COALESCE(c.nome, 'Outros') AS nome, SUM(t.valor) as total
	FROM transacaos t
	LEFT JOIN estabelecimentos e ON e.id = t.estabelecimento_id
	LEFT JOIN categoria c ON c.id = COALESCE(t.categoria_id, e.categoria_id)
	WHERE t.usuario_id = ? AND t.tipo = 'despesa'
	AND t.data >= ? AND t.data <= ?
	GROUP BY COALESCE(c.id, 9999999), COALESCE(c.nome, 'Outros')
	HAVING total > 0
	ORDER BY total DESC`

	if err := r.db.Raw(query, usuarioID, start, end).Scan(&result).Error; err != nil {
		return nil, err
	}
	return result, nil
}
