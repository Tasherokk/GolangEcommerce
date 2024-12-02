package test

import (
    "fmt"
    "github.com/Tasherokk/GolangProject.git/models"
)


func Display(p models.Product) {
    fmt.Printf("Product ID: %d\nName: %s\nPrice: $%.2f\n", p.ProductID, p.Name, p.Price)
}

func UpdateStock(p *models.Product, newStock int) {
    p.Stock = newStock
}


func oopDemo() {
    product := models.Product{
        ProductID:   1,
        Name:        "Laptop",
        Description: "A high-end laptop",
        Price:       1500.00,
        Stock:       10,
    }

    Display(product)
    UpdateStock(&product, 8)
    fmt.Println("Updated Stock:", product.Stock)

}
