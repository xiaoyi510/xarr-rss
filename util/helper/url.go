package helper

import "net/url"

func ParseUrl(uri string) (error, url.Values) {
	parse, err := url.Parse(uri)
	if err != nil {
		return err, nil
	}
	query, err := url.ParseQuery(parse.RawQuery)
	if err != nil {
		return err, nil
	}
	return nil, query
}
