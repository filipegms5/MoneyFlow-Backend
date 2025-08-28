package controllers

import (
	"fmt"
	"net/http"

	"github.com/filipegms5/MoneyFlow-Backend/models"
	"github.com/filipegms5/MoneyFlow-Backend/repositories"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type EstabelecimentoController struct {
	repo *repositories.EstabelecimentoRepository
}

func NewEstabelecimentoController(db *gorm.DB) *EstabelecimentoController {
	return &EstabelecimentoController{
		repo: repositories.NewEstabelecimentoRepository(db),
	}
}

func (c *EstabelecimentoController) Create(ctx *gin.Context) {
	var estabelecimento models.Estabelecimento
	if err := ctx.ShouldBindJSON(&estabelecimento); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := c.repo.Create(&estabelecimento); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, estabelecimento)
}

func (c *EstabelecimentoController) Update(ctx *gin.Context) {
	var estabelecimento models.Estabelecimento
	if err := ctx.ShouldBindJSON(&estabelecimento); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	idParam := ctx.Param("id")
	var id uint
	if _, err := fmt.Sscan(idParam, &id); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}
	estabelecimento.ID = id

	if err := c.repo.Update(&estabelecimento); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, estabelecimento)
}

func (c *EstabelecimentoController) Delete(ctx *gin.Context) {
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

func (c *EstabelecimentoController) GetByID(ctx *gin.Context) {
	idParam := ctx.Param("id")
	var id uint
	if _, err := fmt.Sscan(idParam, &id); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	estabelecimento, err := c.repo.GetByID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Estabelecimento not found"})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	ctx.JSON(http.StatusOK, estabelecimento)
}

func (c *EstabelecimentoController) GetAll(ctx *gin.Context) {
	estabelecimentos, err := c.repo.GetAll()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, estabelecimentos)
}
