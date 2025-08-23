package models

import "time"

type Pengumuman struct {
    IDPengumuman int       `json:"id_pengumuman"`
    Judul        string    `json:"judul"`
    Isi          string    `json:"isi"`
    CreatedAt    time.Time `json:"created_at"`
    UpdatedAt    time.Time `json:"updated_at"`
}