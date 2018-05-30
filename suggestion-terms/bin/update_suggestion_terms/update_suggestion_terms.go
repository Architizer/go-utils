package main

import (
	"flag"
	"fmt"
	"log"
	"strings"

	"github.com/getsentry/raven-go"

	"github.com/Architizer/go-utils/suggestion-terms"
	"github.com/rtt/Go-Solr"
)

func executeCommand() {
	urlPtr := flag.String("url", "http://localhost:8983", "solr url")
	sourcePtr := flag.String("source", "product_source", "Source collection to facet terms from.")
	targetPtr := flag.String("target", "suggestion_terms", "Target collection to update terms on.")
	fqPtr := flag.String("fq", "django_ct:brands.brand", "'fq' param for facet query.")
	facetFieldPtr := flag.String("facetfield", "pseudotiers_ss", "Field to facet on")
	suggestFieldPtr := flag.String("suggestField", "term_t", "Field to store facet value on")
	weightFieldPtr := flag.String("weightfield", "count_i", "Name of weight field. Must have '_i'.")

	flag.Parse()

	var err error
	// Validate args
	if strings.HasSuffix(*weightFieldPtr, string("_i")) != true {
		err = fmt.Errorf("weightField must have suffix '_i'")
		raven.CaptureErrorAndWait(err, nil)
		log.Panic(err)
	}

	sourceURL := fmt.Sprintf("%s/solr/%s", *urlPtr, *sourcePtr)
	targetURL := fmt.Sprintf("%s/solr/%s", *urlPtr, *targetPtr)

	// Build query object
	q := suggestionterms.NewFacetQuery(
		[]string{*fqPtr},
		[]string{*facetFieldPtr},
	)

	// Init a connection
	conn := solr.Connection{URL: sourceURL}

	// Perform the query, checking for errors
	res, err := conn.Select(q)

	if err != nil {
		raven.CaptureErrorAndWait(err, nil)
		log.Panic(err)
	}

	// Convert facets to SuggestionTerms
	facets := res.Results.Facets
	collection := new(suggestionterms.SuggestionTermCollection)
	for _, facet := range facets {
		collection.AddSuggestionTerms(facet, *weightFieldPtr, *suggestFieldPtr)
	}

	// Init a connection
	conn = solr.Connection{URL: targetURL}

	// Build update document
	doc := collection.ToUpdateDocument()

	// Send off the update
	resp, err := conn.Update(doc, true)

	if err != nil {
		raven.CaptureErrorAndWait(err, nil)
		log.Panic(err)
	} else {
		fmt.Println("resp =>", resp)
	}
}

func main() {
	raven.CapturePanic(executeCommand, nil)
}
