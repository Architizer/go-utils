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

// MakeUpdateDocument creates a payload for a delete operation from a list of solr Documents.
func MakeUpdateDocument(conn *solr.Connection, docs []solr.Document) map[string]interface{} {
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
	return updateDoc
}

func deleteDocuments(params solr.URLParamMap, conn *solr.Connection, done chan bool) {
	q := solr.Query{
		Params: params,
	}

	// Select documents for deletion
	res, err := conn.Select(&q)
	if err != nil {
		fmt.Println("Solr select error")
		fmt.Printf("q: %v\n", q)
		fmt.Println(err)
		log.Fatal(err)
	}

	// Delete first page of documents
	updateDoc := MakeUpdateDocument(conn, res.Results.Collection)
	_, err = conn.Update(updateDoc, true)
	if err != nil {
		fmt.Println("Solr update err")
		fmt.Printf("q: %v\n", q)
		log.Fatal(err)
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
			fmt.Println(err)
			log.Fatal(err)
		}

		updateDoc := MakeUpdateDocument(conn, res.Results.Collection)
		_, err = conn.Update(updateDoc, true)
		if err != nil {
			fmt.Println("Solr update err")
			fmt.Printf("q: %v\n", q)
			log.Fatal(err)
		}
		fmt.Printf("Deleted %v documents\n", len(res.Results.Collection))
	}

	done <- true
}

// Convert lines of text into a list of Solr query params to limit request URI length.
func MakeSolrParamsList(lines []string) []solr.URLParamMap {
	maxLength := 500
	solrParamsList := make([]solr.URLParamMap, 0)
	var builder strings.Builder
	for _, line := range lines {
		val := fmt.Sprintf("\"%s\" ", line)
		if builder.Len() == 0 {
			builder.WriteString(fmt.Sprintf("id:"))
		}
		if builder.Len() < maxLength {
			builder.WriteString(val)
		} else {
			solrParams := solr.URLParamMap{
				"q":  []string{"*:*"},
				"fq": []string{builder.String()},
				"fl": []string{"id"},
			}
			solrParamsList = append(solrParamsList, solrParams)
			builder = strings.Builder{}
		}
	}
	return solrParamsList
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

	// Create list of Solr params
	paramsList := MakeSolrParamsList(lines)

	// Init a connection
	conn := solr.Connection{URL: collectionURL}

	// Delete documents that match params in solrParamsList
	done := make(chan bool, 1)
	for _, params := range paramsList {
		go deleteDocuments(params, &conn, done)
	}

	<-done

}
