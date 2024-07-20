package docker

import (
	"os"
	"os/exec"
)

// Unfortunately, we need this to pull an image, because we can't call pull method directly
// due to unexported parameters of the method, so it is impossible to specify the image outside the docker package.
// And I'd like to use this command because of its clear and informative output.
func runDockerCliCommand(args, env []string) error {
	cmd := exec.Command("docker", args...)
	cmd.Env = append(os.Environ(), env...)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	return cmd.Run()
}
