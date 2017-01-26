package handler

import (
	"encoding/json"
	"github.com/ONSdigital/dp-dd-search-indexer/model"
	"github.com/ONSdigital/dp-dd-search-indexer/search"
	"github.com/ONSdigital/go-ns/log"
	"net/http"
	"io/ioutil"
	"io"
)

// Index - HTTP handler for accepting search index requests.
func Index(w http.ResponseWriter, req *http.Request) {

	decoder := json.NewDecoder(req.Body)
	defer func() {
		io.Copy(ioutil.Discard, req.Body)
		err := req.Body.Close()
		if err != nil {
			log.Error(err, log.Data{"message": "Error closing request body."})
		}
	}()

	var document *model.Document
	err := decoder.Decode(&document)
	if err != nil {
		log.Error(err, log.Data{"message": "Failed to parse request data as a document"})
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	log.Debug("Received HTTP request to index data", log.Data{"Document": document})

	err = search.Client.Index(document)
	if err != nil {
		log.Error(err, log.Data{"message": "Error indexing document."})
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
