package controllers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/Tasherokk/GolangProject.git/database"
	"github.com/Tasherokk/GolangProject.git/middleware"
	"github.com/Tasherokk/GolangProject.git/models"
	"github.com/Tasherokk/GolangProject.git/utils"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
)


var user_validator = validator.New()

var allowedRoles = []string{"Guest", "Seller", "Customer"}

func GetUserHandler(w http.ResponseWriter, r *http.Request) {

    log.Printf("Request received: %s %s\n", r.Method, r.URL.Path)

    vars := mux.Vars(r)
    idStr := vars["id"]
    id, err := strconv.Atoi(idStr)
    if err != nil {
        http.Error(w, "Invalid user ID", http.StatusBadRequest)
        return
    }

    user, err := GetUserByID(id)
    if err != nil {
        http.Error(w, "User not found", http.StatusNotFound)
        return
    }

    log.Println(user)
    response, err := json.Marshal(user)
    if err != nil {
        log.Fatalf("Failed to marshal user: %v", err)
    }
    log.Println("Response JSON:", string(response))
    json.NewEncoder(w).Encode(user)
}

func GetAllUsers(w http.ResponseWriter, r *http.Request) {

    rows, err := database.DB.Query("SELECT user_id, username, email, created_at, role FROM users")
    if err != nil {
        http.Error(w, "Error fetching users: "+err.Error(), http.StatusInternalServerError)
        return
    }
    defer rows.Close()

    var users []models.User
    for rows.Next() {
        var user models.User
        if err := rows.Scan(&user.UserID, &user.Username, &user.Email, &user.CreatedAt, &user.Role); err != nil {
            http.Error(w, "Error scanning user: "+err.Error(), http.StatusInternalServerError)
            return
        }
        err := database.DB.QueryRow("SELECT role_name FROM Roles WHERE role_id = $1", user.Role).Scan(&user.RoleName)
        if err != nil {
            http.Error(w, "Error fetching role name: "+err.Error(), http.StatusInternalServerError)
            return
        }
        users = append(users, user)
    }

    json.NewEncoder(w).Encode(users)


}

// Регистрация пользователя
func RegisterUser(w http.ResponseWriter, r *http.Request) {
    var user models.User
    log.Println("Starting user registration...")

    err := json.NewDecoder(r.Body).Decode(&user)
    if err != nil {
        log.Println("Invalid input:", err)
        http.Error(w, "Invalid input", http.StatusBadRequest)
        return
    }
    log.Printf("Decoded user input: %+v\n", user)

    // Валидация данных пользователя
    err = user_validator.Struct(user)
    if err != nil {
        log.Println("Validation error:", err)
        http.Error(w, "Validation error: "+err.Error(), http.StatusBadRequest)
        return
    }
    log.Println("User input validated successfully")

    // Проверка роли пользователя
    roleName := user.RoleName
    if roleName == "" {
        roleName = "Customer" // Роль по умолчанию
    }

    validRole := false
    for _, role := range allowedRoles {
        if role == roleName {
            validRole = true
            break
        }
    }
    if !validRole {
        log.Println("Invalid role:", roleName)
        http.Error(w, "Invalid role", http.StatusBadRequest)
        return
    }
    log.Printf("Role '%s' is valid\n", roleName)

    // Получаем RoleID из таблицы Roles
    var roleID int
    err = database.DB.QueryRow("SELECT role_id FROM Roles WHERE role_name = $1", roleName).Scan(&roleID)
    if err != nil {
        if err == sql.ErrNoRows {
            log.Println("Role not found in database:", roleName)
            http.Error(w, "Role not found", http.StatusBadRequest)
        } else {
            log.Println("Database error while fetching role ID:", err)
            http.Error(w, "Database error", http.StatusInternalServerError)
        }
        return
    }
    log.Printf("Role ID for '%s': %d\n", roleName, roleID)

    // Хешируем пароль
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.PasswordHash), bcrypt.DefaultCost)
    if err != nil {
        log.Println("Error hashing password:", err)
        http.Error(w, "Error hashing password", http.StatusInternalServerError)
        return
    }
    user.PasswordHash = string(hashedPassword)
    user.CreatedAt = time.Now()
    user.Role = roleID

    log.Printf("User data prepared for insertion: %+v\n", user)

    // Вставляем пользователя в базу данных
    var userID int
    err = database.DB.QueryRow(
        "INSERT INTO Users (username, password_hash, email, created_at, role) VALUES ($1, $2, $3, $4, $5) RETURNING user_id",
        user.Username, user.PasswordHash, user.Email, user.CreatedAt, user.Role).Scan(&userID)
    if err != nil {
        log.Println("Error inserting user into database:", err)
        http.Error(w, "Error inserting user: "+err.Error(), http.StatusInternalServerError)
        return
    }
    log.Println("User inserted into database successfully")

    user.UserID = userID

    // Генерация JWT с использованием RoleName
    token, err := utils.GenerateJWT(user.UserID, roleName)
    if err != nil {
        log.Println("Error generating JWT:", err)
        http.Error(w, "Error generating token", http.StatusInternalServerError)
        return
    }

    // Отправка токена и роли клиенту
    response := struct {
        Token string       `json:"token"`
        User  models.User  `json:"user"`
    }{
        Token: token,
        User:  user,
    }

    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(response)
}




// Вход пользователя
func LoginUser(w http.ResponseWriter, r *http.Request) {
    var credentials struct {
        Username string `json:"username"`
        Password string `json:"password"`
    }
    err := json.NewDecoder(r.Body).Decode(&credentials)
    if err != nil {
        http.Error(w, "Invalid input", http.StatusBadRequest)
        return
    }

    // Получение пользователя из базы данных
    var user models.User
    err = database.DB.QueryRow(
        "SELECT user_id, password_hash, role FROM Users WHERE username=$1",
        credentials.Username).Scan(&user.UserID, &user.PasswordHash, &user.Role)
    if err == sql.ErrNoRows {
        http.Error(w, "User not found", http.StatusUnauthorized)
        return
    } else if err != nil {
        http.Error(w, "Database error", http.StatusInternalServerError)
        return
    }

    err = database.DB.QueryRow("SELECT role_name FROM Roles WHERE role_id = $1", user.Role).Scan(&user.RoleName)
    if err != nil {
        http.Error(w, "Error fetching role", http.StatusInternalServerError)
        return
    }

    // Проверка пароля
    err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(credentials.Password))
    if err != nil {
        http.Error(w, "Invalid password", http.StatusUnauthorized)
        return
    }

    // Генерация JWT с использованием RoleName
    token, err := utils.GenerateJWT(user.UserID, user.RoleName)
    if err != nil {
        http.Error(w, "Error generating token", http.StatusInternalServerError)
        return
    }

    // Отправка токена и роли клиенту
    response := struct {
        Token    string `json:"token"`
        User     models.User `json:"user"`
    }{
        Token:    token,
        User:     user,
    }

    json.NewEncoder(w).Encode(response)
}

// Получение текущего пользователя
func GetCurrentUser(w http.ResponseWriter, r *http.Request) {
    userID, ok := r.Context().Value(middleware.UserIDKey).(int)

    if !ok {
        http.Error(w, "User ID not found in context", http.StatusUnauthorized)
        return
    }

    var user models.User
    err := database.DB.QueryRow(
        "SELECT user_id, username, email, created_at, role FROM Users WHERE user_id=$1",
        userID).Scan(&user.UserID, &user.Username, &user.Email, &user.CreatedAt, &user.Role)
    if err != nil {
        http.Error(w, "User not found", http.StatusNotFound)
        return
    }
    err = database.DB.QueryRow("SELECT role_name FROM Roles WHERE role_id = $1", user.Role).Scan(&user.RoleName)
    if err != nil {
        http.Error(w, "Error fetching role", http.StatusInternalServerError)
        return
    }


    json.NewEncoder(w).Encode(user)
}

// Обновление пользователя
func UpdateUser(w http.ResponseWriter, r *http.Request) {
    userID, ok := r.Context().Value(middleware.UserIDKey).(int)

    if !ok {
        http.Error(w, "User ID not found in context", http.StatusUnauthorized)
        return
    }

    var user models.User
    err := json.NewDecoder(r.Body).Decode(&user)
    if err != nil {
        http.Error(w, "Invalid input", http.StatusBadRequest)
        return
    }

    _, err = database.DB.Exec(
        "UPDATE Users SET username=$1, email=$2 WHERE user_id=$3",
        user.Username, user.Email, userID)
    if err != nil {
        http.Error(w, "Error updating user", http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusOK)
}

// Удаление пользователя
func DeleteUser(w http.ResponseWriter, r *http.Request) {
    userID, ok := r.Context().Value(middleware.UserIDKey).(int)

    if !ok {
        http.Error(w, "User ID not found in context", http.StatusUnauthorized)
        return
    }

    _, err := database.DB.Exec("DELETE FROM Users WHERE user_id=$1", userID)
    if err != nil {
        http.Error(w, "Error deleting user", http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusNoContent)
}



func DeleteUserByID(w http.ResponseWriter, r *http.Request) {
    // Извлечение user_id из URL
    vars := mux.Vars(r)
    userID, err := strconv.Atoi(vars["user_id"])
    if err != nil {
        http.Error(w, "Invalid user ID", http.StatusBadRequest)
        return
    }

    // Удаление пользователя из базы данных
    _, err = database.DB.Exec("DELETE FROM Users WHERE user_id=$1", userID)
    if err != nil {
        log.Printf("Error deleting user: %v", err)
        http.Error(w, "Error deleting user", http.StatusInternalServerError)
        return
    }

    // Возвращаем статус 204 (No Content)
    w.WriteHeader(http.StatusNoContent)
}




func CreateUser(newUser *models.User) error {
	// Check if user exists
	exists, err := UserExists(newUser.Username, newUser.Email)
	if err != nil {
		return fmt.Errorf("error checking user existence: %w", err)
	}
	if exists {
		return fmt.Errorf("user already exists with username or email")
	}

	// Insert new user into the database
	sqlStatement := `
		INSERT INTO users (username, email, password_hash, role)
		VALUES ($1, $2, $3, $4)
		RETURNING user_id, username, email
	`
	err = database.DB.QueryRow(sqlStatement, newUser.Username, newUser.Email, newUser.PasswordHash, 1).Scan(&newUser.UserID, &newUser.Username, &newUser.Email)
	if err != nil {
		return fmt.Errorf("error inserting user: %w", err)
	}

	return nil
}



func GetUserByID(id int) (*models.User, error) {
    user := &models.User{}
    query := "SELECT user_id, username, email, created_at, role FROM users WHERE user_id = $1"
    err := database.DB.QueryRow(query, id).Scan(&user.UserID, &user.Username, &user.Email, &user.CreatedAt, &user.Role)
    if err != nil {
        return nil, err
    }
    return user, nil
}


func DeleteUserByUsername(username string) error {
    sqlStatement := `DELETE FROM users WHERE username=$1;`
    _, err := database.DB.Exec(sqlStatement, username)
    return err
}

func DeleteAllUsers() error {
	sqlStatement := `DELETE FROM users`
    _, err := database.DB.Exec(sqlStatement)
    return err

}

func UserExists(username, email string) (bool, error) {
    var exists bool
    sqlStatement := `SELECT EXISTS (
        SELECT 1 FROM users WHERE username = $1 OR email = $2
    )`
    err := database.DB.QueryRow(sqlStatement, username, email).Scan(&exists)
    if err != nil {
        return false, fmt.Errorf("error checking existence: %w", err)
    }
    return exists, nil
}

