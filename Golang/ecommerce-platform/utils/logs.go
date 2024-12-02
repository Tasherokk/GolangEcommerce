package utils

import (
    "github.com/Tasherokk/GolangProject.git/database"
    "time"
)

func LogAction(userID int, action string) error {
    timestamp := time.Now()
    _, err := database.DB.Exec(
        "INSERT INTO AuditLogs (action, user_id, timestamp) VALUES ($1, $2, $3)",
        action, userID, timestamp)
    return err
}
