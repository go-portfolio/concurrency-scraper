package es

import (
	"bytes"
	"encoding/json"
	"fmt"
	"time"

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
	data, err := json.Marshal(result) // 1) превращаем Go-структуру в JSON
	if err != nil {
		return fmt.Errorf("failed to marshal document: %w", err)
	}

	res, err := e.client.Index(
		e.index,               // 2) индекс, куда пишем
		bytes.NewReader(data), // 3) сам JSON-документ
		e.client.Index.WithDocumentID(fmt.Sprint(result.URLID)), // 4) задаём ID = result.URLID
		e.client.Index.WithRefresh("wait_for"),                  // 5) ждём, пока индекс обновится (док сразу доступен в поиске)
		e.client.Index.WithOpType("index"),                      // 6) тип операции (может быть "create", тогда упадёт если ID уже есть)
		e.client.Index.WithRouting(fmt.Sprint(result.URLID%5)),  // 7) кастомный routing для распределения по шардам
		e.client.Index.WithTimeout(5*time.Second),               // 8) ждём максимум 5 секунд, иначе ошибка
	)
	if err != nil {
		return fmt.Errorf("failed to index document: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("error indexing document [%s]: %s", result.URLID, res.String())
	}

	return nil
}

func (e *elasticClient) Close() error {
	// Elasticsearch клиент не требует закрытия соединений явно
	return nil
}
