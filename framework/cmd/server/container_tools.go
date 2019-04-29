package server

import (
	"context"
	"encoding/json"
	"os"
	"runtime"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/go-connections/nat"
	"github.com/k0kubun/pp"
	"github.com/pkg/errors"
	"github.com/rai-project/docker"
	"github.com/spf13/cobra"
)

var (
	defaultContainerImageNames        map[string]string
	defaultContainerNetworkingConfigs map[string]network.NetworkingConfig
	defaultContainerHostConfigs       map[string]container.HostConfig
	defaultContainerConfigs           map[string]container.Config
	containerName                     string
	containerAll                      bool
)

func startContainer(ctx context.Context, imageName string) error {
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
		c, err := docker.NewContainer(
			client,
			docker.Image(imageName),
			docker.ContainerConfig(defaultContainerConfigs[imageName]),
			docker.HostConfig(defaultContainerHostConfigs[imageName]),
			docker.NetworkDisabled(false),
			docker.NetworkConfig(defaultContainerNetworkingConfigs[imageName]),
		)
		if err != nil {
			return err
		}

		c.Start()
	}
	return errors.Errorf("container %v not found", imageName)
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
	case "amd64":
		defaultContainerImageNames = map[string]string{
			"tracer":   "jaegertracing/all-in-one:latest",
			"registry": "consul",
			"database": "mongo:3.0",
		}
	}

	defaultContainerNetworkingConfigs = map[string]network.NetworkingConfig{
		"tracer": network.NetworkingConfig{
			EndpointsConfig: map[string]*network.EndpointSettings{
				"bridge": &network.EndpointSettings{},
			},
		},
		"registry": network.NetworkingConfig{},
		"database": network.NetworkingConfig{},
	}

	defaultContainerConfigs = map[string]container.Config{
		"tracer": unmarshalContainerConfig(`
    {
      "Hostname": "b5b20d6859db",
      "Domainname": "",
      "User": "",
      "AttachStdin": false,
      "AttachStdout": false,
      "AttachStderr": false,
      "ExposedPorts": {
          "14250/tcp": {},
          "14268/tcp": {},
          "16686/tcp": {},
          "5775/udp": {},
          "5778/tcp": {},
          "6831/udp": {},
          "6832/udp": {},
          "9411/tcp": {}
      },
      "Tty": false,
      "OpenStdin": false,
      "StdinOnce": false,
      "Env": [
          "COLLECTOR_ZIPKIN_HTTP_PORT=9411",
          "PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin"
      ],
      "Cmd": [
          "--sampling.strategies-file=/etc/jaeger/sampling_strategies.json"
      ],
      "ArgsEscaped": true,
      "Image": "jaegertracing/all-in-one:latest",
      "Volumes": null,
      "WorkingDir": "",
      "Entrypoint": [
          "/go/bin/all-in-one-linux"
      ],
      "OnBuild": null,
      "Labels": {}
  }`),
		"registry": container.Config{},
		"database": container.Config{},
	}
	defaultContainerHostConfigs = map[string]container.HostConfig{
		"tracer": unmarshalContainerHostConfig(`
    {
      "Binds": null,
      "ContainerIDFile": "",
      "LogConfig": {
        "Type": "json-file",
        "Config": {}
      },
      "NetworkMode": "default",
      "PortBindings": {
        "14268/tcp": [{
          "HostIp": "",
          "HostPort": "14268"
        }],
        "16686/tcp": [{
          "HostIp": "",
          "HostPort": "16686"
        }],
        "5775/udp": [{
          "HostIp": "",
          "HostPort": "5775"
        }],
        "5778/tcp": [{
          "HostIp": "",
          "HostPort": "5778"
        }],
        "6831/udp": [{
          "HostIp": "",
          "HostPort": "6831"
        }],
        "6832/udp": [{
          "HostIp": "",
          "HostPort": "6832"
        }],
        "9411/tcp": [{
          "HostIp": "",
          "HostPort": "9411"
        }]
      },
      "RestartPolicy": {
        "Name": "no",
        "MaximumRetryCount": 0
      },
      "AutoRemove": false,
      "VolumeDriver": "",
      "VolumesFrom": null,
      "CapAdd": null,
      "CapDrop": null,
      "Dns": [],
      "DnsOptions": [],
      "DnsSearch": [],
      "ExtraHosts": null,
      "GroupAdd": null,
      "IpcMode": "shareable",
      "Cgroup": "",
      "Links": null,
      "OomScoreAdj": 0,
      "PidMode": "",
      "Privileged": false,
      "PublishAllPorts": false,
      "ReadonlyRootfs": false,
      "SecurityOpt": null,
      "UTSMode": "",
      "UsernsMode": "",
      "ShmSize": 67108864,
      "Runtime": "runc",
      "ConsoleSize": [
        0,
        0
      ],
      "Isolation": "",
      "CpuShares": 0,
      "Memory": 0,
      "NanoCpus": 0,
      "CgroupParent": "",
      "BlkioWeight": 0,
      "BlkioWeightDevice": [],
      "BlkioDeviceReadBps": null,
      "BlkioDeviceWriteBps": null,
      "BlkioDeviceReadIOps": null,
      "BlkioDeviceWriteIOps": null,
      "CpuPeriod": 0,
      "CpuQuota": 0,
      "CpuRealtimePeriod": 0,
      "CpuRealtimeRuntime": 0,
      "CpusetCpus": "",
      "CpusetMems": "",
      "Devices": [],
      "DeviceCgroupRules": null,
      "DiskQuota": 0,
      "KernelMemory": 0,
      "MemoryReservation": 0,
      "MemorySwap": 0,
      "MemorySwappiness": null,
      "OomKillDisable": false,
      "PidsLimit": 0,
      "Ulimits": null,
      "CpuCount": 0,
      "CpuPercent": 0,
      "IOMaximumIOps": 0,
      "IOMaximumBandwidth": 0,
      "MaskedPaths": [
        "/proc/asound",
        "/proc/acpi",
        "/proc/kcore",
        "/proc/keys",
        "/proc/latency_stats",
        "/proc/timer_list",
        "/proc/timer_stats",
        "/proc/sched_debug",
        "/proc/scsi",
        "/sys/firmware"
      ],
      "ReadonlyPaths": [
        "/proc/bus",
        "/proc/fs",
        "/proc/irq",
        "/proc/sys",
        "/proc/sysrq-trigger"
      ]
    }
    `),
		"registry": container.HostConfig{},
		"database": container.HostConfig{},
	}

	containerCmd.PersistentFlags().StringVar(&containerName, "name", "", "name of the container")
	containerCmd.PersistentFlags().BoolVar(&containerAll, "all", false, "perform operations on all mlmodelscope containers")
	containerCmd.AddCommand(containerRestartCmd, containerPullCmd)
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
