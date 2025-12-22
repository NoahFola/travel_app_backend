package service

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"time"

	"github.com/NoahFola/travel_app_backend/internal/repository"
)

type MediaService struct {
	Repo *repository.MediaRepository
}

func (s *MediaService) UploadMedia(ctx context.Context, file *multipart.FileHeader, activityID string) (*repository.Media, error) {
	// 1. Ensure upload dir exists
	uploadDir := "./uploads"
	if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
		if err := os.Mkdir(uploadDir, 0755); err != nil {
			return nil, err
		}
	}

	// 2. Save file
	filename := fmt.Sprintf("%d_%s", time.Now().UnixNano(), file.Filename)
	dst := filepath.Join(uploadDir, filename)

	if err := s.saveFile(file, dst); err != nil {
		return nil, err
	}

	// 3. Create DB record
	// Determine type based on extension or header? Simple check for now.
	mediaType := "image" // Default
	// Could add logic here to check extension

	media := &repository.Media{
		URL:        "/uploads/" + filename,
		Type:       mediaType,
		ActivityID: activityID,
	}

	if err := s.Repo.Create(ctx, media); err != nil {
		return nil, err
	}

	return media, nil
}

func (s *MediaService) saveFile(file *multipart.FileHeader, dst string) error {
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, src)
	return err
}

func (s *MediaService) ListByActivityID(ctx context.Context, activityID string) ([]repository.Media, error) {
	return s.Repo.ListByActivityID(ctx, activityID)
}
