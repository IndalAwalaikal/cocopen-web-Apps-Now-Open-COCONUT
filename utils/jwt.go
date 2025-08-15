// utils/jwt.go
package utils

import (
    "os"
    "time"
    "fmt"
    "github.com/golang-jwt/jwt/v5"
)

var Secret []byte

func InitJWTSecret() {
    secret := os.Getenv("JWT_SECRET")
    fmt.Printf(" [InitJWTSecret] JWT_SECRET: '%s' (length: %d)\n", secret, len(secret))

    if secret == "" {
        panic("JWT_SECRET tidak ditemukan")
    }
    if len(secret) < 32 {
        panic("JWT_SECRET terlalu pendek!")
    }
    Secret = []byte(secret)
    fmt.Printf(" [InitJWTSecret] Secret di-set: %v\n", Secret)
}

type Claims struct {
    IDUser   int    `json:"id_user"`
    Username string `json:"username"`
    FullName string `json:"full_name"`
    Role     string `json:"role"`
    ProfilePicture  string `json:"profile_picture"`
    jwt.RegisteredClaims
}

func GenerateToken(userID int, username, fullName, role, profilePicture string) string {
    claims := &Claims{
        IDUser:         userID,
        Username:       username,
        FullName:       fullName,
        Role:           role,
        ProfilePicture: profilePicture,
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
            IssuedAt:  jwt.NewNumericDate(time.Now()),
        },
    }

    fmt.Printf("[GenerateToken] Membuat token dengan claims: %+v\n", claims)

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    signedToken, err := token.SignedString(Secret)
    if err != nil {
        panic("Gagal membuat token JWT: " + err.Error())
    }

    fmt.Printf("[GenerateToken] Token berhasil dibuat (50 karakter pertama): %s...\n", signedToken[:50])
    return signedToken
}