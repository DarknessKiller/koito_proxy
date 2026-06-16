package auth

import (
	"koito_proxy/internal/routeutil"
)

type apiPathBuilder struct{}

func newAPIPathBuilder() apiPathBuilder {
	return apiPathBuilder{}
}

func (p apiPathBuilder) LBAuthorization() routeutil.APIPath {
	return routeutil.APIPath("/apis/listenbrainz/1/validate-token")
}

func (p apiPathBuilder) KoitoAuthorization() routeutil.APIPath {
	return routeutil.APIPath("/apis/web/v1/user")
}
