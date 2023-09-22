package quickdb

import (
	"net/http"
	"sync"
)

var region string
var once sync.Once

func GetClosestRegion() string {
	once.Do(func() {
		resp, err := http.Get("https://find-closest-db-region.sqlc.dev")
		if err == nil {
			region = resp.Header.Get("Region")
		}
	})
	return region
}
