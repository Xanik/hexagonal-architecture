package repository

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/olivere/elastic"
	"log"
)

var (
	cfg = "https://elastic:v9hKmmYyQg3dtht2lz759LGM@31b6486899c74ad095db11db761e5e7c.eu-west-1.aws.found.io:9243/"
)

func (m *mongoRespository) Elastic(id string) (interface{}, error) {
	var account map[string]string

	// Initialize a client with the default settings.
	//
	// An `ELASTICSEARCH_URL` environment variable will be used when exported.
	//
	client, err := elastic.NewClient(elastic.SetURL(cfg))
	if err != nil {
		log.Println("Error creating the client: %s", err)
	}

	// Getting the ES version number is quite common, so there's a shortcut
	esversion, err := client.ElasticsearchVersion(cfg)
	if err != nil {
		// Handle error
		fmt.Println(err)
	}
	fmt.Printf("Elasticsearch version %s\n", esversion)

	return account, nil
}

func (m *mongoRespository) IndexDocument(document interface{}, indexName string, id string) (string, error) {
	// Initialize a client with the default settings.
	//
	// An `ELASTICSEARCH_URL` environment variable will be used when exported.
	//
	client, err := elastic.NewSimpleClient(elastic.SetURL(cfg))

	if err != nil {
		log.Println("Error creating the client: %s", err)
	}

	// Getting the ES version number is quite common, so there's a shortcut
	esversion, err := client.ElasticsearchVersion(cfg)
	if err != nil {
		// Handle error
		log.Println(err)
	}
	fmt.Printf("Elasticsearch version %s\n", esversion)

	// Index a struct (using JSON serialization)
	doc := document
	put1, err := client.Index().
		Index(indexName).
		Id(id).
		BodyJson(doc).
		Type("doc").
		Do(context.Background())
	if err != nil {
		// Handle error
		log.Println(err)
	}
	res := fmt.Sprintf("Indexed document %s to index %s, type %s\n", put1.Id, put1.Index, put1.Type)

	return res, nil
}

func (m *mongoRespository) CreateIndex(name string) (string, error) {

	// Initialize a client with the default settings.
	//
	// An `ELASTICSEARCH_URL` environment variable will be used when exported.
	//
	client, err := elastic.NewSimpleClient(elastic.SetURL(cfg))

	if err != nil {
		log.Println("Error creating the client: %s", err)
	}

	// Getting the ES version number is quite common, so there's a shortcut
	esversion, err := client.ElasticsearchVersion(cfg)
	if err != nil {
		// Handle error
		log.Println(err)
	}
	fmt.Printf("Elasticsearch version %s\n", esversion)

	// Use the IndexExists service to check if a specified index exists.
	exists, err := client.IndexExists(name).Do(context.Background())
	if err != nil {
		// Handle error
		log.Println(err)
	}
	if !exists {
		// Create a new index.
		mapping := `
{
	"settings":{
		"number_of_shards":1,
		"number_of_replicas":0
	},
	"mappings":{
		"doc":{
			"properties":{
				"user":{
					"type":"keyword"
				},
				"message":{
					"type":"text",
					"store": true,
					"fielddata": true
				},
            "retweets":{
                "type":"long"
            },
				"tags":{
					"type":"keyword"
				},
				"location":{
					"type":"geo_point"
				},
				"suggest_field":{
					"type":"completion"
				}
			}
		}
	}
}
`
		createIndex, err := client.CreateIndex(name).Body(mapping).Do(context.Background())
		if err != nil {
			// Handle error
			log.Println(err)
			return "", err
		}
		if !createIndex.Acknowledged {
			// Not acknowledged
		}
		return "Index Created", nil
	}
	return "Index Exist", nil

}

func (m *mongoRespository) SearchContent(indexName string, key string, value string) ([]map[string]interface{}, error) {
	// Initialize a client with the default settings.
	//
	// An `ELASTICSEARCH_URL` environment variable will be used when exported.
	//
	client, err := elastic.NewSimpleClient(elastic.SetURL(cfg))

	if err != nil {
		log.Println("Error creating the client: %s", err)
		return nil, err
	}

	// Getting the ES version number is quite common, so there's a shortcut
	esversion, err := client.ElasticsearchVersion(cfg)
	if err != nil {
		// Handle error
		log.Println(err)
		return nil, err
	}
	fmt.Printf("Elasticsearch version %s\n", esversion)

	// Refresh to make sure the documents are searchable.
	_, err = client.Refresh().Index(indexName).Do(context.Background())
	if err != nil {
		log.Println("Error Occured", err)
		return nil, err
	}

	// Search with a term query
	termQuery := elastic.NewWildcardQuery(key, "*"+value+"*")
	searchResult, err := client.Search().
		Index(indexName). // search in index
		Query(termQuery). // specify the query
		From(0).Size(10). // take documents 0-9
		MinScore(0).
		Pretty(true).            // pretty print request and response JSON
		Do(context.Background()) // execute
	if err != nil {
		// Handle error
		log.Println("Error Unhandled", err)
		return nil, err
	}

	// searchResult is of type SearchResult and returns hits, suggestions,
	// and all kinds of other information from Elasticsearch.
	fmt.Printf("Query took %d milliseconds\n", searchResult.TookInMillis)

	// Each is a convenience function that iterates over hits in a search result.
	// It makes sure you don't need to check for nil values in the response.
	// However, it ignores errors in serialization. If you want full control
	// over iterating the hits, see below.
	// TotalHits is another convenience function that works even when something goes wrong.
	fmt.Printf("Found a total of %d docs\n", searchResult.TotalHits())

	// Here's how you iterate through results with full control over each step.
	if searchResult.TotalHits() > 0 {
		fmt.Printf("Found a total of %d docs\n", searchResult.TotalHits())

		var data []map[string]interface{}
		// Iterate through results
		for _, hit := range searchResult.Hits.Hits {
			// hit.Index contains the name of the index

			// Deserialize hit.Source into a struct (could also be just a map[string]interface{}).
			var t map[string]interface{}
			err := json.Unmarshal(*hit.Source, &t)
			if err != nil {
				// Deserialization failed
				return nil, errors.New("Unable To Resolve Search")
			}

			// Work with tweet
			fmt.Printf("Docs by %s: \n", t)
			data = append(data, t)
		}
		return data, nil

	} else {
		// No hits
		fmt.Print("Found no data\n")
		return nil, errors.New("No Data Was Found")
	}
}
