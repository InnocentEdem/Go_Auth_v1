package utils

import (
	"bytes"
	"fmt"
	"html/template"
	"os"

	"github.com/InnocentEdem/Go_Auth_v1/models"
	"gopkg.in/gomail.v2"
)

type FeatureRequest struct {
	FeatureName        string
	FeatureDescription string
	SenderName         string
	SenderEmail        string
}

func SendFeatureRequestEmail(request FeatureRequest) error {
	from := os.Getenv("EMAIL_ADDRESS")
	password := os.Getenv("EMAIL_PASSWORD")

	tmpl, err := template.ParseFiles("templates/feature_request_email.html")
	if err != nil {
		return err
	}

	var body bytes.Buffer
	if err := tmpl.Execute(&body, request); err != nil {
		return err
	}

	mailer := gomail.NewMessage()
	mailer.SetHeader("From", from)
	mailer.SetHeader("To", from)
	mailer.SetHeader("Subject", "New Feature Request")
	mailer.SetBody("text/html", body.String())

	dialer := gomail.NewDialer("smtp.gmail.com", 587, from, password)
	dialer.TLSConfig = nil

	if err := dialer.DialAndSend(mailer); err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}
func SendConfirmationEmail(client models.Client, user models.ClientAppUser) error {
	from := os.Getenv("EMAIL_ADDRESS")
	password := os.Getenv("EMAIL_PASSWORD")

	tmpl, err := template.ParseFiles("templates/user_confirmation_code.html")
	if err != nil {
		return err
	}

	var body bytes.Buffer
	if err := tmpl.Execute(&body, user); err != nil {
		return err
	}

	mailer := gomail.NewMessage()
	mailer.SetHeader("From", from)
	mailer.SetHeader("To", from)
	mailer.SetHeader("Subject", "New Feature Request")
	mailer.SetBody("text/html", body.String())

	dialer := gomail.NewDialer("smtp.gmail.com", 587, from, password)
	dialer.TLSConfig = nil

	if err := dialer.DialAndSend(mailer); err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}
