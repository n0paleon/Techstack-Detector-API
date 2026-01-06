package fetcher

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"TechstackDetectorAPI/internal/core/domain"

	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/assert"
)

func TestHTTPFetcher_Fetch(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/robots.txt" {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("User-agent: *"))
			return
		}
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	client := resty.New()
	fetcher := NewHTTPFetcher(client, 5)

	t.Run("Success fetch multiple requests with caching", func(t *testing.T) {
		plan := domain.FetchPlan{
			BaseURL: server.URL,
			Requests: []domain.FetchRequest{
				{ID: "req1", Path: "/robots.txt", Method: "GET"},
				{ID: "req2", Path: "/robots.txt", Method: "GET"},
				{ID: "req3", Path: "/404", Method: "GET"},
			},
		}

		results := fetcher.Fetch(context.Background(), plan)

		// assertions
		assert.Len(t, results, 3)

		// success
		assert.NotNil(t, results["req1"])
		assert.Equal(t, http.StatusOK, results["req1"].StatusCode)
		assert.Equal(t, []byte("User-agent: *"), results["req1"].Body)

		// cloned from cache
		assert.Equal(t, "req2", results["req2"].RequestID)
		assert.Equal(t, results["req1"].Body, results["req2"].Body)

		// 404
		assert.Equal(t, http.StatusNotFound, results["req3"].StatusCode)
	})

	t.Run("Handle Connection Error", func(t *testing.T) {
		plan := domain.FetchPlan{
			BaseURL: "http://invalid-url-that-does-not-exist.test",
			Requests: []domain.FetchRequest{
				{ID: "fail1", Path: "/", Method: "GET"},
			},
		}

		results := fetcher.Fetch(context.Background(), plan)

		assert.NotNil(t, results["fail1"].Error)
	})
}

func TestHelperFunctions(t *testing.T) {
	t.Run("executionKey should be consistent", func(t *testing.T) {
		req := domain.FetchRequest{Path: "/test", Method: "POST"}
		key := executionKey("http://api.com", req)
		assert.Equal(t, "POST|http://api.com|/test", key)
	})

	t.Run("cloneWithID should not mutate original", func(t *testing.T) {
		original := &domain.HTTPResult{RequestID: "old"}
		cloned := cloneWithID(original, "new")

		assert.Equal(t, "new", cloned.RequestID)
		assert.Equal(t, "old", original.RequestID)
	})
}
