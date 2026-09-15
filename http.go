package cdn_ranges

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"sync"
	"time"
)

const (
	UserAgent = "cdn-ranges (+https://github.com/schniggie/cdn-ranges)"

	// Transient failures (network errors, 429, 5xx) are retried with
	// exponential backoff: 2s, 4s, 8s.
	httpMaxAttempts = 4
	httpBaseBackoff = 2 * time.Second
	httpMaxBackoff  = 60 * time.Second
)

var (
	httpClientOnce sync.Once
	httpClient     *http.Client
)

// HTTPStatusError is returned when a request completes with a non-2xx status
type HTTPStatusError struct {
	URL        string
	StatusCode int
}

func (e *HTTPStatusError) Error() string {
	return fmt.Sprintf("unexpected HTTP status %d from %s", e.StatusCode, e.URL)
}

// Retryable reports whether the status is worth retrying (rate limit or server error)
func (e *HTTPStatusError) Retryable() bool {
	return e.StatusCode == http.StatusTooManyRequests || e.StatusCode >= 500
}

// GetHttpClient returns a global reusable HTTP client
func GetHttpClient() *http.Client {
	httpClientOnce.Do(func() {
		httpClient = &http.Client{
			Transport: &http.Transport{
				MaxIdleConnsPerHost: 20,
			},
			Timeout: 30 * time.Second,
		}
	})

	return httpClient
}

// HttpGet performs an HTTP GET request with context, validating the response
// status and retrying transient failures with exponential backoff
func HttpGet(ctx context.Context, url string) ([]byte, error) {
	var lastErr error

	for attempt := 1; attempt <= httpMaxAttempts; attempt++ {
		body, retryAfter, err := httpGetOnce(ctx, url)
		if err == nil {
			return body, nil
		}
		lastErr = err

		if ctx.Err() != nil || !isRetryable(err) || attempt == httpMaxAttempts {
			break
		}

		delay := httpBaseBackoff << (attempt - 1)
		if retryAfter > delay {
			delay = retryAfter
		}
		if delay > httpMaxBackoff {
			delay = httpMaxBackoff
		}

		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-time.After(delay):
		}
	}

	return nil, fmt.Errorf("%s: %w", url, lastErr)
}

// httpGetOnce performs a single request. On a 429 response the server's
// Retry-After hint (if any) is returned so the caller can honour it.
func httpGetOnce(ctx context.Context, url string) ([]byte, time.Duration, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to create http request: %w", err)
	}
	req.Header.Set("User-Agent", UserAgent)
	req.Header.Set("Accept", "application/json, text/plain, */*")

	resp, err := GetHttpClient().Do(req)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to perform HTTP request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		// Drain so the connection can be reused
		io.Copy(io.Discard, io.LimitReader(resp.Body, 64<<10))

		return nil, parseRetryAfter(resp.Header.Get("Retry-After")), &HTTPStatusError{
			URL:        url,
			StatusCode: resp.StatusCode,
		}
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to read HTTP body: %w", err)
	}

	return body, 0, nil
}

func isRetryable(err error) bool {
	var statusErr *HTTPStatusError
	if errors.As(err, &statusErr) {
		return statusErr.Retryable()
	}

	// Anything else that made it here is a transport level failure
	// (timeout, reset, DNS hiccup) which is worth another attempt.
	return true
}

// parseRetryAfter handles the delay-seconds form of Retry-After. The HTTP-date
// form is ignored and the default backoff is used instead.
func parseRetryAfter(value string) time.Duration {
	if value == "" {
		return 0
	}

	seconds, err := strconv.Atoi(value)
	if err != nil || seconds <= 0 {
		return 0
	}

	return time.Duration(seconds) * time.Second
}
