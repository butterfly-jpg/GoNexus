package model

import (
	"time"

	"gorm.io/gorm"
)

// Session AI会话数据模型
type Session struct {
	ID        string         `gorm:"primaryKey;type:varchar(36)" json:"id"`
	Username  string         `gorm:"index;not null" json:"username"`
	Title     string         `gorm:"type:varchar(100)" json:"title"`
	CreatedAt time.Time      `json:"created_at"` // 自动时间戳
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"_"`
}

// SessionInfo 会话信息
type SessionInfo struct {
	SessionID string    `json:"sessionID"`
	Title     string    `json:"title"`
	CreatedAt time.Time `json:"created_at"`
}
