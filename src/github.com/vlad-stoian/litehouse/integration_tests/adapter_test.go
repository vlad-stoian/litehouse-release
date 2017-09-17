package integration_tests

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
	"github.com/pivotal-cf/on-demand-services-sdk/bosh"
	yaml "gopkg.in/yaml.v2"
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
			generatedManifest *bosh.BoshManifest
			unmarshallErr     error
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
			generatedManifest = new(bosh.BoshManifest)
			unmarshallErr = yaml.Unmarshal(session.Out.Contents(), generatedManifest)
		})

		It("completes successfully", func() {
			Expect(session.ExitCode()).To(Equal(0))
		})

		It("produces valid bosh manifest", func() {
			Expect(unmarshallErr).ToNot(HaveOccurred())
		})

		It("contains bosh instance group", func() {
			Expect(generatedManifest.InstanceGroups).To(HaveLen(1))
			Expect(generatedManifest.InstanceGroups[0].Name).To(Equal("bosh"))
		})

		It("contains correct jobs", func() {
			Expect(generatedManifest.InstanceGroups[0].Jobs).To(HaveLen(8))

			expectedJobNames := []string{"nats", "postgres-9.4", "blobstore", "director", "health_monitor", "warden_cpi", "garden", "user_add"}
			for i, actualJob := range generatedManifest.InstanceGroups[0].Jobs {
				Expect(actualJob.Name).To(Equal(expectedJobNames[i]))
			}
		})

		XIt("contains no unused variables", func() {
			Expect(generatedManifest.Variables).To(HaveLen(0))
		})

		It("has interpolated the external_ip correctly", func() {
			Expect(generatedManifest.InstanceGroups[0].Networks[1].StaticIPs[0]).ToNot(Equal("((external_ip))"))
		})
	})
})
