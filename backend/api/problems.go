package api

import (
	"github.com/gin-gonic/gin"
)

// GetProblems 获取所有问题列表
func GetProblems(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "获取所有问题列表",
		"status":  "success",
	})
}

// GetProblem 获取单个问题详情
func GetProblem(c *gin.Context) {
	id := c.Param("id")
	c.JSON(200, gin.H{
		"message": "获取问题详情",
		"id":      id,
		"status":  "success",
	})
}

// CreateProblem 创建新问题
func CreateProblem(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "创建新问题",
		"status":  "success",
	})
}

// UpdateProblem 更新问题
func UpdateProblem(c *gin.Context) {
	id := c.Param("id")
	c.JSON(200, gin.H{
		"message": "更新问题",
		"id":      id,
		"status":  "success",
	})
}

// DeleteProblem 删除问题
func DeleteProblem(c *gin.Context) {
	id := c.Param("id")
	c.JSON(200, gin.H{
		"message": "删除问题",
		"id":      id,
		"status":  "success",
	})
}

// 添加测试用例
func addTestCase(c *gin.Context) {
	problemID := c.Param("id")
	c.JSON(200, gin.H{
		"message":    "添加测试用例",
		"problem_id": problemID,
		"status":     "success",
	})
}
