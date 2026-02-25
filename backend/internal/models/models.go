package models

import "time"

type User struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	Email        string    `gorm:"size:255;uniqueIndex;not null" json:"email"`
	Username     string    `gorm:"size:255;not null" json:"username"`
	PasswordHash string    `gorm:"size:255;not null" json:"-"`
	Role         string    `gorm:"size:50;not null;default:'student'" json:"role"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type Profile struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	UserID        uint      `gorm:"index;not null" json:"user_id"`
	FullName      string    `gorm:"size:255;not null" json:"full_name"`
	Role          string    `gorm:"size:100;not null" json:"role"`
	PreferredLang string    `gorm:"size:50;default:'en'" json:"preferred_lang"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type TestAttempt struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	UserID     uint      `gorm:"index;not null" json:"user_id"`
	Subject    string    `gorm:"size:255;index;not null" json:"subject"`
	Topic      string    `gorm:"size:255;index;not null" json:"topic"`
	Score      float64   `gorm:"index" json:"score"`
	Language   string    `gorm:"size:50" json:"language"`
	AIFeedback string    `gorm:"type:text" json:"ai_feedback"`
	CreatedAt  time.Time `gorm:"index" json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

