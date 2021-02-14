package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

type GDriveConfig struct {
	ApiKey            string `json:"api_key"`
	URL               string `json:"url"`
	FilenameSubstring string `json:"filename_substring"`
}

type FileList struct {
	Items []ListedFile `json:"items"`
}

type ListedFile struct {
	Kind        string `json:"kind"`
	Title       string `json:"title"`
	DownloadURL string `json:"downloadUrl"`
}

func getID(link string) string {
	chunks := strings.Split(link, "/")
	return chunks[len(chunks)-1]
}

func getCsvUrl(config GDriveConfig, id string) (string, error) {
	req, err := http.NewRequest("GET", config.URL, nil)
	if err != nil {
		return "", fmt.Errorf("gdrive: error constructing list request %w", err)
	}

	q := req.URL.Query()
	q.Add("q", fmt.Sprintf(`'%s' in parents`, id))
	q.Add("key", config.ApiKey)

	req.URL.RawQuery = q.Encode()
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("gdrive: error performing list request %w", err)
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", fmt.Errorf("gdrive: error reading list response %w", err)
	}
	var list FileList
	err = json.Unmarshal(body, &list)
	if err != nil {
		return "", fmt.Errorf("gdrive: error parsing list response %w", err)
	}
	for _, f := range list.Items {
		if strings.Contains(f.Title, config.FilenameSubstring) {
			return f.DownloadURL, nil
		}
	}
	return "", fmt.Errorf("gdrive: file with substring %s not found", config.FilenameSubstring)
}

func downloadCSV(config GDriveConfig, link, downloadPath string) error {
	url, err := getCsvUrl(config, getID(link))
	if err != nil {
		return err
	}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return fmt.Errorf("gdrive: error constructing csv download request %w", err)
	}
	req.Header.Add("Accept", "text/csv")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("gdrive: error performing csv download request %w", err)
	}
	defer res.Body.Close()
	out, err := os.Create(downloadPath)
	if err != nil {
		return fmt.Errorf("gdrive: error allocating csv %w", err)
	}
	defer out.Close()
	_, err = io.Copy(out, res.Body)
	if err != nil {
		return fmt.Errorf("gdrive: error downloading csv %w", err)
	}
	return nil
}
