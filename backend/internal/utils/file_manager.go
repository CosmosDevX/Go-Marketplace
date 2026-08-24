// Package utils
package utils

import (
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"myapp/internal/domain"
)

var allowedExts = map[string]bool{
	".jpg":  true,
	".jpeg": true,
	".png":  true,
	".webp": true,
}

type FileManager struct {
	UploadsDir string
}

func NewFileManager(uploadsDir string) FileManager {
	return FileManager{UploadsDir: uploadsDir}
}

func (m FileManager) SaveFile(file multipart.File, header *multipart.FileHeader) (string, error) {
	if file == nil || header == nil {
		return "", fmt.Errorf("file is required: %w", domain.ErrValidation)
	}

	ext := strings.ToLower(filepath.Ext(header.Filename))
	if !allowedExts[ext] {
		return "", fmt.Errorf("unsupported file type %s: %w", ext, domain.ErrValidation)
	}

	buf := make([]byte, 512)
	n, err := file.Read(buf)
	if err != nil && err != io.EOF {
		return "", fmt.Errorf("read file header: %w", err)
	}
	contentType := http.DetectContentType(buf[:n])
	if !strings.HasPrefix(contentType, "image/") {
		return "", fmt.Errorf("invalid content type %s: %w", contentType, domain.ErrValidation)
	}

	if seeker, ok := file.(io.Seeker); ok {
		if _, err := seeker.Seek(0, io.SeekStart); err != nil {
			return "", fmt.Errorf("seek file: %w", err)
		}
	}

	filename := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)
	savePath := filepath.Join(m.UploadsDir, filename)

	dst, err := os.Create(savePath)
	if err != nil {
		return "", fmt.Errorf("create file %s: %w", filename, err)
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		_ = os.Remove(savePath)
		return "", fmt.Errorf("write file %s: %w", filename, err)
	}

	return filename, nil
}

func (m FileManager) DeleteFile(filename string) error {
	if filename == "" {
		return nil
	}

	filename = filepath.Base(filename)
	filePath := filepath.Join(m.UploadsDir, filename)

	if err := os.Remove(filePath); err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return fmt.Errorf("delete file %s: %w", filename, err)
	}
	return nil
}
