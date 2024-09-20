package utils

import (
	"bytes"
	"fmt"
	"github.com/arifai/zenith/config"
	"html/template"
	"net/smtp"
	"sync"
)

// Mailer is an interface for sending emails.
type Mailer interface {
	// SendMail sends an email to the specified recipients with the given subject and body using the configured SMTP server.
	SendMail(to []string, subject string, body string) error

	// SendMailWithTemplate sends an email to the specified recipients using an HTML template.
	SendMailWithTemplate(to []string, subject string, templateFileName string, data interface{}) error

	// QueueMail enqueues an email to be sent later by a worker.
	QueueMail(to []string, subject string, body string)

	// QueueMailWithTemplate enqueues an email with a template to be sent later by a worker.
	QueueMailWithTemplate(to []string, subject string, templateFileName string, data interface{})

	// Worker processes email requests from the queue, sending each email using the configured SMTP server.
	Worker()

	// Shutdown signals the mailer to stop processing new email requests and waits until all current tasks are completed.
	Shutdown()
}

// MailerImpl provides functionality for sending emails using SMTP. It is initialized with SMTP server configuration details.
type MailerImpl struct {
	config  config.SMTPConfig
	queue   chan emailRequest
	workers int
	wg      sync.WaitGroup
}

// emailRequest represents a request to send an email, including recipient addresses, subject, and email body content.
type emailRequest struct {
	to               []string
	subject          string
	body             string
	templateFileName string
	data             interface{}
}

// NewMailer creates a new MailerImpl instance with the provided SMTPConfig, queue size, and number of worker routines.
func NewMailer(config config.SMTPConfig, queueSize int, workers int) *MailerImpl {
	mailer := &MailerImpl{
		config:  config,
		queue:   make(chan emailRequest, queueSize),
		workers: workers,
	}

	for i := 0; i < workers; i++ {
		go mailer.Worker()
	}

	return mailer
}

func (m *MailerImpl) SendMail(to []string, subject string, body string) error {
	auth := smtp.PlainAuth("", m.config.Username, m.config.Password, m.config.Host)
	msg := "From: " + m.config.Username + "\n" +
		"To: " + fmt.Sprintf("%s", to) + "\n" +
		"Subject: " + subject + "\n\n" +
		body
	return smtp.SendMail(fmt.Sprintf("%s:%d", m.config.Host, m.config.Port), auth, m.config.Username, to, []byte(msg))
}

func (m *MailerImpl) SendMailWithTemplate(to []string, subject string, templateFileName string, data interface{}) error {
	tmpl, err := template.ParseFiles(templateFileName)
	if err != nil {
		return err
	}

	var body bytes.Buffer
	if err := tmpl.Execute(&body, data); err != nil {
		return err
	}

	return m.SendMail(to, subject, body.String())
}

func (m *MailerImpl) QueueMail(to []string, subject string, body string) {
	m.queue <- emailRequest{to: to, subject: subject, body: body}
}

func (m *MailerImpl) QueueMailWithTemplate(to []string, subject string, templateFileName string, data interface{}) {
	m.queue <- emailRequest{to: to, subject: subject, templateFileName: templateFileName, data: data}
}

func (m *MailerImpl) Worker() {
	m.wg.Add(1)
	defer m.wg.Done()
	for email := range m.queue {
		var err error
		if email.templateFileName != "" {
			err = m.SendMailWithTemplate(email.to, email.subject, email.templateFileName, email.data)
		} else {
			err = m.SendMail(email.to, email.subject, email.body)
		}
		if err != nil {
			fmt.Println("Failed to send email:", err)
		} else {
			fmt.Println("Email sent successfully")
		}
	}
}

func (m *MailerImpl) Shutdown() {
	close(m.queue)
	m.wg.Wait()
}
