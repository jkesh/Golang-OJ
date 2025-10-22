package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/jkesh/ojplus/api"
	"github.com/jkesh/ojplus/common/database"
	"github.com/jkesh/ojplus/config"
	"github.com/jkesh/ojplus/models"
	"gorm.io/gorm"
)

// OJPlusServer 表示 OJPlus 后端服务
type OJPlusServer struct {
	config      *config.Config
	router      *gin.Engine
	db          *gorm.DB
	redisClient *redis.Client
}

// NewOJPlusServer 创建新的 OJPlus 服务器实例
func NewOJPlusServer() *OJPlusServer {
	return &OJPlusServer{}
}

// Initialize 初始化服务器
func (s *OJPlusServer) Initialize() error {
	fmt.Println("正在启动 OJPlus 在线评测系统...")

	// 加载配置
	cfg, err := config.LoadConfig("./config/config.json")
	if err != nil {
		return fmt.Errorf("加载配置失败: %v", err)
	}
	s.config = cfg

	// 设置 Gin 模式
	gin.SetMode(gin.DebugMode) // 默认使用调试模式
	fmt.Println("运行在调试模式")
	
	// 初始化 Gin 引擎
	s.router = gin.Default()
	
	// 添加 CORS 中间件
	s.router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})
	
	// 添加数据库中间件，将数据库连接注入到上下文中
	s.router.Use(func(c *gin.Context) {
		c.Set("db", s.db)
		c.Next()
	})

	// 初始化数据库连接
	if err := s.initializeDatabase(); err != nil {
		return err
	}

	// 初始化 Redis 连接
	if err := s.initializeRedis(); err != nil {
		return err
	}
	
	// 初始化路由
	s.initializeRouter()

	return nil
}

// initializeRouter 初始化路由
func (s *OJPlusServer) initializeRouter() {
	// 导入API路由
	api.RegisterRoutes(s.router)
}

// initializeDatabase 初始化数据库连接
func (s *OJPlusServer) initializeDatabase() error {
	fmt.Println("正在初始化 PostgreSQL 数据库连接...")
	db, err := database.InitPostgres(s.config)
	if err != nil {
		return fmt.Errorf("初始化 PostgreSQL 失败: %v", err)
	}
	s.db = db
	fmt.Println("PostgreSQL 数据库连接成功")

	// 自动迁移数据库模型
	fmt.Println("正在进行数据库模型迁移...")
	if err := s.migrateDatabase(); err != nil {
		return fmt.Errorf("数据库迁移失败: %v", err)
	}
	fmt.Println("数据库模型迁移成功")

	return nil
}

// migrateDatabase 迁移数据库模型
func (s *OJPlusServer) migrateDatabase() error {
	// 导入模型包
	if err := s.db.AutoMigrate(
		&models.User{},
		&models.Problem{},
		&models.TestCase{},
		&models.Submission{},
	); err != nil {
		return err
	}
	return nil
}

// initializeRedis 初始化 Redis 连接
func (s *OJPlusServer) initializeRedis() error {
	fmt.Println("正在初始化 Redis 连接...")
	redisClient, err := database.InitRedis(s.config)
	if err != nil {
		return fmt.Errorf("初始化 Redis 失败: %v", err)
	}
	s.redisClient = redisClient
	fmt.Println("Redis 连接成功")
	return nil
}

// Run 启动服务器
func (s *OJPlusServer) Run() error {
	fmt.Printf("服务器启动在端口 %d\n", s.config.Server.Port)

	// 在 goroutine 中启动服务器
	serverErr := make(chan error, 1)
	go func() {
		serverErr <- s.router.Run(fmt.Sprintf(":%d", s.config.Server.Port))
	}()

	// 等待中断信号
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	select {
	case err := <-serverErr:
		return fmt.Errorf("服务器启动失败: %v", err)
	case <-quit:
		fmt.Println("正在关闭服务器...")
	}

	// 清理资源
	s.cleanup()

	fmt.Println("服务器已正常退出")
	return nil
}

// cleanup 清理资源
func (s *OJPlusServer) cleanup() {
	// 关闭 Redis 连接
	if s.redisClient != nil {
		if err := s.redisClient.Close(); err != nil {
			log.Printf("关闭 Redis 连接时出错: %v", err)
		} else {
			fmt.Println("Redis 连接已关闭")
		}
	}

	// 关闭数据库连接
	if s.db != nil {
		sqlDB, err := s.db.DB()
		if err != nil {
			log.Printf("获取数据库实例时出错: %v", err)
		} else {
			if err := sqlDB.Close(); err != nil {
				log.Printf("关闭数据库连接时出错: %v", err)
			} else {
				fmt.Println("数据库连接已关闭")
			}
		}
	}
}

func main() {
	// 创建并初始化服务器
	server := NewOJPlusServer()
	if err := server.Initialize(); err != nil {
		log.Fatalf("服务器初始化失败: %v", err)
	}

	// 运行服务器
	if err := server.Run(); err != nil {
		log.Fatalf("服务器运行失败: %v", err)
	}
}
