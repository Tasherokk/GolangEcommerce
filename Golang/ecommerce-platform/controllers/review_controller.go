package controllers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/Tasherokk/GolangProject.git/database"
	"github.com/Tasherokk/GolangProject.git/middleware"
	"github.com/Tasherokk/GolangProject.git/models"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
)

var review_validator = validator.New()

// Получение отзывов по продукту
func GetReviewsByProduct(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    productIDStr := vars["product_id"]
    productID, err := strconv.Atoi(productIDStr)
    if err != nil {
        http.Error(w, "Invalid product ID", http.StatusBadRequest)
        return
    }

    rows, err := database.DB.Query(`
        SELECT review_id, product_id, user_id, rating, comment, created_at
        FROM Reviews WHERE product_id=$1`, productID)
    if err != nil {
        http.Error(w, "Error fetching reviews", http.StatusInternalServerError)
        return
    }
    defer rows.Close()

    var reviews []models.Review
    for rows.Next() {
        var review models.Review
        err := rows.Scan(&review.ReviewID, &review.ProductID, &review.UserID, &review.Rating, &review.Comment, &review.CreatedAt)
        if err != nil {
            http.Error(w, "Error scanning review", http.StatusInternalServerError)
            return
        }
        reviews = append(reviews, review)
    }

    json.NewEncoder(w).Encode(reviews)
}


func CreateReview(w http.ResponseWriter, r *http.Request) {
    log.Println("Starting review creation...")

    // Получаем ID пользователя из контекста
    userID, ok := r.Context().Value(middleware.UserIDKey).(int)
    if !ok {
        log.Println("User not authenticated")
        http.Error(w, "User not authenticated", http.StatusUnauthorized)
        return
    }
    log.Printf("User ID: %d\n", userID)

    // Получаем ID продукта из URL
    vars := mux.Vars(r)
    productIDStr := vars["product_id"]
    productID, err := strconv.Atoi(productIDStr)
    if err != nil {
        log.Printf("Invalid product ID: %s\n", productIDStr)
        http.Error(w, "Invalid product ID", http.StatusBadRequest)
        return
    }
    log.Printf("Product ID: %d\n", productID)

    // Декодируем тело запроса
    var review models.Review
    err = json.NewDecoder(r.Body).Decode(&review)
    if err != nil {
        log.Println("Invalid input:", err)
        http.Error(w, "Invalid input", http.StatusBadRequest)
        return
    }
    log.Printf("Decoded review input: %+v\n", review)

    // Устанавливаем идентификаторы пользователя и продукта
    review.UserID = userID
    review.ProductID = productID

    // Используем review_validator для валидации
    err = review_validator.Struct(review)
    if err != nil {
        log.Println("Validation error:", err)
        http.Error(w, "Validation error: "+err.Error(), http.StatusBadRequest)
        return
    }
    log.Println("Review input validated successfully")

    // Проверяем, оставлял ли пользователь уже отзыв на этот продукт
    var existingReviewID int
    err = database.DB.QueryRow("SELECT review_id FROM Reviews WHERE user_id = $1 AND product_id = $2", userID, productID).Scan(&existingReviewID)
    if err != sql.ErrNoRows {
        if err == nil {
            log.Printf("User %d has already reviewed product %d\n", userID, productID)
            http.Error(w, "You have already reviewed this product", http.StatusBadRequest)
            return
        } else {
            log.Println("Error checking existing reviews:", err)
            http.Error(w, "Error checking existing reviews", http.StatusInternalServerError)
            return
        }
    }
    log.Println("No existing review found for this product by the user")

    // Вставляем отзыв в базу данных
    _, err = database.DB.Exec(
        "INSERT INTO Reviews (product_id, user_id, rating, comment) VALUES ($1, $2, $3, $4)",
        review.ProductID, review.UserID, review.Rating, review.Comment)
    if err != nil {
        log.Println("Error saving review:", err)
        http.Error(w, "Error saving review", http.StatusInternalServerError)
        return
    }
    log.Println("Review saved successfully")

    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(map[string]string{"message": "Review submitted successfully"})
}
