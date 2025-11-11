package main

import (
	"fmt"
	"log"

	"backend/api"
	"backend/common/database"
	"backend/config"
	"backend/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func main() {
	// 加载配置
	cfg, err := config.LoadConfig("./config.json")
	if err != nil {
		log.Fatalf("加载配置失败: %v", err)
	}

	// 初始化数据库
	db, err := database.InitPostgres(cfg)
	if err != nil {
		log.Fatalf("初始化数据库失败: %v", err)
	}

	// 自动迁移数据库模型
	if err := db.AutoMigrate(
		&models.User{},
		&models.Problem{},
		&models.Submission{},
		&models.TestCase{},
		&models.TestCaseResult{},
		&models.Judge{},
	); err != nil {
		log.Fatalf("数据库迁移失败: %v", err)
	}

	// 重置 submissions 表的自增序列
	db.Exec("SELECT setval(pg_get_serial_sequence('submissions', 'id'), COALESCE(MAX(id), 0) + 1, false) FROM submissions")

	// 初始化Kafka
	if err := api.InitKafka(cfg); err != nil {
		log.Printf("初始化Kafka失败: %v", err)
		log.Printf("系统将在无Kafka模式下运行")
	}

	// 初始化判题结果消费者
	if err := api.InitJudgeResultConsumer(db); err != nil {
		log.Printf("初始化判题结果消费者失败: %v", err)
		log.Printf("系统将在无判题结果消费模式下运行")
	}

	// 初始化路由
	r := initRouter(db)

	// 启动服务器
	addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	log.Printf("服务器启动在 %s", addr)
	if err := r.Run(addr); err != nil {
		log.Fatalf("服务器启动失败: %v", err)
	}
}

// initRouter 初始化路由
func initRouter(db *gorm.DB) *gin.Engine {
	r := gin.Default()

	// 添加数据库连接中间件
	r.Use(func(c *gin.Context) {
		c.Set("db", db)
		c.Next()
	})

	// 添加CORS中间件
	r.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Authorization, X-Requested-With")
		c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Cache-Control, Content-Language, Content-Type")
		c.Header("Access-Control-Allow-Credentials", "true")

		// 如果是OPTIONS请求，直接返回204
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// 注册路由
	api.RegisterRoutes(r)

	return r
}
