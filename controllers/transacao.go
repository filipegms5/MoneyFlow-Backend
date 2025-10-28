package controllers

import (
	"fmt"
	"net/http"

	"github.com/filipegms5/MoneyFlow-Backend/models"
	"github.com/filipegms5/MoneyFlow-Backend/repositories"
	"github.com/filipegms5/MoneyFlow-Backend/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type TransacaoController struct {
	repo *repositories.TransacaoRepository
}

func NewTransacaoController(db *gorm.DB) *TransacaoController {
	return &TransacaoController{
		repo: repositories.NewTransacaoRepository(db),
	}
}

func (c *TransacaoController) Create(ctx *gin.Context) {
	var transacao models.Transacao
	if err := ctx.ShouldBindJSON(&transacao); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Attach authenticated user to the transacao
	if uid, ok := utils.GetUserIDFromContext(ctx); ok {
		transacao.UsuarioID = uid
		transacao.Usuario = nil
	} else {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "user_id not found in token"})
		return
	}

	// Accept nested { "forma_pagamento": { "id": 1 } } by copying ID and niling nested object
	if transacao.FormaPagamento != nil && transacao.FormaPagamento.ID != 0 {
		transacao.FormaPagamentoID = transacao.FormaPagamento.ID
		transacao.FormaPagamento = nil
	}
	if transacao.Estabelecimento != nil && transacao.Estabelecimento.ID != 0 {
		transacao.EstabelecimentoID = transacao.Estabelecimento.ID
		transacao.Estabelecimento = nil
	}

	if err := c.repo.Create(&transacao); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	created, _ := c.repo.GetByID(transacao.ID) // repo should preload associations
	ctx.JSON(http.StatusCreated, created)
}

func (c *TransacaoController) Update(ctx *gin.Context) {
	var id uint64
	if _, err := fmt.Sscan(ctx.Param("id"), &id); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var transacao models.Transacao
	if err := ctx.ShouldBindJSON(&transacao); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	transacao.ID = uint(id)

	if transacao.FormaPagamento != nil && transacao.FormaPagamento.ID != 0 {
		transacao.FormaPagamentoID = transacao.FormaPagamento.ID
		transacao.FormaPagamento = nil
	}
	if transacao.Estabelecimento != nil && transacao.Estabelecimento.ID != 0 {
		transacao.EstabelecimentoID = transacao.Estabelecimento.ID
		transacao.Estabelecimento = nil
	}

	if err := c.repo.Update(&transacao); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	updated, _ := c.repo.GetByID(transacao.ID)
	ctx.JSON(http.StatusOK, updated)
}

func (c *TransacaoController) GetAll(ctx *gin.Context) {
	transacoes, err := c.repo.GetAll()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, transacoes)
}

func (c *TransacaoController) GetByID(ctx *gin.Context) {
	var id uint64
	if _, err := fmt.Sscan(ctx.Param("id"), &id); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	transacao, err := c.repo.GetByID(uint(id))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "transacao not found"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, transacao)
}

func (c *TransacaoController) Delete(ctx *gin.Context) {
	var id uint64
	if _, err := fmt.Sscan(ctx.Param("id"), &id); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	if err := c.repo.Delete(uint(id)); err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "transacao not found"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.Status(http.StatusNoContent)
}

func (c *TransacaoController) GetByTipo(ctx *gin.Context) {
	tipo := ctx.Param("tipo")
	transacoes, err := c.repo.GetByTipo(tipo)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, transacoes)

}

func (c *TransacaoController) GetByPeriodo(ctx *gin.Context) {
	start := ctx.Query("start")
	end := ctx.Query("end")
	if start == "" || end == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "start and end are required (ISO strings)"})
		return
	}
	// Only return transactions for the authenticated user
	uid, ok := utils.GetUserIDFromContext(ctx)
	if !ok || uid == 0 {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "user_id not found in token"})
		return
	}
	transacoes, err := c.repo.GetByPeriodoAndUsuarioIDComRecorrentes(start, end, uid)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, transacoes)
}

func (c *TransacaoController) GetByUserID(ctx *gin.Context) {
	uid, ok := utils.GetUserIDFromContext(ctx)
	if !ok || uid == 0 {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "user_id not found in token"})
		return
	}
	transacoes, err := c.repo.GetByUsuarioID(uid)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, transacoes)
}
