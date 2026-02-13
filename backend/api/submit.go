package api

import (
	"encoding/json"
	"log"
	"net/http"

	"backend/config"
	"backend/models"
	"github.com/IBM/sarama"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// KafkaProducer Kafka生产者实例
var KafkaProducer sarama.SyncProducer

// InitKafka 初始化Kafka生产者
func InitKafka(cfg *config.Config) error {
	kafkaConfig := sarama.NewConfig()
	kafkaConfig.Producer.Return.Successes = true
	kafkaConfig.Producer.Return.Errors = true
	kafkaConfig.ClientID = "oj-submit-service"

	p, err := sarama.NewSyncProducer(cfg.Kafka.Brokers, kafkaConfig)
	if err != nil {
		log.Printf("警告: 初始化Kafka失败: %v", err)
		log.Printf("系统将在无Kafka模式下运行")
		KafkaProducer = nil
		return nil
	}

	KafkaProducer = p

	// 启动事件处理协程
	go handleKafkaEvents()

	return nil
}

// handleKafkaEvents 处理Kafka事件
func handleKafkaEvents() {
	// SyncProducer不需要手动处理事件，这里留空或者可以移除
}

// Submit 提交代码并分配提交ID
func Submit(db *gorm.DB, submission *models.Submission) error {
	// 开启事务
	tx := db.Begin()
	if tx.Error != nil {
		return tx.Error
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 创建提交记录
	// 明确指定不使用已存在的ID值，让数据库自动生成主键
	submission.ID = 0 // 强制清零ID以确保使用自增
	if err := tx.Create(submission).Error; err != nil {
		tx.Rollback()
		return err
	}

	// 如果Kafka可用，则发送消息
	if KafkaProducer != nil {
		// 构造发送到Kafka的消息
		submissionData := map[string]interface{}{
			"submission_id": submission.ID,
			"problem_id":    submission.ProblemID,
			"user_id":       submission.UserID,
			"language":      submission.Language,
			"code":          submission.Code,
			"submitted_at":  submission.SubmittedAt,
		}

		message, err := json.Marshal(submissionData)
		if err != nil {
			tx.Rollback()
			return err
		}

		// 发送到Kafka
		cfg := config.GetConfig()

		// 使用同步发送方式
		msg := &sarama.ProducerMessage{
			Topic: cfg.Kafka.Topic,
			Value: sarama.StringEncoder(message),
		}

		_, _, err = KafkaProducer.SendMessage(msg)
		if err != nil {
			log.Printf("警告: 发送消息到Kafka失败: %v", err)
			// 注意：即使Kafka发送失败，我们仍然提交数据库事务
			// 因为提交记录已经保存到数据库中
		}
	} else {
		log.Println("信息: Kafka不可用，跳过消息发送")
	}

	// 提交事务
	return tx.Commit().Error
}

// SubmitHandler 处理提交请求的HTTP处理函数
func SubmitHandler(c *gin.Context) {
	// 从上下文中获取数据库实例
	db := c.MustGet("db").(*gorm.DB)

	// 解析请求体中的提交数据
	var submission models.Submission
	if err := c.ShouldBindJSON(&submission); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request data",
		})
		return
	}

	// 设置默认状态
	if submission.Status == "" {
		submission.Status = "pending"
	}

	// 从认证中间件中获取用户ID
	userID, err := getCurrentUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "User not authenticated",
		})
		return
	}

	// 设置提交的用户ID
	submission.UserID = userID

	// 调用Submit函数处理提交
	if err := Submit(db, &submission); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to submit code: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Code submitted successfully",
		"status":  "success",
		"id":      submission.ID,
	})
}
