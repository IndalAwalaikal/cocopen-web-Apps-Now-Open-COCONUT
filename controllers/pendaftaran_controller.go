package controllers

import (
	"cocopen-backend/dto"
	"cocopen-backend/models"
	"cocopen-backend/services"
	"cocopen-backend/utils"
	"database/sql"
	"errors"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
)

func CreatePendaftar(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			utils.Error(w, http.StatusMethodNotAllowed, "Metode tidak diizinkan")
			return
		}

		err := r.ParseMultipartForm(10 << 20)
		if err != nil {
			utils.Error(w, http.StatusBadRequest, "Gagal memproses form")
			return
		}

		var req dto.CreatePendaftarRequest
		req.NamaLengkap = r.FormValue("nama_lengkap")
		req.AsalKampus = r.FormValue("asal_kampus")
		req.Prodi = r.FormValue("prodi")
		req.Semester = r.FormValue("semester")
		req.NoWA = r.FormValue("no_wa")
		req.Domisili = r.FormValue("domisili")
		req.AlamatSekarang = r.FormValue("alamat_sekarang")
		req.TinggalDengan = r.FormValue("tinggal_dengan")
		req.AlasanMasuk = r.FormValue("alasan_masuk")
		req.PengetahuanCoconut = r.FormValue("pengetahuan_coconut")


		file, header, err := r.FormFile("foto")
		if err != nil {
    		if err == http.ErrMissingFile {
        	utils.Error(w, http.StatusBadRequest, "Foto wajib diunggah")
        	return
    		} else {
        	utils.Error(w, http.StatusBadRequest, "Gagal membaca file foto: " + err.Error())
        	return
   			}
		}

		defer file.Close()

		ext := strings.ToLower(filepath.Ext(header.Filename))
		if ext != ".jpg" && ext != ".jpeg" && ext != ".png" && ext != ".gif" {
    		utils.Error(w, http.StatusBadRequest, "Format file tidak didukung. Gunakan JPG, PNG, atau GIF.")
    		return
		}

		if header.Size > 2<<20 {
    		utils.Error(w, http.StatusBadRequest, "Ukuran file maksimal 2 MB")
    		return
		}

		fotoName, err := utils.UploadFoto(file, header, utils.FotoPendaftarPath)
		if err != nil {
    		utils.Error(w, http.StatusInternalServerError, "Gagal upload foto")
    		return
		}
		req.FotoPath = fotoName

		if req.NamaLengkap == "" || req.AsalKampus == "" || req.Prodi == "" || req.NoWA == "" {
			utils.Error(w, http.StatusBadRequest, "Field wajib diisi")
			return
		}

		if err := services.CreatePendaftar(db, req); err != nil {
			panic(errors.New("gagal menambahkan pendaftar: " + err.Error()))
		}

		utils.JSONResponse(w, http.StatusCreated, map[string]interface{}{
			"success": true,
			"status":  http.StatusCreated,
			"message": "Pendaftar berhasil dibuat",
		})
	}
}

func GetAllPendaftar(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rows, err := services.GetAllPendaftar(db)
		if err != nil {
			panic(err)
		}
		defer rows.Close()

		var result []models.Pendaftar
		for rows.Next() {
			var p models.Pendaftar
			if err := rows.Scan(
				&p.IDPendaftar, &p.NamaLengkap, &p.AsalKampus, &p.Prodi, &p.Semester, &p.NoWA,
				&p.Domisili, &p.AlamatSekarang, &p.TinggalDengan, &p.AlasanMasuk,
				&p.PengetahuanCoconut, &p.FotoPath, &p.CreatedAt, &p.UpdatedAt, &p.Status,
			); err != nil {
				panic(err)
			}
			result = append(result, p)
		}

		utils.JSONResponse(w, http.StatusOK, result)
	}
}

func GetPendaftarByID(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
            utils.Error(w, http.StatusMethodNotAllowed, "Hanya metode GET yang diizinkan")
            return
        }
		
		idStr := r.URL.Query().Get("id")
		if idStr == "" {
			utils.Error(w, http.StatusBadRequest, "Parameter id wajib diisi")
			return
		}

		idPendaftar, err := strconv.Atoi(idStr)
		if err != nil {
			utils.Error(w, http.StatusBadRequest, "ID tidak valid")
			return
		}

		p, err := services.GetPendaftarByID(db, idPendaftar)
		if err != nil {
			if err == sql.ErrNoRows {
				utils.Error(w, http.StatusNotFound, "Pendaftar tidak ditemukan")
				return
			}
			panic(err)
		}

		utils.JSONResponse(w, http.StatusOK, p)
	}
}

func UpdatePendaftar(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			utils.Error(w, http.StatusMethodNotAllowed, "Metode tidak diizinkan")
			return
		}

		err := r.ParseMultipartForm(10 << 20) 
        if err != nil {
            utils.Error(w, http.StatusBadRequest, "Gagal memproses form")
            return
        }

		idStr := r.FormValue("id_pendaftar")
		id, err := strconv.Atoi(idStr)
		if err != nil || id == 0 {
			utils.Error(w, http.StatusBadRequest, "ID pendaftar tidak valid")
			return
		}

		status := r.FormValue("status")
		if status == "" {
			utils.Error(w, http.StatusBadRequest, "Status wajib diisi")
			return
		}

		if err := services.UpdatePendaftar(db, id, status); 
		err != nil {
            panic(errors.New("gagal memperbarui pendaftar: " + err.Error()))
        }

		utils.JSONResponse(w, http.StatusOK, map[string]interface{}{
			"success": true,
			"status":  http.StatusOK,
			"message": "Pendaftar berhasil diperbarui",
		})
	}
}

func DeletePendaftar(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			utils.Error(w, http.StatusMethodNotAllowed, "Metode tidak diizinkan")
			return
		}

		idStr := r.URL.Query().Get("id")
		if idStr == "" {
			utils.Error(w, http.StatusBadRequest, "Parameter id wajib diisi")
			return
		}

		id, err := strconv.Atoi(idStr)
		if err != nil {
			utils.Error(w, http.StatusBadRequest, "ID tidak valid")
			return
		}

		oldData, err := services.GetPendaftarByID(db, id)
		if err == nil && oldData.FotoPath != "" {
			utils.HapusFoto(utils.FotoPendaftarPath, oldData.FotoPath)
		}

		if err := services.DeletePendaftar(db, id); err != nil {
			if err == sql.ErrNoRows {
				utils.Error(w, http.StatusNotFound, "Pendaftar tidak ditemukan")
				return
			}
			panic(err)
		}

		utils.JSONResponse(w, http.StatusOK, map[string]interface{}{
			"success": true,
			"message": "Pendaftar berhasil dihapus",
		})
	}
}
