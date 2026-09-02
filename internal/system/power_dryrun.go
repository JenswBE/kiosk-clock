//go:build dryrun

package system

import (
	"log"
)

// Shutdown powers off the machine via systemctl.
func Shutdown() error {
	log.Println("[OFFLINE MODE] Shutdown not executed")
	return nil
}

// Reboot restarts the machine via systemctl.
func Reboot() error {
	log.Println("[OFFLINE MODE] Reboot not executed")
	return nil
}
