// utils/jwt.go
package utils

import (
    "os"
    "time"
    "github.com/golang-jwt/jwt/v5"
)

var Secret []byte

func InitJWTSecret() {
    secret := os.Getenv("JWT_SECRET")
    if secret == "" {
        panic("Environment variable JWT_SECRET tidak ditemukan. Harap isi di file .env")
    }
    if len(secret) < 32 {
        panic("JWT_SECRET terlalu pendek! Minimal 32 karakter untuk keamanan.")
    }
    Secret = []byte(secret)
}

type Claims struct {
    ID       int    `json:"id"`
    Username string `json:"username"`
    Role     string `json:"role"`
    jwt.RegisteredClaims
}

func GenerateToken(userID int, username, role string) string {
    claims := &Claims{
        ID:       userID,
        Username: username,
        Role:     role,
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
            IssuedAt:  jwt.NewNumericDate(time.Now()),
        },
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    signedToken, err := token.SignedString(Secret)
    if err != nil {
        panic("Gagal membuat token JWT: " + err.Error())
    }
    return signedToken
}