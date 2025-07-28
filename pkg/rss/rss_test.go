package rss

import (
	"goRssNews/pkg/storage"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestParse(t *testing.T) {
	// тестовый RSS-фид
	validRSS := `<?xml version="1.0" encoding="UTF-8"?>
				<rss version="2.0">
				<channel>
					<title>Test Feed</title>
					<description>Test Description</description>
					<link>http://example.com</link>
					<item>
						<title>Test Item</title>
						<description><![CDATA[<p>Test content</p>]]></description>
						<pubDate>Mon 1 Jan 2023 12:00:00 +0300</pubDate>
						<link>http://example.com/item1</link>
					</item>
				</channel>
				</rss>`

	// Создаем тестовый сервер
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(validRSS))
	}))
	defer ts.Close()
	
	// Тест успешного парсинга
	t.Run("valid RSS feed", func(t *testing.T) {
		news, err := Parse(ts.URL)
		if err != nil {
			t.Fatalf("Parse() error = %v, want nil", err)
		}

		if len(news) != 1 {
			t.Fatalf("Expected 1 news item, got %d", len(news))
		}

		expected := storage.News{
			Title:       "Test Item",
			Content:     "Test content",
			Link:        "http://example.com/item1",
			PublishedAt: time.Date(2023, time.January, 1, 9, 0, 0, 0, time.UTC).Unix(),
		}

		if news[0].Title != expected.Title {
			t.Errorf("Title = %v, want %v", news[0].Title, expected.Title)
		}
		if news[0].Content != expected.Content {
			t.Errorf("Content = %v, want %v", news[0].Content, expected.Content)
		}
		if news[0].Link != expected.Link {
			t.Errorf("Link = %v, want %v", news[0].Link, expected.Link)
		}
		if news[0].PublishedAt != expected.PublishedAt {
			t.Errorf("PublishedAt = %v, want %v", news[0].PublishedAt, expected.PublishedAt)
		}
	})
}
