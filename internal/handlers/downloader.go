package handlers

import (
	"encoding/json"
	"errors"
	"io"
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

func DownloadVideoHandler(w http.ResponseWriter, r *http.Request) {
	videoURL := r.URL.Query().Get("url")
	if videoURL == "" {
		http.Error(w, "Missing video URL", http.StatusBadRequest)
		return
	}

	// Fetch the video
	resp, err := http.Get(videoURL)
	if err != nil || resp.StatusCode != http.StatusOK {
		http.Error(w, "Failed to fetch video", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// Set headers for file download
	w.Header().Set("Content-Disposition", "attachment; filename=video.mp4")
	w.Header().Set("Content-Type", "video/mp4")

	// Stream the video to the client
	_, err = io.Copy(w, resp.Body)
	if err != nil {
		http.Error(w, "Failed to stream video", http.StatusInternalServerError)
		return
	}
}

// Helper function to respond with a JSON-formatted error
func respondWithError(w http.ResponseWriter, statusCode int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	response := map[string]string{"error": message}
	json.NewEncoder(w).Encode(response)
}
