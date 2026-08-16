package service

import (
	"bytes"
	"context"
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"log"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"deepphotos/backend/internal/model"
	"deepphotos/backend/internal/repository"
	"deepphotos/backend/internal/storage"
	"github.com/google/uuid"
)

type PhotoService struct {
	photoRepo     *repository.PhotoRepository
	storageClient *storage.StorageClient
}

func NewPhotoService(photoRepo *repository.PhotoRepository, storageClient *storage.StorageClient) *PhotoService {
	return &PhotoService{
		photoRepo:     photoRepo,
		storageClient: storageClient,
	}
}

func DetectCategory(filename string, contentType string) string {
	mime := strings.ToLower(contentType)
	ext := strings.ToLower(filepath.Ext(filename))

	if strings.HasPrefix(mime, "video/") || ext == ".mp4" || ext == ".mov" || ext == ".webm" || ext == ".mkv" || ext == ".avi" {
		return "video"
	}
	if strings.HasPrefix(mime, "application/") || strings.HasPrefix(mime, "text/") || ext == ".pdf" || ext == ".txt" || ext == ".docx" || ext == ".doc" || ext == ".zip" {
		return "document"
	}
	return "image"
}

func (s *PhotoService) UploadPhoto(ctx context.Context, filename string, contentType string, reader io.Reader, size int64) (*model.Photo, error) {
	buf := new(bytes.Buffer)
	if _, err := io.Copy(buf, reader); err != nil {
		return nil, fmt.Errorf("failed to read file bytes: %w", err)
	}

	photoID := "med_" + uuid.New().String()[:12]
	ext := filepath.Ext(filename)

	category := DetectCategory(filename, contentType)
	objectKey, thumbnailKey := storage.GenerateObjectKey(category, filename)

	if err := s.storageClient.UploadBytes(ctx, objectKey, buf.Bytes(), contentType); err != nil {
		return nil, fmt.Errorf("failed to upload object to MinIO: %w", err)
	}

	imgWidth, imgHeight := 0, 0
	if category == "image" {
		imgConfig, _, err := image.DecodeConfig(bytes.NewReader(buf.Bytes()))
		if err == nil {
			imgWidth = imgConfig.Width
			imgHeight = imgConfig.Height
		}
		go func() {
			if err := s.storageClient.UploadBytes(context.Background(), thumbnailKey, buf.Bytes(), "image/webp"); err != nil {
				log.Printf("Warning: failed to store WebP thumbnail %s: %v", thumbnailKey, err)
			}
		}()
	}

	title := strings.TrimSuffix(filename, ext)

	photo := &model.Photo{
		ID:           photoID,
		Filename:     filename,
		ObjectKey:    objectKey,
		ThumbnailKey: thumbnailKey,
		MimeType:     contentType,
		FileType:     category,
		Size:         int64(buf.Len()),
		Width:        imgWidth,
		Height:       imgHeight,
		ExifModel:    "DeepPhotos Ingest",
		Title:        title,
		TakenAt:      "Recently",
		IsFavorite:   false,
		IsDeleted:    false,
	}

	if err := s.photoRepo.Create(photo); err != nil {
		return nil, fmt.Errorf("failed to save media metadata to SQLite: %w", err)
	}

	return photo, nil
}

func (s *PhotoService) UploadPhotoFromURL(ctx context.Context, imageURL string) (*model.Photo, error) {
	resp, err := http.Get(imageURL)
	if err != nil {
		return nil, fmt.Errorf("failed to download image from URL: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP request returned status %d", resp.StatusCode)
	}

	contentType := resp.Header.Get("Content-Type")
	if contentType == "" {
		contentType = "image/jpeg"
	}

	filename := filepath.Base(imageURL)
	if strings.Contains(filename, "?") {
		filename = strings.Split(filename, "?")[0]
	}
	if filename == "" || filename == "." {
		filename = fmt.Sprintf("web_ingest_%d.jpg", time.Now().Unix())
	}

	return s.UploadPhoto(ctx, filename, contentType, resp.Body, resp.ContentLength)
}

func (s *PhotoService) DeleteSinglePhoto(ctx context.Context, id string) error {
	photo, err := s.photoRepo.DeleteSingle(id)
	if err != nil || photo == nil {
		return fmt.Errorf("media not found or delete failed: %w", err)
	}

	go func() {
		_ = s.storageClient.DeleteObject(context.Background(), photo.ObjectKey)
		_ = s.storageClient.DeleteObject(context.Background(), photo.ThumbnailKey)
	}()

	return nil
}

func (s *PhotoService) DeleteBatchPhotos(ctx context.Context, ids []string) (int, error) {
	photos, err := s.photoRepo.DeleteBatch(ids)
	if err != nil {
		return 0, fmt.Errorf("batch delete failed in SQLite: %w", err)
	}

	var keysToDelete []string
	for _, p := range photos {
		if p.ObjectKey != "" {
			keysToDelete = append(keysToDelete, p.ObjectKey)
		}
		if p.ThumbnailKey != "" {
			keysToDelete = append(keysToDelete, p.ThumbnailKey)
		}
	}

	go func() {
		if len(keysToDelete) > 0 {
			_ = s.storageClient.DeleteObjects(context.Background(), keysToDelete)
		}
	}()

	return len(photos), nil
}
