package storage

import (
	"bytes"
	"context"
	"deepphotos/backend/internal/config"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type StorageClient struct {
	client   *minio.Client
	bucket   string
	IsOnline bool
	diskDir  string
}

// GenerateObjectKey formats object path as: {category}/{uuid4_1}/{uuid4_2}/{uuid4_3}/{filename}
func GenerateObjectKey(category, filename string) (objectKey string, thumbnailKey string) {
	cleanCategory := strings.ToLower(strings.TrimSpace(category))
	if cleanCategory == "" {
		cleanCategory = "image"
	}

	u1 := uuid.New().String()
	u2 := uuid.New().String()
	u3 := uuid.New().String()

	cleanFilename := filepath.Base(filename)
	if cleanFilename == "" || cleanFilename == "." {
		cleanFilename = "file.dat"
	}

	objectKey = fmt.Sprintf("%s/%s/%s/%s/%s", cleanCategory, u1, u2, u3, cleanFilename)
	thumbnailKey = fmt.Sprintf("thumbnails/%s/%s/%s/%s/thumb_%s.webp", cleanCategory, u1, u2, u3, cleanFilename)
	return objectKey, thumbnailKey
}

// GenerateAvatarKey formats user profile photo path under dedicated avatars/ bucket path: avatars/{userID}/{uuid}/{filename}
func GenerateAvatarKey(userID, filename string) string {
	cleanFilename := filepath.Base(filename)
	if cleanFilename == "" || cleanFilename == "." {
		cleanFilename = "avatar.jpg"
	}
	u := uuid.New().String()[:8]
	return fmt.Sprintf("avatars/%s/%s_%s", userID, u, cleanFilename)
}

func InitMinio(cfg *config.Config) (*StorageClient, error) {
	client, err := minio.New(cfg.MinioEndpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.MinioAccessKey, cfg.MinioSecretKey, ""),
		Secure: cfg.MinioUseSSL,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to initialize MinIO client: %w", err)
	}

	diskDir := filepath.Join("data", "storage")
	_ = os.MkdirAll(filepath.Join(diskDir, "originals"), 0755)
	_ = os.MkdirAll(filepath.Join(diskDir, "thumbnails"), 0755)

	s := &StorageClient{
		client:  client,
		bucket:  cfg.MinioBucket,
		diskDir: diskDir,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	exists, err := client.BucketExists(ctx, cfg.MinioBucket)
	if err != nil {
		s.IsOnline = false
		log.Printf("ℹ️ MinIO storage offline at %s (Local disk fallback active at data/storage/)", cfg.MinioEndpoint)
		return s, nil
	}

	if !exists {
		err = client.MakeBucket(context.Background(), cfg.MinioBucket, minio.MakeBucketOptions{})
		if err != nil {
			s.IsOnline = false
			log.Printf("ℹ️ MinIO storage offline (Local disk fallback active at data/storage/)")
			return s, nil
		}
		log.Printf("Created MinIO bucket: %s", cfg.MinioBucket)
	} else {
		log.Printf("Connected to MinIO bucket: %s", cfg.MinioBucket)
	}

	s.IsOnline = true
	return s, nil
}

func (s *StorageClient) UploadObject(ctx context.Context, objectKey string, reader io.Reader, size int64, contentType string) error {
	if s.IsOnline {
		_, err := s.client.PutObject(ctx, s.bucket, objectKey, reader, size, minio.PutObjectOptions{
			ContentType: contentType,
		})
		if err == nil {
			return nil
		}
		log.Printf("MinIO upload failed, falling back to disk: %v", err)
	}

	localPath := filepath.Join(s.diskDir, objectKey)
	_ = os.MkdirAll(filepath.Dir(localPath), 0755)
	outFile, err := os.Create(localPath)
	if err != nil {
		return err
	}
	defer outFile.Close()

	_, err = io.Copy(outFile, reader)
	return err
}

func (s *StorageClient) UploadBytes(ctx context.Context, objectKey string, data []byte, contentType string) error {
	reader := bytes.NewReader(data)
	return s.UploadObject(ctx, objectKey, reader, int64(len(data)), contentType)
}

func (s *StorageClient) GetObject(ctx context.Context, objectKey string) (io.ReadCloser, int64, string, error) {
	if s.IsOnline {
		obj, err := s.client.GetObject(ctx, s.bucket, objectKey, minio.GetObjectOptions{})
		if err == nil {
			stat, statErr := obj.Stat()
			if statErr == nil {
				return obj, stat.Size, stat.ContentType, nil
			}
			obj.Close()
		}
	}

	localPath := filepath.Join(s.diskDir, objectKey)
	f, err := os.Open(localPath)
	if err != nil {
		return nil, 0, "", err
	}

	stat, err := f.Stat()
	if err != nil {
		f.Close()
		return nil, 0, "", err
	}

	ext := strings.ToLower(filepath.Ext(objectKey))
	contentType := "application/octet-stream"
	switch ext {
	case ".jpg", ".jpeg":
		contentType = "image/jpeg"
	case ".png":
		contentType = "image/png"
	case ".webp":
		contentType = "image/webp"
	case ".gif":
		contentType = "image/gif"
	case ".mp4":
		contentType = "video/mp4"
	case ".mov":
		contentType = "video/quicktime"
	case ".webm":
		contentType = "video/webm"
	case ".pdf":
		contentType = "application/pdf"
	case ".txt":
		contentType = "text/plain"
	}

	return f, stat.Size(), contentType, nil
}

func (s *StorageClient) DeleteObject(ctx context.Context, objectKey string) error {
	if s.IsOnline {
		_ = s.client.RemoveObject(ctx, s.bucket, objectKey, minio.RemoveObjectOptions{})
	}
	localPath := filepath.Join(s.diskDir, objectKey)
	return os.Remove(localPath)
}

func (s *StorageClient) DeleteObjects(ctx context.Context, objectKeys []string) error {
	for _, key := range objectKeys {
		if key != "" {
			_ = s.DeleteObject(ctx, key)
		}
	}
	return nil
}
