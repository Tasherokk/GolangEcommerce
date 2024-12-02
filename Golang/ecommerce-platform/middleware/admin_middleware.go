// middleware/admin_middleware.go
package middleware

import (
	"log"
	"net/http"
)

func AdminOnly(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        role, ok := r.Context().Value(RoleKey).(string)
        if !ok {
            log.Println("Role not found in context")
            http.Error(w, "Access denied: Role not found", http.StatusForbidden)
            return
        }
        log.Printf("Role from context: %s", role)
        if role != "Admin" {
            http.Error(w, "Access denied: Insufficient permissions", http.StatusForbidden)
            return
        }
        next.ServeHTTP(w, r)
    })
}
