package storage

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func CreateNotePath(root string, inputPath string) (string, error) {
	cleanPath := strings.TrimSpace(inputPath) // Limpia todos los espacios
	if !strings.HasSuffix(cleanPath, ".md") {
		cleanPath += ".md"
	} // y de paso asegura la extencion .md

	fullPath := filepath.Join(root, cleanPath)

	dir := filepath.Dir(fullPath)

	err := os.MkdirAll(dir, 0o755)
	if err != nil {
		return "", fmt.Errorf("Error a; crear carpetas: %w", err)
	}

	return fullPath, nil
}
