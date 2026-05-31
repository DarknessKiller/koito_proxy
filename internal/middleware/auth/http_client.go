package auth

import (
	"net/http"
	"sync"
	"time"
)

var (
	once   sync.Once
	client *http.Client
)

func httpClient() *http.Client {
	once.Do(func() {
		client = &http.Client{
			Timeout: 10 * time.Second,
			Transport: &http.Transport{
				MaxIdleConns:        100,
				MaxIdleConnsPerHost: 10,
				IdleConnTimeout:     90 * time.Second,
			},
		}
	})
	return client
}
