package koito

import (
	"fmt"
	"net/url"
	"strings"
)

type APIPath string

type pathBuilder struct{}

func newPathBuilder() pathBuilder {
	return pathBuilder{}
}

func (p pathBuilder) MergeEntity() APIPath {
	return APIPath("/apis/web/v1/:entity/:id/merge")
}

func (p pathBuilder) Track(id string) APIPath {
	return APIPath(fmt.Sprintf("/apis/web/v1/track/%s", id))
}

func (p pathBuilder) Artist(id string) APIPath {
	return APIPath(fmt.Sprintf("/apis/web/v1/artist/%s", id))
}

func (p pathBuilder) Album(id string) APIPath {
	return APIPath(fmt.Sprintf("/apis/web/v1/album/%s", id))
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

// URLWithParams resolves an APIPath containing ":param" tokens against baseURL,
// replacing tokens using the provided params map.
func (p APIPath) URLWithParams(baseURL string, params map[string]string) (*url.URL, error) {
	s := p.String()
	for k, v := range params {
		s = strings.ReplaceAll(s, ":"+k, v)
	}

	base, err := url.Parse(baseURL)
	if err != nil {
		return nil, err
	}
	rel, err := url.Parse(s)
	if err != nil {
		return nil, err
	}
	return base.ResolveReference(rel), nil
}
