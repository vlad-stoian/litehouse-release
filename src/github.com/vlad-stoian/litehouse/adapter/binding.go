package adapter

import (
	"fmt"
	"strings"

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

	storeHack, ok := manifest.Properties["store_hack"]
	if !ok {
		return serviceadapter.Binding{}, fmt.Errorf("store_hack key not found in manifest properties: %v", manifest.Properties)
	}

	credentialsMap, ok := storeHack.(map[interface{}]interface{})
	if !ok {
		return serviceadapter.Binding{}, fmt.Errorf("store_hack is not of type map[interface{}]interface{}: %v", manifest.Properties)
	}

	binding := serviceadapter.Binding{
		Credentials: map[string]interface{}{},
	}

	for credentialName, credentialValue := range credentialsMap {
		upperCredentialName := strings.ToUpper(credentialName.(string))
		binding.Credentials[upperCredentialName] = credentialValue
	}

	return binding, nil
}

func (b Binder) DeleteBinding(
	_ string,
	_ bosh.BoshVMs,
	_ bosh.BoshManifest,
	_ serviceadapter.RequestParameters,
) error {
	return nil
}
