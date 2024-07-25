package uri_util

import (
	"net/url"
	"strings"
)

func UriToUrlValues(uri string) *url.Values {
	query := &url.Values{}

	parse, _ := url.Parse(uri)
	if parse != nil {
		queryArr := strings.Split(parse.RawQuery, "&")
		for _, v := range queryArr {
			vArr := strings.Split(v, "=")
			if len(vArr) == 2 {
				unescape, err := url.QueryUnescape(vArr[1])
				if err != nil {
					query.Set(vArr[0], vArr[1])
				} else {
					query.Set(vArr[0], unescape)
				}
			}
		}
	}
	return query
}
