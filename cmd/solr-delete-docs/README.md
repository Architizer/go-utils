Delete Solr documents from a collection given a field name and a text file containing field values.

## Installation
```
go install -i github.com/Architizer/go-utils/cmd/solr-delete-docs
```

## Usage
```
Usage of solr-delete-docs:
  -collection string
        Solr collection to delete documents from
  -field string
        Solr field to query against.
  -file string
        Location of file containing list of values to query against.
  -fq string
        fq param for Select action
  -url string
        Solr host url (default "http://localhost:8983")
```