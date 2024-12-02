// middleware/auth_middleware.go
package middleware

import (
	"context"
	"log"
	"net/http"
	"strings"

	"github.com/Tasherokk/GolangProject.git/utils"
)


type contextKey string

const (
    UserIDKey contextKey = "user_id"
    RoleKey   contextKey = "role"
)


func AuthMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Retrieve the Authorization header
        authHeader := r.Header.Get("Authorization")
        if authHeader == "" {
            http.Error(w, "Missing Authorization header", http.StatusUnauthorized)
            return
        }

        log.Printf("Authorization header: %s", authHeader)
        // Extract the token
        tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
        claims, err := utils.ValidateToken(tokenStr)
        if err != nil {
            http.Error(w, "Invalid token: "+err.Error(), http.StatusUnauthorized)
            return
        }
        
        log.Printf("Token validated. UserID: %d, Role: %s", claims.UserID, claims.Role)
        // Store values in context using the custom key type
        ctx := context.WithValue(r.Context(), UserIDKey, claims.UserID)
        ctx = context.WithValue(ctx, RoleKey, claims.Role)
        next.ServeHTTP(w, r.WithContext(ctx))
    })
}
