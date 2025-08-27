package httpclient

import (
	"io"
	"net/http"
	"time"
)

// Client — интерфейс HTTP клиента для тестирования
type Client interface {
	Fetch(url string) (string, error)
}

// real implementation
type httpClient struct {
	client *http.Client
}

func New() Client {
	return &httpClient{
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (c *httpClient) Fetch(url string) (string, error) {
	resp, err := c.client.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(data), nil
}
