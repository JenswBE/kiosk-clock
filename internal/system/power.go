//go:build !dryrun

package system

import (
	"fmt"
	"os/exec"
	"strings"
)

// Shutdown powers off the machine via systemctl.
func Shutdown() error {
	return runSystemctl("poweroff")
}

// Reboot restarts the machine via systemctl.
func Reboot() error {
	return runSystemctl("reboot")
}

func runSystemctl(action string) error {
	cmd := exec.Command("systemctl", action)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf(
			"systemctl %s failed: %w: %s",
			action,
			err,
			strings.TrimSpace(string(output)),
		)
	}

	return nil
}
