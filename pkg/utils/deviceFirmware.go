package utils

import "strings"

func DetectDeviceFirmware(deviceFirmware string) string {
	if deviceFirmware == "" {
		return "web"
	}

	firmware := strings.ToLower(deviceFirmware)

	// Check for Android (handles cases like "Android 14", "android", etc.)
	if strings.Contains(firmware, "android") {
		return "android"
	}

	// Check for iOS (handles cases like "iOS 18.0.2", "ios", etc.)
	if strings.Contains(firmware, "ios") {
		return "ios"
	}

	return "web"
}
