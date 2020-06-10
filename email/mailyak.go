package email

import (
	"io"
	"net/smtp"

	"github.com/domodwyer/mailyak"
)

type MailYakSmtpConnection struct {
}

func NewMailYakSmtpConnection() *MailYakSmtpConnection {
	return &MailYakSmtpConnection{}
}

func (*MailYakSmtpConnection) New(smtpServer string, smtpUser string, smtpPassword string, smtpFrom string, smtpPort string) Mail {
	return &MailYak{
		mailyak.New(smtpServer+smtpPort, smtp.PlainAuth("", smtpUser, smtpPassword, smtpServer)),
	}
}

type MailYak struct {
	*mailyak.MailYak
}

func (m *MailYak) Html() io.Writer {
	return m.HTML()
}

func (m *MailYak) SetPlain(body string) {
	m.Plain().Set(body)
}
