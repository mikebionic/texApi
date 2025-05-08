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
