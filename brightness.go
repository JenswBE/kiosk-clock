package main

import (
	"errors"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
)

func readBacklightPercentage() (uint8, error) {
	cmd := exec.Command(
		"brightnessctl",
		"--list",
		"--class", "backlight",
		"--machine-readable",
	)

	output, err := cmd.Output()
	if err != nil {
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) {
			return 0, fmt.Errorf(
				"brightnessctl exited with code %d: %w",
				exitErr.ExitCode(),
				err,
			)
		}

		return 0, fmt.Errorf("failed to execute brightnessctl: %w", err)
	}

	line := strings.TrimSpace(string(output))
	if line == "" {
		return 0, fmt.Errorf("brightnessctl returned no backlight devices")
	}

	// --list can theoretically return multiple devices, so use the
	// first non-empty line.
	lines := strings.Split(line, "\n")
	var fields []string

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" {
			fields = strings.Split(line, ",")
			break
		}
	}

	if len(fields) < 5 {
		return 0, fmt.Errorf(
			"unexpected brightnessctl output %q",
			line,
		)
	}

	// Expected:
	//
	// amdgpu_bl1,backlight,28214,44%,64764
	//
	// Field 3 is the current brightness percentage.
	percentage := strings.TrimSpace(fields[3])
	percentage = strings.TrimSuffix(percentage, "%")

	value, err := strconv.ParseInt(percentage, 10, 64)
	if err != nil {
		return 0, fmt.Errorf(
			"invalid brightness percentage %q: %w",
			fields[3],
			err,
		)
	}

	if value < 0 || value > 100 {
		return 0, fmt.Errorf(
			"brightness percentage out of range: %d",
			value,
		)
	}

	return uint8(value), nil
}

func setBacklightBrightness(percentage uint8) error {
	if percentage > 100 {
		return fmt.Errorf(
			"brightness percentage out of range: %d (must be between 0 and 100)",
			percentage,
		)
	}

	value := fmt.Sprintf("%d%%", percentage)

	cmd := exec.Command(
		"brightnessctl",
		"set",
		value,
	)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf(
			"brightnessctl set %s failed: %w: %s",
			value,
			err,
			strings.TrimSpace(string(output)),
		)
	}

	return nil
}
