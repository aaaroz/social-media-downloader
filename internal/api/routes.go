package api

import (
	"github.com/gorilla/mux"
	"net/http"

	"github.com/aaaroz/social-media-downloader/internal/handlers"
)

func SetupRoutes() *mux.Router {
	r := mux.NewRouter()

	// Public routes
	r.HandleFunc("/", homeHandler).Methods("GET")
	r.HandleFunc("/api/platform", handlers.SupportedPlatformsHandler).Methods("GET")
	r.HandleFunc("/api/download", handlers.DownloadHandler).Methods("POST")

	r.PathPrefix("/").Handler(http.FileServer(http.Dir("./public")))

	return r
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hello, World!"))
}
