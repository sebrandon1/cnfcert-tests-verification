package parameters

import (
	"fmt"
	"time"
)

const (
	WaitingTime   = 5 * time.Minute
	RetryInterval = 5
)

var (
	TestNetworkingNameSpace       = "networking-tests"
	AdditionalNetworkingNamespace = "net-tests"
	testPodLabelPrefixName        = "networking-test/test"
	testPodLabelValue             = "testing"
	TestPodLabel                  = fmt.Sprintf("%s: %s", testPodLabelPrefixName, testPodLabelValue)
	TestDeploymentLabels          = map[string]string{
		testPodLabelPrefixName: testPodLabelValue,
		"app":                  "networkingput"}
	TestNadNameA                                 = "networking-nada"
	TestIPamIPNetworkA                           = "10.255.255.0/25"
	TestDeploymentAName                          = "networkingputa"
	TestNadNameB                                 = "networking-nadb"
	TestIPamIPNetworkB                           = "10.255.128.0/25"
	TestDeploymentBName                          = "networkingputb"
	CertsuiteDefaultNetworkTcName                = "networking-icmpv4-connectivity"
	CertsuiteMultusIpv4TcName                    = "networking-icmpv4-connectivity-multus"
	CertsuiteNodePortTcName                      = "access-control-service-type"
	CertsuiteNetworkPolicyDenyAllTcName          = "networking-network-policy-deny-all"
	CertsuiteOcpReservedPortsUsageTcName         = "networking-ocp-reserved-ports-usage"
	CertsuiteUndeclaredContainerPortsUsageTcName = "networking-undeclared-container-ports-usage"
	CertsuiteReservedPartnerPortsTcName          = "networking-reserved-partner-ports"
	CertsuiteDualStackServiceTcName              = "networking-dual-stack-service"
	CertsuiteDpdkCPUPinningExecProbe             = "networking-dpdk-cpu-pinning-exec-probe"

	DpdkPodName                   = "dpdk-pod"
	NetworkingTestSkipLabel       = map[string]string{"redhat-best-practices-for-k8s.com/skip_connectivity_tests": ""}
	NetworkingTestMultusSkipLabel = map[string]string{"redhat-best-practices-for-k8s.com/skip_multus_connectivity_tests": ""}
)

type IPOutputInterface struct {
	IfIndex   uint     `json:"ifindex"`
	IfName    string   `json:"ifname"`
	Flags     []string `json:"flags"`
	Mtu       uint     `json:"mtu"`
	Qdisc     string   `json:"qdisc"`
	Master    string   `json:"master"`
	Operstate string   `json:"operstate"`
	Linkmode  string   `json:"linkmode"`
	Group     string   `json:"group"`
	Txqlen    int      `json:"txqlen"`
	LinkType  string   `json:"link_type"`
	Address   string   `json:"address"`
	Broadcast string   `json:"broadcast"`
}
