package recon

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"slices"
	"sort"
	"strconv"
	"strings"
	"time"

	"cerber/internal/i18n"
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
	scanHeaders    map[string]string
	specHeaders    map[string]string
}

func NewAPIScanner(requestTimeoutSec int, scanHeaders map[string]string, specHeaders map[string]string) *APIScanner {
	reqTimeout := time.Duration(requestTimeoutSec) * time.Second
	if reqTimeout <= 0 {
		reqTimeout = 10 * time.Second
	}
	hScan := make(map[string]string, len(scanHeaders))
	for k, v := range scanHeaders {
		hScan[k] = v
	}
	hSpec := make(map[string]string, len(specHeaders))
	for k, v := range specHeaders {
		hSpec[k] = v
	}
	return &APIScanner{
		client:         &http.Client{},
		requestTimeout: reqTimeout,
		scanHeaders:    hScan,
		specHeaders:    hSpec,
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
		return nil, fmt.Errorf(i18n.T("err_api_key_requires_header"))
	}
	if apiKey == "" && apiKeyHeader != "" {
		return nil, fmt.Errorf(i18n.T("err_api_key_header_requires_key"))
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
		return OpenAPISpec{}, fmt.Errorf(i18n.T("err_invalid_openapi_json"), err)
	}
	if len(spec.Paths) == 0 {
		return OpenAPISpec{}, fmt.Errorf(i18n.T("err_openapi_paths_missing"))
	}
	return spec, nil
}

func (s *APIScanner) Scan(baseURL string, spec OpenAPISpec) []APIScanResult {
	results := make([]APIScanResult, 0)
	paths := make([]string, 0, len(spec.Paths))
	for path := range spec.Paths {
		paths = append(paths, path)
	}
	sort.Strings(paths)

	for _, path := range paths {
		methods := spec.Paths[path]
		methodNames := make([]string, 0, len(methods))
		for method := range methods {
			methodNames = append(methodNames, method)
		}
		sort.Strings(methodNames)

		for _, method := range methodNames {
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
			return nil, fmt.Errorf(i18n.T("err_invalid_status_code"), part)
		}
		if code < 100 || code > 999 {
			return nil, fmt.Errorf(i18n.T("err_status_code_out_of_range"), code)
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

	httpMethod := strings.ToUpper(method)
	var body io.Reader
	if isWriteMethod(httpMethod) {
		body = bytes.NewBufferString("{}")
	}

	req, err := http.NewRequestWithContext(ctx, httpMethod, target, body)
	if err != nil {
		return 0, err
	}
	if isWriteMethod(httpMethod) {
		req.Header.Set("Content-Type", "application/json")
	}
	for k, v := range s.scanHeaders {
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
	for k, v := range s.specHeaders {
		req.Header.Set(k, v)
	}
	resp, err := s.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf(i18n.T("err_spec_load_status"), resp.StatusCode)
	}
	return io.ReadAll(resp.Body)
}

func isWriteMethod(method string) bool {
	switch method {
	case http.MethodPost, http.MethodPut, http.MethodPatch:
		return true
	default:
		return false
	}
}
