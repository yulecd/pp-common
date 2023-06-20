package client

import (
	"fmt"
	urlpkg "net/url"
)

func IsValidHttpUrl(uri string) bool {
	url, err := urlpkg.Parse(uri)
	if err != nil {
		panic(fmt.Sprintf("parse url err: %v", err))
	}

	scheme := url.Scheme
	host := url.Host
	if "http" != scheme && "https" != scheme {
		return false
	}

	if host == "" {
		return false
	}

	return true
}
