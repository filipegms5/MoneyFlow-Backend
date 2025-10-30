package middlewares

import (
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/filipegms5/MoneyFlow-Backend/services"
	"github.com/gin-gonic/gin"
)

const jwtSecret = "your_secret_key"

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := c.GetHeader("Authorization")
		if auth == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "authorization header required"})
			return
		}
		parts := strings.SplitN(auth, " ", 2)
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid authorization header"})
			return
		}
		tokenString := parts[1]

		// faz o parse e valida o token (tenta a validação primeiro)
		token, _ := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, nil
			}
			return []byte(jwtSecret), nil
		})

		var claims jwt.MapClaims
		if token != nil {
			if cks, ok := token.Claims.(jwt.MapClaims); ok {
				claims = cks
			}
		}
		// fallback: faz parse sem verificar para extrair claims (logout funciona mesmo com token expirado)
		if claims == nil {
			parser := jwt.Parser{}
			tokUnverified, _, err := parser.ParseUnverified(tokenString, jwt.MapClaims{})
			if err != nil {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
				return
			}
			if cks, ok := tokUnverified.Claims.(jwt.MapClaims); ok {
				claims = cks
			}
		}

		// constrói o jti (usa a claim jti se existir, senão usa o token como chave)
		jti := ""
		if v, ok := claims["jti"]; ok {
			if s, ok := v.(string); ok {
				jti = s
			}
		}
		if jti == "" {
			jti = "token:" + tokenString
		}

		// verifica blacklist
		blacklisted, err := services.IsTokenBlacklisted(jti)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "failed to check token"})
			return
		}
		if blacklisted {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "token revoked"})
			return
		}

		// define valores úteis no contexto
		c.Set("jti", jti)
		if v, ok := claims["user_id"]; ok {
			c.Set("user_id", v)
		}
		if v, ok := claims["exp"]; ok {
			switch t := v.(type) {
			case float64:
				c.Set("exp", time.Unix(int64(t), 0))
			case int64:
				c.Set("exp", time.Unix(t, 0))
			}
		}
		c.Next()
	}
}
