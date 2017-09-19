package adapter

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"

	"github.com/cloudfoundry/bosh-utils/errors"
	"github.com/pivotal-cf/on-demand-services-sdk/bosh"
	"github.com/pivotal-cf/on-demand-services-sdk/serviceadapter"
	"gopkg.in/yaml.v2"
)

type ManifestGenerator struct {
	boshLiteTemplatePath string
	boshCLIPath          string
}

func NewManifestGenerator(boshLiteTemplatePath, boshCLIPath string) ManifestGenerator {
	return ManifestGenerator{
		boshLiteTemplatePath: boshLiteTemplatePath,
		boshCLIPath:          boshCLIPath,
	}
}

func (mg ManifestGenerator) GenerateManifest(
	deployment serviceadapter.ServiceDeployment,
	plan serviceadapter.Plan,
	params serviceadapter.RequestParameters,
	previousManifest *bosh.BoshManifest,
	previousPlan *serviceadapter.Plan,
) (bosh.BoshManifest, error) {
	if _, err := os.Stat(mg.boshLiteTemplatePath); os.IsNotExist(err) {
		return bosh.BoshManifest{}, errors.WrapErrorf(err, "bosh lite template does not exist: %v", mg.boshLiteTemplatePath)
	}

	if _, err := os.Stat(mg.boshCLIPath); os.IsNotExist(err) {
		return bosh.BoshManifest{}, errors.WrapErrorf(err, "bosh cli does not exist: %v", mg.boshCLIPath)
	}

	arbitraryParams := params.ArbitraryParams()
	externalIp, ok := arbitraryParams["ip"]
	if !ok {
		return bosh.BoshManifest{}, fmt.Errorf("ip key not found in request parameters: %v", arbitraryParams)
	}

	var stdoutBuffer bytes.Buffer
	var stderrBuffer bytes.Buffer

	tempFile, err := ioutil.TempFile("", "example")
	if err != nil {
		return bosh.BoshManifest{}, errors.WrapErrorf(err, "error creating temp file")
	}
	defer os.Remove(tempFile.Name())

	cmd := exec.Command(mg.boshCLIPath,
		"interpolate",
		mg.boshLiteTemplatePath,
		fmt.Sprintf("--var=%s=%s", "director_name", "\"Bosh Lite Director\""),
		fmt.Sprintf("--var=%s=%s", "internal_ip", "127.0.0.1"),
		fmt.Sprintf("--var=%s=%s", "external_ip", externalIp),
		fmt.Sprintf("--vars-store=%s", tempFile.Name()),
	)

	cmd.Stdout = &stdoutBuffer
	cmd.Stderr = &stderrBuffer

	err = cmd.Run()
	if err != nil {
		return bosh.BoshManifest{}, errors.WrapErrorf(err, "error running bosh cli: %v", stderrBuffer.String())
	}

	boshLiteManifestBytes := stdoutBuffer.Bytes()

	var boshLiteManifest bosh.BoshManifest

	err = yaml.Unmarshal(boshLiteManifestBytes, &boshLiteManifest)
	if err != nil {
		return bosh.BoshManifest{}, fmt.Errorf("error unmarshalling: %v", boshLiteManifestBytes)
	}

	return boshLiteManifest, nil
}
