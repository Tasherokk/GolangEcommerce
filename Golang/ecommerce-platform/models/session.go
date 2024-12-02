package models

import "time"

type Session struct {
    SessionID int       `json:"session_id"`
    UserID    int       `json:"user_id"`
    CreatedAt time.Time `json:"created_at"`
    ExpiresAt time.Time `json:"expires_at"`
}
