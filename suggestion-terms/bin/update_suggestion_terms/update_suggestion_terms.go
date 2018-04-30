package main

import (
	"flag"
	"fmt"
	"log"
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

	// Build query object
	q := suggestionterms.NewFacetQuery(
		[]string{*fqPtr},
		[]string{*facetFieldPtr},
	)

	// Init a connection
	conn, err := solr.Init(*hostPtr, *portPtr, *sourcePtr)

	if err != nil {
		fmt.Println(err)
		return
	}

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
	s, err := solr.Init(*hostPtr, *portPtr, *targetPtr)

	if err != nil {
		log.Fatal(err)
	}

	// Build update document
	doc := collection.ToUpdateDocument()

	// Send off the update
	resp, err := s.Update(doc, true)

	if err != nil {
		fmt.Println("error =>", err)
	} else {
		fmt.Println("resp =>", resp)
	}

}
