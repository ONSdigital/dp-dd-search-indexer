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
	"github.com/bsm/sarama-cluster"
	"github.com/gorilla/pat"
	"github.com/justinas/alice"
	"io"
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

	var kafkaConsumer *cluster.Consumer
	exitCh := make(chan struct{})

	if len(config.KafkaBrokers) > 0 && config.KafkaBrokers[0] != "" {
		log.Debug("Creating Kafka consumer.", nil)
		consumerConfig := cluster.NewConfig()
		kafkaConsumer, err = cluster.NewConsumer(config.KafkaBrokers, config.KafkaConsumerTopic, []string{config.KafkaConsumerTopic}, consumerConfig)
		if err != nil {
			log.Error(err, log.Data{"message": "An error occured creating the Kafka consumer"})
			os.Exit(1)
		}

		listenForKafkaMessages(kafkaConsumer, searchClient, config.ElasticSearchIndex, exitCh)
	}

	listenForHTTPRequests(exitCh)
	waitForInterrupt(kafkaConsumer, searchClient, exitCh)
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

func listenForKafkaMessages(kafkaConsumer *cluster.Consumer,
	searchClient search.IndexingClient,
	searchIndex string,
	exitCh chan struct{}) {

	go func() {
		for message := range kafkaConsumer.Messages() {
			search.ProcessIndexRequest(message.Value, searchClient, searchIndex)
		}

		log.Debug("Kafka consumer has stopped.", nil)
		exitCh <- struct{}{}
	}()

}

//func bootstrap(client search.IndexingClient) {
//	armedForcesDataset := &model.Dataset{
//		ID:    "AF001EW",
//		Title: "AF001EW  Members of the Armed Forces by residence type by sex by age",
//		URL:   "http://localhost:20099/datasets/AF001EW",
//		Metadata: &model.Metadata{
//			Description:        "This dataset provides 2011 Census estimates that classify usual residents aged 16 and over who are members of the armed forces by residence type (household or communal resident), by sex and by age. The estimates are as at census day, 27 March 2011.",
//			NationalStatistics: true,
//			Contact: &model.Contact{
//				Name:  "Alexa Bradley",
//				Phone: "+44 (0)1329 444972",
//				Email: "pop.info@ons.gsi.gov.uk",
//			},
//			ReleaseDate: "2017-01-18T00:00:00.000Z",
//			NextRelease: "To be announced",
//			Publications: []string{
//				"http://localhost:20099/datasets/AF002",
//				"http://localhost:20099/datasets/AF003",
//			},
//			TermsAndConditions: "N/A",
//		},
//		Dimensions: []*model.Dimension{{
//			ID:   "D000125",
//			Name: "Sex",
//			Options: []*model.DimensionOption{
//				{ID: "DO000153", Name: "Male"},
//				{ID: "DO000154", Name: "Female"},
//			},
//		}, {
//			ID:   "D000124",
//			Name: "Residence Type",
//			Options: []*model.DimensionOption{
//				{ID: "DO000161", Name: "All categories: Residence Type"},
//				{ID: "DO000162", Name: "Lives in a household"},
//				{ID: "DO000163", Name: "Lives in a communal establishment"},
//			},
//		}, {
//			ID:   "D000123",
//			Name: "Age",
//			Options: []*model.DimensionOption{
//				{ID: "DO000265", Name: "All categories: Age 16 and over"},
//				{ID: "DO000266", Name: "Age 16 to 24"},
//				{ID: "DO000267", Name: "Age 25 to 34"},
//				{ID: "DO000268", Name: "Age 35 to 49"},
//				{ID: "DO000269", Name: "Age 50 and over"},
//			},
//		}, {
//			ID:   "D000126",
//			Name: "2011 Statistical Geography Hierarchy",
//			Type: "geography",
//			Options: []*model.DimensionOption{
//				{ID: "K04000001", Name: "England and Wales"},
//				{ID: "E92000001", Name: "England"},
//				{ID: "W92000004", Name: "Wales"},
//			},
//		}},
//	}
//
//	cpiDataset := &model.Dataset{
//		ID:    "CPI15",
//		Title: "CPI15 Consumer Prices Index (COICOP).",
//		URL:   "http://localhost:20099/datasets/CPI15",
//		Metadata: &model.Metadata{
//			Description:        "Consumer Price Index statistics by Time and Special Aggregate (type of economic activity). This dataset shows the movement of prices over the last five years within the UK economy, broken down by month and various classifications of economic activity. The economic classifications are derived from the COICOP (Classification Of Individual Consumption by Purpose) list.",
//			NationalStatistics: true,
//			Contact: &model.Contact{
//				Name:  "James Tucker",
//				Email: "cpi@ons.gsi.gov.uk",
//				Phone: "+44 (0)1633 456900",
//			},
//			ReleaseDate: "2017-02-19T00:00:00.000Z",
//			NextRelease: "14 February 2017",
//			Methodology: []*model.Methodology{
//				{Title: "Consumer Price Inflation (includes all 4 indices—CPI, CPIH, RPI and RPIJ)", URL: "https://www.ons.gov.uk/economy/inflationandpriceindices/qmis/consumerpriceinflationqmi"},
//			},
//			TermsAndConditions: "",
//		},
//		Dimensions: []*model.Dimension{{
//			ID:   "SP00001",
//			Name: "Special aggregate",
//			Type: "classification",
//			Options: []*model.DimensionOption{{
//				ID:   "FOOD0001",
//				Name: "(01) Food and non-alcoholic beverages",
//				Options: []*model.DimensionOption{
//					{ID: "OPT02", Name: "Food"},
//					{ID: "OPT03", Name: "Bread and cereals"},
//					{ID: "OPT04", Name: "Meat"},
//					{ID: "OPT05", Name: "Fish"},
//					{ID: "OPT06", Name: "Milk, cheese and eggs"},
//					{ID: "OPT07", Name: "Oils and fats"},
//					{ID: "OPT08", Name: "Fruit"},
//				},
//			}, {
//				ID:   "HEALTH0002",
//				Name: "(02) Health",
//				Options: []*model.DimensionOption{
//					{ID: "OPT02", Name: "Medical products, appliances and equipment"},
//					{ID: "OPT03", Name: "Pharmaceutical products"},
//					{ID: "OPT04", Name: "Other medical and therapeutic equipment"},
//					{ID: "OPT05", Name: "Out-patient services"},
//					{ID: "OPT06", Name: "Medical services and paramedical services"},
//					{ID: "OPT07", Name: "Dental services"},
//					{ID: "OPT08", Name: "In-patient service"},
//					{ID: "OPT09", Name: "Medical and paramedic services"},
//				},
//			}},
//		}, {
//			ID:   "T000111",
//			Name: "Time",
//			Type: "time",
//			Options: []*model.DimensionOption{
//				{ID: "OPT02", Name: "February 2013"},
//				{ID: "OPT02", Name: "March 2013"},
//				{ID: "OPT02", Name: "April 2013"},
//				{ID: "OPT02", Name: "May 2013"},
//				{ID: "OPT02", Name: "June 2013"},
//				{ID: "OPT02", Name: "July 2013"},
//				{ID: "OPT02", Name: "August 2013"},
//				{ID: "OPT02", Name: "September 2013"},
//				{ID: "OPT02", Name: "October 2013"},
//				{ID: "OPT02", Name: "November 2013"},
//				{ID: "OPT02", Name: "December 2013"},
//				{ID: "OPT02", Name: "January 2014"},
//				{ID: "OPT02", Name: "February 2014"},
//				{ID: "OPT02", Name: "March 2014"},
//				{ID: "OPT02", Name: "April 2014"},
//				{ID: "OPT02", Name: "May 2014"},
//				{ID: "OPT02", Name: "June 2014"},
//				{ID: "OPT02", Name: "July 2014"},
//				{ID: "OPT02", Name: "August 2014"},
//				{ID: "OPT02", Name: "September 2014"},
//				{ID: "OPT02", Name: "October 2014"},
//				{ID: "OPT02", Name: "November 2014"},
//				{ID: "OPT02", Name: "December 2014"},
//				{ID: "OPT02", Name: "January 2015"},
//				{ID: "OPT02", Name: "February 2015"},
//				{ID: "OPT02", Name: "March 2015"},
//				{ID: "OPT02", Name: "April 2015"},
//				{ID: "OPT02", Name: "May 2015"},
//				{ID: "OPT02", Name: "June 2015"},
//				{ID: "OPT02", Name: "July 2015"},
//				{ID: "OPT02", Name: "August 2015"},
//				{ID: "OPT02", Name: "September 2015"},
//				{ID: "OPT02", Name: "October 2015"},
//				{ID: "OPT02", Name: "November 2015"},
//				{ID: "OPT02", Name: "December 2015"},
//				{ID: "OPT02", Name: "January 2016"},
//				{ID: "OPT02", Name: "February 2016"},
//				{ID: "OPT02", Name: "March 2016"},
//				{ID: "OPT02", Name: "April 2016"},
//				{ID: "OPT02", Name: "May 2016"},
//				{ID: "OPT02", Name: "June 2016"},
//				{ID: "OPT02", Name: "July 2016"},
//				{ID: "OPT02", Name: "August 2016"},
//				{ID: "OPT02", Name: "September 2016"},
//				{ID: "OPT02", Name: "October 2016"},
//			},
//		}},
//	}
//
//	err := client.Index(&model.Document{
//		ID:   "15a8c249-802d-4de4-9a73-7220e4c3b0e0",
//		Type: "dataset",
//		Body: armedForcesDataset,
//	})
//	if err != nil {
//		log.Error(err, nil)
//	}
//
//	err = client.Index(&model.Document{
//		ID:   "12999277-47c4-4c77-a641-81107d9c3201",
//		Type: "dataset",
//		Body: cpiDataset,
//	})
//	if err != nil {
//		log.Error(err, nil)
//	}
//
//}
