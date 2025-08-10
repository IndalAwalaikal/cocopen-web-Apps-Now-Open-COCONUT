// controllers/dashboard_controller.go
package controllers

import (
    "net/http"
    "cocopen-backend/dto"
    "cocopen-backend/utils"
)

func AdminDashboard() http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        if r.Method != http.MethodGet {
            panic(dto.ErrorResponse{
                Success: false,
                Status:  http.StatusMethodNotAllowed,
                Message: "Metode tidak diizinkan",
            })
        }

        claims, ok := r.Context().Value("user_claims").(*utils.Claims)
        if !ok {
            panic(dto.ErrorResponse{
                Success: false,
                Status:  http.StatusUnauthorized,
                Message: "Akses ditolak: tidak terautentikasi",
            })
        }

        if claims.Role != "admin" {
            panic(dto.ErrorResponse{
                Success: false,
                Status:  http.StatusForbidden,
                Message: "Akses ditolak: Anda bukan admin",
            })
        }

        response := dto.DashboardResponse{
            Message: "Selamat datang di Dashboard Admin!",
            User: dto.UserSummary{
                ID:       claims.ID,
                Username: claims.Username,
                Role:     claims.Role,
            },
            Access: "admin",
            Data:   "Anda memiliki akses penuh ke sistem.",
        }

        utils.JSONResponse(w, http.StatusOK, dto.SuccessResponse{
            Success: true,
            Status:  http.StatusOK,
            Message: "Berhasil memuat dashboard admin",
            Data:    response,
        })
    }
}

func UserDashboard() http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        if r.Method != http.MethodGet {
            panic(dto.ErrorResponse{
                Success: false,
                Status:  http.StatusMethodNotAllowed,
                Message: "Metode tidak diizinkan",
            })
        }

        claims, ok := r.Context().Value("user_claims").(*utils.Claims)
        if !ok {
            panic(dto.ErrorResponse{
                Success: false,
                Status:  http.StatusUnauthorized,
                Message: "Akses ditolak: tidak terautentikasi",
            })
        }

        if claims.Role != "user" {
            panic(dto.ErrorResponse{
                Success: false,
                Status:  http.StatusForbidden,
                Message: "Akses ditolak: halaman ini hanya untuk user",
            })
        }

        response := dto.DashboardResponse{
            Message: "Halo, selamat datang di Dashboard Anda!",
            User: dto.UserSummary{
                ID:       claims.ID,
                Username: claims.Username,
                Role:     claims.Role,
            },
            Access: "user",
            Data:   "Anda memiliki akses sebagai pengguna biasa.",
        }

        utils.JSONResponse(w, http.StatusOK, dto.SuccessResponse{
            Success: true,
            Status:  http.StatusOK,
            Message: "Berhasil memuat dashboard user",
            Data:    response,
        })
    }
}