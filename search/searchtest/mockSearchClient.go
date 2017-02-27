package searchtest

import (
	"github.com/ONSdigital/dp-dd-search-indexer/model"
	"github.com/ONSdigital/dp-dd-search-indexer/search"
)

// Checks the MockSearchClient satisfies the IndexingClient interface
var _ search.IndexingClient = (*MockSearchClient)(nil)

// NewMockSearchClient creates a new instance of MockSearchClient
func NewMockSearchClient() *MockSearchClient {
	return &MockSearchClient{}
}

// MockSearchClient provides a mock implementation of IndexingClient
type MockSearchClient struct {
	IndexRequests   []IndexRequest
	CustomIndexFunc func(document *model.Document) error
}

// IndexRequest holds the parameters passed to the index function, allowing them to be captured in tests.
type IndexRequest struct {
	Document *model.Document
}

// Index does not index anything, but captures the request for assertions in tests.
func (elasticSearch *MockSearchClient) Index(document *model.Document, indexName string) error {

	if elasticSearch.CustomIndexFunc != nil {
		return elasticSearch.CustomIndexFunc(document)
	}

	elasticSearch.IndexRequests = append(elasticSearch.IndexRequests, IndexRequest{
		Document: document,
	})
	return nil
}

// Stop - mock implementation does nothing.
func (elasticSearch *MockSearchClient) Stop() {}
