package email

import (
	"bytes"
	"fmt"
	"html/template"
	"net/smtp"
	"os"
	"strconv"
)

type SMTPConfig struct {
	Host     string
	Port     int
	Username string
	Password string
	From     string
}

func LoadSMTPConfig() *SMTPConfig {
	port, _ := strconv.Atoi(getEnv("SMTP_PORT", "587"))
	return &SMTPConfig{
		Host:     getEnv("SMTP_HOST", ""),
		Port:     port,
		Username: getEnv("SMTP_USERNAME", ""),
		Password: getEnv("SMTP_PASSWORD", ""),
		From:     getEnv("SMTP_FROM", "no-reply@arnela.es"),
	}
}

func (c *SMTPConfig) IsConfigured() bool {
	return c.Host != "" && c.Username != "" && c.Password != ""
}

type Mailer struct {
	config *SMTPConfig
}

func NewMailer(cfg *SMTPConfig) *Mailer {
	return &Mailer{config: cfg}
}

func (m *Mailer) Send(to, subject, htmlBody string) error {
	if !m.config.IsConfigured() {
		return fmt.Errorf("SMTP not configured, skipping email to %s", to)
	}

	headers := map[string]string{
		"From":         m.config.From,
		"To":           to,
		"Subject":      subject,
		"MIME-Version": "1.0",
		"Content-Type": "text/html; charset=\"utf-8\"",
	}

	var msg bytes.Buffer
	for k, v := range headers {
		fmt.Fprintf(&msg, "%s: %s\r\n", k, v)
	}
	msg.WriteString("\r\n")
	msg.WriteString(htmlBody)

	auth := smtp.PlainAuth("", m.config.Username, m.config.Password, m.config.Host)
	addr := fmt.Sprintf("%s:%d", m.config.Host, m.config.Port)

	return smtp.SendMail(addr, auth, m.config.From, []string{to}, msg.Bytes())
}

// RenderTemplate renders an HTML template with the given data.
func RenderTemplate(tmplStr string, data interface{}) (string, error) {
	tmpl, err := template.New("email").Parse(tmplStr)
	if err != nil {
		return "", fmt.Errorf("failed to parse template: %w", err)
	}
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("failed to execute template: %w", err)
	}
	return buf.String(), nil
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
