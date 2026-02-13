package api

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

// getCurrentUserID 从Gin上下文中安全提取当前用户ID。
func getCurrentUserID(c *gin.Context) (uint, error) {
	v, exists := c.Get("userID")
	if !exists {
		return 0, fmt.Errorf("user id missing in context")
	}

	switch id := v.(type) {
	case uint:
		return id, nil
	case int:
		if id < 0 {
			return 0, fmt.Errorf("invalid negative user id")
		}
		return uint(id), nil
	case int64:
		if id < 0 {
			return 0, fmt.Errorf("invalid negative user id")
		}
		return uint(id), nil
	case float64:
		if id < 0 {
			return 0, fmt.Errorf("invalid negative user id")
		}
		return uint(id), nil
	case string:
		var parsed uint
		if _, err := fmt.Sscanf(id, "%d", &parsed); err != nil {
			return 0, fmt.Errorf("invalid user id string")
		}
		return parsed, nil
	default:
		return 0, fmt.Errorf("unsupported user id type %T", v)
	}
}
