package main

import (
	"os"

	"path/filepath"

	"github.com/pivotal-cf/on-demand-services-sdk/serviceadapter"
	"github.com/vlad-stoian/litehouse/adapter"
)

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func main() {
	boshCLIPath := filepath.FromSlash(getEnv("BOSH_CLI_PATH", "/var/vcap/packages/go/bin/"))
	boshLiteTemplatePath := filepath.FromSlash(getEnv("BOSH_LITE_TEMPLATE_PATH", "/not/yet/set"))

	manifestGenerator := adapter.NewManifestGenerator(
		boshLiteTemplatePath,
		boshCLIPath,
	)

	binder := new(adapter.Binder)

	serviceadapter.HandleCommandLineInvocation(os.Args, manifestGenerator, binder, nil)
}
