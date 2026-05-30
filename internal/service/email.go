package service

import (
	"bytes"
	"embed"
	"html/template"
	"net/smtp"
	"strconv"

	"BBingyan/internal/config"
)

//go:embed../../templates/*.html
var templateFS embed.FS

var captchaTmpl *template.Template

func init() {
	captchaTmpl = template.Must(template.ParseFS(templateFS, "templates/captcha.html"))
}

func SendValidationCode(email, code string) error {
	var buf bytes.Buffer
	if err := captchaTmpl.Execute(&buf, map[string]any{
		"Code":   code,
		"Expire": config.Conf.Captcha.Expire,
	}); err != nil {
		return err
	}
	return sendMail(email, config.Conf.Captcha.Title, buf.String())
}

func sendMail(to, subject, body string) error {
	cfg := config.Conf.Mail
	auth := smtp.PlainAuth("", cfg.Username, cfg.Password, cfg.Host)
	msg := []byte(
		"To: " + to + "\r\n" +
			"From: " + cfg.Nickname + " <" + cfg.Username + ">\r\n" +
			"Subject: " + subject + "\r\n" +
			"Content-Type: text/html; charset=UTF-8\r\n\r\n" +
			body,
	)
	return smtp.SendMail(cfg.Host+":"+strconv.Itoa(cfg.Port), auth, cfg.Username, []string{to}, msg)
}
