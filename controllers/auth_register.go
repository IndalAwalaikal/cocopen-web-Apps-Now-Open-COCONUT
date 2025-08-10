package controllers

import (
    "cocopen-backend/dto"
    "cocopen-backend/services"
    "cocopen-backend/utils"
    "database/sql"
    "encoding/json"
    "errors"
    "net/http"
    "time"
)

func Register(db *sql.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        if r.Method != http.MethodPost {
            panic(errors.New("method tidak diizinkan"))
        }

        var req dto.RegisterRequest
        if err := json.NewDecoder(r.Body).Decode(&req); 
		err != nil {
            panic(errors.New("format JSON tidak valid"))
        }

        if req.Username == "" || req.Email == "" || req.Password == "" {
            panic(errors.New("semua field wajib diisi"))
        }

        if !utils.IsValidUsername(req.Username) {
            panic(errors.New("username tidak valid (minimal 3 karakter, huruf/angka/_)"))
        }
        if !utils.IsValidEmail(req.Email) {
            panic(errors.New("format email tidak valid"))
        }
        if !utils.IsValidPassword(req.Password) {
            panic(errors.New("password harus minimal 8 karakter, mengandung huruf besar, angka, dan simbol"))
        }

        _, _, _, _, _, err := services.GetUserByUsername(db, req.Username)
        if err == nil {
            panic(errors.New("username sudah digunakan"))
        } else if err != sql.ErrNoRows {
            panic(err)
        }

        _, err = services.GetUserByEmail(db, req.Email)
        if err == nil {
            panic(errors.New("email sudah digunakan"))
        } else if err != sql.ErrNoRows {
            panic(err)
        }

        hashedPassword, err := utils.HashPassword(req.Password)
        if err != nil {
            panic(err)
        }

        userID, err := services.Register(db, req.Username, req.Email, hashedPassword, "user")
        if err != nil {
            panic(err)
        }

        verificationToken := utils.GenerateRandomToken(32)
        expiresAt := time.Now().Add(24 * time.Hour)

        err = services.GenerateVerificationToken(db, userID, verificationToken, expiresAt)
        if err != nil {
            panic(errors.New("gagal membuat token verifikasi: " + err.Error()))
        }

        err = utils.SendVerificationEmail(req.Email, verificationToken)
        if err != nil {
            panic(errors.New("gagal mengirim email verifikasi: " + err.Error()))
        }

        utils.JSONResponse(w, http.StatusCreated, map[string]interface{}{
            "success": true,
            "status":  http.StatusCreated,
            "message": "Akun berhasil dibuat. Silakan cek email untuk verifikasi.",
            "user": map[string]string{
                "username": req.Username,
                "email":    req.Email,
            },
        })
    }
}
