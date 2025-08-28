package httpclient

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// Тест успешного запроса: сервер возвращает "hello world"
func TestFetchSuccess(t *testing.T) {
	// Поднимаем тестовый сервер
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello world"))
	}))
	defer ts.Close()

	c := New()
	body, err := c.Fetch(ts.URL)

	assert.NoError(t, err)
	assert.Equal(t, "hello world", body)
}

// Тест ошибки: сервер возвращает 500
func TestFetchServerError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "internal error", http.StatusInternalServerError)
	}))
	defer ts.Close()

	c := New()
	body, err := c.Fetch(ts.URL)

	// Важно: http.Get не считает 500 ошибкой, поэтому err == nil
	assert.NoError(t, err)
	assert.Contains(t, body, "internal error")
}

// Тест ошибки: невалидный URL
func TestFetchInvalidURL(t *testing.T) {
	c := New()
	_, err := c.Fetch("http://invalid_host")

	assert.Error(t, err)
}

// Тест таймаута
func TestFetchTimeout(t *testing.T) {
	// Сервер, который не отвечает
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(2 * time.Second)
	}))
	defer ts.Close()

	// Клиент с коротким таймаутом
	c := &httpClient{
		client: &http.Client{
			Timeout: 500 * time.Millisecond,
		},
	}

	_, err := c.Fetch(ts.URL)
	assert.Error(t, err)
	// Обычно ошибка будет вида context deadline exceeded
	assert.True(t, errors.Is(err, http.ErrHandlerTimeout) || err != nil)
}
