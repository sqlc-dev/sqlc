package quickdb

import (
	"net/http"
	"sync"
)

var region string
var once sync.Once

func GetClosestRegion() string {
	once.Do(func() {
		resp, err := http.Get("https://debug.fly.dev")
		if err == nil {
			region = resp.Header.Get("Fly-Region")
		}
	})
	return region
}
