package mail

import (
	"fmt"
	"net/smtp"

	"github.com/jordan-wright/email"
)

const (
	smtpAuthAddress   = "smtp.gmail.com"
	smtpServerAddress = "smtp.gmail.com:587"
)

type EmailSender interface {
	Send(subject string, content string, to []string, cc []string, bcc []string, attachFiles []string) error
}

type GmailSender struct {
	name              string
	fromEmailAddr     string
	fromEmailPassword string
}

// Send implements EmailSender.
func (g *GmailSender) Send(
	subject string,
	content string,
	to []string,
	cc []string,
	bcc []string,
	attachFiles []string,
) error {
	letter := email.NewEmail()
	letter.From = fmt.Sprintf("%s <%s>", g.name, g.fromEmailAddr)
	letter.To = to
	letter.Subject = subject
	letter.HTML = []byte(content)
	letter.Cc = cc
	letter.Bcc = bcc
	for _, filename := range attachFiles {
		_, err := letter.AttachFile(filename)
		if err != nil {
			return fmt.Errorf("failed to attach file %s: %w", filename, err)
		}
	}

	smtpAuth := smtp.PlainAuth("", g.fromEmailAddr, g.fromEmailPassword, smtpAuthAddress)
	err := letter.Send(smtpServerAddress, smtpAuth)
	if err != nil {
		return fmt.Errorf("failed to send an email via gmail: %w", err)
	}

	return nil
}

func NewGmailSender(name, fromEmailAddr, fromEmailPassword string) EmailSender {
	return &GmailSender{name, fromEmailAddr, fromEmailPassword}
}
