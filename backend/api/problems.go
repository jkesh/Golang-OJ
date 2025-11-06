package api

import (
	"net/http"

	"backend/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// GetProblems 获取所有问题列表
func GetProblems(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	var problems []models.Problem
	if err := db.Order("id").Find(&problems).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch problems",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"problems": problems,
	})
}

// GetProblem 获取单个问题详情
func GetProblem(c *gin.Context) {
	id := c.Param("id")
	db := c.MustGet("db").(*gorm.DB)
	var problem models.Problem
	if err := db.First(&problem, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Failed to fetch problem",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"problem": problem,
	})
}

// CreateProblem 创建新问题
func CreateProblem(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "创建新问题",
		"status":  "success",
	})
}

// UpdateProblem 更新问题
func UpdateProblem(c *gin.Context) {
	id := c.Param("id")
	c.JSON(http.StatusOK, gin.H{
		"message": "更新问题",
		"id":      id,
		"status":  "success",
	})
}

// DeleteProblem 删除问题
func DeleteProblem(c *gin.Context) {
	id := c.Param("id")
	c.JSON(http.StatusOK, gin.H{
		"message": "删除问题",
		"id":      id,
		"status":  "success",
	})
}

// 添加测试用例
func addTestCase(c *gin.Context) {
	problemID := c.Param("id")
	c.JSON(http.StatusOK, gin.H{
		"message":    "添加测试用例",
		"problem_id": problemID,
		"status":     "success",
	})
}
