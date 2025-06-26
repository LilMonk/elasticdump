package transfer

import (
	"testing"
)

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
