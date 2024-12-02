package models

import "time"

type User struct {
    UserID       int       `json:"user_id"`
    Username     string    `json:"username" validate:"required,min=3"`
    PasswordHash string    `json:"password"` // Используется только при регистрации
    Email        string    `json:"email" validate:"required,email"`
    CreatedAt    time.Time `json:"created_at"`
    Role         int       `json:"role"`
    RoleName     string    `json:"role_name"`
}
