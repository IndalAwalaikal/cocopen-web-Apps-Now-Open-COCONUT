package dto

type CreatePendaftarRequest struct {
    NamaLengkap         string `json:"nama_lengkap"`
    AsalKampus          string `json:"asal_kampus"`
    Prodi               string `json:"prodi"`
    Semester            int    `json:"semester"`
    NoWA                string `json:"no_wa"`
    Domisili            string `json:"domisili"`
    AlamatSekarang      string `json:"alamat_sekarang"`
    TinggalDengan       string `json:"tinggal_dengan"`
    AlasanMasuk         string `json:"alasan_masuk"`
    PengetahuanCoconut  string `json:"pengetahuan_coconut"`
    FotoPath            string `json:"foto_path"`
}

type UpdatePendaftarStatusRequest struct {
    IDPendaftar int    `json:"id_pendaftar" validate:"required"`
    Status      string `json:"status" validate:"required"`
}
