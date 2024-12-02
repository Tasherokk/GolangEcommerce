package controllers

import (
    "encoding/json"
    "net/http"
    "time"
    "github.com/Tasherokk/GolangProject.git/database"
	"github.com/Tasherokk/GolangProject.git/models"
    "strconv"
    "github.com/gorilla/mux"
)

// Получение изображений продукта
func GetProductImages(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    productIDStr := vars["product_id"]
    productID, err := strconv.Atoi(productIDStr)
    if err != nil {
        http.Error(w, "Invalid product ID", http.StatusBadRequest)
        return
    }

    rows, err := database.DB.Query("SELECT image_id, product_id, image_url, created_at FROM ProductImages WHERE product_id=$1", productID)
    if err != nil {
        http.Error(w, "Error fetching product images", http.StatusInternalServerError)
        return
    }
    defer rows.Close()

    var images []models.ProductImage
    for rows.Next() {
        var image models.ProductImage
        err := rows.Scan(&image.ImageID, &image.ProductID, &image.ImageURL, &image.CreatedAt)
        if err != nil {
            http.Error(w, "Error scanning image", http.StatusInternalServerError)
            return
        }
        images = append(images, image)
    }

    json.NewEncoder(w).Encode(images)
}

// Добавление изображения к продукту
func AddProductImage(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    productIDStr := vars["product_id"]
    productID, err := strconv.Atoi(productIDStr)
    if err != nil {
        http.Error(w, "Invalid product ID", http.StatusBadRequest)
        return
    }

    var image models.ProductImage
    err = json.NewDecoder(r.Body).Decode(&image)
    if err != nil {
        http.Error(w, "Invalid input", http.StatusBadRequest)
        return
    }

    image.ProductID = productID
    image.CreatedAt = time.Now()

    err = database.DB.QueryRow(`
        INSERT INTO ProductImages (product_id, image_url, created_at)
        VALUES ($1, $2, $3) RETURNING image_id`,
        image.ProductID, image.ImageURL, image.CreatedAt).
        Scan(&image.ImageID)
    if err != nil {
        http.Error(w, "Error adding product image", http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(image)
}
