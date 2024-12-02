package models

import "time"

type ShoppingCart struct {
    CartID    int       `json:"cart_id"`
    UserID    int       `json:"user_id"`
    CreatedAt time.Time `json:"created_at"`
}
