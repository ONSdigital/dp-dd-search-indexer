package search

import (
	"github.com/ONSdigital/dp-dd-search-indexer/model"
	"gopkg.in/olivere/elastic.v3"
)

// Checks the elasticSearchClient satisfies the IndexingClient interface
var _ IndexingClient = (*elasticSearchClient)(nil)

// IndexingClient - interface for the indexing functions of a search client.
type IndexingClient interface {
	Index(document *model.Document, indexName string) error
	Stop()
}

// elasticSearchClient - Elastic Search specific implementation of IndexingClient
type elasticSearchClient struct {
	client *elastic.Client
}

// Index the given document.
func (elasticSearch *elasticSearchClient) Index(document *model.Document, indexName string) error {
	_, err := elasticSearch.client.Index().
		Index(indexName).
		Type(document.Type).
		Id(document.ID).
		BodyJson(document).
		Refresh(true).
		Do()

	return err
}

// Stop the elastic search client.
func (elasticSearch *elasticSearchClient) Stop() {
	elasticSearch.client.Stop()
}

// NewClient create a new instance of elasticSearchClient.
func NewClient(nodes []string) (IndexingClient, error) {
	client, err := elastic.NewClient(
		elastic.SetURL(nodes...),
		elastic.SetMaxRetries(5))
	if err != nil {
		return nil, err
	}

	return &elasticSearchClient{client}, nil
}
