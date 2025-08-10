package controllers

import (
    "cocopen-backend/dto"
    "cocopen-backend/services"
    "cocopen-backend/utils"
    "encoding/json"
    "errors"
    "net/http"
    "database/sql"
    "time"
)

func ResetPassword(db *sql.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        if r.Method != http.MethodPut {
            panic(errors.New("method tidak diizinkan"))
        }

        var req dto.ResetPasswordRequest
        if err := json.NewDecoder(r.Body).Decode(&req); 
		err != nil {
            panic(errors.New("format JSON tidak valid"))
        }

        if req.NewPassword == "" || req.Token == "" {
            panic(errors.New("token dan password baru harus diisi"))
        }

        userID, expiresAt, err := services.ValidatePasswordResetToken(db, req.Token)
        if err != nil {
            if err == sql.ErrNoRows {
                panic(errors.New("token reset password tidak valid"))
            }
            panic(err)
        }

        if expiresAt.Before(time.Now()) {
            panic(errors.New("token reset password sudah kedaluwarsa"))
        }

        hashedPassword, err := utils.HashPassword(req.NewPassword)
        if err != nil {
            panic(err)
        }

        err = services.ResetPassword(db, userID, hashedPassword)
        if err != nil {
            panic(err)
        }

        err = services.DeletePasswordResetToken(db, req.Token)
        if err != nil {
            panic(err)
        }

        utils.JSONResponse(w, http.StatusOK, map[string]interface{}{
            "success": true,
            "status":  http.StatusOK,
            "message": "Password berhasil direset.",
        })
    }
}
