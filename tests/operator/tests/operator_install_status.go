package operator

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/operator-framework/api/pkg/operators/v1alpha1"

	"github.com/test-network-function/cnfcert-tests-verification/tests/globalhelper"
	"github.com/test-network-function/cnfcert-tests-verification/tests/globalparameters"
	tshelper "github.com/test-network-function/cnfcert-tests-verification/tests/operator/helper"
	tsparams "github.com/test-network-function/cnfcert-tests-verification/tests/operator/parameters"
)

var _ = Describe("Operator install-source,", func() {
	var randomNamespace string
	var origReportDir string
	var origTnfConfigDir string

	BeforeEach(func() {
		// Create random namespace and keep original report and TNF config directories
		randomNamespace, origReportDir, origTnfConfigDir = globalhelper.BeforeEachSetupWithRandomNamespace(tsparams.OperatorNamespace)

		By("Define TNF config file")
		err := globalhelper.DefineTnfConfig(
			[]string{randomNamespace},
			[]string{tsparams.TestPodLabel},
			[]string{},
			[]string{},
			[]string{})
		Expect(err).ToNot(HaveOccurred())

		By("Deploy operator group")
		err = tshelper.DeployTestOperatorGroup(randomNamespace)
		Expect(err).ToNot(HaveOccurred(), "Error deploying operator group")

		By("Deploy cloudbees-ci operator for testing")
		err = tshelper.DeployOperatorSubscription(
			tsparams.OperatorPrefixCloudbees,
			"alpha",
			randomNamespace,
			tsparams.CertifiedOperatorGroup,
			tsparams.OperatorSourceNamespace,
			"",
			v1alpha1.ApprovalAutomatic,
		)
		Expect(err).ToNot(HaveOccurred(), "Error deploying operator "+
			tsparams.OperatorPrefixCloudbees)

		err = tshelper.WaitUntilOperatorIsReady(tsparams.OperatorPrefixCloudbees,
			randomNamespace)
		Expect(err).ToNot(HaveOccurred(), "Operator "+tsparams.OperatorPrefixCloudbees+
			" is not ready")
	})

	AfterEach(func() {
		globalhelper.AfterEachCleanupWithRandomNamespace(randomNamespace, origReportDir, origTnfConfigDir, tsparams.WaitingTime)
	})

	It("one operator that reports Succeeded as its installation status", func() {
		By("Label operator")
		Eventually(func() error {
			return tshelper.AddLabelToInstalledCSV(
				tsparams.OperatorPrefixCloudbees,
				randomNamespace,
				tsparams.OperatorLabel)
		}, tsparams.TimeoutLabelCsv, tsparams.PollingInterval).Should(Not(HaveOccurred()),
			"Error labeling operator "+tsparams.OperatorPrefixCloudbees)

		By("Start test")
		err := globalhelper.LaunchTests(
			tsparams.TnfOperatorInstallStatus,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()))
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(
			tsparams.TnfOperatorInstallStatus,
			globalparameters.TestCasePassed)
		Expect(err).ToNot(HaveOccurred())
	})

	It("two operators, one does not reports Succeeded as its installation status (quick failure) [negative]", func() {
		By("Deploy openvino operator for testing")
		// The OpenVINO operator fails quickly due to the fact that it does not support the install mode type
		// that it is used (OwnNamespace).
		err := tshelper.DeployOperatorSubscription(
			"ovms-operator",
			"alpha",
			randomNamespace,
			tsparams.CertifiedOperatorGroup,
			tsparams.OperatorSourceNamespace,
			"",
			v1alpha1.ApprovalAutomatic,
		)
		Expect(err).ToNot(HaveOccurred(), "Error deploying operator "+
			tsparams.OperatorPrefixOpenvino)

		err = tshelper.WaitUntilOperatorIsReady(tsparams.OperatorPrefixOpenvino,
			randomNamespace)
		Expect(err).ToNot(HaveOccurred(), "Operator "+tsparams.OperatorPrefixOpenvino+
			" is not ready")

		defer func() {
			err := tshelper.DeleteLabelFromInstalledCSV(
				tsparams.OperatorPrefixOpenvino,
				randomNamespace,
				tsparams.OperatorLabel)
			Expect(err).ToNot(HaveOccurred(), "Error removing label from operator "+tsparams.OperatorPrefixOpenvino)
		}()

		By("Label operators")
		Eventually(func() error {
			return tshelper.AddLabelToInstalledCSV(
				tsparams.OperatorPrefixCloudbees,
				randomNamespace,
				tsparams.OperatorLabel)
		}, tsparams.TimeoutLabelCsv, tsparams.PollingInterval).Should(Not(HaveOccurred()),
			"Error labeling operator "+tsparams.OperatorPrefixCloudbees)

		Eventually(func() error {
			return tshelper.AddLabelToInstalledCSV(
				tsparams.OperatorPrefixOpenvino,
				randomNamespace,
				tsparams.OperatorLabel)
		}, tsparams.TimeoutLabelCsv, tsparams.PollingInterval).Should(Not(HaveOccurred()),
			"Error labeling operator "+tsparams.OperatorPrefixOpenvino)

		By("Start test")
		err = globalhelper.LaunchTests(
			tsparams.TnfOperatorInstallStatus,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()))
		Expect(err).To(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(
			tsparams.TnfOperatorInstallStatus,
			globalparameters.TestCaseFailed)
		Expect(err).ToNot(HaveOccurred())
	})

	It("two operators, one does not reports Succeeded as its installation status (delayed failure) [negative]", func() {
		By("Deploy anchore-engine operator for testing")
		// The node selector prevents the operator's installation from succeeding but it stays some time in the Installing phase
		// before failing. This way the failure is delayed, allowing the testing of the CNF Certification Suite timeout mechanism
		// for operator readiness.
		nodeSelector := map[string]string{"target": "none"}
		err := tshelper.DeployOperatorSubscriptionWithNodeSelector(
			"anchore-engine",
			"alpha",
			randomNamespace,
			tsparams.CertifiedOperatorGroup,
			tsparams.OperatorSourceNamespace,
			"",
			v1alpha1.ApprovalAutomatic,
			nodeSelector,
		)
		Expect(err).ToNot(HaveOccurred(), "Error deploying operator "+
			tsparams.OperatorPrefixAnchore)

		// Do not wait until the operator is ready. This time the CNF Certification suite must handle the situation.

		defer func() {
			err := tshelper.DeleteLabelFromInstalledCSV(
				tsparams.OperatorPrefixAnchore,
				randomNamespace,
				tsparams.OperatorLabel)
			Expect(err).ToNot(HaveOccurred(), "Error removing label from operator "+tsparams.OperatorPrefixAnchore)
		}()

		By("Label operators")
		Eventually(func() error {
			return tshelper.AddLabelToInstalledCSV(
				tsparams.OperatorPrefixCloudbees,
				randomNamespace,
				tsparams.OperatorLabel)
		}, tsparams.TimeoutLabelCsv, tsparams.PollingInterval).Should(Not(HaveOccurred()),
			"Error labeling operator "+tsparams.OperatorPrefixCloudbees)

		Eventually(func() error {
			return tshelper.AddLabelToInstalledCSV(
				tsparams.OperatorPrefixAnchore,
				randomNamespace,
				tsparams.OperatorLabel)
		}, tsparams.TimeoutLabelCsv, tsparams.PollingInterval).Should(Not(HaveOccurred()),
			"Error labeling operator "+tsparams.OperatorPrefixAnchore)

		By("Start test")
		err = globalhelper.LaunchTests(
			tsparams.TnfOperatorInstallStatus,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()))
		Expect(err).To(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(
			tsparams.TnfOperatorInstallStatus,
			globalparameters.TestCaseFailed)
		Expect(err).ToNot(HaveOccurred())
	})

})
