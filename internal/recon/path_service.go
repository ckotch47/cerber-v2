package recon

import (
	"context"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type PathResult struct {
	Path       string
	URL        string
	StatusCode int
	Err        error
}

type PathScanner struct {
	client         *http.Client
	requestTimeout time.Duration
	delay          time.Duration
	exclude        map[string]bool
	allowFallback  bool
}

func NewPathScanner(requestTimeoutSec int, delaySec int, exclude []string, allowFallback bool) *PathScanner {
	reqTimeout := time.Duration(requestTimeoutSec) * time.Second
	if reqTimeout <= 0 {
		reqTimeout = 10 * time.Second
	}
	delay := time.Duration(delaySec) * time.Second
	if delay < 0 {
		delay = 0
	}

	return &PathScanner{
		client:         &http.Client{},
		requestTimeout: reqTimeout,
		delay:          delay,
		exclude:        ToStatusSet(exclude),
		allowFallback:  allowFallback,
	}
}

func NormalizeBaseURL(input string) string {
	trimmed := strings.TrimSpace(strings.TrimSuffix(input, "/"))
	if trimmed == "" {
		return ""
	}
	if HasHTTPPrefix(trimmed) {
		return trimmed
	}
	return "https://" + trimmed
}

func HasHTTPPrefix(input string) bool {
	lowered := strings.ToLower(input)
	return strings.HasPrefix(lowered, "http://") || strings.HasPrefix(lowered, "https://")
}

func JoinURL(baseURL, path string) string {
	return strings.TrimRight(baseURL, "/") + "/" + strings.TrimLeft(path, "/")
}

func ToStatusSet(exclude []string) map[string]bool {
	set := make(map[string]bool)
	for _, code := range exclude {
		if code != "" {
			set[code] = true
		}
	}
	return set
}

func (s *PathScanner) IsExcluded(statusCode int) bool {
	return s.exclude[strconv.Itoa(statusCode)]
}

func (s *PathScanner) Scan(baseURL string, paths []string) []PathResult {
	results := make([]PathResult, 0, len(paths))
	for _, path := range paths {
		target := JoinURL(baseURL, path)
		resp, actualURL, err := s.requestWithFallback(target)
		if err != nil {
			results = append(results, PathResult{
				Path: path,
				URL:  actualURL,
				Err:  err,
			})
		} else {
			_ = resp.Body.Close()
			results = append(results, PathResult{
				Path:       path,
				URL:        actualURL,
				StatusCode: resp.StatusCode,
			})
		}

		if s.delay > 0 {
			time.Sleep(s.delay)
		}
	}
	return results
}

func (s *PathScanner) requestWithFallback(target string) (*http.Response, string, error) {
	resp, err := s.requestOnce(target)
	if err == nil {
		return resp, target, nil
	}
	if !s.allowFallback || !strings.HasPrefix(target, "https://") {
		return nil, target, err
	}

	fallbackTarget := "http://" + strings.TrimPrefix(target, "https://")
	fallbackResp, fallbackErr := s.requestOnce(fallbackTarget)
	if fallbackErr != nil {
		return nil, fallbackTarget, fallbackErr
	}
	return fallbackResp, fallbackTarget, nil
}

func (s *PathScanner) requestOnce(target string) (*http.Response, error) {
	ctx, cancel := context.WithTimeout(context.Background(), s.requestTimeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, target, nil)
	if err != nil {
		return nil, err
	}
	return s.client.Do(req)
}
