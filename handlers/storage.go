package handlers

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"
)

const uploadDir = "uploads"

var (
	usersMu sync.RWMutex
	users   = map[string]string{}
)

func InitStorage() error {
	if err := os.MkdirAll(uploadDir, 0o755); err != nil {
		return err
	}
	return nil
}

func storeUser(username, hashedPassword string) {
	usersMu.Lock()
	defer usersMu.Unlock()
	users[strings.ToLower(username)] = hashedPassword
}

func getUserHash(username string) (string, bool) {
	usersMu.RLock()
	defer usersMu.RUnlock()
	hash, ok := users[strings.ToLower(username)]
	return hash, ok
}

func usernameTaken(username string) bool {
	_, exists := getUserHash(username)
	return exists
}

func safeUploadName(originalName string) (string, error) {
	base := filepath.Base(originalName)
	if base == "." || base == string(filepath.Separator) || base == "" {
		return "", errors.New("invalid file name")
	}
	if len(base) > 120 {
		base = base[:120]
	}
	timestamp := time.Now().UTC().Format("20060102150405")
	return timestamp + "_" + strings.ReplaceAll(base, " ", "_"), nil
}

func listUploadedFiles() ([]PdfResource, error) {
	entries, err := os.ReadDir(uploadDir)
	if err != nil {
		return nil, err
	}

	items := make([]PdfResource, 0, len(entries))
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		info, err := entry.Info()
		if err != nil {
			return nil, err
		}
		items = append(items, PdfResource{
			Title:     entry.Name(),
			PdfFile:   "/uploads/" + entry.Name(),
			CreatedAt: info.ModTime(),
		})
	}

	sort.Slice(items, func(i, j int) bool {
		return items[i].CreatedAt.After(items[j].CreatedAt)
	})

	return items, nil
}

func saveUsersSnapshot(path string) error {
	usersMu.RLock()
	defer usersMu.RUnlock()
	b, err := json.MarshalIndent(users, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, b, 0o600)
}
