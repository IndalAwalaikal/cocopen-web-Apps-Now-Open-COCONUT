package models

import (
	"time"
)

type SoalTest struct {
	IDSoal       int       `json:"id_soal"`
	Nomor        int       `json:"nomor"`
	Pertanyaan   string    `json:"pertanyaan"`
	PilihanA     string    `json:"pilihan_a"`
	PilihanB     string    `json:"pilihan_b"`
	PilihanC     string    `json:"pilihan_c"`
	PilihanD     string    `json:"pilihan_d"`
	JawabanBenar string    `json:"jawaban_benar"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type Test struct {
	IDTest       int        `json:"id_test"`
	Judul        string     `json:"judul"`
	Deskripsi    *string    `json:"deskripsi,omitempty"`
	DurasiMenit  int        `json:"durasi_menit"`
	WaktuMulai   *time.Time `json:"waktu_mulai,omitempty"`
	WaktuSelesai *time.Time `json:"waktu_selesai,omitempty"`
	Aktif        bool       `json:"aktif"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}

type HasilTest struct {
	IDHasil      int        `json:"id_hasil"`
	UserID       int        `json:"user_id"`
	PendaftarID  int        `json:"pendaftar_id"`
	IDTest       int        `json:"id_test"`
	SkorBenar    int        `json:"skor_benar"`
	SkorSalah    int        `json:"skor_salah"`
	Nilai        float64    `json:"nilai"`
	WaktuMulai   time.Time  `json:"waktu_mulai"`
	WaktuSelesai *time.Time `json:"waktu_selesai,omitempty"`
	DurasiMenit  *int       `json:"durasi_menit,omitempty"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}

type JawabanUser struct {
	IDJawaban   int       `json:"id_jawaban"`
	IDHasil     int       `json:"id_hasil"`
	IDSoal      int       `json:"id_soal"`
	JawabanUser string    `json:"jawaban_user"`
	IsBenar     bool      `json:"is_benar"`
	CreatedAt   time.Time `json:"created_at"`
}
