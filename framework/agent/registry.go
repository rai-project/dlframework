package agent

import (
	dl "github.com/rai-project/dlframework"
	context "golang.org/x/net/context"
)

type Registry struct {
	base
}

func (c *Registry) FrameworkManifests(context.Context, *dl.FrameworkRequest) (*dl.FrameworkManifestsResponse, error) {
	panic("FrameworkManifests")
	return nil, nil
}
func (c *Registry) FrameworkAgents(context.Context, *dl.FrameworkRequest) (*dl.Agents, error) {
	panic("FrameworkAgents")
	return nil, nil
}
func (c *Registry) ModelManifests(context.Context, *dl.ModelRequest) (*dl.ModelManifestsResponse, error) {
	panic("ModelManifests")
	return nil, nil
}
func (c *Registry) ModelAgents(context.Context, *dl.ModelRequest) (*dl.Agents, error) {
	panic("ModelAgents")
	return nil, nil
}

// func (c *Registry) GetFrameworkManifests(ctx context.Context, ignore *dl.Null) (*dl.GetFrameworkManifestsResponse, error) {
// 	frameworks, err := dl.Frameworks()
// 	if err != nil {
// 		return nil, err
// 	}
// 	pframeworks := make([]*dl.FrameworkManifest, len(frameworks))
// 	for ii := range frameworks {
// 		pframeworks[ii] = &frameworks[ii]
// 	}
// 	return &dl.GetFrameworkManifestsResponse{
// 		Manifests: pframeworks,
// 	}, nil
// }

// func (c *Registry) GetFrameworkManifest(ctx context.Context, req *dl.GetFrameworkManifestRequest) (*dl.FrameworkManifest, error) {
// 	f, err := dl.FindFramework(req.GetFrameworkName() + ":" + req.GetFrameworkVersion())
// 	if err != nil {
// 		return nil, err
// 	}
// 	return f, nil
// }

// func (c *Registry) GetFrameworkModels(ctx context.Context, req *dl.GetFrameworkManifestRequest) (*dl.GetModelManifestsResponse, error) {
// 	frameworks, err := dl.Frameworks()
// 	if err != nil {
// 		return nil, err
// 	}
// 	models := []*dl.ModelManifest{}
// 	for _, framework := range frameworks {
// 		if framework.GetName() != req.GetFrameworkName() ||
// 			framework.GetVersion() != req.GetFrameworkVersion() {
// 			continue
// 		}
// 		ms := framework.Models()
// 		for ii := range ms {
// 			models = append(models, &ms[ii])
// 		}
// 	}
// 	return &dl.GetModelManifestsResponse{
// 		Manifests: models,
// 	}, nil
// }

// func (c *Registry) GetModelManifests(ctx context.Context, ignore *dl.Null) (*dl.GetModelManifestsResponse, error) {
// 	models, err := dl.Models()
// 	if err != nil {
// 		return nil, err
// 	}
// 	pmodels := make([]*dl.ModelManifest, len(models))
// 	for ii := range models {
// 		pmodels[ii] = &models[ii]
// 	}
// 	return &dl.GetModelManifestsResponse{
// 		Manifests: pmodels,
// 	}, nil
// }

// func (c *Registry) GetFrameworkModelManifest(ctx context.Context, req *dl.GetFrameworkModelManifestRequest) (*dl.ModelManifest, error) {
// 	f, err := dl.FindFramework(req.GetFrameworkName() + ":" + req.GetFrameworkVersion())
// 	if err != nil {
// 		return nil, err
// 	}
// 	m, err := f.FindModel(req.GetModelName() + ":" + req.GetModelVersion())
// 	if err != nil {
// 		return nil, err
// 	}
// 	return m, nil
// }

// func (c *Registry) GetModelManifest(ctx context.Context, req *dl.GetModelManifestRequest) (*dl.ModelManifest, error) {
// 	m, err := dl.FindModel(req.GetModelName() + ":" + req.GetModelVersion())
// 	if err != nil {
// 		return nil, err
// 	}
// 	return m, nil
// }

func (c *Registry) PublishInRegistery() error {
	return c.Base.PublishInRegistery("registry")
}
