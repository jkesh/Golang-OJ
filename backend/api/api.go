package api

import (
	"backend/middleware"
	"github.com/gin-gonic/gin"
)

// RegisterRoutes 注册所有API路由
func RegisterRoutes(router *gin.Engine) {
	// API路由组
	apiGroup := router.Group("/api")

	// 认证相关路由 - 无需认证
	apiGroup.POST("/auth/login", Login)
	apiGroup.POST("/auth/register", Register)

	// 需要认证的路由
	authRequired := apiGroup.Group("/")
	authRequired.Use(middleware.JWT())
	{
		// 用户相关
		authRequired.GET("/user/:id", GetCurrentUser)
		authRequired.PUT("/user/me", UpdateCurrentUser)

		// 问题相关
		authRequired.GET("/problems", GetProblems)
		authRequired.GET("/problems/:id", GetProblem)
		authRequired.POST("/problems", CreateProblem)
		authRequired.PUT("/problems/:id", UpdateProblem)
		authRequired.DELETE("/problems/:id", DeleteProblem)
		authRequired.POST("/problems/:id/testcases", addTestCase)

		// 提交相关
		authRequired.POST("/submit", SubmitHandler)
		authRequired.GET("/submissions", GetSubmissions)
		authRequired.GET("/submissions/:id", GetSubmission)
		authRequired.GET("/submissions/:id/result", GetSubmissionResult)
		authRequired.GET("/:id/submit-state", GetUserSubmitState)

		// 管理员路由
		adminRequired := authRequired.Group("/admin")
		adminRequired.Use(middleware.AdminRequired())
		{
			// 用户管理
			adminRequired.GET("/users", GetUsers)
			adminRequired.GET("/users/:id", GetUser)
			adminRequired.PUT("/users/:id", UpdateUser)
			adminRequired.DELETE("/users/:id", DeleteUser)
		}
	}
}
