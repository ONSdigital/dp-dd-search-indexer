package main

import (
	"fmt"
	"github.com/ONSdigital/dp-dd-search-indexer/config"
	"github.com/ONSdigital/dp-dd-search-indexer/handler"
	"github.com/ONSdigital/dp-dd-search-indexer/search"
	"github.com/ONSdigital/go-ns/handlers/healthcheck"
	"github.com/ONSdigital/go-ns/handlers/requestID"
	"github.com/ONSdigital/go-ns/handlers/timeout"
	"github.com/ONSdigital/go-ns/log"
	"github.com/gorilla/pat"
	"github.com/justinas/alice"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func main() {

	fmt.Println("Starting search indexer")
	config.Load()

	log.Debug("Creating search client.", nil)
	searchClient, err := search.NewClient(config.ElasticSearchNodes)
	if err != nil {
		log.Error(err, log.Data{"message": "Failed to create Elastic Search client."})
		os.Exit(1)
	}
	handler.SearchClient = searchClient
	handler.DocumentIndex = config.ElasticSearchIndex
	handler.AreaIndex = config.ElasticSearchAreaIndex

	exitCh := make(chan struct{})

	listenForHTTPRequests(exitCh)
	waitForInterrupt(searchClient, exitCh)
}

func listenForHTTPRequests(exitCh chan struct{}) {

	go func() {
		router := pat.New()
		router.Get("/healthcheck", healthcheck.Handler)
		router.Post("/index-area", handler.IndexGeographicArea)
		router.Post("/index", handler.Index)

		log.Debug("Starting HTTP server", log.Data{"bind_addr": config.BindAddr})

		middleware := []alice.Constructor{
			requestID.Handler(16),
			log.Handler,
			timeout.Handler(10 * time.Second),
		}
		alice := alice.New(middleware...).Then(router)

		server := &http.Server{
			Addr:         config.BindAddr,
			Handler:      alice,
			ReadTimeout:  5 * time.Second,
			WriteTimeout: 10 * time.Second,
		}
		if err := server.ListenAndServe(); err != nil {
			log.Error(err, nil)
		}

		log.Debug("HTTP server has stopped.", nil)
		exitCh <- struct{}{}
	}()
}

func waitForInterrupt(searchClient search.IndexingClient, exitCh chan struct{}) {

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, os.Kill)

	select {
	case <-signals:
		log.Debug("OS signal receieved.", nil)
		shutdown(searchClient)
	case <-exitCh:
		log.Debug("Notification received on exit channel.", nil)
		shutdown(searchClient)
	}
}

func shutdown(searchClient search.IndexingClient) {
	log.Debug("Shutting down.", nil)
	searchClient.Stop()
	log.Debug("Service stopped", nil)
}
