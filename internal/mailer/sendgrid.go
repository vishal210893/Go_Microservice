package mailer

import (
	"fmt"
	"strings"
	"sync"
	"text/template"
	"time"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

// SendGridMailer sends email via SendGrid. Safe for concurrent use.
type SendGridMailer struct {
	fromEmail string
	apiKey    string
	client    *sendgrid.Client
}

// NewSendgrid creates a SendGrid mailer with the given API key and from address.
func NewSendGrid(apiKey, fromEmail string) *SendGridMailer {
	client := sendgrid.NewSendClient(apiKey)

	return &SendGridMailer{
		fromEmail: fromEmail,
		apiKey:    apiKey,
		client:    client,
	}
}

// template cache to avoid reparsing on every send.
var tmplCache sync.Map // map[string]*template.Template

func getTemplate(templateFile string) (*template.Template, error) {
	if v, ok := tmplCache.Load(templateFile); ok {
		if t, ok := v.(*template.Template); ok && t != nil {
			return t, nil
		}
	}

	t, err := template.ParseFS(FS, "template/"+templateFile)
	if err != nil {
		return nil, fmt.Errorf("parse template %q: %w", templateFile, err)
	}
	tmplCache.Store(templateFile, t)
	return t, nil
}

// Send renders the template (expects blocks "subject" and "body") and sends the email.
// Retries on transient failures (HTTP 5xx/429) with exponential backoff.
func (m *SendGridMailer) Send(templateFile, username, email string, data any, isSandbox bool) (int, error) {
	from := mail.NewEmail(FromName, m.fromEmail)
	to := mail.NewEmail(username, email)

	tmpl, err := getTemplate(templateFile)
	if err != nil {
		return -1, err
	}

	var subjBuilder, bodyBuilder strings.Builder
	if err := tmpl.ExecuteTemplate(&subjBuilder, "subject", data); err != nil {
		return -1, fmt.Errorf("execute subject template: %w", err)
	}
	if err := tmpl.ExecuteTemplate(&bodyBuilder, "body", data); err != nil {
		return -1, fmt.Errorf("execute body template: %w", err)
	}

	message := mail.NewSingleEmail(from, strings.TrimSpace(subjBuilder.String()), to, "", bodyBuilder.String())
	message.SetMailSettings(&mail.MailSettings{
		SandboxMode: &mail.Setting{
			Enable: &isSandbox,
		},
	})

	var lastErr error

	for attempt := 0; attempt < maxRetries; attempt++ {
		response, sendErr := m.client.Send(message)
		if sendErr != nil {
			lastErr = sendErr
		} else if response != nil && response.StatusCode >= 200 && response.StatusCode < 300 {
			return response.StatusCode, nil
		} else {
			// Non-success response without a client error
			var status, body string
			if response != nil {
				status = fmt.Sprintf("%d", response.StatusCode)
				body = response.Body
			} else {
				status = "unknown"
			}
			lastErr = fmt.Errorf("sendgrid non-success status %s: %s", status, body)

			// Only retry for 429 or 5xx
			if response == nil || (response.StatusCode != 429 && response.StatusCode < 500) {
				return -1, lastErr
			}
		}

		// Exponential backoff: 1s, 2s, then 4s
		if attempt < maxRetries-1 {
			time.Sleep(time.Second * time.Duration(1<<attempt))
		}
	}

	return -1, fmt.Errorf("failed to send email after %d attempts: %w", maxRetries, lastErr)
}