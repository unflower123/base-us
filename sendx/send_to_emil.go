package sendx

import (
	"context"
	"fmt"
	"gopkg.in/gomail.v2"
)

type gmailConfig struct {
	host     string
	port     int
	from     string
	password string
	title    string
}

func WithGmailConfig(from, pwd, title string) MsgConfigOption {
	return func(c *MsgConfig) {
		c.Mail = &gmailConfig{
			host:     "smtp.gmail.com",
			port:     587,
			from:     from,
			password: pwd,
			title:    title,
		}
	}
}

func WithQmailConfig(from, pwd, title string) MsgConfigOption {
	return func(c *MsgConfig) {
		c.Mail = &gmailConfig{
			host:     "smtp.qq.com",
			port:     587,
			from:     from,
			password: pwd,
			title:    title,
		}
	}
}

func (m gmailConfig) SendMsg(ctx context.Context, to, subject, bodyParam string, args ...any) error {
	body := fmt.Sprintf(bodyParam, args)

	gm := gomail.NewMessage()
	gm.SetAddressHeader("From", m.from, m.title)
	gm.SetHeader("To", to)
	gm.SetHeader("Subject", subject)
	gm.SetBody("text/html", body)

	d := gomail.NewDialer(m.host, m.port, m.from, m.password)
	if err := d.DialAndSend(gm); err != nil {
		return err
	}

	return nil
}
