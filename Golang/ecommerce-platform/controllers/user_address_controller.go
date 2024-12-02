package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/Tasherokk/GolangProject.git/database"
	"github.com/Tasherokk/GolangProject.git/middleware"
	"github.com/Tasherokk/GolangProject.git/models"
	"github.com/go-playground/validator/v10"
)

var address_validator = validator.New()

// Получение адресов пользователя
func GetUserAddresses(w http.ResponseWriter, r *http.Request) {
    userID, ok := r.Context().Value(middleware.UserIDKey).(int)

    if !ok {
        http.Error(w, "User ID not found in context", http.StatusUnauthorized)
        return
    }

    rows, err := database.DB.Query("SELECT address_id, user_id, street, city, state, zip_code FROM UserAddresses WHERE user_id=$1", userID)
    if err != nil {
        http.Error(w, "Error fetching addresses", http.StatusInternalServerError)
        return
    }
    defer rows.Close()

    var addresses []models.UserAddress
    for rows.Next() {
        var address models.UserAddress
        err := rows.Scan(&address.AddressID, &address.UserID, &address.Street, &address.City, &address.State, &address.ZipCode)
        if err != nil {
            http.Error(w, "Error scanning address", http.StatusInternalServerError)
            return
        }
        addresses = append(addresses, address)
    }

    json.NewEncoder(w).Encode(addresses)
}

// Добавление адреса пользователя
func AddUserAddress(w http.ResponseWriter, r *http.Request) {
    userID, ok := r.Context().Value(middleware.UserIDKey).(int)

    if !ok {
        http.Error(w, "User ID not found in context", http.StatusUnauthorized)
        return
    }
    var address models.UserAddress
    err := json.NewDecoder(r.Body).Decode(&address)
    if err != nil {
        http.Error(w, "Invalid input", http.StatusBadRequest)
        return
    }

    address.UserID = userID

    // Валидация данных
    err = address_validator.Struct(address)
    if err != nil {
        http.Error(w, "Validation error: "+err.Error(), http.StatusBadRequest)
        return
    }

    err = database.DB.QueryRow(`
        INSERT INTO UserAddresses (user_id, street, city, state, zip_code)
        VALUES ($1, $2, $3, $4, $5) RETURNING address_id`,
        address.UserID, address.Street, address.City, address.State, address.ZipCode).
        Scan(&address.AddressID)
    if err != nil {
        http.Error(w, "Error adding address", http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(address)
}
