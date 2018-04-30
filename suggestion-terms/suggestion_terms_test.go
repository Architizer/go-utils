package suggestionterms

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	"github.com/rtt/Go-Solr"
)

func TestNewSuggestionTerm(t *testing.T) {
	type TestCase struct {
		Got                 *SuggestionTerm
		ExpectedValue       string
		ExpectedWeightField string
	}
	testCases := []TestCase{
		TestCase{
			Got: NewSuggestionTerm(
				"pseudotiers_ss",
				solr.FacetCount{Value: "Furniture - Contract", Count: 300},
				"brand_count_i",
			),
			ExpectedValue:       "furniture - contract",
			ExpectedWeightField: "brand_count_i",
		},
		TestCase{
			Got: NewSuggestionTerm(
				"product_type_ent_ss",
				solr.FacetCount{Value: "emergency lights", Count: 300},
				"project_count_i",
			),
			ExpectedValue:       "emergency lights",
			ExpectedWeightField: "project_count_i",
		},
	}
	for _, tc := range testCases {
		if tc.Got.Value != tc.ExpectedValue {
			t.Error(
				"FOR:\n", tc, "\n",
				"EXPECTED:\n", tc.ExpectedValue, "\n",
				"GOT:", tc.Got.Value, "\n",
			)
		}
		if tc.Got.WeightField != tc.ExpectedWeightField {
			t.Error(
				"FOR:\n", tc, "\n",
				"EXPECTED:\n", tc.ExpectedWeightField, "\n",
				"GOT:", tc.Got.WeightField, "\n",
			)
		}
	}

}

func TestDocumentID(t *testing.T) {
	type TestCase struct {
		Got      SuggestionTerm
		Expected string
	}

	testCases := []TestCase{
		TestCase{
			Got: *NewSuggestionTerm(
				"pseudotiers_ss",
				solr.FacetCount{Value: "Furniture - Contract", Count: 300},
				"count_i",
			),
			Expected: "pseudotiers_ss.furniture - contract",
		},
		TestCase{
			Got: *NewSuggestionTerm(
				"product_type_ent_ss",
				solr.FacetCount{Value: "emergency lights", Count: 300},
				"count_i",
			),
			Expected: "product_type_ent_ss.emergency lights",
		},
	}

	for _, tc := range testCases {
		got := tc.Got.DocumentID()

		if got != tc.Expected {
			t.Error(
				"For", tc.Got,
				"expected", tc.Expected,
				"got", got,
			)
		}
	}
}

func TestNewDocument(t *testing.T) {
	type TestCase struct {
		Got      SuggestionTerm
		Expected *solr.Document
	}

	testCases := []TestCase{
		TestCase{
			Got: SuggestionTerm{"pseudotiers_ss", "lighting", 100, ""},
			Expected: &solr.Document{
				Fields: map[string]interface{}{
					"id":      "pseudotiers_ss.lighting",
					"term_s":  "lighting",
					"count_i": 100,
				},
			},
		},
		TestCase{
			Got: SuggestionTerm{"product_type_ent_ss", "wood doors", 100, ""},
			Expected: &solr.Document{
				Fields: map[string]interface{}{
					"id":     "product_type_ent_ss.wood doors",
					"term_s": "wood doors",
				},
			},
		},
	}

	for _, tc := range testCases {
		got := tc.Got.NewDocument()
		if got.Fields["id"] != tc.Expected.Fields["id"] {
			t.Error(
				"For", tc.Got,
				"expected", tc.Expected,
				"got", got,
			)
		}
	}
}

func TestAddSuggestionTerms(t *testing.T) {
	type Arguments struct {
		Facet       solr.Facet
		WeightField string
	}
	type TestCase struct {
		Collection *SuggestionTermCollection
		Arguments  *Arguments
		Expected   int
	}

	testCases := []TestCase{
		TestCase{
			Collection: new(SuggestionTermCollection),
			Arguments: &Arguments{
				Facet: solr.Facet{
					Name: "pseudotiers_ss",
					Counts: []solr.FacetCount{
						solr.FacetCount{Value: "Lighting", Count: 100},
						solr.FacetCount{Value: "Stone", Count: 100},
					},
				},
				WeightField: "brand_count_i",
			},
			Expected: 2,
		},
		TestCase{
			Collection: new(SuggestionTermCollection),
			Arguments: &Arguments{
				Facet: solr.Facet{
					Name: "product_type_ent_ss",
					Counts: []solr.FacetCount{
						solr.FacetCount{Value: "wood ceiling", Count: 100},
						solr.FacetCount{Value: "folding doors", Count: 100},
						solr.FacetCount{Value: "window trim", Count: 100},
					},
				},
				WeightField: "project_count_i",
			},
			Expected: 3,
		},
	}

	for _, tc := range testCases {
		tc.Collection.AddSuggestionTerms(tc.Arguments.Facet, tc.Arguments.WeightField)
		got := len(tc.Collection.SuggestionTerms)
		if got != tc.Expected {
			t.Error(
				"For", tc.Arguments,
				"expected", tc.Expected,
				"got", got,
			)
		}
	}
}

func TestToUpdateDocuments(t *testing.T) {
	type TestCase struct {
		Collection    *SuggestionTermCollection
		ExpectedBytes []byte
	}

	expectedBytes, err := json.Marshal(map[string]interface{}{
		"add": []interface{}{
			map[string]interface{}{
				"id":            "pseudotiers_ss.lighting",
				"term_s":        map[string]string{"set": "lighting"},
				"brand_count_i": map[string]int{"set": 100},
			},
			map[string]interface{}{
				"id":            "pseudotiers_ss.plumbing",
				"term_s":        map[string]string{"set": "plumbing"},
				"brand_count_i": map[string]int{"set": 100},
			},
		},
	})

	if err != nil {
		fmt.Println(err)
		return
	}

	stCollection := new(SuggestionTermCollection)
	_ = stCollection.AddSuggestionTerms(
		solr.Facet{
			Name: "pseudotiers_ss",
			Counts: []solr.FacetCount{
				solr.FacetCount{
					Value: "lighting",
					Count: 100,
				},
				solr.FacetCount{
					Value: "plumbing",
					Count: 100,
				},
			},
		},
		"brand_count_i",
	)

	testCases := []TestCase{
		TestCase{
			Collection:    stCollection,
			ExpectedBytes: expectedBytes,
		},
	}

	for _, tc := range testCases {
		payload := tc.Collection.ToUpdateDocument()
		gotBytes, err := solr.JSONToBytes(payload)
		if err != nil {
			fmt.Println(err)
			return
		}
		got := string(*gotBytes)
		expected := string(tc.ExpectedBytes)
		if bytes.Compare(tc.ExpectedBytes, *gotBytes) != 0 {
			t.Error(
				"EXPECTED:\n", expected, "\n",
				"GOT:\n", got, "\n",
			)
		}
	}
}

func TesNewFacetQuery(t *testing.T) {
	query := NewFacetQuery([]string{"django_ct:brands.brand"}, []string{"pseudotiers_ss"})

	type TestCase struct {
		Got      string
		Expected string
	}

	testCases := []TestCase{
		TestCase{
			Got:      query.Params["q"][0],
			Expected: "*:*",
		},
		TestCase{
			Got:      query.Params["fq"][0],
			Expected: "django_ct:brands.brand",
		},
		TestCase{
			Got:      query.Params["facet.field"][0],
			Expected: "pseudotiers_ss",
		},
		TestCase{
			Got:      query.Params["facet"][0],
			Expected: "true",
		},
	}

	for _, tc := range testCases {
		if strings.Compare(tc.Got, tc.Expected) != 0 {
			t.Error(
				"EXPECTED:", tc.Expected,
				"GOT:", tc.Got,
			)
		}
	}
}
