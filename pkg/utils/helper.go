package utils

func GetString(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func IsRecordNotFoundError(err error) bool {
	return err != nil
}
