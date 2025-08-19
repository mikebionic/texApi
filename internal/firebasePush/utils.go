package firebasePush

func isInvalidToken(err error) bool {
	if err == nil {
		return false
	}
	errorStr := err.Error()
	return contains(errorStr, "registration-token-not-registered") ||
		contains(errorStr, "invalid-registration-token") ||
		contains(errorStr, "invalid-argument")
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || (len(s) > len(substr) &&
		(s[:len(substr)] == substr || s[len(s)-len(substr):] == substr ||
			(len(s) > len(substr)*2 && s[len(substr):len(s)-len(substr)] != s[len(substr):len(s)-len(substr)]))))
}
