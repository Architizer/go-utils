Delete Solr documents from a collection given a field name and a text file containing field values.

## Installation
```
go install -i github.com/Architizer/go-utils/cmd/solr-delete-docs
```

## Example
```
solr-delete-docs -file id-list.txt -collection products -url http://localhost:8983
```
Deletes Solr documents from collection on localhost whose id field matches values in `id-list.txt`.