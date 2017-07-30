package framework

import (
	"github.com/rai-project/dlframework"
	context "golang.org/x/net/context"
)

type Registry struct {
}

func (c *Registry) FrameworkManifests(context.Context, *dlframework.FrameworkRequest) (*dlframework.FrameworkManifestsResponse, error) {
	panic("FrameworkManifests")
	return nil, nil
}
func (c *Registry) FrameworkAgents(context.Context, *dlframework.FrameworkRequest) (*dlframework.Agents, error) {
	panic("FrameworkAgents")
	return nil, nil
}
func (c *Registry) ModelManifests(context.Context, *dlframework.ModelRequest) (*dlframework.ModelManifestsResponse, error) {
	panic("ModelManifests")
	return nil, nil
}
func (c *Registry) ModelAgents(context.Context, *dlframework.ModelRequest) (*dlframework.Agents, error) {
	panic("ModelAgents")
	return nil, nil
}
