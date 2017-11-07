package server

import (
	"context"
	"os"
	"runtime"

	"github.com/docker/docker/api/types"
	"github.com/pkg/errors"
	"github.com/rai-project/docker"
	"github.com/spf13/cobra"
)

var (
	defaultContainerImageNames map[string]string
	containerName              string
	containerAll               bool
)

func restartContainer(ctx context.Context, imageName string) error {
	client, err := docker.NewClient(
		docker.ClientContext(ctx),
		docker.Stdout(os.Stdout),
		docker.Stderr(os.Stderr),
	)
	if err != nil {
		return err
	}
	defer client.Close()

	containers, err := client.ContainerList(ctx, types.ContainerListOptions{
		All: true,
	})
	if err != nil {
		return err
	}
	for _, container := range containers {
		if container.Image != imageName {
			continue
		}
		return client.ContainerRestart(ctx, container.ID, nil)
	}
	return errors.Errorf("container %v not found", imageName)
}

func pullContainer(ctx context.Context, imageName string) error {
	client, err := docker.NewClient(
		docker.ClientContext(ctx),
		docker.Stdout(os.Stdout),
		docker.Stderr(os.Stderr),
	)
	if err != nil {
		return err
	}
	defer client.Close()

	return client.PullImage(imageName)
}

var containerCmd = &cobra.Command{
	Use: "container",
	Aliases: []string{
		"docker",
	},
	Short: "Administer CarML containers",
}

var containerRestartCmd = &cobra.Command{
	Use:   "restart",
	Short: "Restarts a CarML container",
	RunE: func(c *cobra.Command, args []string) error {
		containers := []string{}
		if containerAll {
			for k, _ := range defaultContainerImageNames {
				containers = append(containers, k)
			}
		}
		if containerName != "" {
			containers = append(containers, containerName)
		}
		if len(containers) == 0 {
			return errors.New("no containers specified")
		}

		ctx := context.Background()
		for _, c := range containers {
			name, ok := defaultContainerImageNames[c]
			if !ok {
				log.WithField("name", c).
					WithField("containers", defaultContainerImageNames).
					Error("unable to find container name")
				continue
			}
			err := restartContainer(ctx, name)
			if err != nil {
				log.WithError(err).Error("failed to restart container")
			}
		}
		return nil
	},
}

var containerPullCmd = &cobra.Command{
	Use:   "pull",
	Short: "Pulls a CarML container",
	RunE: func(c *cobra.Command, args []string) error {
		containers := []string{}
		if containerAll {
			for k, _ := range defaultContainerImageNames {
				containers = append(containers, k)
			}
		}
		if containerName != "" {
			containers = append(containers, containerName)
		}
		if len(containers) == 0 {
			return errors.New("no containers specified")
		}

		ctx := context.Background()
		for _, c := range containers {
			name, ok := defaultContainerImageNames[c]
			if !ok {
				log.WithField("name", c).
					WithField("containers", defaultContainerImageNames).
					Error("unable to find container name")
				continue
			}
			err := pullContainer(ctx, name)
			if err != nil {
				log.WithError(err).Error("failed to pull container")
			}
		}
		return nil
	},
}

func init() {
	switch runtime.GOARCH {
	case "ppc64le":
		defaultContainerImageNames = map[string]string{
			"tracer":   "carml/jaeger:ppc64le-latest",
			"registry": "carml/consul:ppc64le-latest",
			"database": "c3sr/mongodb:latest",
		}
	case "amd64":
		defaultContainerImageNames = map[string]string{
			"tracer":   "jaegertracing/all-in-one:latest",
			"registry": "consul",
			"database": "mongo:3.0",
		}
	}

	containerCmd.PersistentFlags().StringVar(&containerName, "name", "", "name of the container")
	containerCmd.PersistentFlags().BoolVar(&containerAll, "all", false, "perform operations on all carml containers")
	containerCmd.AddCommand(containerRestartCmd, containerPullCmd)
}
