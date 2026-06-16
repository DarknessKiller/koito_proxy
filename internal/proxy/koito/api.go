package koito

import (
	"koito_proxy/internal/routeutil"
)

type apiPathBuilder struct{}

func newAPIPathBuilder() apiPathBuilder {
	return apiPathBuilder{}
}

func (apiPathBuilder) MergeEntity() routeutil.APIPath {
	return routeutil.APIPath("/apis/web/v1/:entity/:id/merge")
}
