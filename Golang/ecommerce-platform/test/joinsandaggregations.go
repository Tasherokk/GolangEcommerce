package test

import (
    "github.com/Tasherokk/GolangProject.git/database"
)


type OrderWithUser struct {
    OrderID     int
    OrderDate   string
    Status      string
    TotalAmount float64
    UserID      int
    Username    string
    Email       string
}


func getOrdersWithUsers() ([]OrderWithUser, error) {
    sqlStatement := `
    SELECT o.order_id, o.order_date, o.status, o.total_amount,
           u.user_id, u.username, u.email
    FROM orders o
    INNER JOIN users u ON o.user_id = u.user_id;`

    rows, err := database.DB.Query(sqlStatement)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var orders []OrderWithUser
    for rows.Next() {
        var order OrderWithUser
        err := rows.Scan(&order.OrderID, &order.OrderDate, &order.Status, &order.TotalAmount,
            &order.UserID, &order.Username, &order.Email)
        if err != nil {
            return nil, err
        }
        orders = append(orders, order)
    }
    return orders, nil
}


func getTotalSales() (float64, error) {
    sqlStatement := `SELECT SUM(total_amount) FROM orders;`
    var totalSales float64
    err := database.DB.QueryRow(sqlStatement).Scan(&totalSales)
    if err != nil {
        return 0, err
    }
    return totalSales, nil
}


type UserOrderCount struct {
    UserID    int
    Username  string
    Email     string
    OrderCount int
}

func getUsersWithOrderCount() ([]UserOrderCount, error) {
    sqlStatement := `
    SELECT u.user_id, u.username, u.email, COUNT(o.order_id) AS order_count
    FROM users u
    LEFT JOIN orders o ON u.user_id = o.user_id
    GROUP BY u.user_id, u.username, u.email;`

    rows, err := database.DB.Query(sqlStatement)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var userOrderCounts []UserOrderCount
    for rows.Next() {
        var uoc UserOrderCount
        err := rows.Scan(&uoc.UserID, &uoc.Username, &uoc.Email, &uoc.OrderCount)
        if err != nil {
            return nil, err
        }
        userOrderCounts = append(userOrderCounts, uoc)
    }
    return userOrderCounts, nil
}



type RecentOrder struct {
    OrderID   int
    OrderDate string
    TotalAmount float64
}

func getRecentOrders(limit int) ([]RecentOrder, error) {
    sqlStatement := `
    SELECT order_id, order_date, total_amount
    FROM orders
    ORDER BY order_date DESC
    LIMIT $1;`

    rows, err := database.DB.Query(sqlStatement, limit)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var orders []RecentOrder
    for rows.Next() {
        var order RecentOrder
        err := rows.Scan(&order.OrderID, &order.OrderDate, &order.TotalAmount)
        if err != nil {
            return nil, err
        }
        orders = append(orders, order)
    }
    return orders, nil
}
