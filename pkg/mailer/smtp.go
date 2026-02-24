package mailer

import (
	"bytes"
	"html/template"
	"path/filepath"

	"gopkg.in/gomail.v2"
)

type EmailSender interface {
	SendOTP(toEmail string, userName string, otp string) error
	SendInviteCode(toEmail, inviterName, resName, positionName, inviteCode string) error
}
type EmailService struct {
	Dialer *gomail.Dialer
	User   string
}

func NewGoMail(host string, port int, user, pass string) *EmailService {
	return &EmailService{
		Dialer: gomail.NewDialer(host, port, user, pass),
		User:   user,
	}
}

type OTPData struct {
	Name string
	OTP  string
}

func parseTemplate(templateFileName string, data interface{}) (string, error) {
	t, err := template.ParseFiles(filepath.Join("templates", templateFileName))
	if err != nil {
		return "", err
	}

	buf := new(bytes.Buffer)
	if err = t.Execute(buf, data); err != nil {
		return "", err
	}
	return buf.String(), nil
}

func (s *EmailService) SendOTP(toEmail, userName, otp string) error {
	data := OTPData{
		Name: userName,
		OTP:  otp,
	}

	bodyContent, err := parseTemplate("otp_mail.html", data)
	if err != nil {
		return err
	}

	m := gomail.NewMessage()
	m.SetHeader("From", s.User)
	m.SetHeader("To", toEmail)
	m.SetHeader("Subject", "Verify your account - Shifty")
	m.SetBody("text/html", bodyContent)

	return s.Dialer.DialAndSend(m)
}

type inviteTemplateData struct {
	InviterName    string
	RestaurantName string
	PositionName   string
	InviteCode     string
}

func (s *EmailService) SendInviteCode(toEmail, inviterName, resName, positionName, inviteCode string) error {
	data := inviteTemplateData{
		InviterName:    inviterName,
		RestaurantName: resName,
		PositionName:   positionName,
		InviteCode:     inviteCode,
	}

	bodyContent, err := parseTemplate("invite_template.html", data)

	if err != nil {
		return err
	}

	m := gomail.NewMessage()
	m.SetHeader("From", s.User)
	m.SetHeader("To", toEmail)
	m.SetHeader("Subject", "Restaurant invitation - Shifty")
	m.SetBody("text/html", bodyContent)

	return s.Dialer.DialAndSend(m)
}
