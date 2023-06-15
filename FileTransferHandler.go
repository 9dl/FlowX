package FlowX

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

func DownloadFile(url string, outputPath string) error {
	response, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("failed to download file: %w", err)
	}
	defer response.Body.Close()

	file, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	// Create a new io.TeeReader to copy the response body and track progress
	progressReader := &ProgressReader{Reader: response.Body}

	_, err = io.Copy(file, progressReader)
	if err != nil {
		return fmt.Errorf("failed to copy response body to file: %w", err)
	}

	return nil
}

type ProgressReader struct {
	Reader   io.Reader
	Progress int64
}

func (pr *ProgressReader) Read(p []byte) (n int, err error) {
	n, err = pr.Reader.Read(p)
	pr.Progress += int64(n)
	fmt.Printf("Downloaded %d bytes\n", pr.Progress)
	return
}

var client = &http.Client{}

func UploadFile(filePath string, uploadURL string, headers http.Header) error {
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	request, err := http.NewRequest("POST", uploadURL, file)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	// Add custom headers to the request
	request.Header = headers

	response, err := client.Do(request)
	if err != nil {
		return fmt.Errorf("failed to perform request: %w", err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", response.StatusCode)
	}

	return nil
}
