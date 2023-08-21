package tests

import (
	"fmt"

	"github.com/golang/glog"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/test-network-function/cnfcert-tests-verification/tests/globalhelper"
	"github.com/test-network-function/cnfcert-tests-verification/tests/globalparameters"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/config"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/namespaces"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/nodes"

	tshelper "github.com/test-network-function/cnfcert-tests-verification/tests/networking/helper"
	tsparams "github.com/test-network-function/cnfcert-tests-verification/tests/networking/parameters"
)

var _ = Describe("Networking network-policy-deny-all,", func() {

	configSuite, err := config.NewConfig()
	if err != nil {
		glog.Fatal(fmt.Errorf("can not load config file: %w", err))
	}

	var randomNamespace string
	var origReportDir string

	BeforeEach(func() {
		randomNamespace = tsparams.TestNetworkingNameSpace + "-" + globalhelper.GenerateRandomString(10)

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
			[]string{tsparams.TestPodLabel},
			[]string{},
			[]string{},
			[]string{})
		Expect(err).ToNot(HaveOccurred())

		By("Ensure all nodes are labeled with 'worker-cnf' label")
		err = nodes.EnsureAllNodesAreLabeled(globalhelper.GetAPIClient().CoreV1Interface, configSuite.General.CnfNodeLabel)
		Expect(err).ToNot(HaveOccurred())
	})

	AfterEach(func() {
		By("Remove networking test namespaces")
		err := namespaces.DeleteAndWait(
			globalhelper.GetAPIClient().CoreV1Interface,
			randomNamespace,
			tsparams.WaitingTime,
		)
		Expect(err).ToNot(HaveOccurred())

		By("Restore default report directory")
		globalhelper.GetConfiguration().General.TnfReportDir = origReportDir
	})

	// 59740
	It("one deployment, one pod in a namespace with deny all ingress and egress network policy", func() {

		By("Define deployment and create it on cluster")
		err := tshelper.DefineAndCreateDeploymentOnCluster(1, randomNamespace)
		Expect(err).ToNot(HaveOccurred())

		By("Define and create network policy")
		err = tshelper.DefineAndCreateNetworkPolicy("netpolicy1",
			randomNamespace, []string{"Ingress", "Egress"}, tsparams.TestDeploymentLabels)
		Expect(err).ToNot(HaveOccurred())

		By("Start tests")
		err = globalhelper.LaunchTests(
			tsparams.TnfNetworkPolicyDenyAllTcName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()))
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(
			tsparams.TnfNetworkPolicyDenyAllTcName,
			globalparameters.TestCasePassed)
		Expect(err).ToNot(HaveOccurred())

	})

	// 59741
	It("one deployment, one pod in a namespace with only deny all ingress network policy [negative]", func() {

		By("Define deployment and create it on cluster")
		err := tshelper.DefineAndCreateDeploymentOnCluster(1, randomNamespace)
		Expect(err).ToNot(HaveOccurred())

		By("Define and create network policy")
		err = tshelper.DefineAndCreateNetworkPolicy("netpolicy1",
			randomNamespace, []string{"Ingress"}, tsparams.TestDeploymentLabels)
		Expect(err).ToNot(HaveOccurred())

		By("Start tests")
		err = globalhelper.LaunchTests(
			tsparams.TnfNetworkPolicyDenyAllTcName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()))
		Expect(err).To(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(
			tsparams.TnfNetworkPolicyDenyAllTcName,
			globalparameters.TestCaseFailed)
		Expect(err).ToNot(HaveOccurred())

	})

	// 59742
	It("one deployment, one pod in a namespace with only deny all egress network policy [negative]", func() {

		By("Define deployment and create it on cluster")
		err := tshelper.DefineAndCreateDeploymentOnCluster(1, randomNamespace)
		Expect(err).ToNot(HaveOccurred())

		By("Define and create network policy")
		err = tshelper.DefineAndCreateNetworkPolicy("netpolicy1",
			randomNamespace, []string{"Egress"}, tsparams.TestDeploymentLabels)
		Expect(err).ToNot(HaveOccurred())

		By("Start tests")
		err = globalhelper.LaunchTests(
			tsparams.TnfNetworkPolicyDenyAllTcName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()))
		Expect(err).To(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(
			tsparams.TnfNetworkPolicyDenyAllTcName,
			globalparameters.TestCaseFailed)
		Expect(err).ToNot(HaveOccurred())
	})

	// 59743
	It("one deployment, one pod in a namespace with neither deny all ingress or egress network policy [negative]", func() {

		By("Define deployment and create it on cluster")
		err := tshelper.DefineAndCreateDeploymentOnCluster(1, randomNamespace)
		Expect(err).ToNot(HaveOccurred())

		By("Start tests")
		err = globalhelper.LaunchTests(
			tsparams.TnfNetworkPolicyDenyAllTcName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()))
		Expect(err).To(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(
			tsparams.TnfNetworkPolicyDenyAllTcName,
			globalparameters.TestCaseFailed)
		Expect(err).ToNot(HaveOccurred())

	})

	// 59744
	It("two deployments in different namespaces, one pod each, namespaces have deny all ingress and egress network policy",
		func() {

			By("Define first deployment and create it on cluster")
			err := tshelper.DefineAndCreateDeploymentOnCluster(1, randomNamespace)
			Expect(err).ToNot(HaveOccurred())

			By("Define and create first network policy")
			err = tshelper.DefineAndCreateNetworkPolicy("netpolicy1",
				randomNamespace, []string{"Ingress", "Egress"}, tsparams.TestDeploymentLabels)
			Expect(err).ToNot(HaveOccurred())

			By("Create additional namespaces for testing")
			randomSecondaryNamespace := tsparams.AdditionalNetworkingNamespace + "-" + globalhelper.GenerateRandomString(5)
			err = namespaces.Create(randomSecondaryNamespace, globalhelper.GetAPIClient())
			Expect(err).ToNot(HaveOccurred())

			By("Define second deployment and create it on cluster")
			err = tshelper.DefineAndCreateDeploymentWithNamespace(randomSecondaryNamespace, 1)
			Expect(err).ToNot(HaveOccurred())

			By("Define and create second network policy")
			err = tshelper.DefineAndCreateNetworkPolicy("netpolicy2",
				randomSecondaryNamespace, []string{"Ingress", "Egress"}, tsparams.TestDeploymentLabels)
			Expect(err).ToNot(HaveOccurred())

			By("Define TNF config file")
			err = globalhelper.DefineTnfConfig(
				[]string{randomNamespace, randomSecondaryNamespace},
				[]string{tsparams.TestPodLabel},
				[]string{},
				[]string{},
				[]string{})
			Expect(err).ToNot(HaveOccurred())

			By("Start tests")
			err = globalhelper.LaunchTests(
				tsparams.TnfNetworkPolicyDenyAllTcName,
				globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()))
			Expect(err).ToNot(HaveOccurred())

			By("Verify test case status in Junit and Claim reports")
			err = globalhelper.ValidateIfReportsAreValid(
				tsparams.TnfNetworkPolicyDenyAllTcName,
				globalparameters.TestCasePassed)
			Expect(err).ToNot(HaveOccurred())

			By("Delete additional namespaces")
			err = namespaces.DeleteAndWait(
				globalhelper.GetAPIClient().CoreV1Interface,
				randomSecondaryNamespace,
				tsparams.WaitingTime,
			)
			Expect(err).ToNot(HaveOccurred())
		})

	// 59745
	It("two deployments in different namespaces, one pod each, one namespace has only deny all egress network policy [negative]",
		func() {

			By("Define first deployment and create it on cluster")
			err := tshelper.DefineAndCreateDeploymentOnCluster(1, randomNamespace)
			Expect(err).ToNot(HaveOccurred())

			By("Define and create first network policy")
			err = tshelper.DefineAndCreateNetworkPolicy("netpolicy1",
				randomNamespace, []string{"Ingress", "Egress"}, tsparams.TestDeploymentLabels)
			Expect(err).ToNot(HaveOccurred())

			By("Create additional namespaces for testing")
			randomSecondaryNamespace := tsparams.AdditionalNetworkingNamespace + "-" + globalhelper.GenerateRandomString(5)
			err = namespaces.Create(randomSecondaryNamespace, globalhelper.GetAPIClient())
			Expect(err).ToNot(HaveOccurred())

			By("Define second deployment and create it on cluster")
			err = tshelper.DefineAndCreateDeploymentWithNamespace(randomSecondaryNamespace, 1)
			Expect(err).ToNot(HaveOccurred())

			By("Define and create second network policy")
			err = tshelper.DefineAndCreateNetworkPolicy("netpolicy2",
				randomSecondaryNamespace, []string{"Egress"}, tsparams.TestDeploymentLabels)
			Expect(err).ToNot(HaveOccurred())

			By("Start tests")
			err = globalhelper.LaunchTests(
				tsparams.TnfNetworkPolicyDenyAllTcName,
				globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()))
			Expect(err).To(HaveOccurred())

			By("Verify test case status in Junit and Claim reports")
			err = globalhelper.ValidateIfReportsAreValid(
				tsparams.TnfNetworkPolicyDenyAllTcName,
				globalparameters.TestCaseFailed)
			Expect(err).ToNot(HaveOccurred())

			By("Delete additional namespaces")
			err = namespaces.DeleteAndWait(
				globalhelper.GetAPIClient().CoreV1Interface,
				randomSecondaryNamespace,
				tsparams.WaitingTime,
			)
			Expect(err).ToNot(HaveOccurred())
		})
})
