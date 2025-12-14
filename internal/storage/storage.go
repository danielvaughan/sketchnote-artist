package storage

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"cloud.google.com/go/storage"
)

// Store defines the interface for saving and retrieving files.
type Store interface {
	// Save saves the content to the storage.
	// folder: "sketchnotes" or "visual-briefs" (used as subdirectory locally)
	// filename: the name of the file
	// data: the content bytes
	Save(ctx context.Context, folder, filename string, data []byte) error

	// GetPublicURL returns the public URL for the file.
	// folder: "sketchnotes" or "visual-briefs"
	// filename: the name of the file
	GetPublicURL(folder, filename string) string

	// Exists checks if a file exists in the storage.
	Exists(ctx context.Context, folder, filename string) (bool, error)

	// Get returns a reader for the file content.
	Get(ctx context.Context, folder, filename string) (io.ReadCloser, error)
}

// DiskStore implements Store for local file system.
type DiskStore struct{}

func (s *DiskStore) Save(ctx context.Context, folder, filename string, data []byte) error {
	// Ensure directory exists
	if err := os.MkdirAll(folder, 0755); err != nil {
		return fmt.Errorf("failed to create directory %s: %w", folder, err)
	}
	fullPath := filepath.Join(folder, filename)
	return os.WriteFile(fullPath, data, 0644)
}

func (s *DiskStore) GetPublicURL(folder, filename string) string {
	// For local dev, we assume the server maps /images/ to sketchnotes/ or images/
	if folder == "sketchnotes" || folder == "images" {
		return fmt.Sprintf("/images/%s", filename)
	}
	// No public URL for visual briefs locally by default in this requirements scope
	return ""
}

func (s *DiskStore) Exists(ctx context.Context, folder, filename string) (bool, error) {
	fullPath := filepath.Join(folder, filename)
	_, err := os.Stat(fullPath)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func (s *DiskStore) Get(ctx context.Context, folder, filename string) (io.ReadCloser, error) {
	fullPath := filepath.Join(folder, filename)
	return os.Open(fullPath)
}

// GCSStore implements Store for Google Cloud Storage.
type GCSStore struct {
	Client       *storage.Client
	BriefsBucket string
	ImagesBucket string
}

func NewGCSStore(ctx context.Context, briefsBucket, imagesBucket string) (*GCSStore, error) {
	client, err := storage.NewClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create storage client: %w", err)
	}
	return &GCSStore{
		Client:       client,
		BriefsBucket: briefsBucket,
		ImagesBucket: imagesBucket,
	}, nil
}

func (s *GCSStore) getBucketName(folder string) (string, error) {
	switch folder {
	case "visual-briefs":
		return s.BriefsBucket, nil
	case "sketchnotes", "images":
		return s.ImagesBucket, nil
	default:
		return "", fmt.Errorf("unknown storage folder type: %s", folder)
	}
}

func (s *GCSStore) Save(ctx context.Context, folder, filename string, data []byte) error {
	bucketName, err := s.getBucketName(folder)
	if err != nil {
		return err
	}

	wc := s.Client.Bucket(bucketName).Object(filename).NewWriter(ctx)
	if _, err := wc.Write(data); err != nil {
		return fmt.Errorf("failed to write object to bucket %s: %w", bucketName, err)
	}
	if err := wc.Close(); err != nil {
		return fmt.Errorf("failed to close writer for object %s: %w", filename, err)
	}
	return nil
}

func (s *GCSStore) GetPublicURL(folder, filename string) string {
	bucketName, err := s.getBucketName(folder)
	if err != nil {
		return ""
	}
	// Assuming public access or consistent naming
	return fmt.Sprintf("https://storage.googleapis.com/%s/%s", bucketName, filename)
}

func (s *GCSStore) Exists(ctx context.Context, folder, filename string) (bool, error) {
	bucketName, err := s.getBucketName(folder)
	if err != nil {
		return false, err
	}
	_, err = s.Client.Bucket(bucketName).Object(filename).Attrs(ctx)
	if err == storage.ErrObjectNotExist {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}

func (s *GCSStore) Get(ctx context.Context, folder, filename string) (io.ReadCloser, error) {
	bucketName, err := s.getBucketName(folder)
	if err != nil {
		return nil, err
	}
	return s.Client.Bucket(bucketName).Object(filename).NewReader(ctx)
}
