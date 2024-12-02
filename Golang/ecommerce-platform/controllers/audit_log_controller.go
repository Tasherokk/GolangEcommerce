package controllers

import (
    "encoding/json"
    "net/http"

    "github.com/Tasherokk/GolangProject.git/database"
	"github.com/Tasherokk/GolangProject.git/models"
)

func GetAuditLogs(w http.ResponseWriter, r *http.Request) {
    role := r.Context().Value("role").(string)
    if role != "admin" {
        http.Error(w, "Access denied", http.StatusForbidden)
        return
    }

    rows, err := database.DB.Query("SELECT log_id, action, user_id, timestamp FROM AuditLogs")
    if err != nil {
        http.Error(w, "Error fetching audit logs", http.StatusInternalServerError)
        return
    }
    defer rows.Close()

    var logs []models.AuditLog
    for rows.Next() {
        var log models.AuditLog
        err := rows.Scan(&log.LogID, &log.Action, &log.UserID, &log.Timestamp)
        if err != nil {
            http.Error(w, "Error scanning audit log", http.StatusInternalServerError)
            return
        }
        logs = append(logs, log)
    }

    json.NewEncoder(w).Encode(logs)
}
