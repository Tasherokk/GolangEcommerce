package models

import "time"

type Order struct {
    OrderID     int       `json:"order_id"`
    UserID      int       `json:"user_id"`
    OrderDate   time.Time `json:"order_date"`
    Status      string    `json:"status"`
    TotalAmount float64   `json:"total_amount"`
}
