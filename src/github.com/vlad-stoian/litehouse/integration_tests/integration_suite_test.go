package integration_tests

import (
	"os/exec"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"

	"testing"
)

var serviceAdapterBinPath string

var _ = BeforeSuite(func() {
	var err error
	serviceAdapterBinPath, err = gexec.Build("github.com/vlad-stoian/litehouse/cmd/service-adapter")

	Expect(err).ToNot(HaveOccurred())
})

var _ = AfterSuite(func() {
	gexec.CleanupBuildArtifacts()
})

func execBin(args ...string) *gexec.Session {
	cmd := exec.Command(serviceAdapterBinPath, args...)

	session, err := gexec.Start(cmd, GinkgoWriter, GinkgoWriter)
	Expect(err).NotTo(HaveOccurred())
	Eventually(session, "5s").Should(gexec.Exit())

	return session
}

func TestAdapter(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Adapter Integration Tests Suite")
}
