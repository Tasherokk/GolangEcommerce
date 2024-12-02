package controllers

import (
    "encoding/json"
    "net/http"
    "strconv"
	"database/sql"

    "github.com/gorilla/mux"
    "github.com/Tasherokk/GolangProject.git/database"
	"github.com/Tasherokk/GolangProject.git/models"
)


// Получение всех ролей
func GetAllRoles(w http.ResponseWriter, r *http.Request) {
    rows, err := database.DB.Query("SELECT role_id, role_name FROM Roles")
    if err != nil {
        http.Error(w, "Error fetching roles", http.StatusInternalServerError)
        return
    }
    defer rows.Close()

    var roles []models.Role
    for rows.Next() {
        var role models.Role
        err := rows.Scan(&role.RoleID, &role.RoleName)
        if err != nil {
            http.Error(w, "Error scanning role", http.StatusInternalServerError)
            return
        }
        roles = append(roles, role)
    }

    json.NewEncoder(w).Encode(roles)
}

// Создание новой роли
func CreateRole(w http.ResponseWriter, r *http.Request) {
    var role models.Role
    err := json.NewDecoder(r.Body).Decode(&role)
    if err != nil {
        http.Error(w, "Invalid input", http.StatusBadRequest)
        return
    }

    err = database.DB.QueryRow(
        "INSERT INTO Roles (role_name) VALUES ($1) RETURNING role_id",
        role.RoleName).Scan(&role.RoleID)
    if err != nil {
        http.Error(w, "Error creating role", http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(role)
}

// Обновление роли пользователя
func UpdateUserRole(w http.ResponseWriter, r *http.Request) {
    // Получение user_id из URL
    vars := mux.Vars(r)
    userIDStr := vars["user_id"]
    userID, err := strconv.Atoi(userIDStr)
    if err != nil {
        http.Error(w, "Invalid user ID", http.StatusBadRequest)
        return
    }

    // Декодирование тела запроса
    var roleUpdate struct {
        Role string `json:"role_name"`
    }
    err = json.NewDecoder(r.Body).Decode(&roleUpdate)
    if err != nil {
        http.Error(w, "Invalid input", http.StatusBadRequest)
        return
    }

    // Получение RoleID по названию роли
    var roleID int
    err = database.DB.QueryRow("SELECT role_id FROM Roles WHERE role_name = $1", roleUpdate.Role).Scan(&roleID)
    if err != nil {
        if err == sql.ErrNoRows {
            http.Error(w, "Invalid role", http.StatusBadRequest)
        } else {
            http.Error(w, "Database error", http.StatusInternalServerError)
        }
        return
    }

    // Обновление роли пользователя в базе данных
    _, err = database.DB.Exec("UPDATE Users SET role=$1 WHERE user_id=$2", roleID, userID)
    if err != nil {
        http.Error(w, "Error updating user role", http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(map[string]string{"message": "User role updated"})
}