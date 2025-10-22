package client

import (
	"loadsharder/metric"
	"net/http"
	"net/http/httptest"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func createTestServer(responseCode int, delay time.Duration) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if delay > 0 {
			time.Sleep(delay)
		}
		w.WriteHeader(responseCode)
		w.Write([]byte(`{"status": "ok"}`))
	}))
}

func TestClientBaseFunctionality(t *testing.T) {

	server := createTestServer(http.StatusOK, 10*time.Millisecond)
	defer server.Close()

	options := Options{
		Client:            http.DefaultClient,
		RequestTimeout:    20 * time.Second,
		ShardCount:        2,
		WorkersPerShard:   2,
		QueueSizePerShard: 25,
		MaxRPS:            10,
		CapacityFactor:    1.0,
	}

	t.Run("New client", func(t *testing.T) {
		client, err := NewClient(
			options,
			&metric.EmptyMetrics{},
		)
		if err != nil {
			t.Fatalf("Got error: %s", err)
		}
		if client == nil {
			t.Fatal("Got nil client")
		}

		empty, err := NewClient(
			Options{},
			&metric.EmptyMetrics{},
		)

		if err != ErrInvalidParam && empty != nil {
			t.Fatal("Expected nil and err")
		}
	})

	t.Run("Base functionality", func(t *testing.T) {
		client, _ := NewClient(
			options,
			&metric.EmptyMetrics{},
		)

		var wg sync.WaitGroup
		var successCount atomic.Int32

		client.Start()
		client.Start()

		req, _ := http.NewRequest("GET", server.URL, nil)

		for range 10 {
			wg.Add(1)
			client.Add(req, func(resp *http.Response, err error) {
				defer wg.Done()
				if err != nil {
					t.Errorf("Request failed: %v", err)
					return
				}
				if resp.StatusCode != http.StatusOK {
					t.Errorf("Expected status 200, got %d", resp.StatusCode)
					return
				}
				successCount.Add(1)
			})
		}

		wg.Wait()

		if successCount.Load() != 10 {
			t.Fatalf("Expected 10 successful requests, got %d", successCount.Load())
		}

		client.Stop()
		client.Stop()
	})
}
