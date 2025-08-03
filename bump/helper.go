package bump

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/andreimerlescu/checkfs"
	"github.com/andreimerlescu/checkfs/file"
)

// compareInt is used to determine if a > b then 1 or a < b then -1 or 0
func compareInt(a, b int) int {
	if a > b {
		return 1
	}
	if a < b {
		return -1
	}
	return 0
}

func currentIgoVersion() (string, error) {
	lookingFor := filepath.Join(os.Getenv("HOME"), "go", "version")
	err := checkfs.File(lookingFor, file.Options{Exists: true})
	if err != nil {
		return "", err
	}
	b, err := os.ReadFile(lookingFor)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(b)), nil
}
