package options

import "github.com/rai-project/config"
import "github.com/rai-project/nvidia-smi"

var (
	DefaultBatchSize uint32 = 1
	DefaultDevice    device
)

func init() {
	config.BeforeInit(func() {
		nvidiasmi.Wait()
		if nvidiasmi.HasGPU {
			DefaultDevice = device{deviceType: CUDA_DEVICE, id: 0}
		} else {
			DefaultDevice = device{deviceType: CPU_DEVICE, id: 0}
		}
	})
}
