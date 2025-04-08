package helpers

import (
	"context"
	"io"
	"log"
	"net/http"
	"strings"

	"golang.org/x/time/rate"
)

const (
	apiKey = "323b909a-3d84-47be-8c80-8057fce536cd"
)

type Transport struct {
	APIKey string
	Base   http.RoundTripper
}

type RateLimiter struct {
	limiter *rate.Limiter
}

func NewRateLimiter() *RateLimiter {
	return &RateLimiter{
		limiter: rate.NewLimiter(rate.Limit(100), 6000),
	}
}

func (rl *RateLimiter) Allow() bool {
	return rl.limiter.Allow()
}

func (t *Transport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Add("X-API-Key", apiKey)
	req.Header.Add("apiKey", apiKey)
	req.Header.Add("Content-Type", "application/json")
	return t.Base.RoundTrip(req)
}

func newClient() *http.Client {
	return &http.Client{
		Transport: &Transport{
			APIKey: apiKey,
			Base:   http.DefaultTransport,
		},
	}
}

func NewRequest(limiter *RateLimiter, method string, url string, body io.Reader) *http.Response {
	if strings.TrimSpace(url) == "" || strings.TrimSpace(method) == "" {
		log.Println("no url or method provided")
		return nil
	}

	if !limiter.Allow() {
		log.Println("rate limit exceeded")
		return nil
	}

	req, err := http.NewRequestWithContext(context.Background(), method, url, body)
	if err != nil {
		log.Printf("error creating request: %v\n", err)
		return nil
	}
	client := newClient()

	resp, err := client.Do(req)
	if err != nil {
		log.Printf("error making request: %v\n", err)
		return nil
	}
	return resp
}
