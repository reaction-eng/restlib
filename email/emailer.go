// Copyright 2019 Reaction Engineering International. All rights reserved.
// Use of this source code is governed by the MIT license in the file LICENSE.txt.

package email

//go:generate mockgen -destination=../mocks/mock_emailer.go -package=mocks github.com/reaction-eng/restlib/email Emailer,TableInfo

import "github.com/reaction-eng/restlib/utils"

type HeaderInfo struct {
	To      []string
	Bcc     []string
	Subject string
	ReplyTo string
}

type Emailer interface {
	Send(email *HeaderInfo, body string, attachments map[string][]*utils.Base64File) error

	SendTemplateString(email *HeaderInfo, templateString string, data interface{}, attachments map[string][]*utils.Base64File) error
	SendTemplateFile(email *HeaderInfo, templateString string, data interface{}, attachments map[string][]*utils.Base64File) error

	SendTable(email *HeaderInfo, tableData TableInfo, attachments map[string][]*utils.Base64File) error
}

type TableInfo interface {

	//Check to see if it node
	IsNode() bool

	//Get the title
	GetTitle() string

	//Get value
	GetValue() string

	//Get the children
	GetChildren() []TableInfo
}
