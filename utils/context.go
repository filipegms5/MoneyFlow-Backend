package utils

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetUserIDFromContext(c *gin.Context) (uint, bool) {
	userIDVal, exists := c.Get("user_id")
	if !exists {
		return 0, false
	}
	switch v := userIDVal.(type) {
	case float64:
		return uint(v), true
	case int:
		return uint(v), true
	case int64:
		return uint(v), true
	case uint:
		return v, true
	case string:
		if id, err := strconv.Atoi(v); err == nil {
			return uint(id), true
		}
	}
	return 0, false
}
