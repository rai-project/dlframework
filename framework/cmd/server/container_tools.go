package server

// import (
// 	"context"
// 	"os"
// 	"runtime"

// 	"github.com/docker/docker/api/types"
// 	"github.com/pkg/errors"
// 	"github.com/rai-project/docker"
// )

// var defaultContainerImageNames map[string]string

// func restartContainer(ctx context.Context, imageName string) error {
// 	client, err := docker.NewClient(
// 		docker.ClientContext(ctx),
// 		docker.Stdout(os.Stdout),
// 		docker.Stderr(os.Stderr),
// 	)
// 	if err != nil {
// 		return err
// 	}
// 	defer client.Close()

// 	containers, err := client.ContainerList(ctx, types.ContainerListOptions{
// 		All: true,
// 	})
// 	if err != nil {
// 		return err
// 	}
// 	for _, container := range containers {
// 		if container.Image != imageName {
// 			continue
// 		}
// 		return client.ContainerRestart(ctx, container.ID, nil)
// 	}
// 	return errors.Errorf("container %v not found", imageName)
// }

// func pullContainer(ctx context.Context, imageName string) error {
// 	client, err := docker.NewClient(
// 		docker.ClientContext(ctx),
// 		docker.Stdout(os.Stdout),
// 		docker.Stderr(os.Stderr),
// 	)
// 	if err != nil {
// 		return err
// 	}
// 	defer client.Close()

// 	return client.PullImage(imageName)
// }

// // func startContainer(ctx context.Context, imageName string) error {
// // }

// // func stopContainer(ctx context.Context, imageName string) error {
// // }

// func init() {
// 	switch runtime.GOARCH {
// 	case "ppc64le":
// 		defaultContainerImageNames = map[string]string{
// 			"tracer":   "carml/jaeger:ppc64le-latest",
// 			"registry": "carml/consul:ppc64le-latest",
// 			"database": "c3sr/mongodb:latest",
// 		}
// 	case "amd64":
// 		defaultContainerImageNames = map[string]string{
// 			"tracer":   "jaegertracing/all-in-one:latest",
// 			"registry": "consul",
// 			"database": "mongo:3.0",
// 		}
// 	}
// }
