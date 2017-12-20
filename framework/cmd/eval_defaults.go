package cmd

import "strings"

var (
	DefaultEvaulationModels = []string{
		"SqueezeNet_1.0",
		"SqueezeNet_1.1",
		"BVLC-AlexNet_1.0",
		"BVLC-Reference-CaffeNet_1.0",
		"BVLC-GoogLeNet_1.0",
		"ResNet101_1.0",
		"ResNet101_2.0",
		"WRN50_2.0",
		"BVLC-Reference-RCNN-ILSVRC13_1.0",
		"Inception_3.0",
		"Inception_4.0",
		"ResNeXt50-32x4d_1.0",
		"VGG16_1.0",
		"VGG19_1.0",
	}

	DefaultEvaluationFrameworks = []string{
		"mxnet",
		"cntk",
		"caffe2",
		"tensorflow",
		"tensorrt",
		"caffe",
	}

	FrameworkNames = map[string]string{
		"tensorflow": "Tensorflow",
		"tensorrt":   "TensorRT",
		"mxnet":      "MXNet",
		"caffe":      "Caffe",
		"caffe2":     "Caffe2",
		"cntk":       "CNTK",
		"pytorch":    "PyTorch",
	}
)

func ParseModelName(model string) (modelName, version string) {
	splt := strings.Split(model, "_")
	modelName, modelVersion := splt[0], splt[1]
	return modelName, modelVersion
}
