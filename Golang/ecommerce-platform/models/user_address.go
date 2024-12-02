package models

type UserAddress struct {
    AddressID int    `json:"address_id"`
    UserID    int    `json:"user_id"`
    Street    string `json:"street"`
    City      string `json:"city"`
    State     string `json:"state"`
    ZipCode   string `json:"zip_code"`
}
