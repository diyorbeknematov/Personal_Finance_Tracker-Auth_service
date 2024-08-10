package helper

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"net/smtp"
)

func SendPasswordResetEmail(email string, code string) error {
	from := "sanjarbekabdurahmonov123@gmail.com"
	password := "j e h d i a p p c b a b q n x p"

	to := []string{
		email,
	}

	smtpHost := "smtp.gmail.com"
	smtpPort := "587"

	auth := smtp.PlainAuth("", from, password, smtpHost)

	t, err := template.ParseFiles("pkg/email/template.html")
	if err != nil {
		return fmt.Errorf("error parsing template: %v", err)
	}

	var body bytes.Buffer

	mimeHeaders := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	body.Write([]byte(fmt.Sprintf("Subject: Your verification code \n%s\n\n", mimeHeaders)))
	t.Execute(&body, struct {
		Passwd string
	}{

		Passwd: code,
	})

	err = smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, body.Bytes())
	if err != nil {
		return fmt.Errorf("error sending email: %v", err)
	}

	log.Println("Email sended to:", email)
	return nil
}
