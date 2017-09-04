package adapter

import (
	"github.com/pivotal-cf/on-demand-services-sdk/bosh"
	"github.com/pivotal-cf/on-demand-services-sdk/serviceadapter"
)

type Binder struct{}

func (b Binder) CreateBinding(_ string, deploymentTopology bosh.BoshVMs, manifest bosh.BoshManifest, _ serviceadapter.RequestParameters) (serviceadapter.Binding, error) {
	return nil, nil
}
