package tests

import (
	"fmt"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/test-network-function/cnfcert-tests-verification/tests/globalhelper"
	"github.com/test-network-function/cnfcert-tests-verification/tests/globalparameters"
	tshelper "github.com/test-network-function/cnfcert-tests-verification/tests/observability/helper"
	tsparams "github.com/test-network-function/cnfcert-tests-verification/tests/observability/parameters"

	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/deployment"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/namespaces"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/poddisruptionbudget"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/statefulset"
	"k8s.io/apimachinery/pkg/util/intstr"
)

var _ = Describe(tsparams.TnfPodDisruptionBudgetTcName, func() {

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

	const tnfTestCaseName = tsparams.TnfPodDisruptionBudgetTcName

	// 56635
	It("One deployment, pod disruption budget minAvailable value meet requirements", func() {
		qeTcFileName := globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText())

		By("Create deployment")
		dep := deployment.DefineDeployment(tsparams.TestDeploymentBaseName, randomNamespace,
			globalhelper.GetConfiguration().General.TestImage, tsparams.TnfTargetPodLabels)

		deployment.RedefineWithReplicaNumber(dep, 1)

		err := globalhelper.CreateAndWaitUntilDeploymentIsReady(dep, tsparams.DeploymentDeployTimeoutMins)
		Expect(err).ToNot(HaveOccurred())

		By("Create pod disruption budget")
		pdb := poddisruptionbudget.DefinePodDisruptionBudgetMinAvailable(tsparams.TestPdbBaseName, randomNamespace,
			intstr.FromInt(1), tsparams.TnfTargetPodLabels)

		err = globalhelper.CreatePodDisruptionBudget(pdb, tsparams.PdbDeployTimeoutMins)
		Expect(err).ToNot(HaveOccurred())

		By("Start TNF " + tnfTestCaseName + " test case")
		err = globalhelper.LaunchTests(tnfTestCaseName, qeTcFileName)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(tnfTestCaseName, globalparameters.TestCasePassed)
		Expect(err).ToNot(HaveOccurred())
	})

	// 56636
	It("One deployment, pod disruption budget maxUnavailable value meet requirements", func() {
		qeTcFileName := globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText())

		By("Create deployment")
		dep := deployment.DefineDeployment(tsparams.TestDeploymentBaseName, randomNamespace,
			globalhelper.GetConfiguration().General.TestImage, tsparams.TnfTargetPodLabels)

		deployment.RedefineWithReplicaNumber(dep, 2)

		err := globalhelper.CreateAndWaitUntilDeploymentIsReady(dep, tsparams.DeploymentDeployTimeoutMins)
		Expect(err).ToNot(HaveOccurred())

		By("Create pod disruption budget")
		pdb := poddisruptionbudget.DefinePodDisruptionBudgetMaxUnAvailable(tsparams.TestPdbBaseName, randomNamespace,
			intstr.FromInt(1), tsparams.TnfTargetPodLabels)

		err = globalhelper.CreatePodDisruptionBudget(pdb, tsparams.PdbDeployTimeoutMins)
		Expect(err).ToNot(HaveOccurred())

		By("Start TNF " + tnfTestCaseName + " test case")
		err = globalhelper.LaunchTests(tnfTestCaseName, qeTcFileName)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(tnfTestCaseName, globalparameters.TestCasePassed)
		Expect(err).ToNot(HaveOccurred())
	})

	// 56637
	It("One statefulSet, pod disruption budget minAvailable value is zero [negative]", func() {
		qeTcFileName := globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText())

		By("Create statefulSet")
		sf := statefulset.DefineStatefulSet(tsparams.TestStatefulSetBaseName, randomNamespace,
			globalhelper.GetConfiguration().General.TestImage, tsparams.TnfTargetPodLabels)

		statefulset.RedefineWithReplicaNumber(sf, 1)

		err := globalhelper.CreateAndWaitUntilStatefulSetIsReady(sf, tsparams.StatefulSetDeployTimeoutMins)
		Expect(err).ToNot(HaveOccurred())

		By("Create pod disruption budget")
		pdb := poddisruptionbudget.DefinePodDisruptionBudgetMinAvailable(tsparams.TestPdbBaseName, randomNamespace,
			intstr.FromInt(0), tsparams.TnfTargetPodLabels)

		err = globalhelper.CreatePodDisruptionBudget(pdb, tsparams.PdbDeployTimeoutMins)
		Expect(err).ToNot(HaveOccurred())

		By("Start TNF " + tnfTestCaseName + " test case")
		err = globalhelper.LaunchTests(tnfTestCaseName, qeTcFileName)
		Expect(err).To(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(tnfTestCaseName, globalparameters.TestCaseFailed)
		Expect(err).ToNot(HaveOccurred())
	})

	// 56638
	It("One deployment, pod disruption budget maxUnavailable equals to replica number [negative]", func() {
		qeTcFileName := globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText())

		By("Create deployment")
		dep := deployment.DefineDeployment(tsparams.TestDeploymentBaseName, randomNamespace,
			globalhelper.GetConfiguration().General.TestImage, tsparams.TnfTargetPodLabels)

		deployment.RedefineWithReplicaNumber(dep, 2)

		err := globalhelper.CreateAndWaitUntilDeploymentIsReady(dep, tsparams.DeploymentDeployTimeoutMins)
		Expect(err).ToNot(HaveOccurred())

		By("Create pod disruption budget")
		pdb := poddisruptionbudget.DefinePodDisruptionBudgetMaxUnAvailable(tsparams.TestPdbBaseName, randomNamespace,
			intstr.FromInt(2), tsparams.TnfTargetPodLabels)

		err = globalhelper.CreatePodDisruptionBudget(pdb, tsparams.PdbDeployTimeoutMins)
		Expect(err).ToNot(HaveOccurred())

		By("Start TNF " + tnfTestCaseName + " test case")
		err = globalhelper.LaunchTests(tnfTestCaseName, qeTcFileName)
		Expect(err).To(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(tnfTestCaseName, globalparameters.TestCaseFailed)
		Expect(err).ToNot(HaveOccurred())
	})

	// 56746
	It("One deployment, pod disruption budget maxUnavailable is bigger than the replica number [negative]", func() {
		qeTcFileName := globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText())

		By("Create deployment")
		dep := deployment.DefineDeployment(tsparams.TestDeploymentBaseName, randomNamespace,
			globalhelper.GetConfiguration().General.TestImage, tsparams.TnfTargetPodLabels)

		deployment.RedefineWithReplicaNumber(dep, 2)

		err := globalhelper.CreateAndWaitUntilDeploymentIsReady(dep, tsparams.DeploymentDeployTimeoutMins)
		Expect(err).ToNot(HaveOccurred())

		By("Create pod disruption budget")
		pdb := poddisruptionbudget.DefinePodDisruptionBudgetMaxUnAvailable(tsparams.TestPdbBaseName, randomNamespace,
			intstr.FromInt(3), tsparams.TnfTargetPodLabels)

		err = globalhelper.CreatePodDisruptionBudget(pdb, tsparams.PdbDeployTimeoutMins)
		Expect(err).ToNot(HaveOccurred())

		By("Start TNF " + tnfTestCaseName + " test case")
		err = globalhelper.LaunchTests(tnfTestCaseName, qeTcFileName)
		Expect(err).To(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(tnfTestCaseName, globalparameters.TestCaseFailed)
		Expect(err).ToNot(HaveOccurred())
	})
})
