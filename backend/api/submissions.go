package api

import (
	"github.com/gin-gonic/gin"
)

// GetSubmissions 获取所有提交列表
func GetSubmissions(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "获取所有提交列表",
		"status": "success",
	})
}

// GetSubmission 获取单个提交详情
func GetSubmission(c *gin.Context) {
	id := c.Param("id")
	c.JSON(200, gin.H{
		"message": "获取提交详情",
		"id": id,
		"status": "success",
	})
}

// SubmitCode 提交代码
func SubmitCode(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "提交代码",
		"status": "success",
	})
}

// GetSubmissionResult 获取提交结果
func GetSubmissionResult(c *gin.Context) {
	id := c.Param("id")
	c.JSON(200, gin.H{
		"message": "获取提交结果",
		"id": id,
		"status": "success",
	})
}

// RegisterSubmissionRoutes 注册提交相关路由
func RegisterSubmissionRoutes(rg *gin.RouterGroup) {
	submissions := rg.Group("/submissions")
	{
		submissions.GET("", GetSubmissions)
		submissions.GET("/:id", GetSubmission)
		submissions.POST("", SubmitCode)
		submissions.GET("/:id/result", GetSubmissionResult)
	}
}