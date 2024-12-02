package controllers

import (
    "time"

    "github.com/Tasherokk/GolangProject.git/database"
)


func CreateSession(userID int) (int, error) {
    var sessionID int
    createdAt := time.Now()
    expiresAt := createdAt.Add(24 * time.Hour)

    err := database.DB.QueryRow(`
        INSERT INTO Sessions (user_id, created_at, expires_at)
        VALUES ($1, $2, $3) RETURNING session_id`,
        userID, createdAt, expiresAt).Scan(&sessionID)
    if err != nil {
        return 0, err
    }

    return sessionID, nil
}
