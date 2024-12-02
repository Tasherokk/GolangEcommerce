package controllers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/Tasherokk/GolangProject.git/database"
	"github.com/Tasherokk/GolangProject.git/middleware"
	"github.com/Tasherokk/GolangProject.git/models"
	// "github.com/go-playground/validator/v10"
)

// var payment_validator = validator.New()

// Получение платежей пользователя
func GetUserPayments(w http.ResponseWriter, r *http.Request) {
    userID, ok := r.Context().Value(middleware.UserIDKey).(int)

    if !ok {
        http.Error(w, "User ID not found in context", http.StatusUnauthorized)
        return
    }

    rows, err := database.DB.Query(`
        SELECT p.payment_id, p.order_id, p.amount, p.payment_date, p.payment_method
        FROM Payments p
        INNER JOIN Orders o ON p.order_id = o.order_id
        WHERE o.user_id = $1`, userID)
    if err != nil {
        http.Error(w, "Error fetching payments", http.StatusInternalServerError)
        return
    }
    defer rows.Close()

    var payments []models.Payment
    for rows.Next() {
        var payment models.Payment
        err := rows.Scan(&payment.PaymentID, &payment.OrderID, &payment.Amount, &payment.PaymentDate, &payment.PaymentMethod)
        if err != nil {
            http.Error(w, "Error scanning payment", http.StatusInternalServerError)
            return
        }
        payments = append(payments, payment)
    }

    json.NewEncoder(w).Encode(payments)
}


// Создание нового платежа
func CreatePayment(w http.ResponseWriter, r *http.Request) {
    log.Println("Starting payment creation...")

    var payment models.Payment
    err := json.NewDecoder(r.Body).Decode(&payment)
    if err != nil {
        log.Println("Invalid input:", err)

        // Проверка, если ошибка связана с неверным форматом order_id
        var temp map[string]interface{}
        if decodeErr := json.NewDecoder(r.Body).Decode(&temp); decodeErr == nil {
            if orderIDStr, ok := temp["order_id"].(string); ok {
                if orderID, convErr := strconv.Atoi(orderIDStr); convErr == nil {
                    payment.OrderID = orderID
                } else {
                    log.Println("Error converting order_id to int:", convErr)
                    http.Error(w, "Invalid order ID format", http.StatusBadRequest)
                    return
                }
            }
        } else {
            http.Error(w, "Invalid input format", http.StatusBadRequest)
            return
        }
    }

    log.Printf("Decoded payment data: %+v\n", payment)

    // Вставка платежа в таблицу Payments
    log.Println("Inserting payment into database...")
    _, err = database.DB.Exec(
        "INSERT INTO Payments (order_id, amount, payment_method) VALUES ($1, $2, $3)",
        payment.OrderID, payment.Amount, payment.PaymentMethod,
    )
    if err != nil {
        log.Println("Error recording payment:", err)
        http.Error(w, "Error recording payment", http.StatusInternalServerError)
        return
    }
    log.Println("Payment recorded successfully in Payments table")

    // Обновляем статус заказа на "Paid"
    log.Printf("Updating order status to 'Paid' for order ID: %d\n", payment.OrderID)
    _, err = database.DB.Exec("UPDATE Orders SET status = $1 WHERE order_id = $2", "Paid", payment.OrderID)
    if err != nil {
        log.Println("Error updating order status:", err)
        http.Error(w, "Error updating order status", http.StatusInternalServerError)
        return
    }
    log.Println("Order status updated successfully to 'Paid'")

    w.WriteHeader(http.StatusCreated)
    log.Println("Payment creation process completed")
    json.NewEncoder(w).Encode(map[string]string{"message": "Payment recorded successfully"})
}