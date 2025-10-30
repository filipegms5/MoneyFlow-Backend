package controllers

import (
	"fmt"
	"net/http"

	"github.com/filipegms5/MoneyFlow-Backend/models"
	"github.com/filipegms5/MoneyFlow-Backend/repositories"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type CategoriaController struct {
	repo *repositories.CategoriaRepository
}

func NewCategoriaController(db *gorm.DB) *CategoriaController {
	return &CategoriaController{repo: repositories.NewCategoriaRepository(db)}
}

func (c *CategoriaController) Create(ctx *gin.Context) {
	var categoria models.Categoria
	if err := ctx.ShouldBindJSON(&categoria); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := c.repo.Create(&categoria); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, categoria)
}

func (c *CategoriaController) GetAll(ctx *gin.Context) {
	list, err := c.repo.GetAll()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, list)
}

func (c *CategoriaController) GetByID(ctx *gin.Context) {
	var id uint
	if _, err := fmt.Sscan(ctx.Param("id"), &id); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}
	cat, err := c.repo.GetByID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Categoria not found"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, cat)
}

func (c *CategoriaController) Update(ctx *gin.Context) {
	var id uint
	if _, err := fmt.Sscan(ctx.Param("id"), &id); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}
	var categoria models.Categoria
	if err := ctx.ShouldBindJSON(&categoria); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	categoria.ID = id
	if err := c.repo.Update(&categoria); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, categoria)
}

func (c *CategoriaController) Delete(ctx *gin.Context) {
	var id uint
	if _, err := fmt.Sscan(ctx.Param("id"), &id); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}
	if err := c.repo.Delete(id); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.Status(http.StatusNoContent)
}

