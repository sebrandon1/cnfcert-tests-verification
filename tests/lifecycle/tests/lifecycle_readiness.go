package tests

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/test-network-function/cnfcert-tests-verification/tests/globalhelper"
	"github.com/test-network-function/cnfcert-tests-verification/tests/globalparameters"
	"github.com/test-network-function/cnfcert-tests-verification/tests/lifecycle/lifehelper"
	"github.com/test-network-function/cnfcert-tests-verification/tests/lifecycle/lifeparameters"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/daemonset"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/deployment"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/namespaces"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/pod"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/statefulset"
)

var _ = Describe("lifecycle-readiness", func() {

	stringOfSkipTc := globalhelper.GetStringOfSkipTcs(lifeparameters.TnfTestCases,
		lifeparameters.TnfReadinessTcName)

	BeforeEach(func() {
		err := lifehelper.WaitUntilClusterIsStable()
		Expect(err).ToNot(HaveOccurred())

		By("Clean namespace before each test")
		err = namespaces.Clean(lifeparameters.LifecycleNamespace, globalhelper.APIClient)
		Expect(err).ToNot(HaveOccurred())
	})

	// 50145
	It("One deployment, one pod with a readiness probe", func() {
		By("Define deployment with a readiness probe")
		deployment := deployment.RedefineWithReadinessProbe(
			lifehelper.DefineDeployment(1, 1, "lifecycledp"))
		err := globalhelper.CreateAndWaitUntilDeploymentIsReady(deployment, lifeparameters.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Start lifecycle-readiness test")
		err = globalhelper.LaunchTests(
			lifeparameters.LifecycleTestSuiteName,
			globalhelper.ConvertSpecNameToFileName(CurrentGinkgoTestDescription().FullTestText),
			stringOfSkipTc)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(
			lifeparameters.TnfReadinessTcName,
			globalparameters.TestCasePassed)
		Expect(err).ToNot(HaveOccurred())
	})

	// 50146
	It("Two deployments, multiple pods each, all have a readiness probe", func() {
		By("Define first deployment with a readiness probe")
		deploymenta := deployment.RedefineWithReadinessProbe(
			lifehelper.DefineDeployment(3, 1, "lifecycledpa"))
		err := globalhelper.CreateAndWaitUntilDeploymentIsReady(deploymenta, lifeparameters.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Define second deployment with a readiness probe")
		deploymentb := deployment.RedefineWithReadinessProbe(
			lifehelper.DefineDeployment(3, 1, "lifecycledpb"))
		err = globalhelper.CreateAndWaitUntilDeploymentIsReady(deploymentb, lifeparameters.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Start lifecycle-readiness test")
		err = globalhelper.LaunchTests(
			lifeparameters.LifecycleTestSuiteName,
			globalhelper.ConvertSpecNameToFileName(CurrentGinkgoTestDescription().FullTestText),
			stringOfSkipTc)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(
			lifeparameters.TnfReadinessTcName,
			globalparameters.TestCasePassed)
		Expect(err).ToNot(HaveOccurred())
	})

	// 50147
	It("One statefulSet, one pod with a readiness probe", func() {
		By("Define statefulSet with a readiness probe")
		statefulset := statefulset.RedefineWithReadinessProbe(
			lifehelper.DefineStatefulSet("lifecycle-sf"))
		err := lifehelper.CreateAndWaitUntilStatefulSetIsReady(statefulset, lifeparameters.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Start lifecycle-readiness test")
		err = globalhelper.LaunchTests(
			lifeparameters.LifecycleTestSuiteName,
			globalhelper.ConvertSpecNameToFileName(CurrentGinkgoTestDescription().FullTestText),
			stringOfSkipTc)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(
			lifeparameters.TnfReadinessTcName,
			globalparameters.TestCasePassed)
		Expect(err).ToNot(HaveOccurred())
	})

	// 50148
	It("One pod with a readiness probe", func() {
		By("Define pod with a readiness probe")
		pod := pod.RedefineWithReadinessProbe(pod.RedefinePodWithLabel(
			lifehelper.DefinePod("lifecycleput"), lifeparameters.TestDeploymentLabels))
		err := lifehelper.CreateAndWaitUntilPodIsReady(pod, lifeparameters.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Start lifecycle-readiness test")
		err = globalhelper.LaunchTests(
			lifeparameters.LifecycleTestSuiteName,
			globalhelper.ConvertSpecNameToFileName(CurrentGinkgoTestDescription().FullTestText),
			stringOfSkipTc)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(
			lifeparameters.TnfReadinessTcName,
			globalparameters.TestCasePassed)
		Expect(err).ToNot(HaveOccurred())
	})

	// 50149
	It("One daemonSet without a readiness probe [negative]", func() {
		By("Define daemonSet without a readiness probe")
		daemonSet := daemonset.DefineDaemonSet(lifeparameters.LifecycleNamespace,
			globalhelper.Configuration.General.TnfImage,
			lifeparameters.TestDeploymentLabels, "lifecycleds")
		err := globalhelper.CreateAndWaitUntilDaemonSetIsReady(daemonSet, lifeparameters.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Start lifecycle-readiness test")
		err = globalhelper.LaunchTests(
			lifeparameters.LifecycleTestSuiteName,
			globalhelper.ConvertSpecNameToFileName(CurrentGinkgoTestDescription().FullTestText),
			stringOfSkipTc)
		Expect(err).To(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(
			lifeparameters.TnfReadinessTcName,
			globalparameters.TestCaseFailed)
		Expect(err).ToNot(HaveOccurred())
	})

	// 50150
	It("Two deployments, one pod each, one without a readiness probe [negative]", func() {
		By("Define first deployment with a readiness probe")
		deploymenta := deployment.RedefineWithReadinessProbe(
			lifehelper.DefineDeployment(1, 1, "lifecycledpa"))
		err := globalhelper.CreateAndWaitUntilDeploymentIsReady(deploymenta, lifeparameters.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Define second deployment without a readiness probe")
		deploymentb := lifehelper.DefineDeployment(1, 1, "lifecycledpb")
		err = globalhelper.CreateAndWaitUntilDeploymentIsReady(deploymentb, lifeparameters.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Start lifecycle-readiness test")
		err = globalhelper.LaunchTests(
			lifeparameters.LifecycleTestSuiteName,
			globalhelper.ConvertSpecNameToFileName(CurrentGinkgoTestDescription().FullTestText),
			stringOfSkipTc)
		Expect(err).To(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(
			lifeparameters.TnfReadinessTcName,
			globalparameters.TestCaseFailed)
		Expect(err).ToNot(HaveOccurred())
	})
})