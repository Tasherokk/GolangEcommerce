package models

type Category struct {
    CategoryID  int    `json:"category_id"`
    Name        string `json:"name" validate:"required"`
    Description string `json:"description"`
}
