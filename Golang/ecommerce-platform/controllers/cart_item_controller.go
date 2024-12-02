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
	"github.com/gorilla/mux"
)

// Получение элементов корзины
func GetCartItems(w http.ResponseWriter, r *http.Request) {
    userID, ok := r.Context().Value(middleware.UserIDKey).(int)
    if !ok {
        http.Error(w, "User ID not found in context", http.StatusUnauthorized)
        return
    }

    var cartID int
    err := database.DB.QueryRow("SELECT cart_id FROM ShoppingCart WHERE user_id=$1", userID).Scan(&cartID)
    if err == sql.ErrNoRows {
        log.Println("No cart found for user. Creating a new cart...")
        err = database.DB.QueryRow(
            "INSERT INTO ShoppingCart (user_id, created_at) VALUES ($1, NOW()) RETURNING cart_id",
            userID,
        ).Scan(&cartID)
        if err != nil {
            log.Printf("Error creating cart: %v\n", err)
            http.Error(w, "Error creating cart", http.StatusInternalServerError)
            return
        }
    } else if err != nil {
        log.Printf("Error fetching cart: %v\n", err)
        http.Error(w, "Error fetching cart", http.StatusInternalServerError)
        return
    }

    log.Printf("User ID: %d\n", userID)

    rows, err := database.DB.Query(`
        SELECT ci.cart_item_id, ci.product_id, ci.quantity, p.name AS product_name, p.price
        FROM CartItems ci
        JOIN Products p ON ci.product_id = p.product_id
        WHERE ci.cart_id = $1
    `, cartID)
    if err != nil {
        http.Error(w, "Error fetching cart items", http.StatusInternalServerError)
        return
    }
    defer rows.Close()

    var cartItems []models.CartItem
    for rows.Next() {
        var item models.CartItem
        if err := rows.Scan(&item.CartItemID, &item.ProductID, &item.Quantity, &item.ProductName, &item.Price); err != nil {
            http.Error(w, "Error scanning cart items", http.StatusInternalServerError)
            return
        }
        cartItems = append(cartItems, item)
    }

    json.NewEncoder(w).Encode(cartItems)
}

// Добавление элемента в корзину
func AddCartItem(w http.ResponseWriter, r *http.Request) {
    userID, ok := r.Context().Value(middleware.UserIDKey).(int)
    if !ok {
        log.Println("Error: User ID not found in context")
        http.Error(w, "User ID not found in context", http.StatusUnauthorized)
        return
    }

    log.Printf("User ID: %d\n", userID)

    var cartID int
    err := database.DB.QueryRow("SELECT cart_id FROM ShoppingCart WHERE user_id=$1", userID).Scan(&cartID)
    if err == sql.ErrNoRows {
        log.Println("No cart found for user. Creating a new cart...")
        err = database.DB.QueryRow(
            "INSERT INTO ShoppingCart (user_id, created_at) VALUES ($1, NOW()) RETURNING cart_id",
            userID,
        ).Scan(&cartID)
        if err != nil {
            log.Printf("Error creating cart: %v\n", err)
            http.Error(w, "Error creating cart", http.StatusInternalServerError)
            return
        }
    } else if err != nil {
        log.Printf("Error fetching cart: %v\n", err)
        http.Error(w, "Error fetching cart", http.StatusInternalServerError)
        return
    }

    log.Printf("Cart ID: %d\n", cartID)

    var item models.CartItem
    err = json.NewDecoder(r.Body).Decode(&item)
    if err != nil {
        log.Printf("Error decoding request body: %v\n", err)
        http.Error(w, "Invalid input", http.StatusBadRequest)
        return
    }

    log.Printf("Adding product %d with quantity %d to cart %d\n", item.ProductID, item.Quantity, cartID)

    _, err = database.DB.Exec(
        "INSERT INTO CartItems (cart_id, product_id, quantity) VALUES ($1, $2, $3) ON CONFLICT (cart_id, product_id) DO UPDATE SET quantity = CartItems.quantity + EXCLUDED.quantity",
        cartID, item.ProductID, item.Quantity,
    )
    if err != nil {
        log.Printf("Error adding item to cart: %v\n", err)
        http.Error(w, "Error adding item to cart", http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(map[string]string{"message": "Item added to cart"})
}



func DeleteCartItem(w http.ResponseWriter, r *http.Request) {
	// Извлечение `item_id` из URL
	vars := mux.Vars(r)
	itemID, err := strconv.Atoi(vars["item_id"])
	if err != nil {
		http.Error(w, "Invalid cart item ID", http.StatusBadRequest)
		return
	}

	// Удаление элемента из таблицы `CartItems`
	_, err = database.DB.Exec("DELETE FROM CartItems WHERE cart_item_id = $1", itemID)
	if err != nil {
		http.Error(w, "Error removing cart item", http.StatusInternalServerError)
		return
	}

	// Успешный ответ
	w.WriteHeader(http.StatusNoContent)
}


func DecreaseCartItemQuantity(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    cartItemID, err := strconv.Atoi(vars["item_id"])
    if err != nil {
        http.Error(w, "Invalid cart item ID", http.StatusBadRequest)
        return
    }

    // Уменьшить количество на 1
    result, err := database.DB.Exec(
        "UPDATE CartItems SET quantity = quantity - 1 WHERE cart_item_id = $1 AND quantity > 1",
        cartItemID,
    )
    if err != nil {
        http.Error(w, "Error updating cart item", http.StatusInternalServerError)
        return
    }

    rowsAffected, _ := result.RowsAffected()
    if rowsAffected == 0 {
        http.Error(w, "Cannot decrease quantity below 1", http.StatusBadRequest)
        return
    }

    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(map[string]string{"message": "Item quantity decreased"})
}


func IncreaseCartItemQuantity(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    cartItemID, err := strconv.Atoi(vars["item_id"])
    if err != nil {
        http.Error(w, "Invalid cart item ID", http.StatusBadRequest)
        return
    }

    // Получение текущего количества и доступного запаса
    var currentQuantity, stock int
    err = database.DB.QueryRow(`
        SELECT c.quantity, p.stock
        FROM CartItems c
        JOIN Products p ON c.product_id = p.product_id
        WHERE c.cart_item_id = $1
    `, cartItemID).Scan(&currentQuantity, &stock)
    if err != nil {
        if err == sql.ErrNoRows {
            http.Error(w, "Cart item not found", http.StatusNotFound)
        } else {
            http.Error(w, "Database error", http.StatusInternalServerError)
        }
        return
    }

    // Проверка: не превышает ли количество доступный stock
    if currentQuantity >= stock {
        http.Error(w, "Cannot add more items than available in stock", http.StatusBadRequest)
        return
    }

    // Увеличиваем количество на 1
    _, err = database.DB.Exec(
        "UPDATE CartItems SET quantity = quantity + 1 WHERE cart_item_id = $1",
        cartItemID,
    )
    if err != nil {
        http.Error(w, "Error updating cart item", http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(map[string]string{"message": "Item quantity increased"})
}

