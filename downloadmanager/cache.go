package downloadmanager

import (
	"time"

	gocache "github.com/patrickmn/go-cache"
)

var (
	cache *gocache.Cache
)

func init() {

	// Create a cache with a default expiration time of 5 minutes, and which
	// purges expired items every 10 minutes
	cache = gocache.New(5*time.Minute, 10*time.Minute)

}
