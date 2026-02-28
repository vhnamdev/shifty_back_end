package constants

type SevertyType string

const (
	SeverityInfo     SevertyType = "INFO"
	SeverityMinor    SevertyType = "MINOR"
	SeverityMajor    SevertyType = "MAJOR"
	SeverityCritical SevertyType = "CRITICAL"
)
