package models

import "time"

type Product struct {
    ProductID    int       `json:"product_id"`
    Name         string    `json:"name" validate:"required"`
    Description  string    `json:"description"`
    Price        float64   `json:"price" validate:"required,gt=0"`
    Stock        int       `json:"stock" validate:"gte=0"`
    CategoryID   int       `json:"category_id"`
    CreatedAt    time.Time `json:"created_at"`
    CategoryName string    `json:"category_name"`
    ImageURL     string    `json:"image_url,omitempty"`
    SellerID     int       `json:"seller_id"`
}
