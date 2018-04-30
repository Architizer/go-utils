Go package that builds Solr documents from facet queries for a collection devoted to housing terms to build Suggesters from.

## Installation
```
go install github.com/Architizer/go-utils/suggestion-terms/bin/update_suggestion_terms
```

## Usage
```
  -facetfield string
        Field to facet on (default "pseudotiers_ss")
  -fq string
        'fq' param for facet query. (default "django_ct:brands.brand")
  -host string
        solr host (default "localhost")
  -port int
        solr port (default 8983)
  -source string
        Source colllection to facet terms from. (default "product_source")
  -target string
        Target collection to update terms on. (default "suggestion_terms")
  -weightfield string
        Name of weight field. Must have '_i'. (default "count_i")
```