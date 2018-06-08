// Delete documents from a Solr collection matching ids contained in a specified text file
package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"strings"

	solr "github.com/rtt/Go-Solr"
)

// Create a payload for a delete operation from a list of solr Documents.
func deleteDocuments(conn *solr.Connection, docs []solr.Document) error {
	updateDoc := make(map[string]interface{})
	deleteDocs := make([]interface{}, 0)
	for _, doc := range docs {
		switch documentID := doc.Fields["id"].(type) {
		case string:
			deleteDocs = append(
				deleteDocs,
				map[string]interface{}{"id": documentID},
			)
		}
	}
	updateDoc["delete"] = deleteDocs
	res, err := conn.Update(updateDoc, true)
	if res.Success == false {
		log.Fatal("Delete operation not successful")
	}

	return err
}

// Convert lines of text into a list of Solr query params to limit request URI length.
func makeQueryParamList(lines []string) []string {
	maxLength := 500
	queryParamList := make([]string, 0)
	var builder strings.Builder
	for _, line := range lines {
		val := fmt.Sprintf("\"%s\" ", line)
		if builder.Len() == 0 {
			builder.WriteString(fmt.Sprintf("id:"))
		}
		if builder.Len() < maxLength {
			builder.WriteString(val)
		} else {
			queryParamList = append(queryParamList, builder.String())
			builder = strings.Builder{}
		}
	}
	return queryParamList
}

func main() {
	filePtr := flag.String("file", "", "Location of file containing list of values to query against.")
	urlPtr := flag.String("url", "http://localhost:8983", "Solr host url")
	collectionPtr := flag.String("collection", "", "Solr collection to delete documents from")

	flag.Parse()
	var err error

	// Validate args
	if *filePtr == "" {
		err = fmt.Errorf("-file is required")
	}
	if *collectionPtr == "" {
		err = fmt.Errorf("-collection is required")
	}

	// Create URL to collection
	collectionURL := fmt.Sprintf("%s/solr/%s", *urlPtr, *collectionPtr)

	// Read file
	content, err := ioutil.ReadFile(*filePtr)

	if err != nil {
		log.Fatal(err)
	}
	s := string(content)
	lines := strings.Split(s, "\n")

	// Create query strings
	queryParamList := makeQueryParamList(lines)

	// Init a connection
	conn := solr.Connection{URL: collectionURL}

	deletedDocuments := make([]solr.Document, 0)
	for _, queryParam := range queryParamList {
		q := solr.Query{
			Params: solr.URLParamMap{
				"q":  []string{"*:*"},
				"fq": []string{queryParam},
				"fl": []string{"id"},
			},
		}

		// Select documents for deletion
		res, err := conn.Select(&q)
		if err != nil {
			fmt.Println("Solr select error")
			fmt.Printf("q: %v\n", q)
			log.Fatal(err)
		}

		// Delete first page of documents
		err = deleteDocuments(&conn, res.Results.Collection)

		if err != nil {
			fmt.Println("Solr update err")
			fmt.Printf("q: %v\n", q)
			log.Fatal(err)
		}
		// Update deletedDocuments
		for _, doc := range res.Results.Collection {
			deletedDocuments = append(deletedDocuments, doc)
		}
		fmt.Printf("Deleted %v documents\n", len(res.Results.Collection))

		// Delete rest of documents, if necessary
		for i := 10; i < res.Results.NumFound; i += 10 {
			q.Start = i

			// Select documents for deletion
			res, err := conn.Select(&q)
			if err != nil {
				fmt.Println("Solr select err")
				fmt.Printf("q: %v\n", q)
				log.Fatal(err)
			}

			err = deleteDocuments(&conn, res.Results.Collection)
			if err != nil {
				fmt.Println("Solr update err")
				fmt.Printf("q: %v\n", q)
				log.Fatal(err)
			}
			// Update deletedDocuments
			for _, doc := range res.Results.Collection {
				deletedDocuments = append(deletedDocuments, doc)
			}
			fmt.Printf("Deleted %v documents\n", len(res.Results.Collection))
		}
	}

	fmt.Printf("Deleted %v total documents\n", len(deletedDocuments))
}
