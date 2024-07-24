package config

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"html/template"
	"net/smtp"
	"strconv"

	"github.com/shordem/api.thryvo/lib/constants"
)

var SenderName = "Mazimart"

type EmailInterface interface {
	SendWithTemplate(to, subject, templateFile string, data interface{}) error
}

type email struct {
	Host, Username, Password, From string
	Port                           int
}

func NewEmail(env constants.Env) EmailInterface {
	port, _ := strconv.Atoi(env.SMTP_PORT)

	return &email{
		Host:     env.SMTP_HOST,
		Port:     port,
		Username: env.SMTP_USERNAME,
		Password: env.SMTP_PASSWORD,
		From:     env.FROM_EMAIL,
	}
}

func (e *email) Send(to, subject, body string) error {
	auth := smtp.PlainAuth("", e.Username, e.Password, e.Host)
	msg := []byte("From: " + SenderName + " <" + e.From + ">\r\n" +
		"To: " + to + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"MIME-Version: 1.0\r\n" +
		"Content-Type: text/html; charset=utf-8\r\n" +
		"\r\n" +
		body)
	addr := fmt.Sprintf("%s:%d", e.Host, e.Port)

	conn, err := tls.Dial("tcp", addr, &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         e.Host,
	})
	if err != nil {
		return err
	}

	client, err := smtp.NewClient(conn, e.Host)
	if err != nil {
		return err
	}

	// Authenticating
	if err = client.Auth(auth); err != nil {
		return err
	}

	// Setting the sender and recipient
	if err = client.Mail(e.From); err != nil {
		return err
	}

	if err = client.Rcpt(to); err != nil {
		return err
	}

	// Sending the email body
	wc, err := client.Data()
	if err != nil {
		return err
	}

	_, err = wc.Write([]byte(msg))
	if err != nil {
		return err
	}

	if err = wc.Close(); err != nil {
		return err
	}

	return client.Quit()
}

func (e *email) ParseTemplate(templateFile string, data interface{}) (string, error) {
	t, err := template.ParseFiles(templateFile, "templates/layout.html")
	if err != nil {
		return "", err
	}

	buf := new(bytes.Buffer)

	if err = t.ExecuteTemplate(buf, "layout", data); err != nil {
		return "", err
	}
	return buf.String(), nil
}

func (e *email) SendWithTemplate(to, subject, templateFile string, data interface{}) error {
	body, err := e.ParseTemplate(templateFile, data)
	if err != nil {
		return err
	}

	return e.Send(to, subject, body)
}
