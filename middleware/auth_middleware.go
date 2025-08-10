// middleware/auth.go
package middleware

import (
    "context"
    "net/http"
    "strings"

    "github.com/golang-jwt/jwt/v5"
    "cocopen-backend/utils"
)

type contextKey string
var userContextKey contextKey = "user_claims"

func Auth(next http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        authHeader := r.Header.Get("Authorization")
        if authHeader == "" {
            panic("Authorization header diperlukan")
        }

        if !strings.HasPrefix(authHeader, "Bearer ") {
            panic("format token harus Bearer <token>")
        }

        tokenString := strings.TrimPrefix(authHeader, "Bearer ")
        claims := &utils.Claims{}

        token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (any, error) {
            if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
                return nil, jwt.ErrSignatureInvalid
            }
            return utils.Secret, nil
        })

        if err != nil || !token.Valid {
            panic("token tidak valid")
        }

        ctx := context.WithValue(r.Context(), userContextKey, claims)
        r = r.WithContext(ctx)

        next(w, r)
    }
}