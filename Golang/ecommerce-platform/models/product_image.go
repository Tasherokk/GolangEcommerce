package models

import "time"

type ProductImage struct {
    ImageID   int       `json:"image_id"`
    ProductID int       `json:"product_id"`
    ImageURL  string    `json:"image_url"`
    CreatedAt time.Time `json:"created_at"`
}
