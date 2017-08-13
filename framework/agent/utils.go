package agent

import (
	"encoding/base64"
	"os"
	"path/filepath"

	"github.com/Unknwon/com"
	"github.com/pkg/errors"
	"github.com/rai-project/config"
)

func UploadDir() (string, error) {
	dir := filepath.Join(config.App.TempDir, "uploads")
	if !com.IsDir(dir) {
		if err := os.MkdirAll(dir, 0700); err != nil {
			return "", errors.Wrapf(err, "failed to create upload directory %v", dir)
		}
	}
	return dir, nil
}

// isBase64 tests a string to determine if it is a base64 or not.
func isBase64(toTest string) bool {
	_, err := base64.StdEncoding.DecodeString(toTest)
	return err == nil
}

func tryBase64Decode(input []byte) string {
	decoded, err := base64.StdEncoding.DecodeString(string(input))
	if err != nil {
		return string(input)
	}
	return string(decoded)
}
