package search

import (
	"crypto/tls"
	"fmt"
	"net/http"

	"github.com/opensearch-project/opensearch-go/v2"
)

// Client wraps the opensearch client
type Client struct {
	osClient *opensearch.Client
}

// NewClient creates a new OpenSearch client
func NewClient(url string) (*Client, error) {
	osClient, err := opensearch.NewClient(opensearch.Config{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, // Dev only
		},
		Addresses: []string{url},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create opensearch client: %w", err)
	}
	return &Client{osClient: osClient}, nil
}
