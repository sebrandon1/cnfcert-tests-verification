package globalhelper

import (
	"fmt"
	"time"

	"github.com/golang/glog"
	egiClients "github.com/openshift-kni/eco-goinfra/pkg/clients"
	egiDeployment "github.com/openshift-kni/eco-goinfra/pkg/deployment"
	appsv1 "k8s.io/api/apps/v1"

	. "github.com/onsi/gomega"
)

func IsDeploymentReady(namespace, deploymentName string) (bool, error) {
	dep, err := egiDeployment.Pull(egiClients.New(""), deploymentName, namespace)
	if err != nil {
		return false, err
	}

	return dep.IsReady(1 * time.Second), nil
}

// CreateAndWaitUntilDeploymentIsReady creates deployment and wait until all deployment replicas are up and running.
func CreateAndWaitUntilDeploymentIsReady(deployment *appsv1.Deployment,
	timeout time.Duration) error {
	return createAndWaitUntilDeploymentIsReady(deployment, timeout)
}

// createAndWaitUntilDeploymentIsReady creates deployment and wait until all deployment replicas are up and running.
func createAndWaitUntilDeploymentIsReady(deployment *appsv1.Deployment,
	timeout time.Duration) error {

	// create the deployment with eco-goinfra
	depBuilder, err := egiDeployment.NewBuilder(egiClients.New(""), deployment.Name, deployment.Namespace, deployment.Labels, &deployment.Spec.Template.Spec.Containers[0]).Create()
	if err != nil {
		return fmt.Errorf("failed to create deployment %q (ns %s): %w", deployment.Name, deployment.Namespace, err)
	}

	Eventually(func() bool {
		status, err := IsDeploymentReady(depBuilder.Object.Namespace, depBuilder.Object.Name)
		if err != nil {
			glog.V(5).Info(fmt.Sprintf(
				"deployment %s is not ready, retry in 5 seconds", depBuilder.Object.Name))

			return false
		}

		return status
	}, timeout, retryInterval*time.Second).Should(Equal(true), "Deployment is not ready")

	return nil
}

// GetRunningDeployment returns a running deployment.
func GetRunningDeployment(namespace, deploymentName string) (*appsv1.Deployment, error) {
	runningDeployment, err := egiDeployment.Pull(egiClients.New(""), deploymentName, namespace)
	if err != nil {
		return nil, fmt.Errorf("failed to get deployment %q (ns %s): %w", deploymentName, namespace, err)
	}

	return runningDeployment.Object, nil
}

func DeleteDeployment(name, namespace string) error {
	return egiDeployment.NewBuilder(egiClients.New(""), name, namespace, map[string]string{}, nil).Delete()
}
