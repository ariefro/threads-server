package mailer

import (
	"bytes"
	"fmt"
	"html/template"
	"net/smtp"

	"github.com/jordan-wright/email"
)

const (
	smtpAuthAddress   = "smtp.gmail.com"
	smtpServerAddress = "smtp.gmail.com:587"
)

type EmailSender interface {
	SendEmail(
		subject string,
		templateFile string,
		templateData interface{},
		to []string,
		cc []string,
		bcc []string,
		attachFiles []string,
	) error
}

type GmailSender struct {
	name              string
	fromEmailAddress  string
	fromEmailPassword string
}

func NewGmailSender(name, fromEmailAddress, fromEmailPassword string) EmailSender {
	return &GmailSender{
		name,
		fromEmailAddress,
		fromEmailPassword,
	}
}

func (sender *GmailSender) SendEmail(
	subject string,
	templateFile string,
	templateData interface{},
	to []string,
	cc []string,
	bcc []string,
	attachFiles []string,
) error {
	// Load dan render template
	tmpl, err := template.ParseFiles("internal/mailer/templates/" + templateFile)
	if err != nil {
		return fmt.Errorf("failed to parse template: %w", err)
	}

	var contentBuffer bytes.Buffer
	if err := tmpl.Execute(&contentBuffer, templateData); err != nil {
		return fmt.Errorf("failed to execute template: %w", err)
	}

	// Buat email
	e := email.NewEmail()
	e.From = fmt.Sprintf("%s <%s>", sender.name, sender.fromEmailAddress)
	e.Subject = subject
	e.HTML = contentBuffer.Bytes()
	e.To = to
	e.Cc = cc
	e.Bcc = bcc

	// Attach files jika ada
	for _, f := range attachFiles {
		if _, err := e.AttachFile(f); err != nil {
			return fmt.Errorf("failed to attach file %s: %w", f, err)
		}
	}

	// Kirim email menggunakan SMTP
	smtpAuth := smtp.PlainAuth("", sender.fromEmailAddress, sender.fromEmailPassword, smtpAuthAddress)
	return e.Send(smtpServerAddress, smtpAuth)
}
