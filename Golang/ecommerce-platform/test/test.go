package test

import (
    "fmt"
    "log"
    "github.com/Tasherokk/GolangProject.git/models"
    "github.com/Tasherokk/GolangProject.git/controllers"
)


func Test() {

	fmt.Println("Hello, World!")

	fmt.Println("\nData Types:")
	dataTypesDemo()
	fmt.Println("\nControl Structures:")
	controlStructuresDemo()
	fmt.Println("\nFunctions:")
	functionsDemo()
	fmt.Println("\nOOP:")
	oopDemo()
	fmt.Println()
}


func TestDb() {

    prod, err := controllers.GetProductByIDTest(1)
    if err != nil {
        fmt.Println("Error:", err)
        return
    }
    fmt.Println("Product:", prod)

	newUser := models.User{
        Username:     "john_doe",
        PasswordHash: "hashed_password",
        Email:        "john@example.com",
        Role:         1,
    }

	err = controllers.DeleteUserByUsername("john_doe")
	if err != nil {
		log.Fatalf("Failed to delete user: %v", err)
	}

    err = controllers.CreateUser(&newUser)
    if err != nil {
        log.Fatalf("Failed to create user: %v", err)
    }

    user, err := controllers.GetUserByID(newUser.UserID)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println("Retrieved User:", user)
}


func TestDb2() {

	orders, err := getOrdersWithUsers()
    if err != nil {
        log.Fatalf("Error fetching orders with users: %v", err)
    }

    fmt.Println("Orders with User Details:")
    for _, order := range orders {
        fmt.Printf("Order ID: %d, Total: %.2f, User: %s (%s)\n", 
            order.OrderID, order.TotalAmount, order.Username, order.Email)
    }


	totalSales, err := getTotalSales()
    if err != nil {
        log.Fatalf("Error calculating total sales: %v", err)
    }

    fmt.Printf("Total Sales Amount: $%.2f\n", totalSales)

	userOrderCounts, err := getUsersWithOrderCount()
    if err != nil {
        log.Fatalf("Error fetching users with order count: %v", err)
    }

    fmt.Println("Users and Their Order Counts:")
    for _, uoc := range userOrderCounts {
        fmt.Printf("User: %s (%s), Orders: %d\n", uoc.Username, uoc.Email, uoc.OrderCount)
    }


	recentOrders, err := getRecentOrders(5)
    if err != nil {
        log.Fatalf("Error fetching recent orders: %v", err)
    }

    fmt.Println("Recent Orders:")
    for _, order := range recentOrders {
        fmt.Printf("Order ID: %d, Date: %s, Total: %.2f\n", 
            order.OrderID, order.OrderDate, order.TotalAmount)
    }


}