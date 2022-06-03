package repository

import (
	"context"
	"fmt"
	"github.com/olivere/elastic"
	"log"
)

var (
	cfg = "https://elastic:v9hKmmYyQg3dtht2lz759LGM@31b6486899c74ad095db11db761e5e7c.eu-west-1.aws.found.io:9243/"
)

func (m *mongoInteractionRespository) IndexDocument(document interface{}, indexName string, id string) {
	// Delete Index
	m.DeleteIndex(indexName, id)
	// Create Index
	m.CreateIndex(indexName)
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

	fmt.Println(res)

}

func (m *mongoInteractionRespository) CreateIndex(name string) {

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
		}
		if !createIndex.Acknowledged {
			// Not acknowledged
		}
	}

}

func (m *mongoInteractionRespository) DeleteIndex(name string, ID string) {

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

	ctx := context.Background()
	res, err := client.Delete().
		Index(name).
		Type("doc").
		Id(ID).
		Do(ctx)
	if err != nil {
		// Handle error
		log.Println(err)
	}

	fmt.Println(res)

}
