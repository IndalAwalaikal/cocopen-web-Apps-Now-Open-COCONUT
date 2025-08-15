package controllers

import (
    "cocopen-backend/dto"
    "cocopen-backend/services"
    "cocopen-backend/utils"
    "database/sql"
    "encoding/json"
    "net/http"
)

func Login(db *sql.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        if r.Method != http.MethodPost {
            utils.Error(w, http.StatusMethodNotAllowed, "Metode tidak diizinkan")
            return
        }

        var req dto.LoginRequest
        if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
            utils.Error(w, http.StatusBadRequest, "Format JSON tidak valid")
            return
        }

        user, hashed, err := services.Login(db, req.Username)
        if err != nil {
            if err == sql.ErrNoRows {
                utils.Error(w, http.StatusNotFound, "Username tidak ditemukan")
                return
            }
            panic(err)
        }

        if !utils.CheckPassword(req.Password, hashed) {
            utils.Error(w, http.StatusUnauthorized, "Password salah")
            return
        }

        token := utils.GenerateToken(user.IDUser, user.Username, user.FullName, user.Role, user.ProfilePicture)

        utils.JSONResponse(w, http.StatusOK, map[string]interface{}{
            "success": true,
            "status":  http.StatusOK,
            "message": "Anda berhasil login",
            "token":   token,
        })
    }
}
