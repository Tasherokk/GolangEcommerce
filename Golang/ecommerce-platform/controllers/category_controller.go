package controllers

import (
    "database/sql"
    "encoding/json"
    "net/http"
    "strconv"

    "github.com/gorilla/mux"
    "github.com/Tasherokk/GolangProject.git/database"
	"github.com/Tasherokk/GolangProject.git/models"
    "github.com/go-playground/validator/v10"

)

var category_validator = validator.New()

// Получение всех категорий
func GetAllCategories(w http.ResponseWriter, r *http.Request) {
    rows, err := database.DB.Query("SELECT category_id, name, description FROM Categories")
    if err != nil {
        http.Error(w, "Error fetching categories", http.StatusInternalServerError)
        return
    }
    defer rows.Close()

    var categories []models.Category
    for rows.Next() {
        var category models.Category
        err := rows.Scan(&category.CategoryID, &category.Name, &category.Description)
        if err != nil {
            http.Error(w, "Error scanning category", http.StatusInternalServerError)
            return
        }
        categories = append(categories, category)
    }

    json.NewEncoder(w).Encode(categories)
}



// Получение категории по ID
func GetCategoryByID(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    idStr := vars["category_id"]
    id, err := strconv.Atoi(idStr)
    if err != nil {
        http.Error(w, "Invalid category ID", http.StatusBadRequest)
        return
    }

    var category models.Category
    err = database.DB.QueryRow("SELECT category_id, name, description FROM Categories WHERE category_id=$1", id).
        Scan(&category.CategoryID, &category.Name, &category.Description)
    if err != nil {
        if err == sql.ErrNoRows {
            http.Error(w, "Category not found", http.StatusNotFound)
        } else {
            http.Error(w, "Error fetching category", http.StatusInternalServerError)
        }
        return
    }

    json.NewEncoder(w).Encode(category)
}

// Создание новой категории
func CreateCategory(w http.ResponseWriter, r *http.Request) {
    var category models.Category
    err := json.NewDecoder(r.Body).Decode(&category)
    if err != nil {
        http.Error(w, "Invalid input", http.StatusBadRequest)
        return
    }

    // Валидация данных
    err = category_validator.Struct(category)
    if err != nil {
        http.Error(w, "Validation error: "+err.Error(), http.StatusBadRequest)
        return
    }

    err = database.DB.QueryRow(
        "INSERT INTO Categories (name, description) VALUES ($1, $2) RETURNING category_id",
        category.Name, category.Description).Scan(&category.CategoryID)
    if err != nil {
        http.Error(w, "Error creating category", http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(category)
}

// Обновление категории
func UpdateCategory(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    idStr := vars["category_id"]
    id, err := strconv.Atoi(idStr)
    if err != nil {
        http.Error(w, "Invalid category ID", http.StatusBadRequest)
        return
    }

    var category models.Category
    err = json.NewDecoder(r.Body).Decode(&category)
    if err != nil {
        http.Error(w, "Invalid input", http.StatusBadRequest)
        return
    }

    // Валидация данных
    err = category_validator.Struct(category)
    if err != nil {
        http.Error(w, "Validation error: "+err.Error(), http.StatusBadRequest)
        return
    }

    _, err = database.DB.Exec(
        "UPDATE Categories SET name=$1, description=$2 WHERE category_id=$3",
        category.Name, category.Description, id)
    if err != nil {
        http.Error(w, "Error updating category", http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(category)
}

// Удаление категории
func DeleteCategory(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    idStr := vars["category_id"]
    id, err := strconv.Atoi(idStr)
    if err != nil {
        http.Error(w, "Invalid category ID", http.StatusBadRequest)
        return
    }

    _, err = database.DB.Exec("DELETE FROM Categories WHERE category_id=$1", id)
    if err != nil {
        http.Error(w, "Error deleting category", http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusNoContent)
}
