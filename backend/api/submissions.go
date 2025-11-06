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

	// 分页查询，按提交时间倒序排列（最新的在最前）
	offset := (page - 1) * pageSize
	result := db.Offset(offset).Limit(pageSize).Order("submitted_at DESC").Find(&submissions)

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

	db := c.MustGet("db").(*gorm.DB)

	var submission models.Submission
	if err := db.First(&submission, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Submission not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"submission": submission,
		"message":    "获取提交详情",
		"status":     "success",
	})
}

// GetSubmissionResult 获取提交结果
func GetSubmissionResult(c *gin.Context) {
	id := c.Param("id")

	db := c.MustGet("db").(*gorm.DB)

	var submission models.Submission
	if err := db.First(&submission, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Submission not found",
		})
		return
	}

	// 获取题目信息
	var problem models.Problem
	if err := db.First(&problem, submission.ProblemID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Problem not found",
		})
		return
	}

	// 获取测试用例结果
	var testCaseResults []models.TestCaseResult
	if err := db.Where("submission_id = ?", id).Find(&testCaseResults).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch test case results",
		})
		return
	}

	// 如果没有测试用例结果，获取原始测试用例
	if len(testCaseResults) == 0 {
		var testCases []models.TestCase
		if err := db.Where("problem_id = ?", submission.ProblemID).Find(&testCases).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to fetch test cases",
			})
			return
		}

		// 转换为测试用例结果格式
		for _, tc := range testCases {
			testCaseResults = append(testCaseResults, models.TestCaseResult{
				TestCaseID:     tc.ID,
				Input:          tc.Input,
				ExpectedOutput: tc.Output,
			})
		}
	}

	// 构造返回结果
	result := map[string]interface{}{
		"submission": submission,
		"problem":    problem,
		"test_cases": testCaseResults,
	}

	c.JSON(http.StatusOK, gin.H{
		"data":    result,
		"message": "获取提交结果",
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
