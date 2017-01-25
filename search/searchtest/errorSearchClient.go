package searchtest

import (
	"errors"
	"github.com/ONSdigital/dp-dd-search-indexer/model"
	"github.com/ONSdigital/dp-dd-search-indexer/search"
)

// Checks the ErrorSearchClient satisfies the IndexingClient interface
var _ search.IndexingClient = (*ErrorSearchClient)(nil)

// NewErrorSearchClient creates a new instance of ErrorSearchClient
func NewErrorSearchClient() *ErrorSearchClient {
	return &ErrorSearchClient{}
}

// ErrorSearchClient provides a mock implementation of IndexingClient
type ErrorSearchClient struct{}

// Index imitates an error being returned from the index method.
func (elasticSearch *ErrorSearchClient) Index(document *model.Document) error {
	return errors.New("went twang")
}

// Stop - mock implementation does nothing.
func (elasticSearch *ErrorSearchClient) Stop() {}
