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

// Получение заказов пользователя
func GetUserOrders(w http.ResponseWriter, r *http.Request) {
    userID, ok := r.Context().Value(middleware.UserIDKey).(int)
    if !ok {
        http.Error(w, "User ID not found in context", http.StatusUnauthorized)
        return
    }

    // Выполняем запрос для получения заказов пользователя
    rows, err := database.DB.Query(`
        SELECT 
            o.order_id, o.user_id, o.order_date, o.status, o.total_amount,
            COALESCE(p.payment_method, 'N/A') AS payment_method
        FROM Orders o
        LEFT JOIN Payments p ON o.order_id = p.order_id
        WHERE o.user_id=$1`, userID)
    if err != nil {
        http.Error(w, "Error fetching orders", http.StatusInternalServerError)
        return
    }
    defer rows.Close()

    var orders []struct {
        models.Order
        PaymentMethod string `json:"payment_method,omitempty"`
    }

    // Обработка результата
    for rows.Next() {
        var order models.Order
        var paymentMethod string

        err := rows.Scan(&order.OrderID, &order.UserID, &order.OrderDate, &order.Status, &order.TotalAmount, &paymentMethod)
        if err != nil {
            http.Error(w, "Error scanning order", http.StatusInternalServerError)
            return
        }

        orders = append(orders, struct {
            models.Order
            PaymentMethod string `json:"payment_method,omitempty"`
        }{
            Order:         order,
            PaymentMethod: paymentMethod,
        })
    }

    // Проверяем наличие ошибок при обработке строк
    if err = rows.Err(); err != nil {
        http.Error(w, "Error reading rows", http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(orders)
}



// Получение заказа по ID
func GetOrderByID(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    orderIDStr := vars["id"]
    orderID, err := strconv.Atoi(orderIDStr)
    if err != nil {
        http.Error(w, "Invalid order ID", http.StatusBadRequest)
        return
    }

    var order models.Order
    err = database.DB.QueryRow("SELECT order_id, user_id, order_date, status, total_amount FROM Orders WHERE order_id=$1", orderID).
        Scan(&order.OrderID, &order.UserID, &order.OrderDate, &order.Status, &order.TotalAmount)
    if err != nil {
        if err == sql.ErrNoRows {
            http.Error(w, "Order not found", http.StatusNotFound)
        } else {
            http.Error(w, "Error fetching order", http.StatusInternalServerError)
        }
        return
    }

    json.NewEncoder(w).Encode(order)
}

// Создание нового заказа
func CreateOrder(w http.ResponseWriter, r *http.Request) {
    log.Println("Starting checkout process")

    // Получение ID пользователя из контекста
    userID, ok := r.Context().Value(middleware.UserIDKey).(int)
    if !ok {
        log.Println("User ID not found in context")
        http.Error(w, "Unauthorized", http.StatusUnauthorized)
        return
    }
    log.Printf("User ID: %d", userID)

    // Начало транзакции
    tx, err := database.DB.Begin()
    if err != nil {
        log.Printf("Error starting transaction: %v", err)
        http.Error(w, "Failed to start transaction", http.StatusInternalServerError)
        return
    }
    log.Println("Transaction started")

    defer func() {
        if err := recover(); err != nil {
            tx.Rollback()
            log.Printf("Transaction rollback due to panic: %v", err)
            http.Error(w, "Internal server error", http.StatusInternalServerError)
        }
    }()

    // Получение корзины пользователя
    var cartID int
    err = tx.QueryRow("SELECT cart_id FROM ShoppingCart WHERE user_id=$1", userID).Scan(&cartID)
    if err != nil {
        tx.Rollback()
        log.Printf("Error fetching cart: %v", err)
        http.Error(w, "Failed to fetch cart", http.StatusInternalServerError)
        return
    }
    log.Printf("Cart ID: %d", cartID)

    // Создание заказа
    var orderID int
    err = tx.QueryRow(
        "INSERT INTO Orders (user_id, order_date, status, total_amount) VALUES ($1, CURRENT_TIMESTAMP, 'In Progress', 0) RETURNING order_id",
        userID,
    ).Scan(&orderID)
    if err != nil {
        tx.Rollback()
        log.Printf("Error creating order: %v", err)
        http.Error(w, "Failed to create order", http.StatusInternalServerError)
        return
    }
    log.Printf("Order created with ID: %d", orderID)

    // Получение элементов корзины
    rows, err := tx.Query("SELECT product_id, quantity FROM CartItems WHERE cart_id=$1", cartID)
    if err != nil {
        tx.Rollback()
        log.Printf("Error fetching cart items: %v", err)
        http.Error(w, "Failed to fetch cart items", http.StatusInternalServerError)
        return
    }

    // Чтение всех данных из rows
    var items []struct {
        ProductID int
        Quantity  int
    }
    for rows.Next() {
        var item struct {
            ProductID int
            Quantity  int
        }
        if err := rows.Scan(&item.ProductID, &item.Quantity); err != nil {
            rows.Close()
            tx.Rollback()
            log.Printf("Error scanning cart item: %v", err)
            http.Error(w, "Failed to scan cart item", http.StatusInternalServerError)
            return
        }
        items = append(items, item)
    }
    rows.Close() // Закрываем rows сразу после обработки
    if err := rows.Err(); err != nil {
        tx.Rollback()
        log.Printf("Error during rows iteration: %v", err)
        http.Error(w, "Error iterating over rows", http.StatusInternalServerError)
        return
    }

    // Добавление элементов в заказ
    totalAmount := 0.0
    for _, item := range items {
        var price float64
        err = tx.QueryRow("SELECT price FROM Products WHERE product_id=$1", item.ProductID).Scan(&price)
        if err != nil {
            tx.Rollback()
            log.Printf("Error fetching product price for product %d: %v", item.ProductID, err)
            http.Error(w, "Failed to fetch product price", http.StatusInternalServerError)
            return
        }

        totalAmount += price * float64(item.Quantity)
        log.Printf("Adding product %d with quantity %d and price %.2f to order %d", item.ProductID, item.Quantity, price, orderID)

        _, err = tx.Exec(
            "INSERT INTO OrderItems (order_id, product_id, quantity, price) VALUES ($1, $2, $3, $4)",
            orderID, item.ProductID, item.Quantity, price,
        )
        if err != nil {
            tx.Rollback()
            log.Printf("Error inserting order item for order %d: %v", orderID, err)
            http.Error(w, "Failed to insert order item", http.StatusInternalServerError)
            return
        }
    }

    // Обновление суммы заказа
    _, err = tx.Exec("UPDATE Orders SET total_amount=$1 WHERE order_id=$2", totalAmount, orderID)
    if err != nil {
        tx.Rollback()
        log.Printf("Error updating order total: %v", err)
        http.Error(w, "Failed to update order total", http.StatusInternalServerError)
        return
    }

    // Очистка корзины
    _, err = tx.Exec("DELETE FROM CartItems WHERE cart_id=$1", cartID)
    if err != nil {
        tx.Rollback()
        log.Printf("Error clearing cart: %v", err)
        http.Error(w, "Failed to clear cart", http.StatusInternalServerError)
        return
    }

    // Завершение транзакции
    if err := tx.Commit(); err != nil {
        log.Printf("Error committing transaction: %v", err)
        http.Error(w, "Failed to commit transaction", http.StatusInternalServerError)
        return
    }

    log.Println("Checkout process completed successfully")
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(map[string]interface{}{
        "message":  "Order placed successfully",
        "order_id": orderID,
    })
}


// Обновление статуса заказа
func UpdateOrderStatus(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    orderIDStr := vars["order_id"]
    orderID, err := strconv.Atoi(orderIDStr)
    if err != nil {
        http.Error(w, "Invalid order ID", http.StatusBadRequest)
        return
    }

    var statusUpdate struct {
        Status string `json:"status"`
    }
    err = json.NewDecoder(r.Body).Decode(&statusUpdate)
    if err != nil {
        http.Error(w, "Invalid input", http.StatusBadRequest)
        return
    }

    _, err = database.DB.Exec("UPDATE Orders SET status=$1 WHERE order_id=$2", statusUpdate.Status, orderID)
    if err != nil {
        http.Error(w, "Error updating order status", http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(map[string]string{"message": "Order status updated"})
}



func GetSellerOrders(w http.ResponseWriter, r *http.Request) {
    // Получение seller_id из токена или контекста
    sellerID, ok := r.Context().Value(middleware.UserIDKey).(int)
    if !ok {
        http.Error(w, "Seller ID not found in context", http.StatusUnauthorized)
        return
    }

    // SQL-запрос для получения заказов через таблицу order_items
    query := `
        SELECT 
            oi.order_item_id, oi.order_id, oi.product_id, oi.quantity, oi.price,
            p.name AS product_name, o.order_date, o.status, o.total_amount
        FROM 
            OrderItems oi
        JOIN 
            Products p ON oi.product_id = p.product_id
        JOIN 
            Orders o ON oi.order_id = o.order_id
        WHERE 
            p.seller_id = $1
        ORDER BY o.order_date DESC
    `

    rows, err := database.DB.Query(query, sellerID)
    if err != nil {
        http.Error(w, "Error fetching orders", http.StatusInternalServerError)
        return
    }
    defer rows.Close()

    // Собираем заказы в массив
    var orders []struct {
        OrderItemID int     `json:"order_item_id"`
        OrderID     int     `json:"order_id"`
        ProductID   int     `json:"product_id"`
        Quantity    int     `json:"quantity"`
        Price       float64 `json:"price"`
        ProductName string  `json:"product_name"`
        OrderDate   string  `json:"order_date"`
        Status      string  `json:"status"`
        TotalAmount float64 `json:"total_amount"`
    }

    for rows.Next() {
        var order struct {
            OrderItemID int     `json:"order_item_id"`
            OrderID     int     `json:"order_id"`
            ProductID   int     `json:"product_id"`
            Quantity    int     `json:"quantity"`
            Price       float64 `json:"price"`
            ProductName string  `json:"product_name"`
            OrderDate   string  `json:"order_date"`
            Status      string  `json:"status"`
            TotalAmount float64 `json:"total_amount"`
        }

        err = rows.Scan(&order.OrderItemID, &order.OrderID, &order.ProductID, &order.Quantity, &order.Price,
            &order.ProductName, &order.OrderDate, &order.Status, &order.TotalAmount)
        if err != nil {
            http.Error(w, "Error scanning order", http.StatusInternalServerError)
            return
        }

        orders = append(orders, order)
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(orders)
}
