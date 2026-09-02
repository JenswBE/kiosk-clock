//go:build offline

package system

// NTPSynced reports whether the system clock is synchronized via NTP,
// as reported by timedatectl.
func NTPSynced() bool {
	return false
}

// WifiConnected reports whether a WiFi device is currently connected
// to a network, as reported by NetworkManager.
func WifiConnected() (string, bool, error) {
	return "disconnected", false, nil
}
