package main

import (
	"path/filepath"
	"strings"
)

func isValidImageType(filename string) bool {
	validTypes := []string{".jpg", ".jpeg", ".png", ".gif"}
	ext := strings.ToLower(filepath.Ext(filename))
	for _, validType := range validTypes {
		if ext == validType {
			return true
		}
	}
	return false
}