package controllers

import (
    "cocopen-backend/services"
    "cocopen-backend/utils"
    "errors"
    "net/http"
    "time"
    "database/sql"
)

func VerifyEmail(db *sql.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        token := r.URL.Query().Get("token")
        if token == "" {
            panic(errors.New("token verifikasi harus disertakan"))
        }

        userID, expiresAt, err := services.ValidateVerificationToken(db, token)
        if err != nil {
            if err == sql.ErrNoRows {
                panic(errors.New("token verifikasi tidak valid"))
            }
            panic(err)
        }

        if expiresAt.Before(time.Now()) {
            panic(errors.New("token verifikasi sudah kedaluwarsa"))
        }

        if err := services.VerifyEmail(db, userID); 
		err != nil {
            panic(err)
        }

        if err := services.DeleteVerificationToken(db, token); 
		err != nil {
            panic(err)
        }

        utils.JSONResponse(w, http.StatusOK, map[string]interface{}{
            "success": true,
            "status":  http.StatusOK,
            "message": "Email berhasil diverifikasi.",
        })
    }
}
