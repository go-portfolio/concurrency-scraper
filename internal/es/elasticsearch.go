package es

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/go-portfolio/concurrency-scraper/internal/models"
)

type elasticClient struct {
	client *elasticsearch.Client
	index  string
}

func New(addresses []string, index string) (Client, error) {
	cfg := elasticsearch.Config{
		Addresses: addresses,
	}
	es, err := elasticsearch.NewClient(cfg)
	if err != nil {
		return nil, err
	}
	return &elasticClient{client: es, index: index}, nil
}

func (e *elasticClient) IndexPage(result models.ScrapeResult) error {
	data, _ := json.Marshal(result)
	res, err := e.client.Index(
		e.index,
		bytes.NewReader(data),
		e.client.Index.WithDocumentID(fmt.Sprint(result.URLID)),
	)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.IsError() {
		return fmt.Errorf("error indexing document: %s", res.String())
	}
	return nil
}

func (e *elasticClient) Close() error {
	// Elasticsearch клиент не требует закрытия соединений явно
	return nil
}
