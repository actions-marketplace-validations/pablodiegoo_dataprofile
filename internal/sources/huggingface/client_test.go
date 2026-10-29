package huggingface

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/pablodiegoo/dataprofile/internal/domain"
)

func TestCollectAggregatesModelsDatasetsAndSpaces(t *testing.T) {
	var requestedAuthors []string

	mux := http.NewServeMux()
	mux.HandleFunc("/api/models", func(w http.ResponseWriter, r *http.Request) {
		requestedAuthors = append(requestedAuthors, "models:"+r.URL.Query().Get("author"))
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode([]map[string]interface{}{
			{"id": "ada/bert-finetuned", "downloads": 1200, "likes": 15},
			{"id": "ada/ tabular-xgb", "downloads": 450, "likes": 8},
		})
	})
	mux.HandleFunc("/api/datasets", func(w http.ResponseWriter, r *http.Request) {
		requestedAuthors = append(requestedAuthors, "datasets:"+r.URL.Query().Get("author"))
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode([]map[string]interface{}{
			{"id": "ada/br-ecommerce", "downloads": 3100, "likes": 22},
			{"id": "ada/climate-series", "downloads": 800, "likes": 9},
			{"id": "ada/nlp-corpus", "downloads": 150, "likes": 4},
		})
	})
	mux.HandleFunc("/api/spaces", func(w http.ResponseWriter, r *http.Request) {
		requestedAuthors = append(requestedAuthors, "spaces:"+r.URL.Query().Get("author"))
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode([]map[string]interface{}{
			{"id": "ada/forecast-app", "likes": 19},
		})
	})

	srv := httptest.NewServer(mux)
	defer srv.Close()

	client := NewClient()
	client.BaseURL = srv.URL

	summary, err := client.Collect(context.Background(), domain.HuggingFaceSourceConfig{
		Username: "ada",
		Enabled:  true,
	})
	if err != nil {
		t.Fatalf("Collect failed: %v", err)
	}

	if summary.Models != 2 || summary.Datasets != 3 || summary.Spaces != 1 {
		t.Errorf("expected Models=2 Datasets=3 Spaces=1, got %+v", summary)
	}
	if summary.Downloads30d != 5700 {
		t.Errorf("expected Downloads30d=5700, got %d", summary.Downloads30d)
	}
	if summary.Likes != 77 {
		t.Errorf("expected Likes=77, got %d", summary.Likes)
	}
	if len(requestedAuthors) != 3 {
		t.Errorf("expected 3 endpoint queries, got %v", requestedAuthors)
	}
}

func TestCollectZeroArtifactsAndDisabledConfig(t *testing.T) {
	calls := 0
	mux := http.NewServeMux()
	for _, ep := range []string{"/api/models", "/api/datasets", "/api/spaces"} {
		mux.HandleFunc(ep, func(w http.ResponseWriter, r *http.Request) {
			calls++
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`[]`))
		})
	}
	srv := httptest.NewServer(mux)
	defer srv.Close()

	client := NewClient()
	client.BaseURL = srv.URL

	// 1. Disabled config should make zero HTTP requests
	disabledOut, err := client.Collect(context.Background(), domain.HuggingFaceSourceConfig{
		Username: "ada",
		Enabled:  false,
	})
	if err != nil {
		t.Fatalf("expected nil error for disabled config, got %v", err)
	}
	if calls != 0 || disabledOut != (Summary{}) {
		t.Errorf("expected 0 calls and zero Summary when disabled, got calls=%d summary=%+v", calls, disabledOut)
	}

	// 2. Enabled user with zero published artifacts returns zero Summary with nil error
	zeroOut, err := client.Collect(context.Background(), domain.HuggingFaceSourceConfig{
		Username: "zero-user",
		Enabled:  true,
	})
	if err != nil {
		t.Fatalf("expected nil error for zero-artifact user, got %v", err)
	}
	if calls != 3 || zeroOut != (Summary{}) {
		t.Errorf("expected 3 calls and zero Summary, got calls=%d summary=%+v", calls, zeroOut)
	}
}

func TestCollectReturnsErrorWhenAllEndpointsFail(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "service unavailable", http.StatusServiceUnavailable)
	}))
	defer srv.Close()

	client := NewClient()
	client.BaseURL = srv.URL

	_, err := client.Collect(context.Background(), domain.HuggingFaceSourceConfig{
		Username: "ada",
		Enabled:  true,
	})
	if err == nil {
		t.Fatalf("expected error when all Hugging Face endpoints return 503, got nil")
	}
}

func TestCollectAuthenticatedPrivateModelsDatasetsAndSpacesWithDeduplication(t *testing.T) {
	var authHeaders []string

	mux := http.NewServeMux()
	mux.HandleFunc("/api/models", func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("Authorization")
		authHeaders = append(authHeaders, auth)
		w.Header().Set("Content-Type", "application/json")
		if auth == "Bearer hf_secret_pat" {
			_ = json.NewEncoder(w).Encode([]map[string]interface{}{
				{"id": "pablodiegoo/public-llm", "downloads": 900, "likes": 12},
				{"id": "pablodiegoo/private-client-classifier", "downloads": 50, "likes": 3},
				{"id": "pablodiegoo/private-client-classifier", "downloads": 50, "likes": 3}, // duplicate entry to verify deduplication
			})
			return
		}
		_ = json.NewEncoder(w).Encode([]map[string]interface{}{
			{"id": "pablodiegoo/public-llm", "downloads": 900, "likes": 12},
		})
	})
	mux.HandleFunc("/api/datasets", func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("Authorization")
		authHeaders = append(authHeaders, auth)
		w.Header().Set("Content-Type", "application/json")
		if auth == "Bearer hf_secret_pat" {
			_ = json.NewEncoder(w).Encode([]map[string]interface{}{
				{"id": "pablodiegoo/public-corpus", "downloads": 400, "likes": 6},
				{"id": "pablodiegoo/private-proprietary-survey", "downloads": 20, "likes": 1},
			})
			return
		}
		_ = json.NewEncoder(w).Encode([]map[string]interface{}{
			{"id": "pablodiegoo/public-corpus", "downloads": 400, "likes": 6},
		})
	})
	mux.HandleFunc("/api/spaces", func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("Authorization")
		authHeaders = append(authHeaders, auth)
		w.Header().Set("Content-Type", "application/json")
		if auth == "Bearer hf_secret_pat" {
			_ = json.NewEncoder(w).Encode([]map[string]interface{}{
				{"id": "pablodiegoo/public-demo", "likes": 8},
				{"id": "pablodiegoo/private-internal-dashboard", "likes": 2},
			})
			return
		}
		_ = json.NewEncoder(w).Encode([]map[string]interface{}{
			{"id": "pablodiegoo/public-demo", "likes": 8},
		})
	})

	srv := httptest.NewServer(mux)
	defer srv.Close()

	authClient := NewClientWithToken("  hf_secret_pat  ")
	authClient.BaseURL = srv.URL

	summary, err := authClient.Collect(context.Background(), domain.HuggingFaceSourceConfig{
		Username: "pablodiegoo",
		Enabled:  true,
	})
	if err != nil {
		t.Fatalf("Authenticated Collect failed: %v", err)
	}

	for i, h := range authHeaders {
		if h != "Bearer hf_secret_pat" {
			t.Errorf("request %d expected Authorization 'Bearer hf_secret_pat', got %q", i, h)
		}
	}
	if summary.Models != 2 || summary.Datasets != 2 || summary.Spaces != 2 {
		t.Errorf("expected deduplicated public+private Models=2 Datasets=2 Spaces=2, got %+v", summary)
	}
	if summary.Downloads30d != 1370 {
		t.Errorf("expected deduplicated Downloads30d=1370, got %d", summary.Downloads30d)
	}
	if summary.Likes != 32 {
		t.Errorf("expected deduplicated Likes=32, got %d", summary.Likes)
	}
}
