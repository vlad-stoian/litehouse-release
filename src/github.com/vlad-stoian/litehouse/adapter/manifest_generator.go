package adapter

import (
	"github.com/pivotal-cf/on-demand-services-sdk/bosh"
	"github.com/pivotal-cf/on-demand-services-sdk/serviceadapter"
)

type ManifestGenerator struct{}

func (mg ManifestGenerator) GenerateManifest(
	deployment serviceadapter.ServiceDeployment,
	plan serviceadapter.Plan,
	params serviceadapter.RequestParameters,
	previousManifest *bosh.BoshManifest,
	previousPlan *serviceadapter.Plan,
) (bosh.BoshManifest, error) {

	return nil, nil
}
