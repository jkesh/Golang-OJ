package api

import (
	"net/http"
	"strconv"

	"backend/models"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// GetCurrentUser 获取当前登录用户信息
func GetCurrentUser(c *gin.Context) {
	userID := c.Param("id")
	db := c.MustGet("db").(*gorm.DB)
	var user models.User
	if err := db.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "用户不存在"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"user": user,
	})
}

// UpdateCurrentUser 更新当前用户信息
func UpdateCurrentUser(c *gin.Context) {
	userID, _ := c.Get("userID")

	var updateData struct {
		Nickname string `json:"nickname"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求参数"})
		return
	}

	db := c.MustGet("db").(*gorm.DB)
	var user models.User
	if err := db.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "用户不存在"})
		return
	}

	// 更新字段
	if updateData.Nickname != "" {
		user.Nickname = updateData.Nickname
	}

	if updateData.Email != "" {
		// 检查邮箱是否已被其他用户使用
		var existingUser models.User
		if err := db.Where("email = ? AND id != ?", updateData.Email, userID).First(&existingUser).Error; err == nil {
			c.JSON(http.StatusConflict, gin.H{"error": "邮箱已被使用"})
			return
		}
		user.Email = updateData.Email
	}

	if updateData.Password != "" {
		// 加密新密码
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(updateData.Password), bcrypt.DefaultCost)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "密码加密失败"})
			return
		}
		user.Password = string(hashedPassword)
	}

	// 保存更新
	if err := db.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新用户信息失败"})
		return
	}

	// 不返回密码
	user.Password = ""

	c.JSON(http.StatusOK, gin.H{
		"message": "用户信息更新成功",
		"user":    user,
	})
}

// GetUsers 获取所有用户（管理员）
func GetUsers(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	var users []models.User

	if err := db.Select("id, username, email, nickname, avatar, role, created_at, last_login_at").Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取用户列表失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"users": users,
	})
}

// GetUser 获取指定用户（管理员）
func GetUser(c *gin.Context) {
	userID := c.Param("id")
	db := c.MustGet("db").(*gorm.DB)
	var user models.User

	if err := db.Select("id, username, email, nickname, avatar, role, created_at, last_login_at").First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "用户不存在"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user": user,
	})
}

// UpdateUser 更新指定用户（管理员）
func UpdateUser(c *gin.Context) {
	userID := c.Param("id")

	// 检查是否尝试更新自己为普通用户
	currentUserID, _ := c.Get("userID")
	if currentUserID == userID {
		// 获取当前用户的角色
		currentRole, _ := c.Get("role")
		if currentRole == "admin" {
			var updateData struct {
				Nickname string `json:"nickname"`
				Email    string `json:"email"`
				Role     string `json:"role"`
			}

			if err := c.ShouldBindJSON(&updateData); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求参数"})
				return
			}

			// 如果试图将自己降级为普通用户，拒绝操作
			if updateData.Role != "" && updateData.Role != "admin" {
				c.JSON(http.StatusForbidden, gin.H{"error": "不能将自己降级为普通用户"})
				return
			}
		}
	}

	var updateData struct {
		Nickname string `json:"nickname"`
		Email    string `json:"email"`
		Role     string `json:"role"`
	}

	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求参数"})
		return
	}

	db := c.MustGet("db").(*gorm.DB)
	var user models.User
	if err := db.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "用户不存在"})
		return
	}

	// 更新字段
	if updateData.Nickname != "" {
		user.Nickname = updateData.Nickname
	}

	if updateData.Email != "" {
		// 检查邮箱是否已被其他用户使用
		var existingUser models.User
		if err := db.Where("email = ? AND id != ?", updateData.Email, userID).First(&existingUser).Error; err == nil {
			c.JSON(http.StatusConflict, gin.H{"error": "邮箱已被使用"})
			return
		}
		user.Email = updateData.Email
	}

	if updateData.Role != "" {
		if updateData.Role != "user" && updateData.Role != "admin" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "无效的角色类型"})
			return
		}
		user.Role = updateData.Role
	}

	// 保存更新
	if err := db.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新用户信息失败"})
		return
	}

	// 不返回密码
	user.Password = ""

	c.JSON(http.StatusOK, gin.H{
		"message": "用户信息更新成功",
		"user":    user,
	})
}

// DeleteUser 删除指定用户（管理员）
func DeleteUser(c *gin.Context) {
	userID := c.Param("id")
	id, err := strconv.Atoi(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的用户ID"})
		return
	}

	// 不允许删除自己
	currentUserID, _ := c.Get("userID")
	currentID, _ := strconv.Atoi(currentUserID.(string))
	if currentID == id {
		c.JSON(http.StatusForbidden, gin.H{"error": "不能删除当前登录的用户"})
		return
	}

	db := c.MustGet("db").(*gorm.DB)
	if err := db.Delete(&models.User{}, userID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除用户失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "用户删除成功",
	})
}

func GetUserSubmitState(c *gin.Context) {
	userID := c.Param("id")
	db := c.MustGet("db").(*gorm.DB)
	var submission []models.Submission
	if err := db.Where("user_id = ?", userID).Order("id").Find(&submission).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "用户提交状态不存在"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"submit_state": submission,
	})
}
