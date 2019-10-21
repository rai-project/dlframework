// +build ignore

package server

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"runtime"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/go-connections/nat"
	"github.com/k0kubun/pp"
	"github.com/mattn/go-shellwords"
	"github.com/pkg/errors"
	"github.com/rai-project/docker"
	"github.com/spf13/cobra"
)

var (
	defaultContainerImageNames    map[string]string
	defaultContainerImageCommands map[string]string
	containerName                 string
	containerAll                  bool
)

func startContainer(ctx context.Context, name string) error {

	shellCmd := defaultContainerImageCommands[name]

	args, err := shellwords.Parse(shellCmd)
	if err != nil {
		log.WithError(err).WithField("cmd", shellCmd).Error("failed to parse shell command")
		return err
	}
	fmt.Println("Running " + name)

	cmd := exec.Command(args[0], args[1:]...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err = cmd.Start()
	if err != nil {
		log.WithError(err).WithField("cmd", shellCmd).Error("failed to run command")
		return err
	}

	done := make(chan error)
	go func() { done <- cmd.Wait() }()

	select {
	case err := <-done:
		if err != nil {
			log.WithError(err).WithField("cmd", shellCmd).Error("failed to wait for command")
		}
	case <-ctx.Done():
		cmd.Process.Kill()
		log.WithError(ctx.Err()).WithField("cmd", shellCmd).Error("command timeout")
	}

	return nil
}

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
		err = client.ContainerRestart(ctx, container.ID, nil)
		if err != nil {
			return err
		}
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
	Short: "Administer MLModelScope containers",
}

var containerRestartCmd = &cobra.Command{
	Use:   "restart",
	Short: "Restarts a MLModelScope container",
	RunE: func(c *cobra.Command, args []string) error {
		containers := []string{}
		if containerAll {
			for k := range defaultContainerImageNames {
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
	Short: "Pulls a MLModelScope container",
	RunE: func(c *cobra.Command, args []string) error {
		containers := []string{}
		if containerAll {
			for k := range defaultContainerImageNames {
				containers = append(containers, k)
			}
		}
		if containerName != "" {
			pp.Println(containerName)

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
		defaultContainerImageCommands = map[string]string{
			"tracer": `docker run -d -e COLLECTOR_ZIPKIN_HTTP_PORT=9411 -p5775:5775/udp -p6831:6831/udp -p6832:6832/udp -p5778:5778 -p16686:16686 -p14268:14268 -p9411:9411 MLModelScope/jaeger:ppc64le-latest`,
		}
	case "amd64":
		defaultContainerImageNames = map[string]string{
			"tracer":   "jaegertracing/all-in-one:latest",
			"registry": "consul",
			"database": "mongo:3.0",
		}
		defaultContainerImageCommands = map[string]string{
			"tracer": `docker run -d -e COLLECTOR_ZIPKIN_HTTP_PORT=9411 -p5775:5775/udp -p6831:6831/udp -p6832:6832/udp -p5778:5778 -p16686:16686 -p14268:14268 -p9411:9411 jaegertracing/all-in-one:latest`,
		}
	}

	containerCmd.PersistentFlags().StringVar(&containerName, "name", "", "name of the container")
	containerCmd.PersistentFlags().BoolVar(&containerAll, "all", false, "perform operations on all mlmodelscope containers")
	containerCmd.AddCommand(containerRestartCmd, containerPullCmd)
}

func addContainerCmd(c *cobra.Command) {
	c.AddCommand(containerCmd)
}

func unmarshalContainerConfig(str string) container.Config {
	res := container.Config{}
	err := json.Unmarshal([]byte(str), &res)
	if err != nil {
		panic(err)
	}
	return res
}

func unmarshalContainerHostConfig(str string) container.HostConfig {
	res := container.HostConfig{}
	err := json.Unmarshal([]byte(str), &res)
	if err != nil {
		panic(err)
	}
	return res
}

func newPortNoError(proto, port string) nat.Port {
	p, _ := nat.NewPort(proto, port)
	return p
}
