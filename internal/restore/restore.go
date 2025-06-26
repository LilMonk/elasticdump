package restore

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
	"sync"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	"github.com/schollz/progressbar/v3"
)

// Config holds the configuration for restore operations
type Config struct {
	Input       string
	Output      string
	Type        string
	Concurrency int
	Verbose     bool
	Username    string
	Password    string
}

// Client wraps Elasticsearch client with additional functionality
type Client struct {
	API ElasticsearchAPI
	URL string
}

// Document represents an Elasticsearch document for restore
type Document struct {
	Index  string                 `json:"_index"`
	Type   string                 `json:"_type,omitempty"`
	ID     string                 `json:"_id"`
	Source map[string]interface{} `json:"_source"`
}

// ElasticsearchAPI defines the interface for Elasticsearch operations
type ElasticsearchAPI interface {
	Index(index string, body io.Reader, o ...func(*esapi.IndexRequest)) (*esapi.Response, error)
	IndicesPutMapping(indices []string, body io.Reader, o ...func(*esapi.IndicesPutMappingRequest)) (*esapi.Response, error)
	IndicesPutSettings(body io.Reader, o ...func(*esapi.IndicesPutSettingsRequest)) (*esapi.Response, error)
}

// ElasticsearchClientWrapper wraps the actual Elasticsearch client to implement our interface
type ElasticsearchClientWrapper struct {
	client *elasticsearch.Client
}

// NewElasticsearchClientWrapper creates a new wrapper
func NewElasticsearchClientWrapper(client *elasticsearch.Client) *ElasticsearchClientWrapper {
	return &ElasticsearchClientWrapper{client: client}
}

// Index implements ElasticsearchAPI
func (w *ElasticsearchClientWrapper) Index(index string, body io.Reader, o ...func(*esapi.IndexRequest)) (*esapi.Response, error) {
	return w.client.Index(index, body, o...)
}

// IndicesPutMapping implements ElasticsearchAPI
func (w *ElasticsearchClientWrapper) IndicesPutMapping(indices []string, body io.Reader, o ...func(*esapi.IndicesPutMappingRequest)) (*esapi.Response, error) {
	return w.client.Indices.PutMapping(indices, body, o...)
}

// IndicesPutSettings implements ElasticsearchAPI
func (w *ElasticsearchClientWrapper) IndicesPutSettings(body io.Reader, o ...func(*esapi.IndicesPutSettingsRequest)) (*esapi.Response, error) {
	return w.client.Indices.PutSettings(body, o...)
}

// Run executes the restore operation
func Run(config Config) error {
	if config.Verbose {
		fmt.Printf("Starting restore from %s to %s\n", config.Input, config.Output)
		fmt.Printf("Type: %s, Concurrency: %d\n", config.Type, config.Concurrency)
	}

	switch config.Type {
	case "data":
		return restoreData(config)
	case "mapping":
		return restoreMapping(config)
	case "settings":
		return restoreSettings(config)
	default:
		return fmt.Errorf("unsupported restore type: %s", config.Type)
	}
}

// getBaseURL extracts the base URL from the input string
func getBaseURL(s string) string {
	// Remove any index name from the URL
	return strings.TrimSuffix(s, "/"+extractIndex(s))
}

// createClient creates an Elasticsearch client from URL with optional authentication
func createClient(url, username, password string) (*Client, error) {
	cfg := elasticsearch.Config{
		Addresses: []string{url},
	}

	// Add authentication if provided
	if username != "" && password != "" {
		cfg.Username = username
		cfg.Password = password
	}

	client, err := elasticsearch.NewClient(cfg)
	if err != nil {
		return nil, err
	}

	wrapper := NewElasticsearchClientWrapper(client)
	return &Client{API: wrapper, URL: url}, nil
}

// restoreData restores documents from file to Elasticsearch
func restoreData(config Config) error {
	destURL := getBaseURL(config.Output)
	// Create destination client
	destClient, err := createClient(destURL, config.Username, config.Password)
	if err != nil {
		return fmt.Errorf("failed to create destination client: %w", err)
	}

	// Open input file
	file, err := os.Open(config.Input)
	if err != nil {
		return fmt.Errorf("failed to open input file: %w", err)
	}
	defer file.Close()

	// Get file info for progress tracking
	fileInfo, err := file.Stat()
	if err != nil {
		return fmt.Errorf("failed to get file info: %w", err)
	}

	var bar *progressbar.ProgressBar
	if !config.Verbose {
		bar = progressbar.DefaultBytes(fileInfo.Size(), "Restoring documents")
	}

	// Create worker pool
	docChan := make(chan Document, config.Concurrency*2)
	var wg sync.WaitGroup

	index := extractIndex(config.Output)
	// Start workers
	for i := 0; i < config.Concurrency; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for doc := range docChan {
				if err := indexDocument(destClient, index, doc); err != nil {
					fmt.Printf("Error indexing document %s: %v\n", doc.ID, err)
				}
			}
		}()
	}

	// Read and process documents
	go func() {
		defer close(docChan)

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			line := scanner.Text()
			if line == "" {
				continue
			}

			var doc Document
			if err := json.Unmarshal([]byte(line), &doc); err != nil {
				fmt.Printf("Error parsing document: %v\n", err)
				continue
			}

			docChan <- doc

			if bar != nil {
				bar.Add(len(line))
			}
		}

		if err := scanner.Err(); err != nil {
			fmt.Printf("Error reading file: %v\n", err)
		}
	}()

	wg.Wait()

	if bar != nil {
		bar.Finish()
	}

	if config.Verbose {
		fmt.Printf("Restore completed to %s\n", config.Output)
	}

	return nil
}

// restoreMapping restores index mapping from file
func restoreMapping(config Config) error {
	// Read mapping from file
	data, err := os.ReadFile(config.Input)
	if err != nil {
		return fmt.Errorf("failed to read mapping file: %w", err)
	}

	var mapping map[string]interface{}
	if err := json.Unmarshal(data, &mapping); err != nil {
		return fmt.Errorf("failed to parse mapping: %w", err)
	}

	// Create destination client
	destURL := getBaseURL(config.Output)
	destClient, err := createClient(destURL, config.Username, config.Password)
	if err != nil {
		return fmt.Errorf("failed to create destination client: %w", err)
	}

	// Extract index from output URL
	index := extractIndex(config.Output)
	if index == "" {
		return fmt.Errorf("could not extract index from output URL")
	}

	return putMapping(destClient, index, mapping)
}

// restoreSettings restores index settings from file
func restoreSettings(config Config) error {
	// Read settings from file
	data, err := os.ReadFile(config.Input)
	if err != nil {
		return fmt.Errorf("failed to read settings file: %w", err)
	}

	var settings map[string]interface{}
	if err := json.Unmarshal(data, &settings); err != nil {
		return fmt.Errorf("failed to parse settings: %w", err)
	}

	// Create destination client
	destURL := getBaseURL(config.Output)
	destClient, err := createClient(destURL, config.Username, config.Password)
	if err != nil {
		return fmt.Errorf("failed to create destination client: %w", err)
	}

	// Extract index from output URL
	index := extractIndex(config.Output)
	if index == "" {
		return fmt.Errorf("could not extract index from output URL")
	}

	return putSettings(destClient, index, settings)
}

// Helper functions

func extractIndex(url string) string {
	// Extract index name from URL like http://localhost:9200/myindex
	parts := strings.Split(url, "/")
	if len(parts) > 3 {
		return parts[len(parts)-1]
	}
	return ""
}

func indexDocument(client *Client, index string, doc Document) error {
	data, err := json.Marshal(doc.Source)
	if err != nil {
		return err
	}

	// Use the index from the destination URL (extracted from client.URL)
	if index == "" {
		return fmt.Errorf("output index cannot be empty")
	}

	// Use Elasticsearch API interface for indexing
	res, err := client.API.Index(
		index,
		strings.NewReader(string(data)),
		func(r *esapi.IndexRequest) {
			r.DocumentID = doc.ID
			r.Refresh = "false"
		},
	)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.IsError() {
		body, _ := io.ReadAll(res.Body)
		return fmt.Errorf("indexing failed: [%s] %s", res.Status(), string(body))
	}

	return nil
}

func putMapping(client *Client, index string, mapping map[string]interface{}) error {
	data, err := json.Marshal(mapping)
	if err != nil {
		return err
	}

	res, err := client.API.IndicesPutMapping(
		[]string{index},
		strings.NewReader(string(data)),
	)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("put mapping failed: %s", res.String())
	}

	return nil
}

func putSettings(client *Client, index string, settings map[string]interface{}) error {
	data, err := json.Marshal(settings)
	if err != nil {
		return err
	}

	res, err := client.API.IndicesPutSettings(
		strings.NewReader(string(data)),
		func(r *esapi.IndicesPutSettingsRequest) {
			r.Index = []string{index}
		},
	)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("put settings failed: %s", res.String())
	}

	return nil
}
