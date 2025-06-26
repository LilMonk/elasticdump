package transfer

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	"github.com/schollz/progressbar/v3"
)

// ElasticsearchAPI defines the interface for Elasticsearch operations
type ElasticsearchAPI interface {
	Count(o ...func(*esapi.CountRequest)) (*esapi.Response, error)
	Search(o ...func(*esapi.SearchRequest)) (*esapi.Response, error)
	Scroll(o ...func(*esapi.ScrollRequest)) (*esapi.Response, error)
	Index(index string, body io.Reader, o ...func(*esapi.IndexRequest)) (*esapi.Response, error)
	IndicesGetMapping(o ...func(*esapi.IndicesGetMappingRequest)) (*esapi.Response, error)
	IndicesPutMapping(indices []string, body io.Reader, o ...func(*esapi.IndicesPutMappingRequest)) (*esapi.Response, error)
	IndicesGetSettings(o ...func(*esapi.IndicesGetSettingsRequest)) (*esapi.Response, error)
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

// Count implements ElasticsearchAPI
func (w *ElasticsearchClientWrapper) Count(o ...func(*esapi.CountRequest)) (*esapi.Response, error) {
	return w.client.Count(o...)
}

// Search implements ElasticsearchAPI
func (w *ElasticsearchClientWrapper) Search(o ...func(*esapi.SearchRequest)) (*esapi.Response, error) {
	return w.client.Search(o...)
}

// Scroll implements ElasticsearchAPI
func (w *ElasticsearchClientWrapper) Scroll(o ...func(*esapi.ScrollRequest)) (*esapi.Response, error) {
	return w.client.Scroll(o...)
}

// Index implements ElasticsearchAPI
func (w *ElasticsearchClientWrapper) Index(index string, body io.Reader, o ...func(*esapi.IndexRequest)) (*esapi.Response, error) {
	return w.client.Index(index, body, o...)
}

// IndicesGetMapping implements ElasticsearchAPI
func (w *ElasticsearchClientWrapper) IndicesGetMapping(o ...func(*esapi.IndicesGetMappingRequest)) (*esapi.Response, error) {
	return w.client.Indices.GetMapping(o...)
}

// IndicesPutMapping implements ElasticsearchAPI
func (w *ElasticsearchClientWrapper) IndicesPutMapping(indices []string, body io.Reader, o ...func(*esapi.IndicesPutMappingRequest)) (*esapi.Response, error) {
	return w.client.Indices.PutMapping(indices, body, o...)
}

// IndicesGetSettings implements ElasticsearchAPI
func (w *ElasticsearchClientWrapper) IndicesGetSettings(o ...func(*esapi.IndicesGetSettingsRequest)) (*esapi.Response, error) {
	return w.client.Indices.GetSettings(o...)
}

// IndicesPutSettings implements ElasticsearchAPI
func (w *ElasticsearchClientWrapper) IndicesPutSettings(body io.Reader, o ...func(*esapi.IndicesPutSettingsRequest)) (*esapi.Response, error) {
	return w.client.Indices.PutSettings(body, o...)
}

// Config holds the configuration for transfer operations
type Config struct {
	Input       string
	Output      string
	Type        string
	Limit       int
	Concurrency int
	Format      string
	ScrollSize  int
	Verbose     bool
	Username    string
	Password    string
}

// Client wraps Elasticsearch client with additional functionality
type Client struct {
	API ElasticsearchAPI
	URL string
}

// Document represents an Elasticsearch document
type Document struct {
	Index  string                 `json:"_index"`
	Type   string                 `json:"_type,omitempty"`
	ID     string                 `json:"_id"`
	Source map[string]interface{} `json:"_source"`
}

// Run executes the transfer operation
func Run(config Config) error {
	if config.Verbose {
		fmt.Printf("Starting transfer from %s to %s\n", config.Input, config.Output)
		fmt.Printf("Type: %s, Concurrency: %d, ScrollSize: %d\n",
			config.Type, config.Concurrency, config.ScrollSize)
	}

	sourceURL := getBaseURL(config.Input)
	sourceClient, err := createClient(sourceURL, config.Username, config.Password)
	if err != nil {
		return fmt.Errorf("failed to create source client: %w", err)
	}

	switch config.Type {
	case "data":
		return transferData(sourceClient, config)
	case "mapping":
		return transferMapping(sourceClient, config)
	case "settings":
		return transferSettings(sourceClient, config)
	default:
		return fmt.Errorf("unsupported transfer type: %s", config.Type)
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

// transferData transfers documents between clusters
func transferData(sourceClient *Client, config Config) error {
	// Parse index from input URL
	index := extractIndex(config.Input)
	if index == "" {
		return fmt.Errorf("could not extract index from input URL")
	}

	// Check if output is a file or Elasticsearch URL
	if isFile(config.Output) {
		return exportToFile(sourceClient, index, config)
	}

	// Transfer to another Elasticsearch cluster
	destURL := getBaseURL(config.Output)
	destClient, err := createClient(destURL, config.Username, config.Password)
	if err != nil {
		return fmt.Errorf("failed to create destination client: %w", err)
	}

	return transferBetweenClusters(sourceClient, destClient, index, config)
}

// exportToFile exports data to a file
func exportToFile(client *Client, index string, config Config) error {
	file, err := os.Create(config.Output)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer file.Close()

	// Get total count for progress bar
	total, err := getDocumentCount(client, index)
	if err != nil {
		return fmt.Errorf("failed to get document count: %w", err)
	}

	if config.Limit > 0 && config.Limit < total {
		total = config.Limit
	}

	var bar *progressbar.ProgressBar
	if !config.Verbose {
		bar = progressbar.DefaultBytes(int64(total), "Exporting documents")
	}

	// Start scrolling
	scrollSize := min(config.ScrollSize, total)
	scrollID, docs, err := startScroll(client, index, scrollSize)
	if err != nil {
		return fmt.Errorf("failed to start scroll: %w", err)
	}

	exported := 0
	for len(docs) > 0 && (config.Limit == 0 || exported < config.Limit) {
		for _, doc := range docs {
			if config.Limit > 0 && exported >= config.Limit {
				break
			}

			if err := writeDocument(file, doc, config.Format); err != nil {
				return fmt.Errorf("failed to write document: %w", err)
			}

			exported++
			if bar != nil {
				bar.Add(1)
			}
		}

		if config.Limit > 0 && exported >= config.Limit {
			break
		}

		// Continue scrolling
		scrollID, docs, err = continueScroll(client, scrollID)
		if err != nil {
			return fmt.Errorf("failed to continue scroll: %w", err)
		}
	}

	if bar != nil {
		bar.Finish()
	}

	if config.Verbose {
		fmt.Printf("Exported %d documents to %s\n", exported, config.Output)
	}

	return nil
}

// transferBetweenClusters transfers data between two Elasticsearch clusters
func transferBetweenClusters(sourceClient, destClient *Client, index string, config Config) error {
	destIndex := extractIndex(config.Output)
	if destIndex == "" {
		destIndex = index
	}

	// Get total count for progress bar
	total, err := getDocumentCount(sourceClient, index)
	if err != nil {
		return fmt.Errorf("failed to get document count: %w", err)
	}

	if config.Limit > 0 && config.Limit < total {
		total = config.Limit
	}

	var bar *progressbar.ProgressBar
	if !config.Verbose {
		bar = progressbar.DefaultBytes(int64(total), "Transferring documents")
	}

	// Create worker pool
	docChan := make(chan Document, config.Concurrency*2)
	var wg sync.WaitGroup

	// Start workers
	for i := 0; i < config.Concurrency; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for doc := range docChan {
				if err := indexDocument(destClient, destIndex, doc); err != nil {
					fmt.Printf("Error indexing document %s: %v\n", doc.ID, err)
				}
				if bar != nil {
					bar.Add(1)
				}
			}
		}()
	}

	// Start scrolling and send documents to workers
	go func() {
		defer close(docChan)

		scrollID, docs, err := startScroll(sourceClient, index, config.ScrollSize)
		if err != nil {
			fmt.Printf("Failed to start scroll: %v\n", err)
			return
		}

		transferred := 0
		for len(docs) > 0 && (config.Limit == 0 || transferred < config.Limit) {
			for _, doc := range docs {
				if config.Limit > 0 && transferred >= config.Limit {
					return
				}

				docChan <- doc
				transferred++
			}

			if config.Limit > 0 && transferred >= config.Limit {
				return
			}

			// Continue scrolling
			scrollID, docs, err = continueScroll(sourceClient, scrollID)
			if err != nil {
				fmt.Printf("Failed to continue scroll: %v\n", err)
				return
			}
		}
	}()

	wg.Wait()

	if bar != nil {
		bar.Finish()
	}

	if config.Verbose {
		fmt.Printf("Transfer completed to %s\n", config.Output)
	}

	return nil
}

// transferMapping transfers index mapping
func transferMapping(sourceClient *Client, config Config) error {
	index := extractIndex(config.Input)
	if index == "" {
		return fmt.Errorf("could not extract index from input URL")
	}

	mapping, err := getMapping(sourceClient, index)
	if err != nil {
		return fmt.Errorf("failed to get mapping: %w", err)
	}

	if isFile(config.Output) {
		return writeToFile(config.Output, mapping)
	}

	destURL := getBaseURL(config.Output)
	destClient, err := createClient(destURL, config.Username, config.Password)
	if err != nil {
		return fmt.Errorf("failed to create destination client: %w", err)
	}

	destIndex := extractIndex(config.Output)
	if destIndex == "" {
		destIndex = index
	}

	return putMapping(destClient, destIndex, mapping)
}

// transferSettings transfers index settings
func transferSettings(sourceClient *Client, config Config) error {
	index := extractIndex(config.Input)
	if index == "" {
		return fmt.Errorf("could not extract index from input URL")
	}

	settings, err := getSettings(sourceClient, index)
	if err != nil {
		return fmt.Errorf("failed to get settings: %w", err)
	}

	if isFile(config.Output) {
		return writeToFile(config.Output, settings)
	}

	destURL := getBaseURL(config.Output)
	destClient, err := createClient(destURL, config.Username, config.Password)
	if err != nil {
		return fmt.Errorf("failed to create destination client: %w", err)
	}

	destIndex := extractIndex(config.Output)
	if destIndex == "" {
		destIndex = index
	}

	return putSettings(destClient, destIndex, settings)
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

func isFile(path string) bool {
	return !strings.HasPrefix(path, "http://") && !strings.HasPrefix(path, "https://")
}

func getDocumentCount(client *Client, index string) (int, error) {
	// Use Elasticsearch API interface for count
	res, err := client.API.Count(
		func(r *esapi.CountRequest) {
			r.Index = []string{index}
		},
	)
	if err != nil {
		return 0, err
	}
	defer res.Body.Close()

	if res.IsError() {
		return 0, fmt.Errorf("count request failed: %s", res.String())
	}

	var result map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return 0, err
	}

	count, ok := result["count"].(float64)
	if !ok {
		return 0, fmt.Errorf("invalid count response")
	}

	return int(count), nil
}

func startScroll(client *Client, index string, size int) (string, []Document, error) {
	// Use Elasticsearch API interface for scroll search
	res, err := client.API.Search(
		func(r *esapi.SearchRequest) {
			r.Index = []string{index}
			r.Scroll = time.Minute * 5
			r.Size = &size
		},
	)
	if err != nil {
		return "", nil, err
	}
	defer res.Body.Close()

	if res.IsError() {
		return "", nil, fmt.Errorf("search failed: %s", res.String())
	}

	return parseScrollResponse(res.Body)
}

func continueScroll(client *Client, scrollID string) (string, []Document, error) {
	// Use Elasticsearch API interface for scroll continuation
	res, err := client.API.Scroll(
		func(r *esapi.ScrollRequest) {
			r.ScrollID = scrollID
			r.Scroll = time.Minute * 5
		},
	)
	if err != nil {
		return "", nil, err
	}
	defer res.Body.Close()

	if res.IsError() {
		return "", nil, fmt.Errorf("scroll failed: %s", res.String())
	}

	return parseScrollResponse(res.Body)
}

func parseScrollResponse(body io.Reader) (string, []Document, error) {
	var result map[string]interface{}
	if err := json.NewDecoder(body).Decode(&result); err != nil {
		return "", nil, err
	}

	scrollID, _ := result["_scroll_id"].(string)

	hits, ok := result["hits"].(map[string]interface{})
	if !ok {
		return scrollID, nil, nil
	}

	hitsList, ok := hits["hits"].([]interface{})
	if !ok {
		return scrollID, nil, nil
	}

	var docs []Document
	for _, hit := range hitsList {
		hitMap, ok := hit.(map[string]interface{})
		if !ok {
			continue
		}

		doc := Document{
			Index:  hitMap["_index"].(string),
			ID:     hitMap["_id"].(string),
			Source: hitMap["_source"].(map[string]interface{}),
		}

		if t, exists := hitMap["_type"]; exists {
			doc.Type = t.(string)
		}

		docs = append(docs, doc)
	}

	return scrollID, docs, nil
}

func writeDocument(writer io.Writer, doc Document, format string) error {
	switch format {
	case "json":
		return json.NewEncoder(writer).Encode(doc)
	case "ndjson":
		data, err := json.Marshal(doc)
		if err != nil {
			return err
		}
		_, err = fmt.Fprintf(writer, "%s\n", data)
		return err
	default:
		return fmt.Errorf("unsupported format: %s", format)
	}
}

func indexDocument(client *Client, index string, doc Document) error {
	data, err := json.Marshal(doc.Source)
	if err != nil {
		return err
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

func getMapping(client *Client, index string) (map[string]interface{}, error) {
	res, err := client.API.IndicesGetMapping(
		func(r *esapi.IndicesGetMappingRequest) {
			r.Index = []string{index}
		},
	)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var result map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, err
	}

	return result, nil
}

func putMapping(client *Client, index string, mapping map[string]interface{}) error {
	// Extract the mapping for the specific index
	indexMapping, ok := mapping[index].(map[string]interface{})
	if !ok {
		return fmt.Errorf("invalid mapping format")
	}

	mappingData, ok := indexMapping["mappings"].(map[string]interface{})
	if !ok {
		return fmt.Errorf("no mappings found")
	}

	data, err := json.Marshal(mappingData)
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

func getSettings(client *Client, index string) (map[string]interface{}, error) {
	res, err := client.API.IndicesGetSettings(
		func(r *esapi.IndicesGetSettingsRequest) {
			r.Index = []string{index}
		},
	)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var result map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, err
	}

	return result, nil
}

func putSettings(client *Client, index string, settings map[string]interface{}) error {
	// Extract the settings for the specific index
	indexSettings, ok := settings[index].(map[string]interface{})
	if !ok {
		return fmt.Errorf("invalid settings format")
	}

	settingsData, ok := indexSettings["settings"].(map[string]interface{})
	if !ok {
		return fmt.Errorf("no settings found")
	}

	data, err := json.Marshal(settingsData)
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

func writeToFile(filename string, data interface{}) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	return json.NewEncoder(file).Encode(data)
}
