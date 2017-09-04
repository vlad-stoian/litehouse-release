package integration_tests

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
)

var _ = Describe("adapter commands", func() {
	var session *gexec.Session

	Describe("generate-manifest", func() {
		var (
			plan              string
			previousPlan      string
			arbitraryParams   string
			previousManifest  string
			serviceDeployment string
			// generatedManifest *bosh.BoshManifest
			// unmarshallErr     error
		)

		BeforeEach(func() {
			serviceDeployment = `{
				"deployment_name": "litehouse",
				"releases": [
				],
				"stemcell": {
					"stemcell_os": "ubuntu",
					"stemcell_version": "3312.15"
				}
			}`
			plan = `{
				"instance_groups": [{
					"name": "bosh-lite",
					"vm_type": "medium",
					"persisten_disk_type": "10GB",
					"networks": [
						"default"
					],
					"azs": [
						"z1"
					],
					"instances": 1
				}]
			}`

			arbitraryParams = "{}"
			previousManifest = ""
			previousPlan = "{}"
		})

		JustBeforeEach(func() {
			session = execBin("generate-manifest", serviceDeployment, plan, arbitraryParams, previousManifest, previousPlan)
		})

		It("completes successfully", func() {
			Expect(session.ExitCode()).To(Equal(0))
		})
	})
})
