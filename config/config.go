package config

import (
	"github.com/ONSdigital/go-ns/log"
	"os"
	"strings"
)

// BindAddr - the port to listen for HTTP requests
var BindAddr string

// KafkaBrokers - the list of Kafka brokers to use.
var KafkaBrokers []string

// KafkaConsumerTopic - the Kafka topic name to consume messages from.
var KafkaConsumerTopic string

// KafkaConsumerGroup - the Kafka consumer group to use.
var KafkaConsumerGroup string

// ElasticSearchNodes - a list of elastic search nodes to use
var ElasticSearchNodes []string

// ElasticSearchIndex - the Elastic Search index name to write to.
var ElasticSearchIndex string

// Load any defined environment variables into the configuration.
func Load() {
	BindAddr = getEnvironmentVariable("BIND_ADDR", ":20050")
	KafkaBrokers = strings.Split(getEnvironmentVariable("KAFKA_ADDR", "localhost:9092"), ",")
	KafkaConsumerTopic = getEnvironmentVariable("KAFKA_CONSUMER_TOPIC", "search-index-request")
	KafkaConsumerGroup = getEnvironmentVariable("KAFKA_CONSUMER_GROUP", "search-index-request")
	ElasticSearchNodes = strings.Split(getEnvironmentVariable("ELASTIC_SEARCH_NODES", "http://127.0.0.1:9200"), ",")
	ElasticSearchIndex = getEnvironmentVariable("ELASTIC_SEARCH_INDEX", "ons")

	log.Debug("Loaded configuration values:", log.Data{
		"BIND_ADDR":            BindAddr,
		"KAFKA_ADDR":           KafkaBrokers,
		"KAFKA_CONSUMER_TOPIC": KafkaConsumerTopic,
		"KAFKA_CONSUMER_GROUP": KafkaConsumerGroup,
		"ELASTIC_SEARCH_NODES": ElasticSearchNodes,
		"ELASTIC_SEARCH_INDEX": ElasticSearchIndex,
	})
}

func getEnvironmentVariable(name string, defaultValue string) string {
	environmentValue := os.Getenv(name)
	if environmentValue != "" {
		return environmentValue
	}

	return defaultValue
}
