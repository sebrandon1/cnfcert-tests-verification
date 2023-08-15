package tests

import (
	"fmt"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/test-network-function/cnfcert-tests-verification/tests/globalhelper"
	"github.com/test-network-function/cnfcert-tests-verification/tests/globalparameters"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/namespaces"
	corev1 "k8s.io/api/core/v1"

	tshelper "github.com/test-network-function/cnfcert-tests-verification/tests/observability/helper"
	tsparams "github.com/test-network-function/cnfcert-tests-verification/tests/observability/parameters"
)

var _ = Describe(tsparams.TnfTerminationMsgPolicyTcName, func() {
	const tnfTestCaseName = tsparams.TnfTerminationMsgPolicyTcName
	qeTcFileName := globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText())

	var randomNamespace string
	var origReportDir string

	BeforeEach(func() {
		randomNamespace = tsparams.TestNamespace + "-" + globalhelper.GenerateRandomString(10)

		By(fmt.Sprintf("Create %s namespace", randomNamespace))
		err := namespaces.Create(randomNamespace, globalhelper.GetAPIClient())
		Expect(err).ToNot(HaveOccurred())

		By("Override default report directory")
		origReportDir = globalhelper.GetConfiguration().General.TnfReportDir
		reportDir := origReportDir + "/" + randomNamespace
		globalhelper.OverrideReportDir(reportDir)

		By("Define TNF config file")
		err = globalhelper.DefineTnfConfig(
			[]string{randomNamespace},
			tshelper.GetTnfTargetPodLabelsSlice(),
			[]string{},
			[]string{},
			[]string{tsparams.CrdSuffix1, tsparams.CrdSuffix2})
		Expect(err).ToNot(HaveOccurred())
	})

	AfterEach(func() {
		By(fmt.Sprintf("Remove %s namespace", randomNamespace))
		err := namespaces.DeleteAndWait(
			globalhelper.GetAPIClient().CoreV1Interface,
			randomNamespace,
			tsparams.NsResourcesDeleteTimeoutMins,
		)
		Expect(err).ToNot(HaveOccurred())

		By("Restore default report directory")
		globalhelper.GetConfiguration().General.TnfReportDir = origReportDir
	})

	// Positive #1.
	It("One deployment one pod one container with terminationMessagePolicy set to FallbackToLogsOnError", func() {

		By("Create deployment in the cluster")
		deployment := tshelper.DefineDeploymentWithTerminationMsgPolicies(tsparams.TestDeploymentBaseName, randomNamespace, 1,
			[]corev1.TerminationMessagePolicy{corev1.TerminationMessageFallbackToLogsOnError})

		err := globalhelper.CreateAndWaitUntilDeploymentIsReady(deployment, tsparams.DeploymentDeployTimeoutMins)
		Expect(err).ToNot(HaveOccurred())

		By("Start TNF " + tnfTestCaseName + " test case")
		err = globalhelper.LaunchTests(tnfTestCaseName, qeTcFileName)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(tnfTestCaseName, globalparameters.TestCasePassed)
		Expect(err).ToNot(HaveOccurred())
	})

	// // Positive #2.
	It("One deployment one pod two containers both with terminationMessagePolicy set to FallbackToLogsOnError", func() {

		By("Create deployment in the cluster")
		deployment := tshelper.DefineDeploymentWithTerminationMsgPolicies(tsparams.TestDeploymentBaseName,
			randomNamespace, 1,
			[]corev1.TerminationMessagePolicy{
				corev1.TerminationMessageFallbackToLogsOnError,
				corev1.TerminationMessageFallbackToLogsOnError,
			})

		err := globalhelper.CreateAndWaitUntilDeploymentIsReady(deployment, tsparams.DeploymentDeployTimeoutMins)
		Expect(err).ToNot(HaveOccurred())

		By("Start TNF " + tnfTestCaseName + " test case")
		err = globalhelper.LaunchTests(tnfTestCaseName, qeTcFileName)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(tnfTestCaseName, globalparameters.TestCasePassed)
		Expect(err).ToNot(HaveOccurred())
	})

	// Positive #3.
	It("One daemonset with two containers, both with terminationMessagePolicy "+
		"set to FallbackToLogsOnError", func() {

		By("Create deployment in the cluster")
		daemonSet := tshelper.DefineDaemonSetWithTerminationMsgPolicies(tsparams.TestDaemonSetBaseName,
			randomNamespace,
			[]corev1.TerminationMessagePolicy{
				corev1.TerminationMessageFallbackToLogsOnError,
				corev1.TerminationMessageFallbackToLogsOnError,
			})

		err := globalhelper.CreateAndWaitUntilDaemonSetIsReady(daemonSet, tsparams.DaemonSetDeployTimeoutMins)
		Expect(err).ToNot(HaveOccurred())

		By("Start TNF " + tnfTestCaseName + " test case")
		err = globalhelper.LaunchTests(tnfTestCaseName, qeTcFileName)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(tnfTestCaseName, globalparameters.TestCasePassed)
		Expect(err).ToNot(HaveOccurred())
	})

	// Positive #4
	It("One deployment and one statefulset, both with one pod with one container, "+
		"all with terminationMessagePolicy set to FallbackToLogsOnError", func() {

		By("Create deployment in the cluster")
		deployment := tshelper.DefineDeploymentWithTerminationMsgPolicies(tsparams.TestDeploymentBaseName,
			randomNamespace, 1,
			[]corev1.TerminationMessagePolicy{corev1.TerminationMessageFallbackToLogsOnError})

		err := globalhelper.CreateAndWaitUntilDeploymentIsReady(deployment, tsparams.DeploymentDeployTimeoutMins)
		Expect(err).ToNot(HaveOccurred())

		By("Create statefulset in the cluster")
		statefulSet := tshelper.DefineStatefulSetWithTerminationMsgPolicies(tsparams.TestStatefulSetBaseName,
			randomNamespace, 1,
			[]corev1.TerminationMessagePolicy{corev1.TerminationMessageFallbackToLogsOnError})

		err = globalhelper.CreateAndWaitUntilStatefulSetIsReady(statefulSet, tsparams.StatefulSetDeployTimeoutMins)
		Expect(err).ToNot(HaveOccurred())

		By("Start TNF " + tnfTestCaseName + " test case")
		err = globalhelper.LaunchTests(tnfTestCaseName, qeTcFileName)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(tnfTestCaseName, globalparameters.TestCasePassed)
		Expect(err).ToNot(HaveOccurred())
	})

	// Negative #1.
	It("One deployment one pod one container without terminationMessagePolicy [negative]", func() {

		By("Create deployment in the cluster")
		deployment := tshelper.DefineDeploymentWithTerminationMsgPolicies(tsparams.TestDeploymentBaseName,
			randomNamespace, 1,
			[]corev1.TerminationMessagePolicy{tsparams.UseDefaultTerminationMsgPolicy})

		err := globalhelper.CreateAndWaitUntilDeploymentIsReady(deployment, tsparams.DeploymentDeployTimeoutMins)
		Expect(err).ToNot(HaveOccurred())

		By("Start TNF " + tnfTestCaseName + " test case")
		err = globalhelper.LaunchTests(tnfTestCaseName, qeTcFileName)
		Expect(err).To(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(tnfTestCaseName, globalparameters.TestCaseFailed)
		Expect(err).ToNot(HaveOccurred())
	})

	// Negative #2.
	It("One deployment one pod two containers, only one container with terminationMessagePolicy "+
		"set to FallbackToLogsOnError [negative]", func() {

		By("Create deployment in the cluster")
		deployment := tshelper.DefineDeploymentWithTerminationMsgPolicies(tsparams.TestDeploymentBaseName,
			randomNamespace, 1,
			[]corev1.TerminationMessagePolicy{
				tsparams.UseDefaultTerminationMsgPolicy,
				corev1.TerminationMessageFallbackToLogsOnError,
			})

		err := globalhelper.CreateAndWaitUntilDeploymentIsReady(deployment, tsparams.DeploymentDeployTimeoutMins)
		Expect(err).ToNot(HaveOccurred())

		By("Start TNF " + tnfTestCaseName + " test case")
		err = globalhelper.LaunchTests(tnfTestCaseName, qeTcFileName)
		Expect(err).To(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(tnfTestCaseName, globalparameters.TestCaseFailed)
		Expect(err).ToNot(HaveOccurred())
	})

	// Negative #3.
	It("One deployment with two pods with one container each without terminationMessagePolicy set [negative]", func() {

		By("Create deployment in the cluster")
		deployment := tshelper.DefineDeploymentWithTerminationMsgPolicies(tsparams.TestDeploymentBaseName,
			randomNamespace, 2,
			[]corev1.TerminationMessagePolicy{
				tsparams.UseDefaultTerminationMsgPolicy,
			})

		err := globalhelper.CreateAndWaitUntilDeploymentIsReady(deployment, tsparams.DeploymentDeployTimeoutMins)
		Expect(err).ToNot(HaveOccurred())

		By("Start TNF " + tnfTestCaseName + " test case")
		err = globalhelper.LaunchTests(tnfTestCaseName, qeTcFileName)
		Expect(err).To(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(tnfTestCaseName, globalparameters.TestCaseFailed)
		Expect(err).ToNot(HaveOccurred())
	})

	// Negative #4.
	It("One deployment and one statefulset, both with one pod with one container, "+
		"only the deployment has terminationMessagePolicy set to FallbackToLogsOnError [negative]", func() {

		By("Create deployment in the cluster")
		deployment := tshelper.DefineDeploymentWithTerminationMsgPolicies(tsparams.TestDeploymentBaseName,
			randomNamespace, 1,
			[]corev1.TerminationMessagePolicy{corev1.TerminationMessageFallbackToLogsOnError})

		err := globalhelper.CreateAndWaitUntilDeploymentIsReady(deployment, tsparams.DeploymentDeployTimeoutMins)
		Expect(err).ToNot(HaveOccurred())

		By("Create statefulset in the cluster")
		statefulSet := tshelper.DefineStatefulSetWithTerminationMsgPolicies(tsparams.TestStatefulSetBaseName,
			randomNamespace, 1,
			[]corev1.TerminationMessagePolicy{tsparams.UseDefaultTerminationMsgPolicy})

		err = globalhelper.CreateAndWaitUntilStatefulSetIsReady(statefulSet, tsparams.StatefulSetDeployTimeoutMins)
		Expect(err).ToNot(HaveOccurred())

		By("Start TNF " + tnfTestCaseName + " test case")
		err = globalhelper.LaunchTests(tnfTestCaseName, qeTcFileName)
		Expect(err).To(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(tnfTestCaseName, globalparameters.TestCaseFailed)
		Expect(err).ToNot(HaveOccurred())
	})

	// Negative #5.
	It("One deployment and one daemonset, both with one pod with one container, "+
		"only the deployment has terminationMessagePolicy set to FallbackToLogsOnError [negative]", func() {

		By("Create deployment in the cluster")
		deployment := tshelper.DefineDeploymentWithTerminationMsgPolicies(tsparams.TestDeploymentBaseName,
			randomNamespace, 1,
			[]corev1.TerminationMessagePolicy{corev1.TerminationMessageFallbackToLogsOnError})

		err := globalhelper.CreateAndWaitUntilDeploymentIsReady(deployment, tsparams.DeploymentDeployTimeoutMins)
		Expect(err).ToNot(HaveOccurred())

		By("Create daemonset in the cluster")
		daemonSet := tshelper.DefineDaemonSetWithTerminationMsgPolicies(tsparams.TestDaemonSetBaseName,
			randomNamespace,
			[]corev1.TerminationMessagePolicy{tsparams.UseDefaultTerminationMsgPolicy})

		err = globalhelper.CreateAndWaitUntilDaemonSetIsReady(daemonSet, tsparams.DaemonSetDeployTimeoutMins)
		Expect(err).ToNot(HaveOccurred())

		By("Start TNF " + tnfTestCaseName + " test case")
		err = globalhelper.LaunchTests(tnfTestCaseName, qeTcFileName)
		Expect(err).To(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(tnfTestCaseName, globalparameters.TestCaseFailed)
		Expect(err).ToNot(HaveOccurred())
	})

	// Skip #1.
	It("One deployment with one pod and one container without TNF target labels [skip]", func() {

		By("Create deployment without TNF target labels in the cluster")
		deployment := tshelper.DefineDeploymentWithoutTargetLabels(tsparams.TestDeploymentBaseName, randomNamespace)

		err := globalhelper.CreateAndWaitUntilDeploymentIsReady(deployment, tsparams.DeploymentDeployTimeoutMins)
		Expect(err).ToNot(HaveOccurred())

		By("Start TNF " + tnfTestCaseName + " test case")
		err = globalhelper.LaunchTests(tnfTestCaseName, qeTcFileName)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(tnfTestCaseName, globalparameters.TestCaseSkipped)
		Expect(err).ToNot(HaveOccurred())
	})
})
