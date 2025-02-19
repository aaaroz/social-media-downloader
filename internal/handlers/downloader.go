package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/aaaroz/social-media-downloader/internal/downloader"
	"github.com/aaaroz/social-media-downloader/internal/models"
	"github.com/aaaroz/social-media-downloader/internal/services"
)

func DownloadHandler(w http.ResponseWriter, r *http.Request) {
	var req models.DownloadRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	downloaderService := services.NewDownloaderService()
	result, err := downloaderService.DownloadMedia(req.URL)
	if err != nil {
		// Check if the error is due to an unsupported platform
		var unsupportedErr *downloader.UnsupportedPlatformError
		if errors.As(err, &unsupportedErr) {
			respondWithError(w, http.StatusUnsupportedMediaType, err.Error())
			return
		}
		// Handle other errors as internal server errors
		respondWithError(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

// Helper function to respond with a JSON-formatted error
func respondWithError(w http.ResponseWriter, statusCode int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	response := map[string]string{"error": message}
	json.NewEncoder(w).Encode(response)
}
