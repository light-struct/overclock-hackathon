package domain

import "time"

type Profile struct {
	ID            int64     `json:"id"`
	FullName      string    `json:"full_name"`
	Role          string    `json:"role"`
	PreferredLang string    `json:"preferred_lang"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

