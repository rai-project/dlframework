package monitors

import (
	"github.com/google/gops/agent"
)

func init() {
	if err := agent.Listen(agent.Options{}); err != nil {
		log.Fatal(err)
	}
}
