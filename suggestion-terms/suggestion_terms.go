package suggestionterms

import (
	"fmt"
	"strings"

	"github.com/rtt/Go-Solr"
)

const defaultWeightField = "count_i"

// SuggestionTerm represents a suggestion_terms document.
type SuggestionTerm struct {
	Field       string
	Value       string
	Count       int
	WeightField string
}

// SuggestionTermCollection represents a collection of suggestion_terms documents
type SuggestionTermCollection struct {
	SuggestionTerms []SuggestionTerm
	WeightField     string
}

// NewSuggestionTerm creates a SuggestionTerm from a facet name and FacentCount object
func NewSuggestionTerm(facetName string, facetCount solr.FacetCount, weightField string) *SuggestionTerm {
	st := new(SuggestionTerm)
	st.Field = facetName
	st.Value = strings.ToLower(facetCount.Value)
	st.Count = facetCount.Count
	st.WeightField = weightField
	return st
}

// DocumentID returns a Solr document id for a SuggestionTerm
func (st *SuggestionTerm) DocumentID() string {
	return fmt.Sprint(st.Field, ".", st.Value)
}

// NewDocument converts a SuggestionTerm into a solr.Document
func (st *SuggestionTerm) NewDocument() *solr.Document {
	return &solr.Document{Fields: map[string]interface{}{
		"id":           st.DocumentID(),
		"term_s":       map[string]string{"set": st.Value},
		st.WeightField: map[string]int{"set": st.Count},
	}}
}

// AddSuggestionTerms adds SuggestionTerms to a SuggestionTermCollection
// facet: solr.Facet to create suggestion terms from
// weightField: Name of weight field on Solr document
func (collection *SuggestionTermCollection) AddSuggestionTerms(facet solr.Facet, weightField string) error {
	var facetCount solr.FacetCount
	for i := 0; i < len(facet.Counts); i++ {
		facetCount = facet.Counts[i]
		st := NewSuggestionTerm(facet.Name, facetCount, weightField)
		collection.SuggestionTerms = append(collection.SuggestionTerms, *st)
	}
	return nil
}

// ToUpdateDocument converts SuggestionTerms into a payload for solr.Connection.Update
func (collection *SuggestionTermCollection) ToUpdateDocument() map[string]interface{} {
	docs := make([]interface{}, len(collection.SuggestionTerms))

	for i, st := range collection.SuggestionTerms {
		docs[i] = st.NewDocument().Doc()
	}

	payload := map[string]interface{}{
		"add": docs,
	}

	return payload
}

// MakeFacetQuery constructs a Solr facet query
func MakeFacetQuery(fq []string, facetField []string) solr.Query {
	return solr.Query{
		Params: solr.URLParamMap{
			"q":           []string{"*:*"},
			"fq":          fq,
			"facet.field": facetField,
			"facet":       []string{"true"},
		},
	}
}
