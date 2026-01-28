package utils

func GetString(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}
