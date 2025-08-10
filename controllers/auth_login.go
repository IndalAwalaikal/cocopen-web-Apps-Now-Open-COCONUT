package controllers

import (
    "cocopen-backend/dto"
    "cocopen-backend/services"
    "cocopen-backend/utils"
    "database/sql"
    "encoding/json"
    "errors"
    "net/http"
)

func Login(db *sql.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        if r.Method != http.MethodPost {
            panic(errors.New("method tidak diizinkan"))
        }

        var req dto.LoginRequest
        if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
            panic(errors.New("format JSON tidak valid"))
        }

        user, hashed, err := services.Login(db, req.Username)
        if err != nil {
            if err == sql.ErrNoRows {
                panic(errors.New("username tidak ditemukan atau belum diverifikasi"))
            }
            panic(err)
        }

        if !utils.CheckPassword(req.Password, hashed) {
            panic(errors.New("password salah"))
        }

        token := utils.GenerateToken(user.IDUser, user.Username, user.Role)

        utils.JSONResponse(w, http.StatusOK, map[string]interface{}{
            "success": true,
            "status":  http.StatusOK,
            "message": "Anda berhasil login",
            "token":   token,
            "user": map[string]interface{}{
                "id":       user.IDUser,
                "username": user.Username,
                "role":     user.Role,
                "email":    user.Email,
            },
        })
    }
}
