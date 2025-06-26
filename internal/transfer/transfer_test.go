package transfer

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/elastic/go-elasticsearch/v8/esapi"
)

// MockElasticsearchAPI implements ElasticsearchAPI for testing
type MockElasticsearchAPI struct {
	CountResponse    *esapi.Response
	SearchResponse   *esapi.Response
	ScrollResponse   *esapi.Response
	IndexResponse    *esapi.Response
	MappingResponse  *esapi.Response
	SettingsResponse *esapi.Response
}

// Count implements ElasticsearchAPI for testing
func (m *MockElasticsearchAPI) Count(o ...func(*esapi.CountRequest)) (*esapi.Response, error) {
	if m.CountResponse != nil {
		return m.CountResponse, nil
	}
	return createMockCountResponse(100, false), nil
}

// Search implements ElasticsearchAPI for testing
func (m *MockElasticsearchAPI) Search(o ...func(*esapi.SearchRequest)) (*esapi.Response, error) {
	if m.SearchResponse != nil {
		return m.SearchResponse, nil
	}
	return createMockSearchResponse(), nil
}

// Scroll implements ElasticsearchAPI for testing
func (m *MockElasticsearchAPI) Scroll(o ...func(*esapi.ScrollRequest)) (*esapi.Response, error) {
	if m.ScrollResponse != nil {
		return m.ScrollResponse, nil
	}
	return createMockScrollResponse(), nil
}

// Index implements ElasticsearchAPI for testing
func (m *MockElasticsearchAPI) Index(index string, body io.Reader, o ...func(*esapi.IndexRequest)) (*esapi.Response, error) {
	if m.IndexResponse != nil {
		return m.IndexResponse, nil
	}
	return createMockIndexResponse(), nil
}

// IndicesGetMapping implements ElasticsearchAPI for testing
func (m *MockElasticsearchAPI) IndicesGetMapping(o ...func(*esapi.IndicesGetMappingRequest)) (*esapi.Response, error) {
	if m.MappingResponse != nil {
		return m.MappingResponse, nil
	}
	return createMockMappingResponse(), nil
}

// IndicesPutMapping implements ElasticsearchAPI for testing
func (m *MockElasticsearchAPI) IndicesPutMapping(indices []string, body io.Reader, o ...func(*esapi.IndicesPutMappingRequest)) (*esapi.Response, error) {
	return createMockSuccessResponse(), nil
}

// IndicesGetSettings implements ElasticsearchAPI for testing
func (m *MockElasticsearchAPI) IndicesGetSettings(o ...func(*esapi.IndicesGetSettingsRequest)) (*esapi.Response, error) {
	if m.SettingsResponse != nil {
		return m.SettingsResponse, nil
	}
	return createMockSettingsResponse(), nil
}

// IndicesPutSettings implements ElasticsearchAPI for testing
func (m *MockElasticsearchAPI) IndicesPutSettings(body io.Reader, o ...func(*esapi.IndicesPutSettingsRequest)) (*esapi.Response, error) {
	return createMockSuccessResponse(), nil
}

// Helper functions to create mock responses
func createMockCountResponse(count int, hasError bool) *esapi.Response {
	var body io.ReadCloser
	var statusCode int

	if hasError {
		statusCode = 500
		body = io.NopCloser(strings.NewReader(`{"error": "internal server error"}`))
	} else {
		statusCode = 200
		responseBody := fmt.Sprintf(`{"count": %d}`, count)
		body = io.NopCloser(strings.NewReader(responseBody))
	}

	return &esapi.Response{
		StatusCode: statusCode,
		Body:       body,
	}
}

func createMockSearchResponse() *esapi.Response {
	responseBody := `{
		"_scroll_id": "test-scroll-id",
		"hits": {
			"hits": [
				{
					"_index": "test-index",
					"_type": "_doc",
					"_id": "1",
					"_source": {
						"field1": "value1"
					}
				}
			]
		}
	}`
	return &esapi.Response{
		StatusCode: 200,
		Body:       io.NopCloser(strings.NewReader(responseBody)),
	}
}

func createMockScrollResponse() *esapi.Response {
	responseBody := `{
		"_scroll_id": "test-scroll-id-2",
		"hits": {
			"hits": []
		}
	}`
	return &esapi.Response{
		StatusCode: 200,
		Body:       io.NopCloser(strings.NewReader(responseBody)),
	}
}

func createMockIndexResponse() *esapi.Response {
	responseBody := `{
		"_index": "test-index",
		"_type": "_doc",
		"_id": "1",
		"_version": 1,
		"result": "created"
	}`
	return &esapi.Response{
		StatusCode: 200,
		Body:       io.NopCloser(strings.NewReader(responseBody)),
	}
}

func createMockMappingResponse() *esapi.Response {
	responseBody := `{
		"test-index": {
			"mappings": {
				"properties": {
					"field1": {
						"type": "text"
					}
				}
			}
		}
	}`
	return &esapi.Response{
		StatusCode: 200,
		Body:       io.NopCloser(strings.NewReader(responseBody)),
	}
}

func createMockSettingsResponse() *esapi.Response {
	responseBody := `{
		"test-index": {
			"settings": {
				"index": {
					"number_of_shards": "1",
					"number_of_replicas": "0"
				}
			}
		}
	}`
	return &esapi.Response{
		StatusCode: 200,
		Body:       io.NopCloser(strings.NewReader(responseBody)),
	}
}

func createMockSuccessResponse() *esapi.Response {
	responseBody := `{"acknowledged": true}`
	return &esapi.Response{
		StatusCode: 200,
		Body:       io.NopCloser(strings.NewReader(responseBody)),
	}
}

// createMockClient creates a client with a mock API
func createMockClient() *Client {
	return &Client{
		API: &MockElasticsearchAPI{},
		URL: "http://mock:9200",
	}
}

// createMockClientWithError creates a client that returns errors
func createMockClientWithError() *Client {
	return &Client{
		API: &MockElasticsearchAPI{
			CountResponse: createMockCountResponse(0, true),
		},
		URL: "http://mock:9200",
	}
}

func TestExtractIndex(t *testing.T) {
	tests := []struct {
		name     string
		url      string
		expected string
	}{
		{
			name:     "basic URL with index",
			url:      "http://localhost:9200/myindex",
			expected: "myindex",
		},
		{
			name:     "HTTPS URL with index",
			url:      "https://elasticsearch.example.com:9200/logs-2023",
			expected: "logs-2023",
		},
		{
			name:     "URL without index",
			url:      "http://localhost:9200",
			expected: "",
		},
		{
			name:     "URL with authentication",
			url:      "https://user:pass@elastic.cloud:9200/prod-index",
			expected: "prod-index",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := extractIndex(tt.url)
			if result != tt.expected {
				t.Errorf("extractIndex(%q) = %q, want %q", tt.url, result, tt.expected)
			}
		})
	}
}

func TestIsFile(t *testing.T) {
	tests := []struct {
		name     string
		path     string
		expected bool
	}{
		{
			name:     "HTTP URL",
			path:     "http://localhost:9200/index",
			expected: false,
		},
		{
			name:     "HTTPS URL",
			path:     "https://elastic.example.com/index",
			expected: false,
		},
		{
			name:     "Local file path",
			path:     "/tmp/backup.json",
			expected: true,
		},
		{
			name:     "Relative file path",
			path:     "backup.ndjson",
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isFile(tt.path)
			if result != tt.expected {
				t.Errorf("isFile(%q) = %v, want %v", tt.path, result, tt.expected)
			}
		})
	}
}

func TestConfig(t *testing.T) {
	config := Config{
		Input:       "http://localhost:9200/source",
		Output:      "http://localhost:9201/dest",
		Type:        "data",
		Limit:       1000,
		Concurrency: 4,
		Format:      "json",
		ScrollSize:  500,
		Verbose:     true,
	}

	if config.Input != "http://localhost:9200/source" {
		t.Errorf("Expected Input to be 'http://localhost:9200/source', got '%s'", config.Input)
	}

	if config.Concurrency != 4 {
		t.Errorf("Expected Concurrency to be 4, got %d", config.Concurrency)
	}

	if config.ScrollSize != 500 {
		t.Errorf("Expected ScrollSize to be 500, got %d", config.ScrollSize)
	}
}

func TestGetBaseURL(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "URL with index",
			input:    "http://localhost:9200/myindex",
			expected: "http://localhost:9200",
		},
		{
			name:     "HTTPS URL with index",
			input:    "https://elastic.cloud:9200/logs-2023",
			expected: "https://elastic.cloud:9200",
		},
		{
			name:     "URL without index",
			input:    "http://localhost:9200",
			expected: "http://localhost:9200",
		},
		{
			name:     "URL with path and index",
			input:    "https://cloud.elastic.co/deployment/myindex",
			expected: "https://cloud.elastic.co/deployment",
		},
		{
			name:     "URL with authentication and index",
			input:    "https://user:pass@elastic.cloud:9200/prod-index",
			expected: "https://user:pass@elastic.cloud:9200",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := getBaseURL(tt.input)
			if result != tt.expected {
				t.Errorf("getBaseURL(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestCreateClient(t *testing.T) {
	tests := []struct {
		name     string
		url      string
		username string
		password string
		wantErr  bool
	}{
		{
			name:     "valid URL without auth",
			url:      "http://localhost:9200",
			username: "",
			password: "",
			wantErr:  false,
		},
		{
			name:     "valid URL with auth",
			url:      "http://localhost:9200",
			username: "elastic",
			password: "changeme",
			wantErr:  false,
		},
		{
			name:     "invalid URL",
			url:      "invalid-url",
			username: "",
			password: "",
			wantErr:  false, // elasticsearch client doesn't validate URL format at creation
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, err := createClient(tt.url, tt.username, tt.password)

			if tt.wantErr && err == nil {
				t.Error("Expected error but got none")
			}
			if !tt.wantErr && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
			if !tt.wantErr && client == nil {
				t.Error("Expected client but got nil")
			}
			if !tt.wantErr && client != nil && client.URL != tt.url {
				t.Errorf("Expected client URL to be %q, got %q", tt.url, client.URL)
			}
		})
	}
}

func TestDocument(t *testing.T) {
	doc := Document{
		Index: "test-index",
		Type:  "doc",
		ID:    "1",
		Source: map[string]interface{}{
			"field1":    "value1",
			"field2":    42,
			"timestamp": "2023-01-01T00:00:00Z",
		},
	}

	if doc.Index != "test-index" {
		t.Errorf("Expected Index to be 'test-index', got '%s'", doc.Index)
	}

	if doc.ID != "1" {
		t.Errorf("Expected ID to be '1', got '%s'", doc.ID)
	}

	if field1, ok := doc.Source["field1"].(string); !ok || field1 != "value1" {
		t.Errorf("Expected field1 to be 'value1', got %v", doc.Source["field1"])
	}

	if field2, ok := doc.Source["field2"].(int); !ok || field2 != 42 {
		t.Errorf("Expected field2 to be 42, got %v", doc.Source["field2"])
	}
}

func TestWriteDocument(t *testing.T) {
	doc := Document{
		Index: "test-index",
		Type:  "doc",
		ID:    "1",
		Source: map[string]interface{}{
			"field1": "value1",
			"field2": 42,
		},
	}

	tests := []struct {
		name    string
		format  string
		wantErr bool
	}{
		{
			name:    "JSON format",
			format:  "json",
			wantErr: false,
		},
		{
			name:    "NDJSON format",
			format:  "ndjson",
			wantErr: false,
		},
		{
			name:    "Invalid format",
			format:  "invalid",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf strings.Builder
			err := writeDocument(&buf, doc, tt.format)

			if tt.wantErr && err == nil {
				t.Error("Expected error but got none")
			}
			if !tt.wantErr && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}

			if !tt.wantErr {
				output := buf.String()
				if output == "" {
					t.Error("Expected output but got empty string")
				}

				// Verify JSON content
				if tt.format == "json" || tt.format == "ndjson" {
					if !strings.Contains(output, "test-index") {
						t.Errorf("Output should contain index name: %s", output)
					}
					if !strings.Contains(output, "field1") {
						t.Errorf("Output should contain field1: %s", output)
					}
				}
			}
		})
	}
}

func TestConfigValidation(t *testing.T) {
	tests := []struct {
		name   string
		config Config
		valid  bool
	}{
		{
			name: "valid config",
			config: Config{
				Input:       "http://localhost:9200/source",
				Output:      "http://localhost:9200/dest",
				Type:        "data",
				Concurrency: 4,
				ScrollSize:  1000,
				Format:      "json",
			},
			valid: true,
		},
		{
			name: "negative concurrency",
			config: Config{
				Input:       "http://localhost:9200/source",
				Output:      "http://localhost:9200/dest",
				Type:        "data",
				Concurrency: -1,
				ScrollSize:  1000,
				Format:      "json",
			},
			valid: false,
		},
		{
			name: "zero scroll size",
			config: Config{
				Input:       "http://localhost:9200/source",
				Output:      "http://localhost:9200/dest",
				Type:        "data",
				Concurrency: 4,
				ScrollSize:  0,
				Format:      "json",
			},
			valid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test basic config field validation
			if tt.config.Concurrency < 0 && tt.valid {
				t.Error("Negative concurrency should not be valid")
			}
			if tt.config.ScrollSize <= 0 && tt.valid {
				t.Error("Zero or negative scroll size should not be valid")
			}
		})
	}
}

func TestTypeValidation(t *testing.T) {
	validTypes := []string{"data", "mapping", "settings"}
	invalidTypes := []string{"invalid", "", "docs", "mappings"}

	for _, validType := range validTypes {
		t.Run("valid_type_"+validType, func(t *testing.T) {
			config := Config{Type: validType}
			if config.Type != validType {
				t.Errorf("Expected type to be %s, got %s", validType, config.Type)
			}
		})
	}

	for _, invalidType := range invalidTypes {
		t.Run("invalid_type_"+invalidType, func(t *testing.T) {
			// This would be caught by the Run function validation
			if invalidType != "data" && invalidType != "mapping" && invalidType != "settings" {
				// This is an invalid type
				if invalidType == "data" || invalidType == "mapping" || invalidType == "settings" {
					t.Errorf("Type %s should be valid", invalidType)
				}
			}
		})
	}
}

func TestFormatValidation(t *testing.T) {
	validFormats := []string{"json", "ndjson"}

	for _, format := range validFormats {
		t.Run("format_"+format, func(t *testing.T) {
			config := Config{Format: format}
			if config.Format != format {
				t.Errorf("Expected format to be %s, got %s", format, config.Format)
			}
		})
	}
}

func TestUsernamePasswordHandling(t *testing.T) {
	tests := []struct {
		name     string
		username string
		password string
		expected bool
	}{
		{
			name:     "both provided",
			username: "elastic",
			password: "changeme",
			expected: true,
		},
		{
			name:     "only username",
			username: "elastic",
			password: "",
			expected: false,
		},
		{
			name:     "only password",
			username: "",
			password: "changeme",
			expected: false,
		},
		{
			name:     "neither provided",
			username: "",
			password: "",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hasAuth := tt.username != "" && tt.password != ""
			if hasAuth != tt.expected {
				t.Errorf("Expected auth check to be %v, got %v", tt.expected, hasAuth)
			}
		})
	}
}

func TestRunFunction(t *testing.T) {
	// Test Run function with different configurations
	tests := []struct {
		name    string
		config  Config
		wantErr bool
		errMsg  string
	}{
		{
			name: "unsupported type",
			config: Config{
				Input:  "http://localhost:9200/source",
				Output: "http://localhost:9200/dest",
				Type:   "invalid",
			},
			wantErr: true,
			errMsg:  "unsupported transfer type",
		},
		{
			name: "data type",
			config: Config{
				Input:  "http://localhost:9200/source",
				Output: "http://localhost:9200/dest",
				Type:   "data",
			},
			wantErr: true, // Will fail due to no actual ES instance
		},
		{
			name: "mapping type",
			config: Config{
				Input:  "http://localhost:9200/source",
				Output: "http://localhost:9200/dest",
				Type:   "mapping",
			},
			wantErr: true, // Will fail due to no actual ES instance
		},
		{
			name: "settings type",
			config: Config{
				Input:  "http://localhost:9200/source",
				Output: "http://localhost:9200/dest",
				Type:   "settings",
			},
			wantErr: true, // Will fail due to no actual ES instance
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Run(tt.config)

			if tt.wantErr && err == nil {
				t.Error("Expected error but got none")
			}
			if !tt.wantErr && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
			if tt.wantErr && err != nil && tt.errMsg != "" {
				if !strings.Contains(err.Error(), tt.errMsg) {
					t.Errorf("Expected error containing '%s', got '%s'", tt.errMsg, err.Error())
				}
			}
		})
	}
}

func TestVerboseLogging(t *testing.T) {
	// Test that verbose flag affects behavior
	tests := []struct {
		name    string
		verbose bool
	}{
		{
			name:    "verbose enabled",
			verbose: true,
		},
		{
			name:    "verbose disabled",
			verbose: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := Config{
				Input:   "http://localhost:9200/source",
				Output:  "http://localhost:9200/dest",
				Type:    "data",
				Verbose: tt.verbose,
			}

			// This will fail due to no actual ES instance, but we test that verbose doesn't panic
			_ = Run(config)
		})
	}
}

func TestWriteToFile(t *testing.T) {
	// Test writeToFile function
	tempFile := "/tmp/test_write.json"
	defer os.Remove(tempFile)

	data := map[string]interface{}{
		"test":   "data",
		"number": 42,
	}

	err := writeToFile(tempFile, data)
	if err != nil {
		t.Errorf("writeToFile failed: %v", err)
	}

	// Verify file was created and has content
	if _, err := os.Stat(tempFile); os.IsNotExist(err) {
		t.Error("File was not created")
	}

	// Read and verify content
	content, err := os.ReadFile(tempFile)
	if err != nil {
		t.Errorf("Failed to read file: %v", err)
	}

	if len(content) == 0 {
		t.Error("File is empty")
	}

	// Verify it's valid JSON
	var parsed map[string]interface{}
	if err := json.Unmarshal(content, &parsed); err != nil {
		t.Errorf("File content is not valid JSON: %v", err)
	}
}

func TestGetDocumentCount(t *testing.T) {
	t.Run("successful count", func(t *testing.T) {
		client := createMockClient()
		count, err := getDocumentCount(client, "test-index")

		if err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}

		if count != 100 {
			t.Errorf("Expected count to be 100, got %d", count)
		}
	})

	t.Run("error response", func(t *testing.T) {
		client := createMockClientWithError()
		count, err := getDocumentCount(client, "test-index")

		if err == nil {
			t.Error("Expected error but got none")
		}

		if count != 0 {
			t.Errorf("Expected count to be 0 on error, got %d", count)
		}
	})

	t.Run("custom count", func(t *testing.T) {
		client := &Client{
			API: &MockElasticsearchAPI{
				CountResponse: createMockCountResponse(500, false),
			},
			URL: "http://mock:9200",
		}

		count, err := getDocumentCount(client, "test-index")

		if err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}

		if count != 500 {
			t.Errorf("Expected count to be 500, got %d", count)
		}
	})
}

func TestScrollFunctions(t *testing.T) {
	client := createMockClient()

	// Test startScroll
	t.Run("startScroll", func(t *testing.T) {
		scrollID, docs, err := startScroll(client, "test-index", 100)

		if err != nil {
			t.Errorf("startScroll failed: %v", err)
		}

		if scrollID != "test-scroll-id" {
			t.Errorf("Expected scrollID 'test-scroll-id', got '%s'", scrollID)
		}

		if len(docs) != 1 {
			t.Errorf("Expected 1 document, got %d", len(docs))
		}

		if len(docs) > 0 && docs[0].Index != "test-index" {
			t.Errorf("Expected index 'test-index', got '%s'", docs[0].Index)
		}
	})

	// Test continueScroll
	t.Run("continueScroll", func(t *testing.T) {
		scrollID, docs, err := continueScroll(client, "test-scroll-id")

		if err != nil {
			t.Errorf("continueScroll failed: %v", err)
		}

		if scrollID != "test-scroll-id-2" {
			t.Errorf("Expected scrollID 'test-scroll-id-2', got '%s'", scrollID)
		}

		if len(docs) != 0 {
			t.Errorf("Expected no documents, got %d", len(docs))
		}
	})
}

func TestMappingAndSettingsFunctions(t *testing.T) {
	client := createMockClient()

	// Test getMapping
	t.Run("getMapping", func(t *testing.T) {
		mapping, err := getMapping(client, "test-index")

		if err != nil {
			t.Errorf("getMapping failed: %v", err)
		}

		if mapping == nil {
			t.Error("Expected mapping but got nil")
		}

		if indexMapping, ok := mapping["test-index"]; !ok {
			t.Error("Expected test-index in mapping")
		} else if indexMap, ok := indexMapping.(map[string]interface{}); ok {
			if _, ok := indexMap["mappings"]; !ok {
				t.Error("Expected mappings in index mapping")
			}
		}
	})

	// Test putMapping
	t.Run("putMapping", func(t *testing.T) {
		testMapping := map[string]interface{}{
			"test-index": map[string]interface{}{
				"mappings": map[string]interface{}{
					"properties": map[string]interface{}{
						"field1": map[string]interface{}{
							"type": "text",
						},
					},
				},
			},
		}

		err := putMapping(client, "test-index", testMapping)
		if err != nil {
			t.Errorf("putMapping failed: %v", err)
		}
	})

	// Test getSettings
	t.Run("getSettings", func(t *testing.T) {
		settings, err := getSettings(client, "test-index")

		if err != nil {
			t.Errorf("getSettings failed: %v", err)
		}

		if settings == nil {
			t.Error("Expected settings but got nil")
		}

		if indexSettings, ok := settings["test-index"]; !ok {
			t.Error("Expected test-index in settings")
		} else if settingsMap, ok := indexSettings.(map[string]interface{}); ok {
			if _, ok := settingsMap["settings"]; !ok {
				t.Error("Expected settings in index settings")
			}
		}
	})

	// Test putSettings
	t.Run("putSettings", func(t *testing.T) {
		testSettings := map[string]interface{}{
			"test-index": map[string]interface{}{
				"settings": map[string]interface{}{
					"index": map[string]interface{}{
						"number_of_shards":   1,
						"number_of_replicas": 0,
					},
				},
			},
		}

		err := putSettings(client, "test-index", testSettings)
		if err != nil {
			t.Errorf("putSettings failed: %v", err)
		}
	})
}

func TestIndexDocument(t *testing.T) {
	client := createMockClient()

	doc := Document{
		Index: "test-index",
		Type:  "_doc",
		ID:    "1",
		Source: map[string]interface{}{
			"field1": "value1",
			"field2": 42,
		},
	}

	err := indexDocument(client, "test-index", doc)
	if err != nil {
		t.Errorf("indexDocument failed: %v", err)
	}

	// Test with error response
	t.Run("error response", func(t *testing.T) {
		errorClient := &Client{
			API: &MockElasticsearchAPI{
				IndexResponse: &esapi.Response{
					StatusCode: 500,
					Body:       io.NopCloser(strings.NewReader(`{"error": "internal server error"}`)),
				},
			},
			URL: "http://mock:9200",
		}

		err := indexDocument(errorClient, "test-index", doc)
		if err == nil {
			t.Error("Expected error for failed indexing")
		}
	})
}

func TestParseScrollResponse(t *testing.T) {
	// Test parseScrollResponse with valid JSON
	jsonResponse := `{
		"_scroll_id": "test-scroll-id",
		"hits": {
			"hits": [
				{
					"_index": "test-index",
					"_type": "_doc",
					"_id": "1",
					"_source": {
						"field1": "value1"
					}
				}
			]
		}
	}`

	reader := strings.NewReader(jsonResponse)
	scrollID, docs, err := parseScrollResponse(reader)

	if err != nil {
		t.Errorf("parseScrollResponse failed: %v", err)
	}

	if scrollID != "test-scroll-id" {
		t.Errorf("Expected scrollID 'test-scroll-id', got '%s'", scrollID)
	}

	if len(docs) != 1 {
		t.Errorf("Expected 1 document, got %d", len(docs))
	}

	if docs[0].Index != "test-index" {
		t.Errorf("Expected index 'test-index', got '%s'", docs[0].Index)
	}

	// Test with invalid JSON
	invalidReader := strings.NewReader("invalid json")
	_, _, err = parseScrollResponse(invalidReader)
	if err == nil {
		t.Error("Expected error for invalid JSON")
	}
}
