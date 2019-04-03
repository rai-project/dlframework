// +build ignore

package main

var (
	DefaultEvaulationModels = []string{
		"SqueezeNet_1.0",
		"SqueezeNet_1.1",
		"BVLC_AlexNet_1.0",
		"BVLC_Reference_CaffeNet_1.0",
		"BVLC_GoogLeNet_1.0",
		"ResNet101_1.0",
		"ResNet101_2.0",
		"WRN50_2.0",
		"BVLC_Reference_RCNN_ILSVRC13_1.0",
		"Inception_3.0",
		"Inception_4.0",
		"ResNeXt50_32x4d_1.0",
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
		"tensorflow": "TensorFlow",
		"tensorrt":   "TensorRT",
		"mxnet":      "MXNet",
		"caffe":      "Caffe",
		"caffe2":     "Caffe2",
		"cntk":       "CNTK",
		"pytorch":    "PyTorch",
	}
)
