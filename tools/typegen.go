package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"sort"

	"gopkg.in/yaml.v2"

	"github.com/newrelic/newrelic-client-go/newrelic"
	"github.com/newrelic/newrelic-client-go/pkg/nerdgraph"

	log "github.com/sirupsen/logrus"
)

// Config is the information keeper for generating go structs from type names.
type Config struct {
	PackageTypeMap map[string][]string `yaml:"package_type_map"`
}

func main() {

	verbose := flag.Bool("v", false, "increase verbosity")

	flag.Parse()

	if *verbose {
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetLevel(log.InfoLevel)
	}

	apiKey := os.Getenv("NEW_RELIC_API_KEY")
	nr, err := newrelic.New(newrelic.ConfigPersonalAPIKey(apiKey))
	if err != nil {
		log.Fatal(err)
	}

	schema, err := nr.NerdGraph.QuerySchema()
	if err != nil {
		log.Fatal(err)
	}

	yamlFile, err := ioutil.ReadFile("configs/typegen.yaml")
	if err != nil {
		log.Fatal(err)
	}

	var config Config

	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		log.Fatal(err)
	}

	for packageName, typeNames := range config.PackageTypeMap {

		types, err := nerdgraph.ResolveSchemaTypes(*schema, typeNames)
		if err != nil {
			log.Error(err)
		}

		f, err := os.Create(fmt.Sprintf("pkg/%s/types.go", packageName))
		if err != nil {
			log.Error(err)
		}

		_, err = f.WriteString(fmt.Sprintf("package %s\n", packageName))
		if err != nil {
			log.Error(err)
		}

		defer f.Close()

		keys := make([]string, 0, len(types))
		for k := range types {
			keys = append(keys, k)
		}
		sort.Strings(keys)

		for _, k := range keys {
			_, err := f.WriteString(types[k])
			if err != nil {
				log.Error(err)
			}
		}

	}

}
