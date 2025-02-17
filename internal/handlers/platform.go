package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/aaaroz/social-media-downloader/pkg/socialmedia"
)

// SupportedPlatformsHandler menampilkan daftar platform yang didukung
func SupportedPlatformsHandler(w http.ResponseWriter, r *http.Request) {
	// Ambil daftar platform yang didukung dari socialmedia.PlatformPatterns
	supportedPlatforms := make([]string, 0, len(socialmedia.PlatformPatterns))
	for platform := range socialmedia.PlatformPatterns {
		supportedPlatforms = append(supportedPlatforms, platform)
	}

	// Buat respons JSON
	response := map[string]interface{}{
		"supported_platforms": supportedPlatforms,
	}

	// Set header dan kirim respons JSON
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
