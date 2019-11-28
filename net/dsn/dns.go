package dsn

import (
	"net/url"
)

//data source name
type DSN struct {
	*url.URL
}


//parse dsn to url
func Parse(dnsUrl string) (*DSN, error) {
	durl, err := url.Parse(dnsUrl)
	return &DSN{durl}, err
}



