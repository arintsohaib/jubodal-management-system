package search

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"

	"github.com/opensearch-project/opensearch-go/v2/opensearchapi"
)

// Service provides high-level search functions
type Service struct {
	client *Client
}

// NewService creates a new search service
func NewService(client *Client) *Service {
	return &Service{client: client}
}

// SearchGlobal searches across multiple entities (mostly users for now)
func (s *Service) SearchGlobal(ctx context.Context, query string, limit, offset int) (interface{}, error) {
	searchQuery := map[string]interface{}{
		"query": map[string]interface{}{
			"multi_match": map[string]interface{}{
				"query":  query,
				"fields": []string{"full_name", "full_name_bn", "position_name", "jurisdiction_name"},
				"fuzziness": "AUTO",
			},
		},
		"from": offset,
		"size": limit,
	}

	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(searchQuery); err != nil {
		return nil, err
	}

	searchReq := opensearchapi.SearchRequest{
		Index: []string{IndexUsers},
		Body:  &buf,
	}

	res, err := searchReq.Do(ctx, s.client.osClient)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.IsError() {
		return nil, fmt.Errorf("search error: %s", res.String())
	}

	var r map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		return nil, err
	}

	return r["hits"], nil
}

// INDEXING

// IndexUser indexes or updates a user in OpenSearch
func (s *Service) IndexUser(ctx context.Context, u map[string]interface{}) error {
	id := fmt.Sprintf("%v", u["id"])
	body, err := json.Marshal(u)
	if err != nil {
		return err
	}

	indexReq := opensearchapi.IndexRequest{
		Index:      IndexUsers,
		DocumentID: id,
		Body:       bytes.NewReader(body),
		Refresh:    "wait_for",
	}

	res, err := indexReq.Do(ctx, s.client.osClient)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("indexing error: %s", res.String())
	}

	return nil
}

// Autocomplete provides suggestions for a partial query
func (s *Service) Autocomplete(ctx context.Context, query string) (interface{}, error) {
	searchQuery := map[string]interface{}{
		"query": map[string]interface{}{
			"bool": map[string]interface{}{
				"should": []interface{}{
					map[string]interface{}{
						"prefix": map[string]interface{}{"full_name": query},
					},
					map[string]interface{}{
						"prefix": map[string]interface{}{"full_name_bn": query},
					},
				},
			},
		},
		"size": 5,
	}

	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(searchQuery); err != nil {
		return nil, err
	}

	searchReq := opensearchapi.SearchRequest{
		Index: []string{IndexUsers},
		Body:  &buf,
	}

	res, err := searchReq.Do(ctx, s.client.osClient)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var r map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		return nil, err
	}

	return r["hits"], nil
}
