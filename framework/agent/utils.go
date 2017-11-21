package agent

import (
	"context"
	"encoding/base64"
	"os"
	"path/filepath"

	"github.com/Unknwon/com"
	"github.com/pkg/errors"
	"github.com/rai-project/config"
	"github.com/rai-project/dlframework/framework/options"
	cupti "github.com/rai-project/go-cupti"
	"github.com/rai-project/tracer"
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

func Partition(in []interface{}, partitionSize int) (out [][]interface{}) {
	cnt := (len(in)-1)/partitionSize + 1
	for i := 0; i < cnt; i++ {
		start := i * partitionSize
		end := (i + 1) * partitionSize
		if end > len(in) {
			end = len(in)
		}
		part := in[start:end]
		out = append(out, part)
	}

	return out
}

func cuptiTrace(ctx context.Context, opts *options.Options) {
	if opts.UsesGPU() && opts.TraceLevel() >= tracer.HARDWARE_TRACE {
		cu, err := cupti.New(cupti.Context(ctx))
		if err == nil {
			defer func() {
				cu.Wait()
				cu.Close()
			}()
		}
	}
}
