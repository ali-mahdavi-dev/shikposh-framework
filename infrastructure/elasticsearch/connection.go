package elasticsearch

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"time"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
)

type Connection interface {
	IndexDocument(ctx context.Context, index string, id string, document interface{}) error
	GetDocument(ctx context.Context, index string, id string) (map[string]interface{}, error)
	DeleteDocument(ctx context.Context, index string, id string) error
	Search(ctx context.Context, index string, query map[string]interface{}) (map[string]interface{}, error)
	HealthCheck(ctx context.Context) error
}

type connection struct {
	client *elasticsearch.Client
}

func NewElasticsearchConnection(cfg Config) (*connection, error) {
	esCfg := elasticsearch.Config{
		Addresses: []string{fmt.Sprintf("http://%s:%s", cfg.Host, cfg.Port)},
	}

	if cfg.Username != "" && cfg.Password != "" {
		esCfg.Username = cfg.Username
		esCfg.Password = cfg.Password
	}

	client, err := elasticsearch.NewClient(esCfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create elasticsearch client: %w", err)
	}

	conn := &connection{client: client}

	// Test connection with retry (Elasticsearch might not be ready immediately)
	ctx := context.Background()
	maxRetries := 5
	retryDelay := 2 * time.Second

	var lastErr error
	for i := 0; i < maxRetries; i++ {
		if err := conn.HealthCheck(ctx); err == nil {
			return conn, nil
		}
		lastErr = err
		if i < maxRetries-1 {
			time.Sleep(retryDelay)
		}
	}

	return nil, fmt.Errorf("failed to connect to elasticsearch after %d retries: %w", maxRetries, lastErr)
}

type Config struct {
	Host     string
	Port     string
	Username string
	Password string
}

func (c *connection) IndexDocument(ctx context.Context, index string, id string, document interface{}) error {
	jsonData, err := json.Marshal(document)
	if err != nil {
		return fmt.Errorf("failed to marshal document: %w", err)
	}

	req := esapi.IndexRequest{
		Index:      index,
		DocumentID: id,
		Body:       bytes.NewReader(jsonData),
		Refresh:    "true",
	}

	res, err := req.Do(ctx, c.client)
	if err != nil {
		return fmt.Errorf("failed to index document: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		body, _ := io.ReadAll(res.Body)
		return fmt.Errorf("elasticsearch error: %s: %s", res.Status(), string(body))
	}

	return nil
}

func (c *connection) GetDocument(ctx context.Context, index string, id string) (map[string]interface{}, error) {
	req := esapi.GetRequest{
		Index:      index,
		DocumentID: id,
	}

	res, err := req.Do(ctx, c.client)
	if err != nil {
		return nil, fmt.Errorf("failed to get document: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		if res.StatusCode == 404 {
			return nil, fmt.Errorf("document not found")
		}
		body, _ := io.ReadAll(res.Body)
		return nil, fmt.Errorf("elasticsearch error: %s: %s", res.Status(), string(body))
	}

	var result map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return result, nil
}

func (c *connection) DeleteDocument(ctx context.Context, index string, id string) error {
	req := esapi.DeleteRequest{
		Index:      index,
		DocumentID: id,
		Refresh:    "true",
	}

	res, err := req.Do(ctx, c.client)
	if err != nil {
		return fmt.Errorf("failed to delete document: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		if res.StatusCode == 404 {
			return nil // Document doesn't exist, consider it deleted
		}
		body, _ := io.ReadAll(res.Body)
		return fmt.Errorf("elasticsearch error: %s: %s", res.Status(), string(body))
	}

	return nil
}

func (c *connection) Search(ctx context.Context, index string, query map[string]interface{}) (map[string]interface{}, error) {
	queryJSON, err := json.Marshal(query)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal query: %w", err)
	}

	req := esapi.SearchRequest{
		Index: []string{index},
		Body:  bytes.NewReader(queryJSON),
	}

	res, err := req.Do(ctx, c.client)
	if err != nil {
		return nil, fmt.Errorf("failed to search: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		body, _ := io.ReadAll(res.Body)
		return nil, fmt.Errorf("elasticsearch error: %s: %s", res.Status(), string(body))
	}

	var result map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return result, nil
}

func (c *connection) HealthCheck(ctx context.Context) error {
	res, err := c.client.Info()
	if err != nil {
		return fmt.Errorf("failed to ping elasticsearch: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		body, _ := io.ReadAll(res.Body)
		return fmt.Errorf("elasticsearch health check failed: %s: %s", res.Status(), string(body))
	}

	return nil
}
