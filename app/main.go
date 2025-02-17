package main

import (
	"log"
	"net/http"

	"github.com/aaaroz/social-media-downloader/configs"
	"github.com/aaaroz/social-media-downloader/internal/api"
)

func main() {
	cfg := configs.LoadConfig()

	router := api.SetupRoutes()

	log.Printf("Server running on port %s in %s mode\n", cfg.Port, cfg.Environment)
	log.Fatal(http.ListenAndServe(":"+cfg.Port, router))
}
