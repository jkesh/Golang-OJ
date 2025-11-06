package api

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"backend/config"
	"backend/models"
	"github.com/IBM/sarama"
	"gorm.io/gorm"
)

// JudgeResult 判题结果结构
type JudgeResult struct {
	SubmissionID uint             `json:"submission_id"`
	ProblemID    uint             `json:"problem_id"`
	Status       string           `json:"status"`   // accepted, wrong_answer, time_limit_exceeded, memory_limit_exceeded, runtime_error, compilation_error
	RunTime      int              `json:"run_time"` // 毫秒
	Memory       int              `json:"memory"`   // KB
	ErrorMessage string           `json:"error_message"`
	TestCases    []TestCaseResult `json:"test_cases"`
}

// TestCaseResult 测试用例结果结构
type TestCaseResult struct {
	TestCaseID     uint   `json:"test_case_id"`
	Input          string `json:"input"`
	ExpectedOutput string `json:"expected_output"`
	UserOutput     string `json:"user_output"`
	Status         string `json:"status"`   // passed, failed, time_limit_exceeded, memory_limit_exceeded, runtime_error
	RunTime        int    `json:"run_time"` // 毫秒
	Memory         int    `json:"memory"`   // KB
	ErrorMessage   string `json:"error_message"`
}

// InitJudgeResultConsumer 初始化判题结果消费者
func InitJudgeResultConsumer(db *gorm.DB) error {
	cfg := config.GetConfig()

	// 创建消费者配置
	kafkaConfig := sarama.NewConfig()
	kafkaConfig.Consumer.Return.Errors = true

	// 创建消费者组
	group, err := sarama.NewConsumerGroup(cfg.Kafka.Brokers, "judge-result-group", kafkaConfig)
	if err != nil {
		return err
	}

	// 启动消费者协程
	go func() {
		defer func() {
			if err := group.Close(); err != nil {
				log.Printf("Error closing consumer group: %v", err)
			}
		}()

		ctx := context.Background()
		for {
			if err := group.Consume(ctx, []string{cfg.Kafka.ResultTopic}, &JudgeResultConsumer{db: db}); err != nil {
				log.Printf("Error consuming judge results: %v", err)
			}
		}
	}()

	return nil
}

// JudgeResultConsumer 判题结果消费者结构
type JudgeResultConsumer struct {
	db *gorm.DB
}

// Setup ConsumerGroupHandler接口实现
func (c *JudgeResultConsumer) Setup(sarama.ConsumerGroupSession) error {
	return nil
}

// Cleanup ConsumerGroupHandler接口实现
func (c *JudgeResultConsumer) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

// ConsumeClaim ConsumerGroupHandler接口实现
func (c *JudgeResultConsumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for message := range claim.Messages() {
		log.Printf("Received judge result: %s", string(message.Value))

		// 解析判题结果
		var result JudgeResult
		if err := json.Unmarshal(message.Value, &result); err != nil {
			log.Printf("Error unmarshaling judge result: %v", err)
			session.MarkMessage(message, "")
			continue
		}

		// 处理判题结果
		if err := c.ProcessJudgeResult(&result); err != nil {
			log.Printf("Error processing judge result: %v", err)
			session.MarkMessage(message, "")
			continue
		}

		// 标记消息已处理
		session.MarkMessage(message, "")
	}

	return nil
}

// ProcessJudgeResult 处理判题结果
func (c *JudgeResultConsumer) ProcessJudgeResult(result *JudgeResult) error {
	// 开启事务
	tx := c.db.Begin()
	if tx.Error != nil {
		return tx.Error
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 更新提交记录
	var submission models.Submission
	if err := tx.First(&submission, result.SubmissionID).Error; err != nil {
		tx.Rollback()
		return err
	}

	// 更新提交状态
	submission.Status = result.Status
	submission.RunTime = result.RunTime
	submission.Memory = result.Memory
	submission.ErrorMessage = result.ErrorMessage
	submission.JudgedAt = time.Now()

	if err := tx.Save(&submission).Error; err != nil {
		tx.Rollback()
		return err
	}

	// 保存测试用例结果
	for _, testCaseResult := range result.TestCases {
		testCaseRes := models.TestCaseResult{
			SubmissionID:   result.SubmissionID,
			TestCaseID:     testCaseResult.TestCaseID,
			Input:          testCaseResult.Input,
			ExpectedOutput: testCaseResult.ExpectedOutput,
			UserOutput:     testCaseResult.UserOutput,
			Status:         testCaseResult.Status,
			RunTime:        testCaseResult.RunTime,
			Memory:         testCaseResult.Memory,
			ErrorMessage:   testCaseResult.ErrorMessage,
		}

		if err := tx.Create(&testCaseRes).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	// 提交事务
	return tx.Commit().Error
}
