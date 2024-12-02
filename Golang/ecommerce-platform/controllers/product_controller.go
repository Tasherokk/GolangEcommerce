// controllers/product_controller.go
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
	"github.com/go-playground/validator/v10"
	"github.com/go-redis/redis/v8"
	"github.com/gorilla/mux"
)



var product_validator = validator.New()


// Получение списка продуктов с кешированием
func GetAllProducts(w http.ResponseWriter, r *http.Request) {
    val, err := database.RDB.Get(database.Ctx, "products").Result()
    if err == redis.Nil {
        // Ключ не найден, получаем данные из базы
        rows, err := database.DB.Query(`
            SELECT 
                p.product_id, p.name, p.description, p.price, 
                p.stock, p.category_id, c.name AS category_name, 
                p.created_at
            FROM Products p
            LEFT JOIN Categories c ON p.category_id = c.category_id
        `)
        if err != nil {
            http.Error(w, "Error fetching products", http.StatusInternalServerError)
            return
        }
        defer rows.Close()

        var products []models.Product
        for rows.Next() {
            var product models.Product
            err := rows.Scan(&product.ProductID, &product.Name, &product.Description, &product.Price,
                &product.Stock, &product.CategoryID, &product.CategoryName, &product.CreatedAt)
            if err != nil {
                http.Error(w, "Error scanning product", http.StatusInternalServerError)
                return
            }
            products = append(products, product)
        }

        productsJSON, _ := json.Marshal(products)
        // Сохраняем в кеш на 10 минут
        database.RDB.Set(database.Ctx, "products", productsJSON, time.Minute*10)

        w.Header().Set("Content-Type", "application/json")
        w.Write(productsJSON)
    } else if err != nil {
        http.Error(w, "Cache error", http.StatusInternalServerError)
        return
    } else {
        // Данные найдены в кеше
        w.Header().Set("Content-Type", "application/json")
        w.Write([]byte(val))
    }
}




// Получение продукта по ID
func GetProductByID(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    idStr := vars["product_id"]
    id, err := strconv.Atoi(idStr)
    if err != nil {
        http.Error(w, "Invalid product ID", http.StatusBadRequest)
        return
    }

    var product models.Product
    err = database.DB.QueryRow(`
        SELECT p.product_id, p.name, p.description, p.price, p.stock, c.name AS category_name
        FROM Products p
        JOIN Categories c ON p.category_id = c.category_id
        WHERE p.product_id = $1
    `, id).Scan(&product.ProductID, &product.Name, &product.Description, &product.Price, &product.Stock, &product.CategoryName)
    if err != nil {
        if err == sql.ErrNoRows {
            http.Error(w, "Product not found", http.StatusNotFound)
        } else {
            http.Error(w, "Error fetching product details", http.StatusInternalServerError)
        }
        return
    }

    // Fetch product image
    var imageURL sql.NullString
    err = database.DB.QueryRow(`
        SELECT image_url
        FROM ProductImages
        WHERE product_id = $1
        LIMIT 1
    `, id).Scan(&imageURL)
    if err != nil && err != sql.ErrNoRows {
        http.Error(w, "Error fetching product image", http.StatusInternalServerError)
        return
    }
    product.ImageURL = imageURL.String

    // Fetch reviews with authors
    rows, err := database.DB.Query(`
        SELECT 
            r.review_id, r.product_id, r.user_id, r.rating, r.comment, r.created_at, 
            u.username AS author_name 
        FROM Reviews r
        JOIN Users u ON r.user_id = u.user_id
        WHERE r.product_id = $1
    `, id)
    if err != nil {
        http.Error(w, "Error fetching reviews", http.StatusInternalServerError)
        return
    }
    defer rows.Close()

    type ReviewWithAuthor struct {
        models.Review
        AuthorName string `json:"author_name"`
    }

    var reviews []ReviewWithAuthor
    for rows.Next() {
        var review ReviewWithAuthor
        err := rows.Scan(
            &review.ReviewID, &review.ProductID, &review.UserID, 
            &review.Rating, &review.Comment, &review.CreatedAt, &review.AuthorName,
        )
        if err != nil {
            http.Error(w, "Error scanning review", http.StatusInternalServerError)
            return
        }
        reviews = append(reviews, review)
    }

    response := struct {
        Product models.Product      `json:"product"`
        Reviews []ReviewWithAuthor `json:"reviews"`
    }{
        Product: product,
        Reviews: reviews,
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(response)
}




// Создание нового продукта
func CreateProduct(w http.ResponseWriter, r *http.Request) {
    log.Println("Starting product creation...")

    // Получение seller_id из контекста
    sellerID, ok := r.Context().Value(middleware.UserIDKey).(int)
    if !ok {
        log.Println("Seller ID not found in context")
        http.Error(w, "Unauthorized", http.StatusUnauthorized)
        return
    }
    log.Printf("Seller ID: %d\n", sellerID)

    var product models.Product
    err := json.NewDecoder(r.Body).Decode(&product)
    if err != nil {
        log.Println("Invalid input:", err)
        http.Error(w, "Invalid input", http.StatusBadRequest)
        return
    }
    log.Printf("Decoded product data: %+v\n", product)

    // Валидация данных
    err = product_validator.Struct(product)
    if err != nil {
        log.Println("Validation error:", err)
        http.Error(w, "Validation error: "+err.Error(), http.StatusBadRequest)
        return
    }
    log.Println("Product data validated successfully")

    product.CreatedAt = time.Now()
    product.SellerID = sellerID // Устанавливаем seller_id

    // Вставка продукта в базу данных
    err = database.DB.QueryRow(
        `INSERT INTO Products (name, description, price, stock, category_id, created_at, seller_id) 
         VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING product_id`,
        product.Name, product.Description, product.Price, product.Stock, product.CategoryID, product.CreatedAt, product.SellerID,
    ).Scan(&product.ProductID)
    if err != nil {
        log.Println("Error creating product:", err)
        http.Error(w, "Error creating product", http.StatusInternalServerError)
        return
    }

    log.Printf("Product created successfully with ID: %d\n", product.ProductID)

    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(product)
}



// Обновление продукта
func UpdateProduct(w http.ResponseWriter, r *http.Request) {
    log.Println("Starting product update...")

    vars := mux.Vars(r)
    idStr := vars["product_id"]
    id, err := strconv.Atoi(idStr)
    if err != nil {
        log.Println("Invalid product ID:", err)
        http.Error(w, "Invalid product ID", http.StatusBadRequest)
        return
    }
    log.Printf("Product ID to update: %d\n", id)

    var product models.Product
    err = json.NewDecoder(r.Body).Decode(&product)
    if err != nil {
        log.Println("Invalid input:", err)
        http.Error(w, "Invalid input", http.StatusBadRequest)
        return
    }
    product.ProductID = id
    log.Printf("Decoded product data: %+v\n", product)

    // Валидация данных
    err = product_validator.Struct(product)
    if err != nil {
        log.Println("Validation error:", err)
        http.Error(w, "Validation error: "+err.Error(), http.StatusBadRequest)
        return
    }
    log.Println("Product data validated successfully")

    err = database.DB.QueryRow("SELECT category_id FROM Categories WHERE name = $1", product.CategoryName).Scan(&product.CategoryID)
    if err != nil {
        if err == sql.ErrNoRows {
            log.Println("Category not found:", product.CategoryName)
            http.Error(w, "Category not found", http.StatusBadRequest)
            return
        }
        log.Println("Error fetching category ID:", err)
        http.Error(w, "Error fetching category ID", http.StatusInternalServerError)
        return
    }
    log.Printf("Resolved category ID: %d for category name: %s\n", product.CategoryID, product.CategoryName)

     // Обновляем продукт
     _, err = database.DB.Exec(
        "UPDATE Products SET name=$1, description=$2, price=$3, stock=$4, category_id=$5 WHERE product_id=$6",
        product.Name, product.Description, product.Price, product.Stock, product.CategoryID, product.ProductID,
    )
    if err != nil {
        log.Println("Error updating product:", err)
        http.Error(w, "Error updating product", http.StatusInternalServerError)
        return
    }

    // Обновляем или добавляем изображение
    if product.ImageURL != "" {
        var exists bool
        err = database.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM ProductImages WHERE product_id=$1)", product.ProductID).Scan(&exists)
        if err != nil {
            log.Println("Error checking image existence:", err)
            http.Error(w, "Error updating product image", http.StatusInternalServerError)
            return
        }

        if exists {
            // Обновляем существующее изображение
            _, err = database.DB.Exec(
                "UPDATE ProductImages SET image_url=$1, created_at=NOW() WHERE product_id=$2",
                product.ImageURL, product.ProductID,
            )
            if err != nil {
                log.Println("Error updating product image:", err)
                http.Error(w, "Error updating product image", http.StatusInternalServerError)
                return
            }
            log.Println("Product image updated successfully")
        } else {
            // Добавляем новое изображение
            _, err = database.DB.Exec(
                "INSERT INTO ProductImages (product_id, image_url, created_at) VALUES ($1, $2, NOW())",
                product.ProductID, product.ImageURL,
            )
            if err != nil {
                log.Println("Error inserting product image:", err)
                http.Error(w, "Error inserting product image", http.StatusInternalServerError)
                return
            }
            log.Println("Product image inserted successfully")
        }
    }

    log.Println("Product updated successfully in Products table")

    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(product)
}



// Удаление продукта
func DeleteProduct(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    idStr := vars["product_id"]
    id, err := strconv.Atoi(idStr)
    if err != nil {
        log.Printf("Error deleting product: %v", err)
        http.Error(w, "Invalid product ID", http.StatusBadRequest)
        return
    }

    _, err = database.DB.Exec("DELETE FROM Products WHERE product_id=$1", id)
    if err != nil {
        http.Error(w, "Error deleting product", http.StatusInternalServerError)
        return
    }

    // Удаление кэша
    cacheKey := "products"
    err = database.RDB.Del(database.Ctx, cacheKey).Err()
    if err != nil {
        log.Printf("Error deleting product from cache: %v", err)
        http.Error(w, "Error clearing cache", http.StatusInternalServerError)
        return
    }

    log.Printf("Product with ID %d deleted successfully, cache cleared", id)

    w.WriteHeader(http.StatusNoContent)
}





func GetProductByIDTest(id int) (models.Product, error) {
    var product models.Product
    cacheKey := fmt.Sprintf("product:%d", id)

    val, err := database.RDB.Get(database.Ctx, cacheKey).Result()
    if err == redis.Nil {
		log.Println("getProductByID Cache miss, fetch from DB")
        sqlStatement := `SELECT product_id, name, description, price FROM products WHERE product_id=$1;`
        row := database.DB.QueryRow(sqlStatement, id)
        err := row.Scan(&product.ProductID, &product.Name, &product.Description, &product.Price)
        if err != nil {
            return product, err
        }
		log.Println("getProductByID Caching the result")
        productJSON, _ := json.Marshal(product)
        database.RDB.Set(database.Ctx, cacheKey, productJSON, time.Minute*10)
    } else if err != nil {
        return product, err
    } else {
		log.Println("getProductByID Cache hit, unmarshal")
        json.Unmarshal([]byte(val), &product)
    }
    return product, nil
}



func GetSellerProducts(w http.ResponseWriter, r *http.Request) {
    log.Println("Starting GetSellerProducts handler")

    // Получение ID текущего пользователя из контекста
    userID, ok := r.Context().Value(middleware.UserIDKey).(int)
    if !ok {
        log.Println("Error: User ID not found in context")
        http.Error(w, "User ID not found in context", http.StatusUnauthorized)
        return
    }
    log.Printf("User ID: %d\n", userID)

    // Запрос на получение продуктов продавца
    log.Println("Fetching products for seller")
    rows, err := database.DB.Query(`
        SELECT DISTINCT ON (p.product_id) 
            p.product_id, p.name, p.description, p.price, p.stock, 
            p.category_id, c.name AS category_name, p.created_at, COALESCE(pi.image_url, '') AS image_url
        FROM 
            Products p
        LEFT JOIN 
            Categories c ON p.category_id = c.category_id
        LEFT JOIN 
            ProductImages pi ON p.product_id = pi.product_id
        WHERE 
            p.seller_id = $1
    `, userID)
    if err != nil {
        log.Printf("Error fetching products from database: %v\n", err)
        http.Error(w, "Error fetching products", http.StatusInternalServerError)
        return
    }
    defer rows.Close()

    log.Println("Query executed successfully, scanning results")

    // Список для хранения продуктов
    var products []models.Product

    // Чтение результатов запроса
    for rows.Next() {
        var product models.Product
        err := rows.Scan(
            &product.ProductID, &product.Name, &product.Description, &product.Price,
            &product.Stock, &product.CategoryID, &product.CategoryName, &product.CreatedAt, &product.ImageURL,
        )
        if err != nil {
            log.Printf("Error scanning product row: %v\n", err)
            http.Error(w, "Error scanning product", http.StatusInternalServerError)
            return
        }
        products = append(products, product)
    }

    if rows.Err() != nil {
        log.Printf("Error iterating through rows: %v\n", rows.Err())
        http.Error(w, "Error processing query results", http.StatusInternalServerError)
        return
    }

    log.Printf("Products fetched successfully. Total products: %d\n", len(products))

    // Отправка данных клиенту
    w.Header().Set("Content-Type", "application/json")
    if err := json.NewEncoder(w).Encode(products); err != nil {
        log.Printf("Error encoding products to JSON: %v\n", err)
        http.Error(w, "Error encoding response", http.StatusInternalServerError)
        return
    }

    log.Printf("Products fetched: %+v", products)
    log.Println("Response sent successfully")
}

