package config_test

import (
	"github.com/ONSdigital/dp-dd-search-indexer/config"
	. "github.com/smartystreets/goconvey/convey"
	"os"
	"reflect"
	"strings"
	"testing"
)

func Test(t *testing.T) {

	Convey("Given some preset environment variables", t, func() {

		bindAddress := ":theport"
		elasticSearchNodes := []string{"elasticNode1", "elasticNode2"}
		elasticSearchIndex := "search-index"

		_ = os.Setenv("BIND_ADDR", bindAddress)
		_ = os.Setenv("ELASTIC_SEARCH_NODES", strings.Join(elasticSearchNodes, ","))
		_ = os.Setenv("ELASTIC_SEARCH_INDEX", elasticSearchIndex)

		Convey("When the config is loaded", func() {

			config.Load()

			Convey("Then the expected environment variable values are loaded into the config", func() {
				So(config.BindAddr, ShouldEqual, bindAddress)
				So(reflect.DeepEqual(config.ElasticSearchNodes, elasticSearchNodes), ShouldBeTrue)
				So(config.ElasticSearchIndex, ShouldEqual, elasticSearchIndex)
			})
		})
	})
}
