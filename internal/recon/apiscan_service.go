package recon

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"slices"
	"strconv"
	"strings"
	"time"
)

var allowedHTTPMethods = []string{"get", "post", "put", "patch", "delete"}

type OpenAPISpec struct {
	Paths map[string]map[string]json.RawMessage `json:"paths"`
}

type APIScanResult struct {
	Method     string
	Path       string
	StatusCode int
	Err        error
}

type APIScanner struct {
	client         *http.Client
	requestTimeout time.Duration
	headers        map[string]string
}

func NewAPIScanner(requestTimeoutSec int, headers map[string]string) *APIScanner {
	reqTimeout := time.Duration(requestTimeoutSec) * time.Second
	if reqTimeout <= 0 {
		reqTimeout = 10 * time.Second
	}
	h := make(map[string]string, len(headers))
	for k, v := range headers {
		h[k] = v
	}
	return &APIScanner{
		client:         &http.Client{},
		requestTimeout: reqTimeout,
		headers:        h,
	}
}

func BuildAuthHeaders(jwt string, apiKeyHeader string, apiKey string) (map[string]string, error) {
	headers := make(map[string]string)
	if strings.TrimSpace(jwt) != "" {
		headers["Authorization"] = "Bearer " + strings.TrimSpace(jwt)
	}

	apiKeyHeader = strings.TrimSpace(apiKeyHeader)
	apiKey = strings.TrimSpace(apiKey)
	if apiKey != "" && apiKeyHeader == "" {
		return nil, fmt.Errorf("для --api-key требуется --api-key-header")
	}
	if apiKey == "" && apiKeyHeader != "" {
		return nil, fmt.Errorf("для --api-key-header требуется --api-key")
	}
	if apiKey != "" && apiKeyHeader != "" {
		headers[apiKeyHeader] = apiKey
	}
	return headers, nil
}

func (s *APIScanner) LoadSpec(specSource string) (OpenAPISpec, error) {
	var data []byte
	var err error

	if HasHTTPPrefix(specSource) {
		data, err = s.readFromURL(specSource)
	} else {
		data, err = os.ReadFile(specSource)
	}
	if err != nil {
		return OpenAPISpec{}, err
	}

	var spec OpenAPISpec
	if err := json.Unmarshal(data, &spec); err != nil {
		return OpenAPISpec{}, fmt.Errorf("невалидный openapi json: %w", err)
	}
	if len(spec.Paths) == 0 {
		return OpenAPISpec{}, fmt.Errorf("в спецификации отсутствует paths")
	}
	return spec, nil
}

func (s *APIScanner) Scan(baseURL string, spec OpenAPISpec) []APIScanResult {
	results := make([]APIScanResult, 0)
	for path, methods := range spec.Paths {
		for method := range methods {
			if !isAllowedMethod(method) {
				continue
			}
			status, err := s.request(method, JoinURL(baseURL, path))
			results = append(results, APIScanResult{
				Method:     strings.ToUpper(method),
				Path:       path,
				StatusCode: status,
				Err:        err,
			})
		}
	}
	return results
}

func ParseStatusCodes(value string) (map[int]bool, error) {
	set := make(map[int]bool)
	if strings.TrimSpace(value) == "" {
		return set, nil
	}
	parts := strings.Split(value, ",")
	for _, part := range parts {
		part = strings.TrimSpace(part)
		code, err := strconv.Atoi(part)
		if err != nil {
			return nil, fmt.Errorf("невалидный статус-код %q", part)
		}
		if code < 100 || code > 999 {
			return nil, fmt.Errorf("статус-код вне диапазона 100..999: %d", code)
		}
		set[code] = true
	}
	return set, nil
}

func ShouldIncludeStatus(code int, show map[int]bool, exclude map[int]bool) bool {
	if len(show) > 0 {
		return show[code]
	}
	if len(exclude) > 0 {
		return !exclude[code]
	}
	return true
}

func isAllowedMethod(method string) bool {
	return slices.Contains(allowedHTTPMethods, strings.ToLower(method))
}

func (s *APIScanner) request(method string, target string) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), s.requestTimeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, strings.ToUpper(method), target, nil)
	if err != nil {
		return 0, err
	}
	for k, v := range s.headers {
		req.Header.Set(k, v)
	}
	resp, err := s.client.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()
	_, _ = io.Copy(io.Discard, resp.Body)
	return resp.StatusCode, nil
}

func (s *APIScanner) readFromURL(source string) ([]byte, error) {
	ctx, cancel := context.WithTimeout(context.Background(), s.requestTimeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, source, nil)
	if err != nil {
		return nil, err
	}
	for k, v := range s.headers {
		req.Header.Set(k, v)
	}
	resp, err := s.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("не удалось загрузить спецификацию: status %d", resp.StatusCode)
	}
	return io.ReadAll(resp.Body)
}
