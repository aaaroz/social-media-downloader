package tiktok

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"
)

type TikTokDownloader struct {
	URL      string
	Scrapers map[string]func() (map[string]interface{}, error)
}

func NewTikTokDownloader(url string) *TikTokDownloader {
	downloader := &TikTokDownloader{
		URL: url,
	}

	downloader.Scrapers = map[string]func() (map[string]interface{}, error){
		"api1": downloader.api1,
		"api2": downloader.api2,
	}
	return downloader
}

func (t *TikTokDownloader) Download() (map[string]interface{}, error) {
	for name, scraper := range t.Scrapers {
		result, err := scraper()
		if err == nil && result != nil {
			return result, nil
		}
		fmt.Printf("Scraper %s error: %v\n", name, err)
	}
	return nil, errors.New("All scraping method is failed!")
}

func (t *TikTokDownloader) api1() (map[string]interface{}, error) {
	apiURL := "https://ssstik.io/abc?url=dl"
	formData := url.Values{}
	formData.Set("id", t.URL)
	formData.Set("locale", "id")
	formData.Set("tt", "NmVQZUpk")

	res, err := http.PostForm(apiURL, formData)
	if err != nil {
		return nil, fmt.Errorf("API 2 Error: failed to send POST request - %v", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API 2 Error: unexpected status code %d", res.StatusCode)
	}

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("API 2 Error: failed to read response body - %v", err)
	}

	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(resBody))
	if err != nil {
		return nil, fmt.Errorf("API 2 Error: failed to parse HTML response - %v", err)
	}

	// Extract download links from the response
	downloadLinks := extractDownloadLinks(doc)

	// Extract the onclick attribute for the HD download button
	var hdDownloadURL string
	doc.Find("#hd_download").Each(func(i int, s *goquery.Selection) {
		onclick, exists := s.Attr("onclick")
		if exists {
			reHD := regexp.MustCompile(`downloadX\('(/abc\?url=[^']+)'\)`)
			hdMatches := reHD.FindStringSubmatch(onclick)
			if len(hdMatches) > 1 {
				hdDownloadURL = "https://ssstik.io" + hdMatches[1]
			}
		}
	})

	// Add the extracted HD download URL to the download links
	if hdDownloadURL != "" {
		downloadLinks = append(downloadLinks, map[string]interface{}{
			"type":     "download_video_hd",
			"url":      hdDownloadURL,
			"filename": fmt.Sprintf("tiktok_%s_hd.mp4", sanitizeFilename("HD")),
		})
	}

	// Extract username from the URL
	username := "unknown"
	if parts := strings.Split(t.URL, "@"); len(parts) > 1 {
		username = strings.Split(parts[1], "/")[0]
	}

	// Format the current date
	date := time.Now().Format("02 January 2006 15:04")

	return map[string]interface{}{
		"platform":  "tiktok",
		"caption":   "TikTok Video",
		"author":    "TikTok User",
		"username":  username,
		"img-thumb": "",
		"like":      0,
		"views":     0,
		"comments":  0,
		"date":      date,
		"downloads": downloadLinks,
	}, nil
}

func (t *TikTokDownloader) api2() (map[string]interface{}, error) {
	apiURL := "https://www.tikwm.com/api/"
	data := url.Values{}
	data.Set("url", t.URL)
	data.Set("count", "12")
	data.Set("cursor", "0")
	data.Set("web", "1")
	data.Set("hd", "1")
	resp, err := http.PostForm(apiURL, data)
	if err != nil {
		return nil, fmt.Errorf("API 1 Error: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API 1 Error: status code %d", resp.StatusCode)
	}
	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("API 1 Error: %v", err)
	}
	if result["code"].(float64) != 0 {
		return nil, fmt.Errorf("API 1 Error: %v", result["msg"])
	}
	dataMap := result["data"].(map[string]interface{})
	author := dataMap["author"].(map[string]interface{})
	createTime := time.Unix(int64(dataMap["create_time"].(float64)), 0)
	date := createTime.Format("02 January 2006 15:04")
	return map[string]interface{}{
		"platform":  "tiktok",
		"caption":   dataMap["title"],
		"author":    author["nickname"],
		"username":  author["unique_id"],
		"img-thumb": dataMap["cover"],
		"like":      int(dataMap["digg_count"].(float64)),
		"views":     int(dataMap["play_count"].(float64)),
		"comments":  int(dataMap["comment_count"].(float64)),
		"date":      date,
		"downloads": []map[string]interface{}{
			{
				"type":     "download_video_hd",
				"url":      "https://www.tikwm.com" + dataMap["hdplay"].(string),
				"filename": fmt.Sprintf("tiktok_%s_hd.mp4", author["unique_id"]),
			},
			{
				"type":     "download_video_480p",
				"url":      "https://www.tikwm.com" + dataMap["wmplay"].(string),
				"filename": fmt.Sprintf("tiktok_%s_watermark.mp4", author["unique_id"]),
			},
			{
				"type":     "download_audio",
				"url":      "https://www.tikwm.com" + dataMap["music"].(string),
				"filename": fmt.Sprintf("tiktok_%s_audio.mp3", author["unique_id"]),
			},
		},
	}, nil
}

// Helper function to extract download links from a goquery document
func extractDownloadLinks(doc *goquery.Document) []map[string]interface{} {
	downloadLinks := []map[string]interface{}{}
	doc.Find(".download_link").Each(func(i int, s *goquery.Selection) {
		link, exists := s.Attr("href")
		if exists {
			downloadType := strings.TrimSpace(s.Text())
			filename := fmt.Sprintf("tiktok_%s_%d.mp4", sanitizeFilename(downloadType), i)
			downloadLinks = append(downloadLinks, map[string]interface{}{
				"type":     downloadType,
				"url":      link,
				"filename": filename,
			})
		}
	})
	return downloadLinks
}

// Helper function to sanitize filenames
func sanitizeFilename(name string) string {
	return strings.ReplaceAll(strings.TrimSpace(name), " ", "_")
}
