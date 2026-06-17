package routeutil

import (
	"testing"
)

func TestAPIPathURL(t *testing.T) {
	tests := []struct {
		name    string
		path    APIPath
		baseURL string
		want    string
		wantErr bool
	}{
		{
			name:    "builds URL from path and base",
			path:    APIPath("/apis/test/endpoint"),
			baseURL: "http://localhost:4110",
			want:    "http://localhost:4110/apis/test/endpoint",
		},
		{
			name:    "handles base URL with trailing slash",
			path:    APIPath("/apis/test/endpoint"),
			baseURL: "http://localhost:4110/",
			want:    "http://localhost:4110/apis/test/endpoint",
		},
		{
			name:    "handles base URL with path prefix",
			path:    APIPath("/relative/path"),
			baseURL: "http://localhost:4110/base",
			want:    "http://localhost:4110/relative/path",
		},
		{
			name:    "returns error for empty base URL",
			path:    APIPath("/test"),
			baseURL: "",
			wantErr: true,
		},
		{
			name:    "returns error for invalid base URL",
			path:    APIPath("/test"),
			baseURL: "://invalid",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.path.URL(tt.baseURL)
			if (err != nil) != tt.wantErr {
				t.Errorf("URL() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if got.String() != tt.want {
					t.Errorf("URL() = %v, want %v", got.String(), tt.want)
				}
			}
		})
	}
}

func TestAPIPathURLWithParams(t *testing.T) {
	tests := []struct {
		name    string
		path    APIPath
		baseURL string
		params  map[string]string
		want    string
		wantErr bool
	}{
		{
			name:    "replaces path parameters",
			path:    APIPath("/apis/web/v1/:entity/:id/merge"),
			baseURL: "http://localhost:4110",
			params:  map[string]string{"entity": "track", "id": "456"},
			want:    "http://localhost:4110/apis/web/v1/track/456/merge",
		},
		{
			name:    "handles no parameters",
			path:    APIPath("/apis/test/endpoint"),
			baseURL: "http://localhost:4110",
			params:  nil,
			want:    "http://localhost:4110/apis/test/endpoint",
		},
		{
			name:    "URL-encodes parameter with slashes",
			path:    APIPath("/apis/web/v1/:entity/:id/merge"),
			baseURL: "http://localhost:4110",
			params:  map[string]string{"entity": "track", "id": "../../admin/rules"},
			want:    "http://localhost:4110/apis/web/v1/track/..%2F..%2Fadmin%2Frules/merge",
		},
		{
			name:    "URL-encodes parameter with special characters",
			path:    APIPath("/apis/web/v1/:entity/:id/merge"),
			baseURL: "http://localhost:4110",
			params:  map[string]string{"entity": "track", "id": "123?inject=true"},
			want:    "http://localhost:4110/apis/web/v1/track/123%3Finject=true/merge",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.path.URLWithParams(tt.baseURL, tt.params)
			if (err != nil) != tt.wantErr {
				t.Errorf("URLWithParams() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if got.String() != tt.want {
					t.Errorf("URLWithParams() = %v, want %v", got.String(), tt.want)
				}
			}
		})
	}
}

func TestAPIPathString(t *testing.T) {
	path := APIPath("/apis/test")
	if path.String() != "/apis/test" {
		t.Errorf("String() = %v, want %v", path.String(), "/apis/test")
	}
}
