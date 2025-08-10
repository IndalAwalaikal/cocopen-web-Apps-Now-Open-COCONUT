// middleware/role.go
package middleware

import (
    "net/http"
    "cocopen-backend/utils"
)

func Role(requiredRole string) func(http.HandlerFunc) http.HandlerFunc {
    return func(next http.HandlerFunc) http.HandlerFunc {
        return func(w http.ResponseWriter, r *http.Request) {
            claims, ok := r.Context().Value("user_claims").(*utils.Claims)
            if !ok {
                panic("Akses ditolak: tidak terautentikasi")
            }

            if claims.Role != requiredRole {
                panic("Akses ditolak: role tidak sesuai")
            }

            next(w, r)
        }
    }
}