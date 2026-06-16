package unit_test

import "github.com/stretchr/testify/mock"

// ===== MOCK MAILER =====
type MockMailer struct{ mock.Mock }

func (m *MockMailer) SendOTP(email, name, otp string) error {
	return m.Called(email, name, otp).Error(0)
}
func (m *MockMailer) SendInviteCode(toEmail, inviterName, resName, positionName, inviteCode string) error {
	return m.Called(toEmail, inviterName, resName, positionName, inviteCode).Error(0)
}
