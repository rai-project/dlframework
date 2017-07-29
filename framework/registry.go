package framework

import (
	"github.com/rai-project/dlframework"
	context "golang.org/x/net/context"
)

type Registry struct {
}

func (c *Registry) GetFrameworkManifests(ctx context.Context, ignore *dlframework.Null) (*dlframework.GetFrameworkManifestsResponse, error) {
	panic("GetFrameworkManifests")
	return nil, nil
}

func (c *Registry) GetFrameworkManifest(ctx context.Context, req *dlframework.GetFrameworkManifestRequest) (*dlframework.FrameworkManifest, error) {
	panic("GetFrameworkManifest")
	return nil, nil
}

func (c *Registry) GetFrameworkModels(ctx context.Context, req *dlframework.GetFrameworkManifestRequest) (*dlframework.GetModelManifestsResponse, error) {
	panic("GetFrameworkModels")
	return nil, nil
}

func (c *Registry) GetModelManifests(ctx context.Context, ignore *dlframework.Null) (*dlframework.GetModelManifestsResponse, error) {
	panic("GetModelManifests")
	return nil, nil
}

func (c *Registry) GetFrameworkModelManifest(ctx context.Context, req *dlframework.GetFrameworkModelManifestRequest) (*dlframework.ModelManifest, error) {
	panic("GetFrameworkModelManifest")
	return nil, nil
}

func (c *Registry) GetModelManifest(ctx context.Context, req *dlframework.GetModelManifestRequest) (*dlframework.ModelManifest, error) {
	panic("GetModelManifest")
	return nil, nil
}
