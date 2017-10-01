package monitors

import "github.com/gbbr/memstats"

func init() {
	memstats.Serve()
}
