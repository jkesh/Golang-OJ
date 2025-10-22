package models

import (
	"gorm.io/gorm"
)

// SubmitState 记录用户提交状态
type SubmitState struct {
	gorm.Model
	UserID            uint   `json:"user_id" gorm:"not null;index"`
	User              User   `json:"user" gorm:"foreignKey:UserID"`
	SubmittedProblems string `json:"submitted_problems" gorm:"type:text"` // 存储已提交题目ID，用逗号分隔
	AcceptedProblems  string `json:"accepted_problems" gorm:"type:text"`  // 存储已通过题目ID，用逗号分隔
	TotalSubmitted    int    `json:"total_submitted" gorm:"default:0"`    // 总提交次数
	TotalAccepted     int    `json:"total_accepted" gorm:"default:0"`     // 总通过次数
}

// TableName 指定表名
func (SubmitState) TableName() string {
	return "submit_states"
}