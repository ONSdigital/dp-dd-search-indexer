package search

import (
	"github.com/ONSdigital/dp-dd-search-indexer/model"
	"gopkg.in/olivere/elastic.v3"
)

// Checks the elasticSearchClient satisfies the IndexingClient interface
var _ IndexingClient = (*elasticSearchClient)(nil)

// IndexingClient - interface for the indexing functions of a search client.
type IndexingClient interface {
	Index(document *model.Document) error
	Stop()
}

// elasticSearchClient - Elastic Search specific implementation of IndexingClient
type elasticSearchClient struct {
	client *elastic.Client
	index  string
}

// Index the given document.
func (elasticSearch *elasticSearchClient) Index(document *model.Document) error {
	_, err := elasticSearch.client.Index().
		Index(elasticSearch.index).
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
func NewClient(nodes []string, index string) (IndexingClient, error) {
	client, err := elastic.NewClient(
		elastic.SetURL(nodes...),
		elastic.SetMaxRetries(5))
	if err != nil {
		return nil, err
	}

	return &elasticSearchClient{client, index}, nil
}
