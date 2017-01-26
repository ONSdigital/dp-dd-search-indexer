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
		return
	}
	handler.SearchClient = searchClient

	log.Debug("Creating Kafka consumer.", nil)
	consumerConfig := cluster.NewConfig()
	kafkaConsumer, err := cluster.NewConsumer(config.KafkaBrokers, config.KafkaConsumerTopic, []string{config.KafkaConsumerTopic}, consumerConfig)
	if err != nil {
		log.Error(err, log.Data{"message": "An error occured creating the Kafka consumer"})
		return
	}

	listenForKafkaMessages(kafkaConsumer, searchClient)
	listenForHTTPRequests()
	waitForInterrupt(kafkaConsumer, searchClient)
}

func listenForHTTPRequests() {

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
	}()
}

func waitForInterrupt(kafkaConsumer io.Closer, searchClient search.IndexingClient) {

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)
	<-signals
	log.Debug("Shutting down...", nil)
	err := kafkaConsumer.Close()
	if err != nil {
		log.Error(err, log.Data{"message": "An error occured closing the Kafka consumer"})
	}
	searchClient.Stop()
	log.Debug("Service stopped", nil)

}

func listenForKafkaMessages(kafkaConsumer *cluster.Consumer, searchClient search.IndexingClient) {

	go func() {
		for message := range kafkaConsumer.Messages() {
			search.ProcessIndexRequest(message.Value, searchClient)
		}
	}()

}
