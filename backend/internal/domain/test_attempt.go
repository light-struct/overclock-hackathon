package domain

import "time"

type TestAttempt struct {
	ID         int64     `json:"id"`
	UserID     int64     `json:"user_id"`
	Subject    string    `json:"subject"`
	Topic      string    `json:"topic"`
	Score      float64   `json:"score"`
	Language   string    `json:"language"`
	AIFeedback string    `json:"ai_feedback"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

