package controllers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/filipegms5/MoneyFlow-Backend/models"
	"github.com/filipegms5/MoneyFlow-Backend/repositories"
	"github.com/filipegms5/MoneyFlow-Backend/services"
	"github.com/filipegms5/MoneyFlow-Backend/utils"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UsuarioController struct {
	repo *repositories.UsuarioRepository
}

func NewUsuarioController(db *gorm.DB) *UsuarioController {
	return &UsuarioController{
		repo: repositories.NewUsuarioRepository(db),
	}
}

func (ctrl *UsuarioController) Create(c *gin.Context) {
	var usuario models.Usuario
	if err := c.ShouldBindJSON(&usuario); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// Hash the password before saving
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(usuario.Senha), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}
	usuario.Senha = string(hashedPassword)

	if err := ctrl.repo.Create(&usuario); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, usuario)
}

func (ctrl *UsuarioController) Update(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}
	var usuario models.Usuario
	if err := c.ShouldBindJSON(&usuario); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	usuario.ID = uint(id)
	if err := ctrl.repo.Update(&usuario); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, usuario)
}

func (ctrl *UsuarioController) Delete(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}
	if err := ctrl.repo.Delete(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

func (ctrl *UsuarioController) Login(c *gin.Context) {
	var loginData struct {
		Email string `json:"email" binding:"required,email"`
		Senha string `json:"senha" binding:"required"`
	}
	if err := c.ShouldBindJSON(&loginData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	token, err := services.Login(loginData.Email, loginData.Senha)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"token": token})
}

func (ctrl *UsuarioController) Logout(c *gin.Context) {
	jti, exists := c.Get("jti")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Token JTI not found"})
		return
	}
	exp, exists := c.Get("exp")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Token expiration not found"})
		return
	}
	if err := services.Logout(jti.(string), exp.(time.Time)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to logout"})
		return
	}
	c.Status(http.StatusOK)
}

func (ctrl *UsuarioController) GetUserID(c *gin.Context) {
	if uid, ok := utils.GetUserIDFromContext(c); ok {
		c.JSON(http.StatusOK, gin.H{"user_id": uid})
		return
	}
	c.JSON(http.StatusUnauthorized, gin.H{"error": "user_id not found in token"})
}
