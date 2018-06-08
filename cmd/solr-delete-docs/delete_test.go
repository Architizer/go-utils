package main

import (
	"math/rand"
	"testing"
	"time"

	"github.com/icrowley/fake"
)

func TestMakeSolrParamsList(t *testing.T) {
	lines := make([]string, 50)
	seed := rand.NewSource(time.Now().UnixNano())
	for i := range lines {
		r := rand.New(seed)
		chars := fake.CharactersN(r.Intn(100))
		lines[i] = chars
	}
	paramsList := MakeSolrParamsList(lines)

	var expected []string
	var expectedLength int
	var got []string
	var gotLength int
	for _, params := range paramsList {
		expected = []string{"*:*"}
		got = params["q"]
		if len(got) != len(expected) {
			t.Error(
				"Expected: ", expected,
				"Got: ", got,
			)
		}
		for _, expectedVal := range expected {
			for _, gotVal := range got {
				if expectedVal != gotVal {
					t.Error(
						"Expected: ", expectedVal,
						"Got: ", gotVal,
					)
				}
			}
		}

		expectedLength = 1
		gotLength = len(params["fq"])
		if expectedLength != gotLength {
			t.Error(
				"Expected: ", expectedLength,
				"Got: ", gotLength,
			)
		}

		expectedLength = 500
		gotLength = len(params["fq"][0])

		if expectedLength >= gotLength {
			t.Error(
				"Expected: ", expectedLength,
				"Got: ", gotLength,
			)
		}
	}
}
