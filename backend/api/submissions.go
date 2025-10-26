package api

import (
	"net/http"
	"strconv"

	"backend/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// GetSubmissions 获取所有提交列表 (带分页)
func GetSubmissions(c *gin.Context) {
	// 获取分页参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	// 确保参数有效
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 || pageSize > 100 {
		pageSize = 10
	}

	db := c.MustGet("db").(*gorm.DB)

	var submissions []models.Submission
	var total int64

	// 获取总数
	db.Model(&models.Submission{}).Count(&total)

	// 分页查询
	offset := (page - 1) * pageSize
	result := db.Offset(offset).Limit(pageSize).Find(&submissions)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch submissions",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"submissions": submissions,
		"total":       total,
		"page":        page,
		"page_size":   pageSize,
		"message":     "获取所有提交列表",
		"status":      "success",
	})
}

// GetSubmission 获取单个提交详情
func GetSubmission(c *gin.Context) {
	id := c.Param("id")
	c.JSON(200, gin.H{
		"message": "获取提交详情",
		"id":      id,
		"status":  "success",
	})
}

// GetSubmissionResult 获取提交结果
func GetSubmissionResult(c *gin.Context) {
	id := c.Param("id")
	c.JSON(200, gin.H{
		"message": "获取提交结果",
		"id":      id,
		"status":  "success",
	})
}

// RegisterSubmissionRoutes 注册提交相关路由
func RegisterSubmissionRoutes(rg *gin.RouterGroup) {
	submissions := rg.Group("/submissions")
	{
		submissions.GET("", GetSubmissions)
		submissions.GET("/:id", GetSubmission)
		submissions.GET("/:id/result", GetSubmissionResult)
	}
}
