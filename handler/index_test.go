package handler_test

import (
	"bytes"
	"github.com/ONSdigital/dp-dd-search-indexer/handler"
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

			handler.SearchClient = searchtest.NewErrorSearchClient()
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

		Convey("When the index handler is called and returns an error", func() {

			handler.SearchClient = searchtest.NewMockSearchClient()
			handler.Index(recorder, request)

			Convey("Then the response code is a 200 - OK", func() {
				So(recorder.Code, ShouldEqual, http.StatusOK)
			})
		})
	})
}
