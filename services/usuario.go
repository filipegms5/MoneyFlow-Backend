package services

import (
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/filipegms5/MoneyFlow-Backend/repositories"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var usuarioRepository *repositories.UsuarioRepository
var tokenBlacklistRepo *repositories.TokenBlacklistRepository

func InitUsuarioService(db *gorm.DB) {
	usuarioRepository = repositories.NewUsuarioRepository(db)
}

func Login(email, senha string) (string, error) {
	usuario, err := usuarioRepository.GetByEmail(email)
	if err != nil {
		return "", err
	}
	err = bcrypt.CompareHashAndPassword([]byte(usuario.Senha), []byte(senha))
	if err != nil {
		return "", err
	}
	// Generate JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": usuario.ID,
		"email":   usuario.Email,
		"exp":     time.Now().Add(time.Hour * 72).Unix(),
	})
	tokenString, err := token.SignedString([]byte("your_secret_key"))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func Logout(jti string, exp time.Time) error {
	return tokenBlacklistRepo.Add(jti, exp)
}

func IsTokenBlacklisted(jti string) (bool, error) {
	return tokenBlacklistRepo.IsRevoked(jti)
}
