// middleware/recovery.go
package middleware

import (
    "log"
    "net/http"
    "runtime"
    "strings"
    "fmt"

    "cocopen-backend/dto"
    "cocopen-backend/utils"
)

func Recovery(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        defer func() {
            if err := recover(); err != nil {
                log.Printf("terjadi kesalahan: %v", err)
                for i := 2; ; i++ {
                    _, file, line, ok := runtime.Caller(i)
                    if !ok {
                        break
                    }
                    log.Printf("  %s:%d", file, line)
                }

                var response dto.ErrorResponse

                if appErr, ok := err.(dto.ErrorResponse); ok {
                    utils.JSONResponse(w, appErr.Status, appErr)
                    return
                }
                errorMsg := strings.ToLower(string(fmt.Sprint(err)))

                switch {
                case strings.Contains(errorMsg, "method tidak diizinkan"):
                    response = dto.ErrorResponse{Success: false, Status: http.StatusMethodNotAllowed, Message: "Metode tidak diizinkan"}
                case strings.Contains(errorMsg, "tidak terautentikasi"), strings.Contains(errorMsg, "authorization header"), strings.Contains(errorMsg, "token tidak valid"):
                    response = dto.ErrorResponse{Success: false, Status: http.StatusUnauthorized, Message: "Akses ditolak: tidak terautentikasi"}
                case strings.Contains(errorMsg, "akses ditolak"), strings.Contains(errorMsg, "bukan admin"), strings.Contains(errorMsg, "hanya untuk user"):
                    response = dto.ErrorResponse{Success: false, Status: http.StatusForbidden, Message: "Akses ditolak: role tidak sesuai"}
                case strings.Contains(errorMsg, "tidak ditemukan"):
                    response = dto.ErrorResponse{Success: false, Status: http.StatusNotFound, Message: "Data tidak ditemukan"}
                case strings.Contains(errorMsg, "wajib"), strings.Contains(errorMsg, "tidak valid"), strings.Contains(errorMsg, "format"), strings.Contains(errorMsg, "sudah digunakan"):
                    response = dto.ErrorResponse{Success: false, Status: http.StatusBadRequest, Message: "Permintaan tidak valid"}
                default:
                    response = dto.ErrorResponse{Success: false, Status: http.StatusInternalServerError, Message: "Terjadi kesalahan internal"}
                }

                utils.JSONResponse(w, response.Status, response)
            }
        }()
        next.ServeHTTP(w, r)
    })
}