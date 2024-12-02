// routes/routes.go
package routes

import (
    "github.com/gorilla/mux"
    "github.com/Tasherokk/GolangProject.git/controllers"
    "github.com/Tasherokk/GolangProject.git/middleware"
)

func RegisterRoutes(router *mux.Router) {
    // Маршруты, не требующие аутентификации
    router.HandleFunc("/api/register", controllers.RegisterUser).Methods("POST")
    router.HandleFunc("/api/login", controllers.LoginUser).Methods("POST")

    router.HandleFunc("/api/products", controllers.GetAllProducts).Methods("GET")
    router.HandleFunc("/api/products/{product_id}", controllers.GetProductByID).Methods("GET")

    router.HandleFunc("/api/categories", controllers.GetAllCategories).Methods("GET")
    router.HandleFunc("/api/categories/{category_id}", controllers.GetCategoryByID).Methods("GET")

    // Подмаршруты с аутентификацией
    api := router.PathPrefix("/api").Subrouter()
    api.Use(middleware.AuthMiddleware)

    // Маршруты для пользователя
    api.HandleFunc("/users/me", controllers.GetCurrentUser).Methods("GET")
    api.HandleFunc("/users/me", controllers.UpdateUser).Methods("PUT")
    api.HandleFunc("/users/me", controllers.DeleteUser).Methods("DELETE")

    // Маршруты для продуктов (создание, обновление, удаление - требуют прав администратора)
    api.HandleFunc("/products", controllers.CreateProduct).Methods("POST")
    api.HandleFunc("/products/{product_id}", controllers.UpdateProduct).Methods("PUT")
    api.HandleFunc("/products/{product_id}", controllers.DeleteProduct).Methods("DELETE")
    api.HandleFunc("/seller/products", controllers.GetSellerProducts).Methods("GET")


    // Маршруты для категорий (создание, обновление, удаление - требуют прав администратора)
    api.HandleFunc("/categories", controllers.CreateCategory).Methods("POST")
    api.HandleFunc("/categories/{category_id}", controllers.UpdateCategory).Methods("PUT")
    api.HandleFunc("/categories/{category_id}", controllers.DeleteCategory).Methods("DELETE")

    // Маршруты для заказов
    api.HandleFunc("/orders", controllers.GetUserOrders).Methods("GET")
    api.HandleFunc("/orders/checkout", controllers.CreateOrder).Methods("POST")
    api.HandleFunc("/orders/{order_id}", controllers.GetOrderByID).Methods("GET")
    api.HandleFunc("/orders/{order_id}/status", controllers.UpdateOrderStatus).Methods("PATCH")
    api.HandleFunc("/seller/orders", controllers.GetSellerOrders).Methods("GET")

    // Маршруты для элементов заказа
    api.HandleFunc("/orders/{order_id}/items", controllers.GetOrderItems).Methods("GET")
    api.HandleFunc("/orders/{order_id}/items", controllers.AddOrderItems).Methods("POST")

    // Маршруты для корзины
    api.HandleFunc("/cart", controllers.GetUserCart).Methods("GET")

    // Маршруты для элементов корзины
    api.HandleFunc("/cart/items", controllers.GetCartItems).Methods("GET")
    api.HandleFunc("/cart/items", controllers.AddCartItem).Methods("POST")
    // api.HandleFunc("/cart/items/{item_id}", controllers.UpdateCartItem).Methods("PUT")
    api.HandleFunc("/cart/items/{item_id}", controllers.DeleteCartItem).Methods("DELETE")
    api.HandleFunc("/cart/items/{item_id}/decrease", controllers.DecreaseCartItemQuantity).Methods("PATCH")
    api.HandleFunc("/cart/items/{item_id}/increase", controllers.IncreaseCartItemQuantity).Methods("PATCH")


    // Маршруты для платежей
    api.HandleFunc("/payments", controllers.GetUserPayments).Methods("GET")
    api.HandleFunc("/payments", controllers.CreatePayment).Methods("POST")

    // Маршруты для отзывов
    api.HandleFunc("/products/{product_id}/reviews", controllers.GetReviewsByProduct).Methods("GET")
    api.HandleFunc("/products/{product_id}/reviews", controllers.CreateReview).Methods("POST")
    // api.HandleFunc("/reviews/{review_id}", controllers.UpdateReview).Methods("PUT")
    // api.HandleFunc("/reviews/{review_id}", controllers.DeleteReview).Methods("DELETE")

    // Маршруты для адресов пользователя
    api.HandleFunc("/addresses", controllers.GetUserAddresses).Methods("GET")
    api.HandleFunc("/addresses", controllers.AddUserAddress).Methods("POST")
    // api.HandleFunc("/addresses/{address_id}", controllers.UpdateUserAddress).Methods("PUT")
    // api.HandleFunc("/addresses/{address_id}", controllers.DeleteUserAddress).Methods("DELETE")

    // Маршруты для изображений продукта
    api.HandleFunc("/products/{product_id}/images", controllers.GetProductImages).Methods("GET")
    api.HandleFunc("/products/{product_id}/images", controllers.AddProductImage).Methods("POST")
    // api.HandleFunc("/products/{product_id}/images/{image_id}", controllers.DeleteProductImage).Methods("DELETE")

    // Маршруты для ролей (требуют прав администратора)
    adminRoutes := api.PathPrefix("").Subrouter()
    adminRoutes.Use(middleware.AdminOnly)

    
    adminRoutes.HandleFunc("/admin/roles", controllers.GetAllRoles).Methods("GET")
    adminRoutes.HandleFunc("/admin/roles", controllers.CreateRole).Methods("POST")
    adminRoutes.HandleFunc("/admin/users/{user_id}/role", controllers.UpdateUserRole).Methods("PUT")
    adminRoutes.HandleFunc("/admin/users/{user_id}", controllers.DeleteUserByID).Methods("DELETE")
    adminRoutes.HandleFunc("/admin/users", controllers.GetAllUsers).Methods("GET")



    // Маршруты для журналов аудита (требуют прав администратора)
    adminRoutes.HandleFunc("/audit-logs", controllers.GetAuditLogs).Methods("GET")

    // Дополнительные маршруты для сессий, если требуется
    // api.HandleFunc("/sessions", controllers.CreateSession).Methods("POST")
    // api.HandleFunc("/sessions", controllers.DeleteSession).Methods("DELETE")
	
}
