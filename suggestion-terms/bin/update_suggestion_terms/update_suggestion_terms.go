package main

import (
	"flag"
	"fmt"
	"strings"

	"github.com/Architizer/go-utils/suggestion-terms"
	"github.com/rtt/Go-Solr"
)

func main() {
	hostPtr := flag.String("host", "localhost", "solr host")
	portPtr := flag.Int("port", 8983, "solr port")
	sourcePtr := flag.String("source", "product_source", "Source colllection to facet terms from.")
	targetPtr := flag.String("target", "suggestion_terms", "Target collection to update terms on.")
	fqPtr := flag.String("fq", "django_ct:brands.brand", "'fq' param for facet query.")
	facetFieldPtr := flag.String("facetfield", "pseudotiers_ss", "Field to facet on")
	weightFieldPtr := flag.String("weightfield", "count_i", "Name of weight field. Must have '_i'.")

	flag.Parse()

	// Validate args
	if strings.HasSuffix(*weightFieldPtr, string("_i")) != true {
		fmt.Println("weightField must have suffix '_i'")
		return
	}
	if len(*hostPtr) == 0 {
		fmt.Println("Invalid hostname (must be length >= 1)")
	}
	if *portPtr <= 0 || *portPtr > 65535 {
		fmt.Println("Invalid port (must be 1..65535")
	}

	var sourceURL string
	var targetURL string
	if *portPtr == 80 {
		sourceURL = fmt.Sprintf("%s/solr/%s", *hostPtr, *sourcePtr)
		targetURL = fmt.Sprintf("%s/solr/%s", *hostPtr, *targetPtr)
	} else {
		sourceURL = fmt.Sprintf("%s:%d/solr/%s", *hostPtr, *portPtr, *sourcePtr)
		targetURL = fmt.Sprintf("%s:%d/solr/%s", *hostPtr, *portPtr, *targetPtr)
	}

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
		fmt.Println(err)
		return
	}

	// Convert facets to SuggestionTerms
	facets := res.Results.Facets
	collection := new(suggestionterms.SuggestionTermCollection)
	for _, facet := range facets {
		collection.AddSuggestionTerms(facet, *weightFieldPtr)
	}

	// Init a connection
	conn = solr.Connection{URL: targetURL}

	// Build update document
	doc := collection.ToUpdateDocument()

	// Send off the update
	resp, err := conn.Update(doc, true)

	if err != nil {
		fmt.Println("error =>", err)
	} else {
		fmt.Println("resp =>", resp)
	}

}
