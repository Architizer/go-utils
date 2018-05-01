FROM golang:1.10

WORKDIR /go/src/github.com/Architizer/go-utils/suggestion-terms
COPY . .

RUN go get -v github.com/rtt/Go-Solr github.com/getsentry/raven-go
RUN go install -v github.com/Architizer/go-utils/suggestion-terms/bin/update_suggestion_terms

CMD ["update_suggestion_terms"]