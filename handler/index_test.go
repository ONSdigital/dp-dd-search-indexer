package handler_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/ONSdigital/dp-dd-search-indexer/handler"
	"github.com/ONSdigital/dp-dd-search-indexer/model"
	"github.com/ONSdigital/dp-dd-search-indexer/search/searchtest"
	. "github.com/smartystreets/goconvey/convey"
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test(t *testing.T) {

	validJSONDocument := "{\"type\":\"testtype\",\"id\":\"234\"}"

	Convey("Given a HTTP request with an invalid document JSON object in the body", t, func() {

		recorder := httptest.NewRecorder()
		requestBodyReader := bytes.NewReader([]byte("{not a valid document}"))
		request, _ := http.NewRequest("POST", "/", requestBodyReader)

		Convey("When the index handler is called", func() {

			handler.Index(recorder, request)

			Convey("Then the response code is a 400 - bad request", func() {
				So(recorder.Code, ShouldEqual, http.StatusBadRequest)
			})
		})
	})

	Convey("Given a valid JSON input, but invalid search request ", t, func() {

		recorder := httptest.NewRecorder()
		requestBodyReader := bytes.NewReader([]byte("[]"))
		request, _ := http.NewRequest("POST", "/", requestBodyReader)

		Convey("When the index handler is called", func() {

			handler.SearchClient = searchtest.NewMockSearchClient()
			handler.Index(recorder, request)

			Convey("Then the response code is a 400 - bad request", func() {
				So(recorder.Code, ShouldEqual, http.StatusBadRequest)
			})
		})
	})

	Convey("Given a valid search request ", t, func() {

		recorder := httptest.NewRecorder()
		requestBodyReader := bytes.NewReader([]byte(validJSONDocument))
		request, _ := http.NewRequest("POST", "/", requestBodyReader)

		Convey("When the index handler is called and returns an error", func() {

			mockSearchClient := searchtest.NewMockSearchClient()
			mockSearchClient.CustomIndexFunc = func(document *model.Document) error {
				return errors.New("went twang")
			}
			handler.SearchClient = mockSearchClient

			handler.Index(recorder, request)

			Convey("Then the response code is a 500 - internal server error", func() {
				So(recorder.Code, ShouldEqual, http.StatusInternalServerError)
			})
		})
	})

	Convey("Given a valid search request ", t, func() {

		recorder := httptest.NewRecorder()
		requestBodyReader := bytes.NewReader([]byte(validJSONDocument))
		request, _ := http.NewRequest("POST", "/", requestBodyReader)
		var expectedDocument model.Document
		_ = json.Unmarshal([]byte(validJSONDocument), &expectedDocument)

		Convey("When the index handler is called and returns an error", func() {

			mockSearchClient := searchtest.NewMockSearchClient()
			handler.SearchClient = mockSearchClient
			handler.Index(recorder, request)

			Convey("Then the response code is a 200 - OK", func() {
				So(recorder.Code, ShouldEqual, http.StatusOK)
				actualDocument := mockSearchClient.IndexRequests[0].Document
				So(actualDocument.ID, ShouldEqual, expectedDocument.ID)
				So(actualDocument.Title, ShouldEqual, expectedDocument.Title)
			})
		})
	})
}
