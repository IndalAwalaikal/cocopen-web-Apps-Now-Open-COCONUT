package controllers

import (
	"cocopen-backend/utils"
	"encoding/json"
	"net/http"
	"time"
)

func NotifikasiStreamHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("X-Accel-Buffering", "no")

	flusher, ok := w.(http.Flusher)
	if !ok {
		utils.Error(w, http.StatusInternalServerError, "Streaming tidak didukung")
		return
	}

	w.Write([]byte(":\n\n"))
	flusher.Flush()

	for {
		select {
		case event := <-utils.NotifikasiChan:
			data, err := json.Marshal(event)
			if err != nil {
				continue
			}
			w.Write([]byte(" " + string(data) + "\n\n"))
			flusher.Flush()

		case <-time.After(30 * time.Second):
			w.Write([]byte(":\n\n"))
			flusher.Flush()
		}
	}
}
