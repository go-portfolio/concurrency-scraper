package scraper

import (
	"errors"
	"testing"
	"time"

	"github.com/go-portfolio/concurrency-scraper/internal/models"	
	"github.com/go-portfolio/concurrency-scraper/pkg/logger"
)

type MockDB struct {
	URLs  []models.URL
	Saved []models.ScrapeResult
}

func (m *MockDB) GetURLs() ([]models.URL, error) {
	if m.URLs == nil {
		return nil, errors.New("no urls")
	}
	return m.URLs, nil
}

func (m *MockDB) SaveResult(urlID int, content string) (int, error) {
	return 42, nil
}

func (m *MockDB) SavePageData(r models.ScrapeResult, resultID int) error {
	m.Saved = append(m.Saved, r)
	return nil
}

func (m *MockDB) Exec(query string, args ...interface{}) (interface{}, error) {
	return nil, nil
}

// MockClient реализует httpclient.Client
type MockClient struct{}

func (m *MockClient) Fetch(url string) (string, error) {
	return "<html><title>Test</title><body>Hello world</body></html>", nil
}

// SyncPool — синхронный пул для теста
type SyncPool struct{}

func (p *SyncPool) Submit(task func()) { task() }
func (p *SyncPool) Close()             {}

func TestScraper_Run(t *testing.T) {
	mockDB := &MockDB{
		URLs: []models.URL{
			{ID: 1, URL: "https://example.com"},
			{ID: 2, URL: "https://test.com"},
		},
	}
	mockClient := &MockClient{}
	mockPool := &SyncPool{}
	mockLogger := &logger.MockLogger{}

	s := New(mockLogger, mockClient, mockDB, mockPool)

	err := s.Run(2)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(mockDB.Saved) != len(mockDB.URLs) {
		t.Fatalf("expected %d saved results, got %d", len(mockDB.URLs), len(mockDB.Saved))
	}

	for _, r := range mockDB.Saved {
		if r.Title != "Test" {
			t.Errorf("expected title 'Test', got %q", r.Title)
		}
		if r.WordCount != 2 { // "Hello world"
			t.Errorf("expected word count 2, got %d", r.WordCount)
		}
		if time.Since(r.FetchedAt) > time.Minute {
			t.Errorf("FetchedAt seems wrong: %v", r.FetchedAt)
		}
	}
}
