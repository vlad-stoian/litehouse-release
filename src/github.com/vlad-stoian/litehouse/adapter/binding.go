package adapter

import (
	"github.com/pivotal-cf/on-demand-services-sdk/bosh"
	"github.com/pivotal-cf/on-demand-services-sdk/serviceadapter"
)

type Binder struct{}

func (b Binder) CreateBinding(
	bindingID string,
	deploymentTopology bosh.BoshVMs,
	manifest bosh.BoshManifest,
	requestParameters serviceadapter.RequestParameters,
) (serviceadapter.Binding, error) {

	return serviceadapter.Binding{}, nil
}

func (b Binder) DeleteBinding(
	_ string,
	_ bosh.BoshVMs,
	_ bosh.BoshManifest,
	_ serviceadapter.RequestParameters,
) error {
	return nil
}
