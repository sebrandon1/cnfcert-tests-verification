package container

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"strings"
)

// SelectEngine check what container engine is present on the machine.
func SelectEngine() (string, error) {
	// If this is a LOCAL_RUN, check to see if the TNF_CONTAINER_CLIENT variable is also set and use that.
	if os.Getenv("LOCAL_RUN") == "true" {
		if os.Getenv("TNF_CONTAINER_CLIENT") != "" {
			return os.Getenv("TNF_CONTAINER_CLIENT"), nil
		}
		return "nil", fmt.Errorf("no container Engine present on host machine")
	}

	for _, containerEngine := range []string{"docker", "podman"} {
		containerEngineCMD := exec.Command(containerEngine)
		directoryName, _ := path.Split(containerEngineCMD.Path)

		if directoryName != "" {
			if strings.Contains(containerEngineCMD.String(), "docker") {
				err := validateDockerDaemonRunning()
				if err != nil {
					return "", err
				}
			}

			return containerEngine, nil
		}
	}

	return "nil", fmt.Errorf("no container Engine present on host machine")
}

func validateDockerDaemonRunning() error {
	isDaemonRunning := exec.Command("systemctl", "is-active", "--quiet", "docker")

	if isDaemonRunning.Run() != nil {
		return fmt.Errorf("docker daemon is not active on host")
	}

	return nil
}
