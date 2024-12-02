package controllers

import (
    "encoding/json"
    "net/http"
    "strconv"

    "github.com/gorilla/mux"
    "github.com/Tasherokk/GolangProject.git/database"
	"github.com/Tasherokk/GolangProject.git/models"
)

// Получение элементов заказа по ID заказа
func GetOrderItems(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    orderIDStr := vars["order_id"]
    orderID, err := strconv.Atoi(orderIDStr)
    if err != nil {
        http.Error(w, "Invalid order ID", http.StatusBadRequest)
        return
    }

    // Получение информации о заказе
    var status string
    err = database.DB.QueryRow("SELECT status FROM Orders WHERE order_id=$1", orderID).Scan(&status)
    if err != nil {
        http.Error(w, "Error fetching order status", http.StatusInternalServerError)
        return
    }

    query := `
        SELECT 
            oi.order_item_id,
            oi.order_id,
            oi.product_id,
            p.name AS product_name,
            oi.quantity,
            oi.price
        FROM OrderItems oi
        JOIN Products p ON oi.product_id = p.product_id
        WHERE oi.order_id = $1
    `

    rows, err := database.DB.Query(query, orderID)
    if err != nil {
        http.Error(w, "Error fetching order items", http.StatusInternalServerError)
        return
    }
    defer rows.Close()

    var orderItems []models.OrderItem
    for rows.Next() {
        var item models.OrderItem
        err := rows.Scan(&item.OrderItemID, &item.OrderID, &item.ProductID, &item.ProductName, &item.Quantity, &item.Price)
        if err != nil {
            http.Error(w, "Error scanning order item", http.StatusInternalServerError)
            return
        }
        orderItems = append(orderItems, item)
    }

    // Создание ответа с данными о статусе и элементах заказа
    response := map[string]interface{}{
        "status": status,
        "items":  orderItems,
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(response)
}

// Добавление элементов к заказу
func AddOrderItems(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    orderIDStr := vars["order_id"]
    orderID, err := strconv.Atoi(orderIDStr)
    if err != nil {
        http.Error(w, "Invalid order ID", http.StatusBadRequest)
        return
    }

    var items []models.OrderItem
    err = json.NewDecoder(r.Body).Decode(&items)
    if err != nil {
        http.Error(w, "Invalid input", http.StatusBadRequest)
        return
    }

    for _, item := range items {
        _, err := database.DB.Exec(
            "INSERT INTO OrderItems (order_id, product_id, quantity, price) VALUES ($1, $2, $3, $4)",
            orderID, item.ProductID, item.Quantity, item.Price)
        if err != nil {
            http.Error(w, "Error adding order item", http.StatusInternalServerError)
            return
        }
    }

    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(map[string]string{"message": "Order items added"})
}
