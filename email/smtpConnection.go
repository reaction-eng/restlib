// Copyright 2019 Reaction Engineering International. All rights reserved.
// Use of this source code is governed by the MIT license in the file LICENSE.txt.

package email

//go:generate mockgen -destination=../mocks/mock_smtpConnection.go -package=mocks github.com/reaction-eng/restlib/email SmtpConnection,Mail

import "io"

type SmtpConnection interface {
	New(smtpServer string, smtpUser string, smtpPassword string, smtpFrom string, smtpPort string) Mail
}

type Mail interface {
	To(addrs ...string)
	Bcc(addrs ...string)
	Subject(sub string)
	From(addr string)
	ReplyTo(addr string)
	Html() io.Writer
	SetPlain(body string)
	Send() error
	Attach(name string, r io.Reader)
}
