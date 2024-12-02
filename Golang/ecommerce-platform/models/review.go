package models

import "time"

type Review struct {
    ReviewID  int       `json:"review_id"`
    ProductID int       `json:"product_id"`
    UserID    int       `json:"user_id"`
    Rating    int       `json:"rating" validate:"required,min=1,max=5"`
    Comment   string    `json:"comment"`
    CreatedAt time.Time `json:"created_at"`
}
