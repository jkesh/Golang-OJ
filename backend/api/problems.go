package api

import (
	"net/http"
	"strconv"

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
	db := c.MustGet("db").(*gorm.DB)

	var problem models.Problem
	if err := c.ShouldBindJSON(&problem); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request data",
		})
		return
	}

	// 获取用户ID
	userID, err := getCurrentUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "User not authenticated",
		})
		return
	}

	problem.CreatedBy = userID

	if err := db.Create(&problem).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create problem: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Problem created successfully",
		"problem": problem,
		"status":  "success",
	})
}

// UpdateProblem 更新问题
func UpdateProblem(c *gin.Context) {
	id := c.Param("id")
	db := c.MustGet("db").(*gorm.DB)

	// 查找题目
	var problem models.Problem
	if err := db.First(&problem, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Problem not found",
		})
		return
	}

	// 绑定更新数据
	var updateData models.Problem
	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request data",
		})
		return
	}

	// 更新字段
	problem.Title = updateData.Title
	problem.Description = updateData.Description
	problem.Difficulty = updateData.Difficulty
	problem.TimeLimit = updateData.TimeLimit
	problem.MemoryLimit = updateData.MemoryLimit
	problem.Tags = updateData.Tags

	// 保存更新
	if err := db.Save(&problem).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to update problem: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Problem updated successfully",
		"problem": problem,
		"status":  "success",
	})
}

// DeleteProblem 删除问题
func DeleteProblem(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid problem ID",
		})
		return
	}

	db := c.MustGet("db").(*gorm.DB)

	// 查找题目
	var problem models.Problem
	if err := db.First(&problem, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Problem not found",
		})
		return
	}

	// 删除题目
	if err := db.Delete(&problem).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to delete problem: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Problem deleted successfully",
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
