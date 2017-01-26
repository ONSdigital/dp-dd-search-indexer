package main

import (
	"github.com/ONSdigital/dp-dd-search-indexer/config"
	"github.com/ONSdigital/dp-dd-search-indexer/handler"
	"github.com/ONSdigital/dp-dd-search-indexer/search"
	"github.com/ONSdigital/go-ns/handlers/healthcheck"
	"github.com/ONSdigital/go-ns/log"
	"github.com/bsm/sarama-cluster"
	"github.com/gorilla/pat"
	"github.com/justinas/alice"
	"io"
	"net/http"
	"os"
	"os/signal"
)

func main() {

	config.Load()

	log.Debug("Creating search client.", nil)
	searchClient, err := search.NewClient(config.ElasticSearchNodes, config.ElasticSearchIndex)
	if err != nil {
		log.Error(err, log.Data{"message": "Failed to create Elastic Search client."})
		os.Exit(1)
	}
	handler.SearchClient = searchClient

	log.Debug("Creating Kafka consumer.", nil)
	consumerConfig := cluster.NewConfig()
	kafkaConsumer, err := cluster.NewConsumer(config.KafkaBrokers, config.KafkaConsumerTopic, []string{config.KafkaConsumerTopic}, consumerConfig)
	if err != nil {
		log.Error(err, log.Data{"message": "An error occured creating the Kafka consumer"})
		os.Exit(1)
	}

	exitCh := make(chan struct{})

	listenForKafkaMessages(kafkaConsumer, searchClient, exitCh)
	listenForHTTPRequests(exitCh)
	waitForInterrupt(kafkaConsumer, searchClient, exitCh)
}

func listenForHTTPRequests(exitCh chan struct{}) {

	go func() {
		router := pat.New()
		alice := alice.New().Then(router)
		router.Get("/healthcheck", healthcheck.Handler)
		router.Post("/index", handler.Index)
		log.Debug("Starting server", log.Data{"bind_addr": config.BindAddr})
		server := &http.Server{
			Addr:    config.BindAddr,
			Handler: alice,
		}
		if err := server.ListenAndServe(); err != nil {
			log.Error(err, nil)
		}

		log.Debug("HTTP server has stopped.", nil)
		exitCh <- struct{}{}
	}()
}

func waitForInterrupt(kafkaConsumer io.Closer, searchClient search.IndexingClient, exitCh chan struct{}) {

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, os.Kill)

	select {
	case <-signals:
		log.Debug("OS signal receieved.", nil)
		shutdown(kafkaConsumer, searchClient)
	case <-exitCh:
		log.Debug("Notification received on exit channel.", nil)
		shutdown(kafkaConsumer, searchClient)
	}
}

func shutdown(kafkaConsumer io.Closer, searchClient search.IndexingClient) {
	log.Debug("Shutting down.", nil)
	err := kafkaConsumer.Close()
	if err != nil {
		log.Error(err, log.Data{"message": "An error occured closing the Kafka consumer"})
	}
	searchClient.Stop()
	log.Debug("Service stopped", nil)
}


func listenForKafkaMessages(kafkaConsumer *cluster.Consumer, searchClient search.IndexingClient, exitCh chan struct{}) {

	go func() {
		for message := range kafkaConsumer.Messages() {
			search.ProcessIndexRequest(message.Value, searchClient)
		}

		log.Debug("Kafka consumer has stopped.", nil)
		exitCh <- struct{}{}
	}()

}
