package listenbrainz

import (
	"fmt"
	"net/url"
)

type APIPath string

type pathBuilder struct{}

func newPathBuilder() pathBuilder {
	return pathBuilder{}
}

func (p pathBuilder) SubmitListen() APIPath {
	return APIPath("/apis/listenbrainz/1/submit-listens")
}

func (p APIPath) String() string {
	return string(p)
}

func (p APIPath) URL(baseURL string) (*url.URL, error) {
	if baseURL == "" {
		return nil, fmt.Errorf("base URL is empty")
	}

	base, err := url.Parse(baseURL)
	if err != nil {
		return nil, err
	}

	rel, err := url.Parse(p.String())
	if err != nil {
		return nil, err
	}

	return base.ResolveReference(rel), nil
}
