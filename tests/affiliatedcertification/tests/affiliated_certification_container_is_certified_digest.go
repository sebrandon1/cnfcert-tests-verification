package tests

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/test-network-function/cnfcert-tests-verification/tests/globalhelper"
	"github.com/test-network-function/cnfcert-tests-verification/tests/globalparameters"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/deployment"

	tsparams "github.com/test-network-function/cnfcert-tests-verification/tests/affiliatedcertification/parameters"

	"github.com/test-network-function/oct/pkg/certdb/onlinecheck"
)

var _ = Describe("Affiliated-certification container-is-certified-digest,", Serial, func() {
	var randomNamespace string
	var randomReportDir string
	var randomTnfConfigDir string

	BeforeEach(func() {
		// Create random namespace and keep original report and TNF config directories
		randomNamespace, randomReportDir, randomTnfConfigDir = globalhelper.BeforeEachSetupWithRandomNamespace(
			tsparams.TestCertificationNameSpace)

		By("Define tnf config file")
		err := globalhelper.DefineTnfConfig(
			[]string{randomNamespace},
			[]string{tsparams.TestPodLabel},
			[]string{},
			[]string{},
			[]string{}, randomTnfConfigDir)
		Expect(err).ToNot(HaveOccurred(), "error defining tnf config file")

		By("Check if the test image is certified prior to deployment")
		// Using the 'oct' repo, we should do a quick assertion to see if the image is available
		// and certified.
		onlineValidator := onlinecheck.NewOnlineValidator()

		// The information for this is gathered from:
		//nolint:lll
		// https://catalog.redhat.com/api/containers/v1/images?filter=image_id==sha256:41bc5b622db8b5e0d608e7524c39928b191270666252578edbf1e0f84a9e3cab
		//nolint:lll
		Expect(onlineValidator.IsContainerCertified("registry.access.redhat.com", "ubi8/nodejs-12", "latest", "sha256:41bc5b622db8b5e0d608e7524c39928b191270666252578edbf1e0f84a9e3cab")).To(BeTrue())
	})

	AfterEach(func() {
		globalhelper.AfterEachCleanupWithRandomNamespace(randomNamespace, randomReportDir, randomTnfConfigDir, tsparams.Timeout)
	})

	// 66765
	It("one container to test, container is certified digest", func() {
		if globalhelper.IsKindCluster() {
			Skip("Skip test due to image pull missing credentials in Kind")
		}

		By("Define deployment with certified container")
		dep := deployment.DefineDeployment("affiliated-cert-deployment", randomNamespace,
			tsparams.CertifiedContainerURLNodeJs, tsparams.TestDeploymentLabels)

		By("Create deployment")
		err := globalhelper.CreateAndWaitUntilDeploymentIsReady(dep, tsparams.Timeout)
		Expect(err).ToNot(HaveOccurred())

		By("Start test")
		err = globalhelper.LaunchTests(
			tsparams.TestCaseNameContainerDigest,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()), randomReportDir, randomTnfConfigDir)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(
			tsparams.TestCaseNameContainerDigest,
			globalparameters.TestCasePassed, randomReportDir)
		Expect(err).ToNot(HaveOccurred())
	})

	// 66766
	It("one container to test, container is not certified digest [negative]", func() {
		By("Define deployment with uncertified container")

		dep := deployment.DefineDeployment("affiliated-cert-deployment", randomNamespace,
			tsparams.UncertifiedContainerURLCnfTest, tsparams.TestDeploymentLabels)

		By("Create deployment")
		err := globalhelper.CreateAndWaitUntilDeploymentIsReady(dep, tsparams.Timeout)
		Expect(err).ToNot(HaveOccurred())

		By("Start test")
		err = globalhelper.LaunchTests(
			tsparams.TestCaseNameContainerDigest,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()), randomReportDir, randomTnfConfigDir)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(
			tsparams.TestCaseNameContainerDigest,
			globalparameters.TestCaseFailed, randomReportDir)
		Expect(err).ToNot(HaveOccurred())
	})

	// 66767
	It("two containers to test, both are certified digest", func() {
		if globalhelper.IsKindCluster() {
			Skip("Skip test due to image pull missing credentials in Kind")
		}

		By("Define deployments with certified containers")
		dep := deployment.DefineDeployment("affiliated-cert-deployment", randomNamespace,
			tsparams.CertifiedContainerURLNodeJs, tsparams.TestDeploymentLabels)

		By("Create deployment 1")
		err := globalhelper.CreateAndWaitUntilDeploymentIsReady(dep, tsparams.Timeout)
		Expect(err).ToNot(HaveOccurred())

		dep2 := deployment.DefineDeployment("affiliated-cert-deployment-2", randomNamespace,
			tsparams.CertifiedContainerURLCockroachDB, tsparams.TestDeploymentLabels)

		By("Create deployment 2")
		err = globalhelper.CreateAndWaitUntilDeploymentIsReady(dep2, tsparams.Timeout)
		Expect(err).ToNot(HaveOccurred())

		By("Start test")
		err = globalhelper.LaunchTests(
			tsparams.TestCaseNameContainerDigest,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()), randomReportDir, randomTnfConfigDir)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(
			tsparams.TestCaseNameContainerDigest,
			globalparameters.TestCasePassed, randomReportDir)
		Expect(err).ToNot(HaveOccurred())
	})

	// 66768
	It("two containers to test, one is certified, one is not digest [negative]", func() {
		if globalhelper.IsKindCluster() {
			Skip("Skip test due to image pull missing credentials in Kind")
		}

		By("Define deployments with different container certification statuses")
		dep := deployment.DefineDeployment("affiliated-cert-deployment", randomNamespace,
			tsparams.UncertifiedContainerURLCnfTest, tsparams.TestDeploymentLabels)

		By("Create deployment 1")
		err := globalhelper.CreateAndWaitUntilDeploymentIsReady(dep, tsparams.Timeout)
		Expect(err).ToNot(HaveOccurred())

		dep2 := deployment.DefineDeployment("affiliated-cert-deployment-2", randomNamespace,
			tsparams.CertifiedContainerURLCockroachDB, tsparams.TestDeploymentLabels)

		By("Create deployment 2")
		err = globalhelper.CreateAndWaitUntilDeploymentIsReady(dep2, tsparams.Timeout)
		Expect(err).ToNot(HaveOccurred())

		By("Start test")
		err = globalhelper.LaunchTests(
			tsparams.TestCaseNameContainerDigest,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()), randomReportDir, randomTnfConfigDir)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(
			tsparams.TestCaseNameContainerDigest,
			globalparameters.TestCaseFailed, randomReportDir)
		Expect(err).ToNot(HaveOccurred())
	})

})
