package client

import "strings"

func cleanNames() {
	frameworkName = strings.ToLower(frameworkName)
	frameworkVersion = strings.ToLower(frameworkVersion)
	modelName = strings.ToLower(modelName)
	modelVersion = strings.ToLower(modelVersion)
}
