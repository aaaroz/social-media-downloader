package services

import (
	"errors"

	"github.com/aaaroz/social-media-downloader/internal/downloader"
)

type DownloaderService struct {
	downloader *downloader.Downloader
}

func NewDownloaderService() *DownloaderService {
	return &DownloaderService{
		downloader: downloader.NewDownloader(),
	}
}

func (ds *DownloaderService) DownloadMedia(url string) (interface{}, error) {
	if url == "" {
		return nil, errors.New("URL cannot be empty")
	}
	return ds.downloader.Download(url)
}
