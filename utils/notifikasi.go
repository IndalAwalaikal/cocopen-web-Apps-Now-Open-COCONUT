// utils/notifikasi.go
package utils

import "cocopen-backend/dto"

var NotifikasiChan = make(chan dto.NotifikasiEvent, 100)

func BroadcastNotifikasi(judul, isi string) {
    event := dto.NotifikasiEvent{
        Judul: judul,
        Isi:   isi,
    }
    select {
    case NotifikasiChan <- event:
    default:
    }
}