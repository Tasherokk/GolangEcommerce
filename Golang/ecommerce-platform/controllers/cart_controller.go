package controllers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"github.com/Tasherokk/GolangProject.git/database"
	"github.com/Tasherokk/GolangProject.git/middleware"
	"github.com/Tasherokk/GolangProject.git/models"
)

// Получение корзины пользователя
func GetUserCart(w http.ResponseWriter, r *http.Request) {
    userID, ok := r.Context().Value(middleware.UserIDKey).(int)

    if !ok {
        http.Error(w, "User ID not found in context", http.StatusUnauthorized)
        return
    }

    var cart models.ShoppingCart
    err := database.DB.QueryRow("SELECT cart_id, user_id, created_at FROM ShoppingCart WHERE user_id=$1", userID).
        Scan(&cart.CartID, &cart.UserID, &cart.CreatedAt)
    if err != nil {
        if err == sql.ErrNoRows {
            // Создаем новую корзину
            cart.UserID = userID
            cart.CreatedAt = time.Now()
            err = database.DB.QueryRow(
                "INSERT INTO ShoppingCart (user_id, created_at) VALUES ($1, $2) RETURNING cart_id",
                cart.UserID, cart.CreatedAt).Scan(&cart.CartID)
            if err != nil {
                http.Error(w, "Error creating cart", http.StatusInternalServerError)
                return
            }
        } else {
            http.Error(w, "Error fetching cart", http.StatusInternalServerError)
            return
        }
    }

    json.NewEncoder(w).Encode(cart)
}
