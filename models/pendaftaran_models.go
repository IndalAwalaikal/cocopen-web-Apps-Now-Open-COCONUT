// model/pendaftar.go
package models

import "time"

type Pendaftar struct {
    IDPendaftar         int       `json:"id_pendaftar"`
    NamaLengkap         string    `json:"nama_lengkap" validate:"required"`
    AsalKampus          string    `json:"asal_kampus" validate:"required"`
    Prodi               string    `json:"prodi" validate:"required"`
    Semester            string    `json:"semester" validate:"required"`
    NoWA                string    `json:"no_wa" validate:"required"`
    Domisili            string    `json:"domisili" validate:"required"`
    AlamatSekarang      string    `json:"alamat_sekarang" validate:"required"`
    TinggalDengan       string    `json:"tinggal_dengan" validate:"required"`
    AlasanMasuk         string    `json:"alasan_masuk" validate:"required"`
    PengetahuanCoconut  string    `json:"pengetahuan_coconut"`
	FotoPath 			string    `json:"foto_path"`
    CreatedAt           time.Time `json:"created_at"`
    UpdatedAt           time.Time `json:"updated_at"`
    Status              string    `json:"status"`
    UserID              *int      `json:"user_id,omitempty"`
}