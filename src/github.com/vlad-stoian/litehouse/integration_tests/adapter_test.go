package integration_tests

import (
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
	"github.com/pivotal-cf/on-demand-services-sdk/bosh"
	"github.com/pivotal-cf/on-demand-services-sdk/serviceadapter"
	yaml "gopkg.in/yaml.v2"
)

var _ = Describe("adapter commands", func() {
	var session *gexec.Session

	Describe("generate-manifest", func() {
		var (
			plan                   string
			previousPlan           string
			arbitraryParams        string
			previousManifest       string
			serviceDeployment      string
			generatedManifest      *bosh.BoshManifest
			generatedManifestBytes []byte
			unmarshallErr          error
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
					"persistent_disk_type": "10GB",
					"networks": [
						"default"
					],
					"azs": [
						"z1"
					],
					"instances": 1
				}]
			}`

			arbitraryParams = "{ \"parameters\": { \"ip\": \"external-ip\" } }"
			previousManifest = ""
			previousPlan = "{}"

			os.Setenv("BOSH_CLI_PATH", "/usr/local/bin/bosh")
			os.Setenv("BOSH_LITE_TEMPLATE_PATH", "/Users/pivotal/workspace/litehouse-release/src/github.com/vlad-stoian/litehouse/assets/bosh-lite-template.yml")
		})

		JustBeforeEach(func() {
			session = execBin("generate-manifest", serviceDeployment, plan, arbitraryParams, previousManifest, previousPlan)
			generatedManifest = new(bosh.BoshManifest)
			generatedManifestBytes = session.Out.Contents()
			unmarshallErr = yaml.Unmarshal(generatedManifestBytes, generatedManifest)
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

		It("contains no unused variables", func() {
			Expect(generatedManifest.Variables).To(HaveLen(0))
		})

		It("has interpolated the external_ip correctly", func() {
			Expect(generatedManifest.InstanceGroups[0].Networks[1].StaticIPs[0]).ToNot(Equal("((external_ip))"))
		})

		It("has all the varibles inside the store_hack", func() {
			By("contains the store_hack property")
			Expect(generatedManifest.Properties).To(HaveKey("store_hack"))
			storeHackProperty := generatedManifest.Properties["store_hack"]

			By("store_hack property is a map")
			Expect(storeHackProperty).To(BeAssignableToTypeOf(map[interface{}]interface{}{}))
			storeHackMap := storeHackProperty.(map[interface{}]interface{})

			By("contains all the keys")
			expectedProperties := []string{
				"bosh_client",
				"bosh_client_secret",
				"bosh_ca_cert",
				"bosh_environment",
				"bosh_gw_user",
				"bosh_gw_private_key",
			}
			for _, expectedProperty := range expectedProperties {
				Expect(storeHackMap).To(HaveKey(expectedProperty))
			}

		})
	})

	FDescribe("create-binding", func() {
		var (
			bindingID             string
			deploymentTopology    string
			manifest              string
			requestParameters     string
			generatedBinding      *serviceadapter.Binding
			generatedBindingBytes []byte
			unmarshallErr         error
		)

		BeforeEach(func() {
			bindingID = "binding-id"
			deploymentTopology = `{
				"service-instance": ["1.1.1.1"]
			}`

			manifest = `{
				"instance_groups": [{
					"name": "bosh-lite",
					"vm_type": "medium",
					"persistent_disk_type": "10GB",
					"networks": [{
						"name": "default"
					}],
					"azs": [
						"z1"
					],
					"instances": 1
				}]
			}`

			requestParameters = "{}"

			os.Setenv("BOSH_CLI_PATH", "/usr/local/bin/bosh")
			os.Setenv("BOSH_LITE_TEMPLATE_PATH", "/Users/pivotal/workspace/litehouse-release/src/github.com/vlad-stoian/litehouse/assets/bosh-lite-template.yml")
		})

		JustBeforeEach(func() {
			session = execBin("create-binding", bindingID, deploymentTopology, manifest, requestParameters)
			generatedBinding = new(serviceadapter.Binding)
			generatedBindingBytes = session.Out.Contents()
			unmarshallErr = yaml.Unmarshal(generatedBindingBytes, generatedBinding)
		})

		It("completes successfully", func() {
			Expect(session.ExitCode()).To(Equal(0))
		})

		It("produces valid binding", func() {
			Expect(unmarshallErr).ToNot(HaveOccurred())
		})
	})
})
