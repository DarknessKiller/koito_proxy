package listenbrainz

import (
	"koito_proxy/internal/routeutil"
)

type apiPathBuilder struct{}

func newAPIPathBuilder() apiPathBuilder {
	return apiPathBuilder{}
}

func (p apiPathBuilder) SubmitListen() routeutil.APIPath {
	return routeutil.APIPath("/apis/listenbrainz/1/submit-listens")
}
