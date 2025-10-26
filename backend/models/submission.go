package models

import (
	"time"

	"gorm.io/gorm"
)

// Submission 代码提交实体模型
type Submission struct {
	ID uint `json:"id" gorm:"primaryKey;autoIncrement:true"`
	gorm.Model
	ProblemID    uint      `json:"problem_id" gorm:"not null"`
	UserID       uint      `json:"user_id" gorm:"not null"`
	Language     string    `json:"language" gorm:"not null"`
	Code         string    `json:"code" gorm:"type:text;not null"`
	Status       string    `json:"status" gorm:"default:'pending'"` // pending, judging, accepted, wrong_answer, etc.
	RunTime      int       `json:"run_time"`                        // 毫秒
	Memory       int       `json:"memory"`                          // KB
	SubmittedAt  time.Time `json:"submitted_at" gorm:"autoCreateTime"`
	JudgedAt     time.Time `json:"judged_at"`
	ErrorMessage string    `json:"error_message" gorm:"type:text"`
}

// TableName 指定表名
func (Submission) TableName() string {
	return "submissions"
}
