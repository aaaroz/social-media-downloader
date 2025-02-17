package models

type DownloadRequest struct {
	URL string `json:"url"`
}

type DownloadResponse struct {
	URL    string `json:"url"`
	Status string `json:"status"`
}
