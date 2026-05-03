package sendto

import (
	"bytes"
	"fmt"
	"go-learning/global"
	"html/template"
	"net/smtp"
	"strings"

	"go.uber.org/zap"
)

/*
1. Vào https://myaccount.google.com/security
2. Bật 2-Step Verification (Xác minh 2 bước dùng SMS hoặc Google Authenticator đều được)
3. Vào https://myaccount.google.com/apppasswords -> Không hỗ trợ
- Tài khoản Google Workspace (email công ty)
- Tài khoản do trường học cấp
- Tài khoản bị admin quản lý
4. Tạo App Mail
5. Password trả về dạng "abcd efgh ijkl mnop" -> SMTPPassword = "abcdefghijklmnop"
*/

const (
	SMTPHost     = "smtp.gmail.com"
	SMTPPort     = "587"
	SMTPUsername = "example@gmail.com"
	SMTPPassword = ""
)

type EmailAddress struct {
	Address string `json:"address"`
	Name    string `json:"name"`
}

type Mail struct {
	From    EmailAddress
	To      []string
	Subject string
	Body    string
}

func SendTextEmailOtp(to []string, from string, otp string) error {
	contentEmail := Mail{
		From: EmailAddress{
			Address: from,
			Name:    "Admin",
		},
		To:      to,
		Subject: "Go-learning OTP Verification",
		Body:    fmt.Sprintf("Your OTP is %s. Please enter it to verify it to your account.", otp),
	}

	messageEmail := buildMessage(contentEmail)

	// send smtp
	auth := smtp.PlainAuth("", SMTPUsername, SMTPPassword, SMTPHost)
	err := smtp.SendMail(SMTPHost+":"+SMTPPort, auth, from, to, []byte(messageEmail))
	if err != nil {
		global.Logger.Error("Email send failed::", zap.Error(err))
		return err
	}

	return nil
}

func SendTemplateEmailOtp(to []string, from string, templateName string, templateData map[string]interface{}) error {
	htmlBody, err := getMailTemplate(templateName, templateData)
	if err != nil {
		return err
	}

	return send(to, from, htmlBody)
}

func getMailTemplate(templateName string, templateData map[string]interface{}) (string, error) {
	htmlTemplate := new(bytes.Buffer)
	t := template.Must(template.New(templateName).ParseFiles("templates-email/" + templateName))
	err := t.Execute(htmlTemplate, templateData)
	if err != nil {
		return "", err
	}

	return htmlTemplate.String(), nil
}

func buildMessage(mail Mail) string {
	msg := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\r\n"
	msg += fmt.Sprintf("From: %s\r\n", mail.From.Address)
	msg += fmt.Sprintf("To: %s\r\n", strings.Join(mail.To, ";"))
	msg += fmt.Sprintf("Subject: %s\r\n", mail.Subject)
	msg += fmt.Sprintf("\r\n%s\r\n", mail.Body)

	return msg
}

func send(to []string, from string, htmlBody string) error {
	contentEmail := Mail{
		From: EmailAddress{
			Address: from,
			Name:    "Admin",
		},
		To:      to,
		Subject: "Go-learning OTP Verification",
		Body:    htmlBody,
	}

	messageEmail := buildMessage(contentEmail)

	// send smtp
	auth := smtp.PlainAuth("", SMTPUsername, SMTPPassword, SMTPHost)
	err := smtp.SendMail(SMTPHost+":"+SMTPPort, auth, from, to, []byte(messageEmail))
	if err != nil {
		global.Logger.Error("Email send failed::", zap.Error(err))
		return err
	}

	return nil
}
