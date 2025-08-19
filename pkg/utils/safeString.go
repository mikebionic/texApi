package utils

func SafeString(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func SafeInt(i *int) int {
	if i == nil {
		return 0
	}
	return *i
}

var (
	EmptyString = ""
)

func SafeBool(i *bool) bool {
	if i == nil {
		return false
	}
	return *i
}

func BoolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}
