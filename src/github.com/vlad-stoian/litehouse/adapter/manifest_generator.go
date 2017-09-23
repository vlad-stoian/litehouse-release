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
	serviceDeployment serviceadapter.ServiceDeployment,
	plan serviceadapter.Plan,
	requestParameters serviceadapter.RequestParameters,
	previousManifest *bosh.BoshManifest,
	previousPlan *serviceadapter.Plan,
) (bosh.BoshManifest, error) {
	if _, err := os.Stat(mg.boshLiteTemplatePath); os.IsNotExist(err) {
		return bosh.BoshManifest{}, errors.WrapErrorf(err, "Bosh lite template does not exist: %v", mg.boshLiteTemplatePath)
	}

	if _, err := os.Stat(mg.boshCLIPath); os.IsNotExist(err) {
		return bosh.BoshManifest{}, errors.WrapErrorf(err, "Bosh cli does not exist: %v", mg.boshCLIPath)
	}

	arbitraryParams := requestParameters.ArbitraryParams()
	externalIpInterface, ok := arbitraryParams["ip"]
	if !ok {
		return bosh.BoshManifest{}, fmt.Errorf("IP key not found in request parameters: %v", arbitraryParams)
	}

	externalIp, ok := externalIpInterface.(string)
	if !ok {
		return bosh.BoshManifest{}, fmt.Errorf("IP key found but it is not of type string: %v", externalIpInterface)
	}

	boshLiteManifest, err := mg.interpolateManifest(externalIp)
	if err != nil {
		return bosh.BoshManifest{}, err
	}

	boshLiteManifest.Name = serviceDeployment.DeploymentName

	return boshLiteManifest, nil
}

func (mg ManifestGenerator) interpolateManifest(externalIp string) (bosh.BoshManifest, error) {
	var stdoutBuffer bytes.Buffer
	var stderrBuffer bytes.Buffer

	tempFile, err := ioutil.TempFile("", "example")
	if err != nil {
		return bosh.BoshManifest{}, errors.WrapErrorf(err, "Error creating temp file")
	}
	defer os.Remove(tempFile.Name())

	cmd := exec.Command(
		mg.boshCLIPath,
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
		return bosh.BoshManifest{}, errors.WrapErrorf(err, "Error running bosh cli, stderr: %v", stderrBuffer.String())
	}

	boshLiteManifestBytes := stdoutBuffer.Bytes()

	var boshLiteManifest bosh.BoshManifest

	err = yaml.Unmarshal(boshLiteManifestBytes, &boshLiteManifest)
	if err != nil {
		return bosh.BoshManifest{}, errors.WrapErrorf(err, "Error unmarshalling: %v", boshLiteManifestBytes)
	}

	return boshLiteManifest, nil
}
