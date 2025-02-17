package downloader

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/aaaroz/social-media-downloader/internal/downloader/tiktok"
	"github.com/aaaroz/social-media-downloader/pkg/socialmedia"
)

type UnsupportedPlatformError struct {
	Platform string
}

func (e *UnsupportedPlatformError) Error() string {
	return fmt.Sprintf("platform %s is not supported yet", e.Platform)
}

type Downloader struct {
	platformPatterns map[string]*regexp.Regexp
	downloaders      map[string]interface{}
}

func NewDownloader() *Downloader {
	return &Downloader{
		platformPatterns: socialmedia.PlatformPatterns,
		downloaders: map[string]interface{}{
			"tiktok": tiktok.NewTikTokDownloader,
		},
	}
}

func (d *Downloader) GetPlatform(url string) (string, error) {
	if url == "" {
		return "", errors.New("URL can't be empty.")
	}
	for platform, pattern := range d.platformPatterns {
		if pattern.MatchString(url) {
			return platform, nil
		}
	}
	return "", &UnsupportedPlatformError{Platform: "unknown"}
}

func (d *Downloader) Download(url string) (interface{}, error) {
	if url == "" {
		return nil, errors.New("URL can't be empty")
	}
	url = strings.TrimSpace(url)
	if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
		url = "https://" + url
	}
	platform, err := d.GetPlatform(url)
	if err != nil {
		var unsupportedErr *UnsupportedPlatformError
		if errors.As(err, &unsupportedErr) {
			return nil, unsupportedErr
		}
		return nil, err
	}
	downloaderFunc, ok := d.downloaders[platform]
	if !ok {
		return nil, &UnsupportedPlatformError{Platform: platform}
	}
	switch platform {
	case "tiktok":
		downloader := downloaderFunc.(func(string) *tiktok.TikTokDownloader)(url)
		return downloader.Download()
	default:
		return nil, &UnsupportedPlatformError{Platform: platform}
	}
}
