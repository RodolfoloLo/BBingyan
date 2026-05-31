package service

import (
	"bytes"
	"html/template"
	"net/smtp"
	"strconv"

	"BBingyan/internal/config"
)

const captchaHTML = `<!DOCTYPE html>
<html>
<head><meta charset="UTF-8"></head>
<body style="font-family: Arial, sans-serif; padding: 20px;">
  <h2>BBingyan 邮箱验证码</h2>
  <p>您的验证码是：</p>
  <p style="font-size: 32px; font-weight: bold; color: #4A90D9; letter-spacing: 8px;">{{.Code}}</p>
  <p>有效期 {{.Expire}} 秒，请勿泄露给他人。</p>
  <hr>
  <p style="color: #999; font-size: 12px;">此邮件由系统自动发送，请勿回复。</p>
</body>
</html>`

var captchaTmpl *template.Template

func init() {
	captchaTmpl = template.Must(template.New("captcha").Parse(captchaHTML))
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
