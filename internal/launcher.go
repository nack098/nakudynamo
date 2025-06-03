package internal

import (
	"os"
	"os/exec"
)

func Start(dynoEnv *DynamoEnvironment) (*exec.Cmd, error) {
	cmd := exec.Command(dynoEnv.JREPath, "-jar", dynoEnv.DynamoJarPath, "-inMemory")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Start(); err != nil {
		return nil, err
	}
	return cmd, nil
}

func StopDynamoDB(cmd *exec.Cmd) error {
	return cmd.Process.Signal(os.Interrupt)
}
