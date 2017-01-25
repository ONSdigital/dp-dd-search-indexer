package search

import (
	"encoding/json"
	"github.com/ONSdigital/dp-dd-search-indexer/model"
	"github.com/ONSdigital/go-ns/log"
)

// ProcessIndexRequest takes a []byte which contains the raw JSON document and indexes it in search.
func ProcessIndexRequest(msg []byte) {

	var document *model.Document
	err := json.Unmarshal(msg, &document)
	if err != nil {
		log.Debug("Failed to parse json request data", log.Data{"message": string(msg)})
		return
	}

	log.Debug("Indexing document", log.Data{
		"Document": document,
	})

	err = Client.Index(document)
	if err != nil {
		log.Error(err, nil)
	}
}
