package agent

import (
	"fmt"
	"os"
	"testing"

	"github.com/rai-project/config"
	dl "github.com/rai-project/dlframework"
	_ "github.com/rai-project/dlframework/frameworks/tensorflow"
	"github.com/stretchr/testify/assert"
)

func XXXTestModelRegistration(t *testing.T) {
	models, err := dl.Models()
	assert.NoError(t, err)
	for _, model := range models {
		fmt.Println(model.GetName() + ":" + model.GetVersion())
	}
}

func TestGRPCRegistration(t *testing.T) {
	RegisterRegistryServer()
}

func TestMain(m *testing.M) {
	config.Init(
		config.AppName("carml"),
		config.DebugMode(true),
		config.VerboseMode(true),
	)
	os.Exit(m.Run())
}
