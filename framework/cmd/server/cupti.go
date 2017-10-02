// +build !linux64,!amd64,!cgo

package server

func init() {
	cuptiHandle = noopCloser{}
}
