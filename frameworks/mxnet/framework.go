package mxnet

import "github.com/rai-project/dlframework"

var thisFramework = dlframework.FrameworkManifest{
	Name: "MXNet",
	DefaultContainer: map[string]*dlframework.ContainerHardware{
		"amd64": &dlframework.ContainerHardware{
			Cpu: "raiproject/carml-mxnet:amd64-cpu",
			Gpu: "raiproject/carml-mxnet:amd64-gpu",
		},
		"ppc64le": &dlframework.ContainerHardware{
			Cpu: "raiproject/carml-mxnet:ppc64le-gpu",
			Gpu: "raiproject/carml-mxnet:ppc64le-gpu",
		},
	},
}

func init() {
	dlframework.AddFramework(thisFramework)
}
