package models

import (
	"gorm.io/gorm"
)

// TestCase 测试用例实体模型
type TestCase struct {
	gorm.Model
	ProblemID  uint   `json:"problem_id" gorm:"not null"`
	Input      string `json:"input" gorm:"type:text"`
	Output     string `json:"output" gorm:"type:text"`
	IsExample  bool   `json:"is_example" gorm:"default:false"` // 是否为示例测试用例
	IsHidden   bool   `json:"is_hidden" gorm:"default:false"`  // 是否为隐藏测试用例
	Weight     int    `json:"weight" gorm:"default:1"`         // 测试用例权重
}

// TableName 指定表名
func (TestCase) TableName() string {
	return "test_cases"
}