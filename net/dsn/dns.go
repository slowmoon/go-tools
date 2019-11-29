package dsn

import (
	"net/url"
	"strings"
)

//data source name
type DSN struct {
	*url.URL
}

func (d *DSN)Binding(v interface{}) (url.Values, error) {
	assigns := make(map[string]assignFunc)
	if d.User != nil {
		username := d.User.Username()
		password , ok := d.User.Password()
		if ok {
			assigns["password"] = stringAssignFunc(password)
		}
		assigns["username"] = stringAssignFunc(username)
	}
	assigns["network"] = stringAssignFunc(d.Scheme)
	assigns["address"] = addressesAssignFunc(d.Address())

}


func (d *DSN)Address()[]string  {
	switch d.Scheme {
	case "unix", "unixgram", "unixpacket":
		return []string{d.Path}
	default:
		return strings.Split(d.Host, ",")
	}
}

//parse dsn to url
func Parse(dnsUrl string) (*DSN, error) {
	durl, err := url.Parse(dnsUrl)
	return &DSN{durl}, err
}



