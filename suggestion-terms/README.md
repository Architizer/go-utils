Go package that builds Solr documents from facet queries for a collection devoted to housing terms to build Suggesters from.

## Installation
```
go install github.com/Architizer/go-utils/suggestion-terms/bin/update_suggestion_terms
```

## Download dependencies
```
go get github.com/rtt/Go-Solr
```

## Usage
Command to update suggestion terms is `update_suggestion_terms`.

Arguments:
```
  -facetfield string
        Field to facet on (default "pseudotiers_ss")
  -fq string
        'fq' param for facet query. (default "django_ct:brands.brand")
  -source string
        Source colllection to facet terms from. (default "product_source")
  -target string
        Target collection to update terms on. (default "suggestion_terms")
  -url string
      solr url (default "http://localhost:8983")
  -weightfield string
        Name of weight field. Must have '_i'. (default "count_i")
```