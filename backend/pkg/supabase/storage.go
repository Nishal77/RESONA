package supabase

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"time"

	"github.com/google/uuid"
)

type StorageClient struct {
	baseURL    string
	serviceKey string
	bucket     string
	httpClient *http.Client
}

func NewStorageClient(supabaseURL, serviceKey, bucket string) *StorageClient {
	return &StorageClient{
		baseURL:    supabaseURL,
		serviceKey: serviceKey,
		bucket:     bucket,
		httpClient: &http.Client{Timeout: 30 * time.Second},
	}
}

func (s *StorageClient) Upload(fileContent []byte, originalFilename, contentType string) (string, error) {
	ext := filepath.Ext(originalFilename)
	objectPath := fmt.Sprintf("uploads/%s%s", uuid.New().String(), ext)

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", objectPath)
	if err != nil {
		return "", fmt.Errorf("create form file: %w", err)
	}
	if _, err = io.Copy(part, bytes.NewReader(fileContent)); err != nil {
		return "", fmt.Errorf("copy file: %w", err)
	}
	writer.Close()

	url := fmt.Sprintf("%s/storage/v1/object/%s/%s", s.baseURL, s.bucket, objectPath)
	req, err := http.NewRequest(http.MethodPost, url, body)
	if err != nil {
		return "", fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+s.serviceKey)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("upload request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		b, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("supabase storage error %d: %s", resp.StatusCode, string(b))
	}

	var result struct {
		Key string `json:"Key"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("decode response: %w", err)
	}

	publicURL := fmt.Sprintf("%s/storage/v1/object/public/%s/%s", s.baseURL, s.bucket, objectPath)
	return publicURL, nil
}
