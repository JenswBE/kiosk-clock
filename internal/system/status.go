//go:build !offline

package system

import (
	"fmt"
	"os/exec"
	"strings"
)

// NTPSynced reports whether the system clock is synchronized via NTP,
// as reported by timedatectl.
func NTPSynced() bool {
	cmd := exec.Command(
		"timedatectl",
		"show",
		"-p", "NTPSynchronized",
		"--value",
	)

	out, err := cmd.Output()
	if err != nil {
		return false
	}

	return parseNTPSynchronized(string(out))
}

func parseNTPSynchronized(output string) bool {
	return strings.TrimSpace(output) == "yes"
}

// WifiConnected reports whether a WiFi device is currently connected
// to a network, as reported by NetworkManager.
func WifiConnected() (string, bool, error) {
	cmd := exec.Command(
		"nmcli",
		"-t",
		"-f", "TYPE,STATE",
		"device", "status",
	)

	out, err := cmd.Output()
	if err != nil {
		return "", false, fmt.Errorf("failed to check wifi status: %w", err)
	}

	status, ok := parseWifiConnected(string(out))
	return status, ok, nil
}

func parseWifiConnected(output string) (string, bool) {
	for _, line := range strings.Split(strings.TrimSpace(output), "\n") {
		fields := strings.SplitN(strings.TrimSpace(line), ":", 2)
		if len(fields) != 2 {
			continue
		}

		if fields[0] == "wifi" {
			isOK := fields[1] == "connected"
			return fields[1], isOK
		}
	}

	return "", false
}
