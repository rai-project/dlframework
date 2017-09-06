package predict

import "github.com/rai-project/dlframework"

type Base struct {
	Framework dlframework.FrameworkManifest
	Model     dlframework.ModelManifest
}
