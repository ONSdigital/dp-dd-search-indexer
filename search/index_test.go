package search_test

import (
	"encoding/json"
	"github.com/ONSdigital/dp-dd-search-indexer/model"
	"github.com/ONSdigital/dp-dd-search-indexer/search"
	"github.com/ONSdigital/dp-dd-search-indexer/search/searchtest"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func Test(t *testing.T) {

	Convey("Given a new index request", t, func() {
		expectedRequest := model.Document{
			ID:   "123",
			Type: "thetype",
		}
		documentJson, _ := json.Marshal(expectedRequest)
		searchClient := searchtest.NewMockSearchClient()
		search.Client = searchClient

		Convey("When the index request is processed", func() {
			search.ProcessIndexRequest(documentJson)

			Convey("Then the search client is called with the expected parameters", func() {
				var actualRequest searchtest.IndexRequest = searchClient.IndexRequests[0]
				So(actualRequest.Document.Type, ShouldEqual, expectedRequest.Type)
				So(actualRequest.Document.ID, ShouldEqual, expectedRequest.ID)
			})
		})
	})
}
