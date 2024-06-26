package main

import (
	"bytes"
	"html/template"
	"log"
	"time"

	"github.com/vanng822/go-premailer/premailer"
	mail "github.com/xhit/go-simple-mail/v2"
)

// used when communicating with mail server
type Mail struct {
	Domain      string
	Host        string
	Port        int
	Username    string
	Password    string
	Encryption  string
	FromAddress string
	FromName    string
}

// individual email message model
type Message struct {
	From        string
	FromName    string
	To          string
	Subject     string
	Attachments []string
	Data        any
	DataMap     map[string]any
}

func (m *Mail) SendSMTPMessage(msg Message) error {
	if msg.From == "" {
		msg.From = m.FromAddress
	}

	if msg.FromName == "" {
		msg.FromName = m.FromName
	}

	data := map[string]any{
		"message": msg.Data,
	}

	msg.DataMap = data

	formattedMessage, err := m.buildHTMLMessage(msg)
	if err != nil {
		log.Printf("Failed to build HTML message: %s\n", err)
		return err
	}

	plainText, err := m.buildPlainTextMessage(msg)
	if err != nil {
		log.Println("Failed to build plain text message")
		return err
	}

	srv := mail.NewSMTPClient()
	srv.Host = m.Host
	srv.Port = m.Port
	srv.Username = m.Username
	srv.Password = m.Password
	srv.Encryption = m.getEncryption(m.Encryption)
	srv.KeepAlive = false
	srv.ConnectTimeout = 10 * time.Second
	srv.SendTimeout = 10 * time.Second

	smtpClient, err := srv.Connect()
	if err != nil {
		log.Println("Failed to connect to SMTP client")
		return err
	}

	email := mail.NewMSG()
	email.SetFrom(msg.From).
		SetSubject(msg.Subject).
		AddTo(msg.To)

	email.SetBody(mail.TextPlain, plainText)
	email.AddAlternative(mail.TextHTML, formattedMessage)

	if len(msg.Attachments) > 0 {
		for _, attachment := range msg.Attachments {
			email.AddAttachment(attachment)
		}
	}

	err = email.Send(smtpClient)
	if err != nil {
		log.Println("Failed to send email to SMTP client")
		return err
	}

	return nil
}

func (m *Mail) buildHTMLMessage(msg Message) (string, error) {
	templateToRender := "templates/mail.html.gohtml"
	t, err := template.New("email-html").ParseFiles(templateToRender)
	if err != nil {
		return "", err
	}

	var tpl bytes.Buffer
	if err = t.ExecuteTemplate(&tpl, "body", msg.DataMap); err != nil {
		return "", err
	}

	formattedMessage := tpl.String()
	formattedMessage, err = m.inlineCSS(formattedMessage)
	if err != nil {
		return "", err
	}

	return formattedMessage, nil
}

func (m *Mail) buildPlainTextMessage(msg Message) (string, error) {
	templateToRender := "templates/mail.plain.gohtml"
	t, err := template.New("email-plain").ParseFiles(templateToRender)
	if err != nil {
		return "", err
	}

	var tpl bytes.Buffer
	if err = t.ExecuteTemplate(&tpl, "body", msg.DataMap); err != nil {
		return "", err
	}

	plainMessage := tpl.String()

	return plainMessage, nil
}

func (m *Mail) inlineCSS(s string) (string, error) {
	opts := premailer.Options{
		RemoveClasses:     false,
		CssToAttributes:   false,
		KeepBangImportant: true,
	}

	prem, err := premailer.NewPremailerFromString(s, &opts)
	if err != nil {
		return "", err
	}

	html, err := prem.Transform()
	if err != nil {
		return "", err
	}

	return html, nil
}

func (m *Mail) getEncryption(s string) mail.Encryption {
	switch s {
	case "TLS":
		return mail.EncryptionSTARTTLS
	case "SSL":
		return mail.EncryptionSSLTLS
	case "None", "":
		return mail.EncryptionNone
	default:
		return mail.EncryptionSTARTTLS
	}
}
