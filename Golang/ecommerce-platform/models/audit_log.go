package models

import "time"

type AuditLog struct {
    LogID    int       `json:"log_id"`
    Action   string    `json:"action"`
    UserID   int       `json:"user_id"`
    Timestamp time.Time `json:"timestamp"`
}
