// Copyright 2019 Reaction Engineering International. All rights reserved.
// Use of this source code is governed by the MIT license in the file LICENSE.txt.

package email

import "github.com/reaction-eng/restlib/utils"

/**
Simple struct for email
*/
type HeaderInfo struct {
	To      []string
	Subject string
	ReplyTo string
}

/**
Simple email message
*/
type Interface interface {

	/**
	Get the specific news istem
	*/
	SendEmail(email *HeaderInfo, body string, attachments map[string][]*utils.Base64File) error

	/**
	Send html email
	*/
	SendEmailTemplateString(email *HeaderInfo, templateString string, data interface{}, attachments map[string][]*utils.Base64File) error
	SendEmailTemplateFile(email *HeaderInfo, templateString string, data interface{}, attachments map[string][]*utils.Base64File) error

	/**
	Send html email
	*/
	SendEmailTable(email *HeaderInfo, tableData TableInfo, attachments map[string][]*utils.Base64File) error
}

/**
Simple email message
*/
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
