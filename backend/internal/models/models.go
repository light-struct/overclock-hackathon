package models

import "time"

// User represents a registered system user (for auth).
type User struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	Email        string    `gorm:"size:255;uniqueIndex;not null" json:"email"`
	Username     string    `gorm:"size:255;not null" json:"username"`
	PasswordHash string    `gorm:"size:255;not null" json:"-"`
	Role         string    `gorm:"size:50;not null" json:"role"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// Profile represents the Profiles table in the database.
type Profile struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	FullName      string    `gorm:"size:255;not null" json:"full_name"`
	Role          string    `gorm:"size:100;not null" json:"role"`
	PreferredLang string    `gorm:"size:50" json:"preferred_lang"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// TestAttempt represents the TestAttempts table in the database.
type TestAttempt struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	UserID     uint      `gorm:"not null" json:"user_id"`
	Subject    string    `gorm:"size:255;not null" json:"subject"`
	Topic      string    `gorm:"size:255;not null" json:"topic"`
	Score      float64   `json:"score"`
	Language   string    `gorm:"size:50" json:"language"`
	AIFeedback string    `gorm:"type:text" json:"ai_feedback"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

