package models

import (
	"gorm.io/gorm"
)

// TestCaseResult 测试用例结果实体模型
type TestCaseResult struct {
	gorm.Model
	SubmissionID   uint   `json:"submission_id" gorm:"not null;index"`
	TestCaseID     uint   `json:"test_case_id" gorm:"not null"`
	Input          string `json:"input" gorm:"type:text"`
	ExpectedOutput string `json:"expected_output" gorm:"type:text"`
	UserOutput     string `json:"user_output" gorm:"type:text"`
	Status         string `json:"status"`   // passed, failed, time_limit_exceeded, memory_limit_exceeded, runtime_error
	RunTime        int    `json:"run_time"` // 毫秒
	Memory         int    `json:"memory"`   // KB
	ErrorMessage   string `json:"error_message" gorm:"type:text"`
}

// TableName 指定表名
func (TestCaseResult) TableName() string {
	return "test_case_results"
}
