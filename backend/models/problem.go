package models

import (
	"time"
	
	"gorm.io/gorm"
)

// Problem 题目实体模型
type Problem struct {
	gorm.Model
	Title       string    `json:"title" gorm:"not null"`
	Description string    `json:"description" gorm:"type:text"`
	Difficulty  string    `json:"difficulty" gorm:"default:'medium'"`
	TimeLimit   int       `json:"time_limit" gorm:"default:1000"` // 毫秒
	MemoryLimit int       `json:"memory_limit" gorm:"default:256"` // MB
	Tags        string    `json:"tags"`
	CreatedBy   uint      `json:"created_by"`
	UpdatedAt   time.Time `json:"updated_at"`
	TestCases   []TestCase `json:"test_cases,omitempty" gorm:"foreignKey:ProblemID"`
	Submissions []Submission `json:"submissions,omitempty" gorm:"foreignKey:ProblemID"`
}

// TableName 指定表名
func (Problem) TableName() string {
	return "problems"
}