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

type MetaFinanceiraController struct {
	repo          *repositories.MetaFinanceiraRepository
	transacaoRepo *repositories.TransacaoRepository
}

func NewMetaFinanceiraController(db *gorm.DB) *MetaFinanceiraController {
	return &MetaFinanceiraController{
		repo:          repositories.NewMetaFinanceiraRepository(db),
		transacaoRepo: repositories.NewTransacaoRepository(db),
	}
}

func (ctrl *MetaFinanceiraController) Create(c *gin.Context) {
	var m models.MetaFinanceira
	if err := c.ShouldBindJSON(&m); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if uid, ok := utils.GetUserIDFromContext(c); ok {
		m.UsuarioID = uid
		m.Usuario = nil
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user_id not found in token"})
		return
	}
	if err := ctrl.repo.Create(&m); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	got, _ := ctrl.repo.GetByID(m.ID)
	c.JSON(http.StatusCreated, got)
}

func (ctrl *MetaFinanceiraController) Update(c *gin.Context) {
	var id uint64
	if _, err := fmt.Sscan(c.Param("id"), &id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	var m models.MetaFinanceira
	if err := c.ShouldBindJSON(&m); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	m.ID = uint(id)
	if err := ctrl.repo.Update(&m); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	got, _ := ctrl.repo.GetByID(m.ID)
	c.JSON(http.StatusOK, got)
}

func (ctrl *MetaFinanceiraController) Delete(c *gin.Context) {
	var id uint64
	if _, err := fmt.Sscan(c.Param("id"), &id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	if err := ctrl.repo.Delete(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

func (ctrl *MetaFinanceiraController) GetByID(c *gin.Context) {
	var id uint64
	if _, err := fmt.Sscan(c.Param("id"), &id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	m, err := ctrl.repo.GetByID(uint(id))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "meta not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, m)
}

func (ctrl *MetaFinanceiraController) GetAll(c *gin.Context) {
	metas, err := ctrl.repo.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, metas)
}

func (ctrl *MetaFinanceiraController) GetByUser(c *gin.Context) {
	uid, ok := utils.GetUserIDFromContext(c)
	if !ok || uid == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user_id not found in token"})
		return
	}
	metas, err := ctrl.repo.GetByUsuarioID(uid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, metas)
}

func (ctrl *MetaFinanceiraController) GetTransacoesPeriodo(c *gin.Context) {
	var id uint64
	if _, err := fmt.Sscan(c.Param("id"), &id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	meta, err := ctrl.repo.GetByID(uint(id))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "meta not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	uid, ok := utils.GetUserIDFromContext(c)
	if !ok || uid == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user_id not found in token"})
		return
	}

	txs, err := ctrl.transacaoRepo.GetByPeriodoAndUsuarioID(meta.DataInicio, meta.DataFim, uid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var despesas []models.Transacao
	var receitas []models.Transacao
	for _, t := range txs {
		if t.Tipo == "despesa" {
			despesas = append(despesas, t)
		} else if t.Tipo == "receita" {
			receitas = append(receitas, t)
		}
	}

	c.JSON(http.StatusOK, gin.H{"despesas": despesas, "receitas": receitas})
}
