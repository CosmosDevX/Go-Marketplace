// Package utils
package utils

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"runtime"
	"time"
)

type FileManager struct {
	RootDir string
}

func NewFileManager() FileManager {
	_, filename, _, _ := runtime.Caller(0)
	currentDir := filepath.Dir(filename)
	rootDir := filepath.Join(currentDir, "../..")

	return FileManager{
		RootDir: rootDir,
	}
}

func (m FileManager) SaveFile(file multipart.File, header *multipart.FileHeader, saveDirectory string) (string, error) {
	if file == nil || header == nil {
		return "", nil
	}

	fileExt := filepath.Ext(header.Filename)
	filename := fmt.Sprintf("%d%s", time.Now().UnixNano(), fileExt)
	savePath := filepath.Join(m.RootDir+fmt.Sprintf("/%s", saveDirectory), filename)

	dst, err := os.Create(savePath)
	if err != nil {
		return "", fmt.Errorf("save file %s in %s: %w", filename, saveDirectory, err)
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		return "", fmt.Errorf("write date to file %s in %s: %w", filename, saveDirectory, err)
	}

	return filename, nil
}

func (m FileManager) DeleteFile(path, filename string) error {
	if filename == "" {
		return nil
	}

	filePath := filepath.Join(m.RootDir+path, filename)
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return fmt.Errorf("file %s not exists in %s: %w", filename, path, err)
	}

	err := os.Remove(filePath)
	if err != nil {
		return fmt.Errorf("delete file %s from %s: %w", filename, path, err)
	}

	return nil
}
