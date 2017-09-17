package adapter

import (
	"fmt"
	"io/ioutil"
	"path/filepath"

	boshtpl "github.com/cloudfoundry/bosh-cli/director/template"
	"github.com/pivotal-cf/on-demand-services-sdk/bosh"
	"github.com/pivotal-cf/on-demand-services-sdk/serviceadapter"
	yaml "gopkg.in/yaml.v2"
)

type ManifestGenerator struct{}

func (mg ManifestGenerator) GenerateManifest(
	deployment serviceadapter.ServiceDeployment,
	plan serviceadapter.Plan,
	params serviceadapter.RequestParameters,
	previousManifest *bosh.BoshManifest,
	previousPlan *serviceadapter.Plan,
) (bosh.BoshManifest, error) {

	boshLiteTemplatePath := filepath.FromSlash("../assets/bosh-lite-template.yml")
	contents, err := ioutil.ReadFile(boshLiteTemplatePath)
	if err != nil {
		return bosh.BoshManifest{}, fmt.Errorf("Error reading from: %v", boshLiteTemplatePath)
	}

	template := boshtpl.NewTemplate(contents)
	variables := boshtpl.StaticVariables{
		"external_ip": "1234",
	}

	// something := cfgtypes.NewValueGeneratorConcrete(cfgtypes.NewVarsCertLoader(vars))
	// https://github.com/cloudfoundry/bosh-cli/blob/06eede24ca808e599a73fd38251ed84d0c506e07/cmd/var_flags.go

	boshLiteManifestBytes, _ := template.Evaluate(variables, nil, boshtpl.EvaluateOpts{})

	var boshLiteManifest bosh.BoshManifest
	err = yaml.Unmarshal(boshLiteManifestBytes, &boshLiteManifest)
	if err != nil {
		return bosh.BoshManifest{}, fmt.Errorf("Error unmarshalling: %v", contents)
	}

	return boshLiteManifest, nil
}
