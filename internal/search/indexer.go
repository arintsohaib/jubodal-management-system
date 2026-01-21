package search

import (
	"context"
	"strings"

	"github.com/opensearch-project/opensearch-go/v2/opensearchapi"
)

const (
	IndexUsers      = "users"
	IndexCommittees = "committees"
	IndexActivities = "activities"
)

// InitIndices creates indices with mappings if they don't exist
func (c *Client) InitIndices(ctx context.Context) error {
	indices := []string{IndexUsers, IndexCommittees, IndexActivities}
	
	for _, index := range indices {
		// Check if exists
		res, err := opensearchapi.IndicesExistsRequest{
			Index: []string{index},
		}.Do(ctx, c.osClient)
		if err != nil {
			return err
		}
		if res.StatusCode == 200 {
			continue // Already exists
		}

		// Create index with mapping (simplified version for now)
		mapping := getMappingFor(index)
		createReq := opensearchapi.IndicesCreateRequest{
			Index: index,
			Body:  strings.NewReader(mapping),
		}
		_, err = createReq.Do(ctx, c.osClient)
		if err != nil {
			return err
		}
	}
	return nil
}

func getMappingFor(index string) string {
	switch index {
	case IndexUsers:
		return `{
			"settings": {
				"analysis": {
					"analyzer": {
						"bangla_analyzer": {
							"type": "custom",
							"tokenizer": "standard",
							"filter": [
								"lowercase",
								"decimal_digit",
								"indic_normalization"
							]
						}
					}
				}
			},
			"mappings": {
				"properties": {
					"id": {"type": "keyword"},
					"full_name": {"type": "text"},
					"full_name_bn": {"type": "text", "analyzer": "bangla_analyzer"},
					"phone": {"type": "keyword"},
					"jurisdiction_name": {"type": "text"},
					"position_name": {"type": "keyword"}
				}
			}
		}`
	default:
		return `{}`
	}
}
