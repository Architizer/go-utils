package suggestionterms

import (
	"fmt"
	"strings"

	"github.com/rtt/Go-Solr"
)

// DefaultSuggestField is the default value to store facet values on.
const DefaultSuggestField = "term_t"

// DefaultWeightField is the default value to store facet count on.
const DefaultWeightField = "count_i"

// SuggestionTerm represents a suggestion_terms document.
type SuggestionTerm struct {
	Field        string
	Value        string
	Count        int
	WeightField  string
	SuggestField string
}

// SuggestionTermCollection represents a collection of suggestion_terms documents
type SuggestionTermCollection struct {
	SuggestionTerms []SuggestionTerm
}

// NewSuggestionTerm creates a SuggestionTerm from a facet name and FacentCount object
func NewSuggestionTerm(facetName string, facetCount solr.FacetCount, weightField string, suggestField string) *SuggestionTerm {
	st := new(SuggestionTerm)
	st.Field = facetName
	st.Value = strings.ToLower(facetCount.Value)
	st.Count = facetCount.Count

	// Set default SuggestField
	if suggestField == "" {
		suggestField = DefaultSuggestField
	}

	// Set default WeightField
	if weightField == "" {
		weightField = DefaultWeightField
	}
	st.WeightField = weightField
	st.SuggestField = suggestField
	return st
}

// DocumentID returns a Solr document id for a SuggestionTerm
func (st *SuggestionTerm) DocumentID() string {
	return fmt.Sprint(st.Field, ".", st.Value)
}

// NewDocument converts a SuggestionTerm into a solr.Document
func (st *SuggestionTerm) NewDocument() *solr.Document {
	return &solr.Document{Fields: map[string]interface{}{
		"id":            st.DocumentID(),
		st.SuggestField: map[string]string{"set": st.Value},
		"term_type_s":   map[string]string{"set": st.Field},
		st.WeightField:  map[string]int{"set": st.Count},
	}}
}

// AddSuggestionTerms adds SuggestionTerms to a SuggestionTermCollection
func (collection *SuggestionTermCollection) AddSuggestionTerms(facet solr.Facet, weightField string, suggestField string) {
	var facetCount solr.FacetCount
	for i := 0; i < len(facet.Counts); i++ {
		facetCount = facet.Counts[i]
		st := NewSuggestionTerm(facet.Name, facetCount, weightField, suggestField)
		collection.SuggestionTerms = append(collection.SuggestionTerms, *st)
	}
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

// NewFacetQuery constructs a Solr facet query
func NewFacetQuery(fq []string, facetField []string) *solr.Query {
	return &solr.Query{
		Params: solr.URLParamMap{
			"q":           []string{"*:*"},
			"fq":          fq,
			"facet.field": facetField,
			"facet":       []string{"true"},
			"facet.limit": []string{"-1"},
			"facet.mincount": []string{"1"},
		},
	}
}
