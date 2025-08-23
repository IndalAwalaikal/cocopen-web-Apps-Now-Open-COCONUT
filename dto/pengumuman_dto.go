package dto

type PengumumanResponse struct {
    IDPengumuman int    `json:"id_pengumuman"`
    Judul        string `json:"judul"`
    Isi          string `json:"isi"`
    CreatedAt    string `json:"created_at"`
    WaktuLalu    string `json:"waktu_lalu"`
}

type PengumumanCreateRequest struct {
    Judul   string `json:"judul" validate:"required,min=3"`
    Isi     string `json:"isi" validate:"required,min=10"`
    IdJadwal *int  `json:"id_jadwal,omitempty"` 
}