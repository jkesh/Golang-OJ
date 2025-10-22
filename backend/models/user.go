package models

import (
	"time"
	
	"gorm.io/gorm"
)

// User 用户实体模型
type User struct {
	gorm.Model
	Username    string    `json:"username" gorm:"unique;not null"`
	Email       string    `json:"email" gorm:"unique;not null"`
	Password    string    `json:"-" gorm:"not null"` // 不在JSON中显示密码
	Nickname    string    `json:"nickname"`
	Avatar      string    `json:"avatar"`
	Role        string    `json:"role" gorm:"default:'user'"` // user, admin
	LastLoginAt time.Time `json:"last_login_at"`
	Submissions []Submission `json:"submissions,omitempty" gorm:"foreignKey:UserID"`
}

// TableName 指定表名
func (User) TableName() string {
	return "users"
}