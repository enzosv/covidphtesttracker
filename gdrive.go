package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
)

type GDriveConfig struct {
	ApiKey          string `json:"api_key"`
	URL             string `json:"url"`
	TestSubstring   string `json:"test_substring"`
	LinkSubstring   string `json:"link_substring"`
	ReadmeSubstring string `json:"readme_substring"`
	DailyURL        string `json:"daily_url"`
}

type FileList struct {
	Items []ListedFile `json:"items"`
}

type ListedFile struct {
	Title       string `json:"title"`
	DownloadURL string `json:"downloadUrl"`
	MimeType    string `json:"mimeType"`
}

func getRedirectURL(link string) (*url.URL, error) {
	req, err := http.NewRequest("GET", link, nil)
	if err != nil {
		return nil, fmt.Errorf("redirect: error consturcting request: %w", err)
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("redirect: error performing request %w", err)
	}
	defer res.Body.Close()
	return res.Request.URL, nil
}

func getFolderID(link url.URL) string {
	chunks := strings.Split(link.Path, "/")
	return chunks[len(chunks)-1]
}

func (c GDriveConfig) getFileUrl(folderID, mime, fileSubstring string) (string, error) {
	req, err := http.NewRequest("GET", c.URL, nil)
	if err != nil {
		return "", fmt.Errorf("gdrive: error constructing list request %w", err)
	}

	q := req.URL.Query()
	q.Add("q", fmt.Sprintf(`'%s' in parents`, folderID))
	q.Add("key", c.ApiKey)

	req.URL.RawQuery = q.Encode()
	fmt.Println(req.URL.String())
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
		if f.MimeType == mime && strings.Contains(f.Title, fileSubstring) {
			return f.DownloadURL, nil
		}
	}

	return "", fmt.Errorf("gdrive: file with substring %s not found", fileSubstring)
}

func (c GDriveConfig) download(folderID, downloadPath, mime, substring string) (string, error) {
	log.Printf("\n\t[INFO] Downloading %s\n\n", substring)
	url, err := c.getFileUrl(folderID, mime, substring)
	if err != nil {
		return "", err
	}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", fmt.Errorf("gdrive: error constructing %s download request %w", substring, err)
	}
	req.Header.Add("Accept", mime)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("gdrive: error performing %s download request %w", substring, err)
	}
	defer res.Body.Close()
	out, err := os.Create(downloadPath)
	if err != nil {
		return "", fmt.Errorf("gdrive: error allocating %s %w", substring, err)
	}
	defer out.Close()
	_, err = io.Copy(out, res.Body)
	if err != nil {
		return "", fmt.Errorf("gdrive: error downloading %s %w", substring, err)
	}
	return url, nil
}

func (c GDriveConfig) GetReadmeLinks(date string) ([]string, error) {
	pdfPath := fmt.Sprintf("%s.pdf", date)
	defer os.Remove(pdfPath)
	redirected, err := getRedirectURL(c.DailyURL)
	if err != nil {
		return nil, err
	}
	folderID := getFolderID(*redirected)
	defer os.Remove(pdfPath)
	_, err = c.download(folderID, pdfPath, "application/pdf", c.ReadmeSubstring)
	if err != nil {
		return nil, err
	}
	return findSubstringInPDF(pdfPath, c.LinkSubstring)
}

func (c GDriveConfig) GetTestFolderLink(links []string, csvPath string) (string, string, error) {
	var err error
	for _, l := range links {
		redirected, err := getRedirectURL(l)
		if err != nil {
			fmt.Println(err)
			continue
		}
		folderID := getFolderID(*redirected)
		fileurl, err := c.download(folderID, csvPath, "text/csv", c.TestSubstring)
		if err != nil {
			fmt.Println(err)
			continue
		}
		return redirected.String(), fileurl, nil
	}
	return "", "", err
}
