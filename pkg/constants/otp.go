package constants

type OTPPurpose string

const (
	PurposeRegister      OTPPurpose = "register"
	PurposeResetPassword OTPPurpose = "reset"
)
