package agent

import (
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
