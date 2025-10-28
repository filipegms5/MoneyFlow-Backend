package controllers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/filipegms5/MoneyFlow-Backend/models"
	"github.com/filipegms5/MoneyFlow-Backend/repositories"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type FormaPagamentoController struct {
	repo *repositories.FormaPagamentoRepository
}

func NewFormaPagamentoController(db *gorm.DB) *FormaPagamentoController {
	return &FormaPagamentoController{
		repo: repositories.NewFormaPagamentoRepository(db),
	}
}

func (c *FormaPagamentoController) Create(ctx *gin.Context) {
	var formaPagamento models.FormaPagamento
	if err := ctx.ShouldBindJSON(&formaPagamento); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := c.repo.Create(&formaPagamento); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, formaPagamento)
}

func (c *FormaPagamentoController) Update(ctx *gin.Context) {
	var formaPagamento models.FormaPagamento
	if err := ctx.ShouldBindJSON(&formaPagamento); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	idParam := ctx.Param("id")
	var id uint
	if _, err := fmt.Sscan(idParam, &id); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}
	formaPagamento.ID = id

	if err := c.repo.Update(&formaPagamento); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, formaPagamento)
}

func (c *FormaPagamentoController) Delete(ctx *gin.Context) {
	idParam := ctx.Param("id")
	var id uint
	if _, err := fmt.Sscan(idParam, &id); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	if err := c.repo.Delete(id); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.Status(http.StatusNoContent)
}

func (c *FormaPagamentoController) GetByID(ctx *gin.Context) {
	idParam := ctx.Param("id")
	var id uint
	if _, err := fmt.Sscan(idParam, &id); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	formaPagamento, err := c.repo.GetByID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Forma de Pagamento not found"})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	ctx.JSON(http.StatusOK, formaPagamento)
}

func (c *FormaPagamentoController) GetAll(ctx *gin.Context) {
	formasPagamento, err := c.repo.GetAll()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, formasPagamento)
}

func (c *FormaPagamentoController) GetFirstQtd(ctx *gin.Context) {
	qtdStr := ctx.Param("qtd")
	qtd, err := strconv.Atoi(qtdStr)
	if err != nil || qtd <= 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "qtd inválida"})
		return
	}

	formasPagamento, err := c.repo.GetFirst(qtd)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, formasPagamento)
}
