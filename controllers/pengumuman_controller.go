package controllers

import (
	"database/sql"
	"cocopen-backend/dto"
	"cocopen-backend/middleware"
	"cocopen-backend/models"
	"cocopen-backend/services"
	"cocopen-backend/utils"
	"net/http"
	"strconv"
	"time"
)

// GetPengumumanHandler mengambil semua pengumuman
func GetPengumumanHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			utils.Error(w, http.StatusMethodNotAllowed, "Hanya GET yang diizinkan")
			return
		}

		pengumumans, err := services.GetAllPengumuman(db)
		if err != nil {
			utils.Error(w, http.StatusInternalServerError, "Gagal ambil pengumuman")
			return
		}

		var results []dto.PengumumanResponse
		now := time.Now()
		for _, p := range pengumumans {
			diff := now.Sub(p.CreatedAt)

			var waktuLalu string
			switch {
			case diff.Seconds() < 60:
				waktuLalu = "baru"
			case diff.Minutes() < 60:
				waktuLalu = strconv.Itoa(int(diff.Minutes())) + " menit yang lalu"
			case diff.Hours() < 24:
				waktuLalu = strconv.Itoa(int(diff.Hours())) + " jam yang lalu"
			case diff.Hours() < 48:
				waktuLalu = "1 hari yang lalu"
			default:
				waktuLalu = strconv.Itoa(int(diff.Hours()/24)) + " hari yang lalu"
			}

			results = append(results, dto.PengumumanResponse{
				IDPengumuman: p.IDPengumuman,
				Judul:        p.Judul,
				Isi:          p.Isi,
				CreatedAt:    p.CreatedAt.Format("2006-01-02 15:04:05"),
				WaktuLalu:    waktuLalu,
			})
		}

		utils.JSONResponse(w, http.StatusOK, results)
	}
}

// CreatePengumumanHandler dibuat oleh admin
func CreatePengumumanHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			utils.Error(w, http.StatusMethodNotAllowed, "Hanya POST yang diizinkan")
			return
		}

		claims, ok := r.Context().Value(middleware.UserContextKey).(*utils.Claims)
		if !ok || claims.Role != "admin" {
			utils.Error(w, http.StatusForbidden, "Akses ditolak")
			return
		}

		var req dto.PengumumanCreateRequest
		if err := utils.ParseAndValidate(r, &req); err != nil {
			utils.Error(w, http.StatusBadRequest, err.Error())
			return
		}

		pengumuman := models.Pengumuman{
			Judul: req.Judul,
			Isi:   req.Isi,
		}

		err := services.CreatePengumuman(db, pengumuman)
		if err != nil {
			utils.Error(w, http.StatusInternalServerError, "Gagal buat pengumuman")
			return
		}

		utils.JSONResponse(w, http.StatusCreated, map[string]interface{}{
			"success": true,
			"message": "Pengumuman berhasil dibuat",
		})
	}
}

// UpdatePengumumanHandler
func UpdatePengumumanHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			utils.Error(w, http.StatusMethodNotAllowed, "Hanya PUT yang diizinkan")
			return
		}

		claims, ok := r.Context().Value(middleware.UserContextKey).(*utils.Claims)
		if !ok || claims.Role != "admin" {
			utils.Error(w, http.StatusForbidden, "Akses ditolak")
			return
		}

		idStr := r.URL.Query().Get("id")
		id, err := strconv.Atoi(idStr)
		if err != nil || id <= 0 {
			utils.Error(w, http.StatusBadRequest, "ID tidak valid")
			return
		}

		var req dto.PengumumanCreateRequest
		if err := utils.ParseAndValidate(r, &req); err != nil {
			utils.Error(w, http.StatusBadRequest, err.Error())
			return
		}

		p, err := services.GetPengumumanByID(db, id)
		if err != nil {
			if err == sql.ErrNoRows {
				utils.Error(w, http.StatusNotFound, "Pengumuman tidak ditemukan")
				return
			}
			utils.Error(w, http.StatusInternalServerError, "Gagal ambil data")
			return
		}

		p.Judul = req.Judul
		p.Isi = req.Isi

		err = services.UpdatePengumuman(db, p)
		if err != nil {
			utils.Error(w, http.StatusInternalServerError, "Gagal update pengumuman")
			return
		}

		utils.JSONResponse(w, http.StatusOK, map[string]interface{}{
			"success": true,
			"message": "Pengumuman berhasil diperbarui",
		})
	}
}

// DeletePengumumanHandler
func DeletePengumumanHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			utils.Error(w, http.StatusMethodNotAllowed, "Hanya DELETE yang diizinkan")
			return
		}

		claims, ok := r.Context().Value(middleware.UserContextKey).(*utils.Claims)
		if !ok || claims.Role != "admin" {
			utils.Error(w, http.StatusForbidden, "Akses ditolak")
			return
		}

		idStr := r.URL.Query().Get("id")
		id, err := strconv.Atoi(idStr)
		if err != nil || id <= 0 {
			utils.Error(w, http.StatusBadRequest, "ID tidak valid")
			return
		}

		err = services.DeletePengumuman(db, id)
		if err != nil {
			utils.Error(w, http.StatusInternalServerError, "Gagal hapus pengumuman")
			return
		}

		utils.JSONResponse(w, http.StatusOK, map[string]interface{}{
			"success": true,
			"message": "Pengumuman berhasil dihapus",
		})
	}
}