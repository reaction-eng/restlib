package email

import (
	"bitbucket.org/reidev/restlib/configuration"
	"bytes"
	"encoding/json"
	"html/template"
	"log"
	"net/smtp"
)

/**
Simple struct for email
*/
type SmtpSender struct {
	smtpServer   string
	smtpUser     string
	smtpPassword string
	smtpFrom     string
	smtpPort     string
}

const (
	MIMEHTML = "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	MIMETEXT = "MIME-version: 1.0;\nContent-Type: text/plain; charset=\"UTF-8\";\n\n"
	MIMEJSON = "MIME-version: 1.0;\nContent-Type: application/json; charset=\"UTF-8\";\n\n"
)

//Provide a method to make a new AnimalRepoSql
func NewSmtpSender(configFile string) *SmtpSender {

	//Load up the config
	config, err := configuration.NewConfiguration(configFile)

	if err != nil {
		log.Fatal(err)
	}

	sender := SmtpSender{
		smtpServer:   config.GetString("smtp_server"),
		smtpPort:     config.GetString("smtp_port"),
		smtpUser:     config.GetString("smtp_user"),
		smtpPassword: config.GetString("smtp_password"),
		smtpFrom:     config.GetString("smtp_from"),
	}

	return &sender

}

/**
Get all of the news
*/
func (repo *SmtpSender) SendEmail(email *HeaderInfo, body string) error {
	//Get the summary
	subject := "Subject: " + email.Subject + "!\n"
	msg := []byte(subject + MIMETEXT + "\n" + body)

	//Send it
	err := smtp.SendMail(
		repo.smtpServer+repo.smtpPort,                                         //smtp address
		smtp.PlainAuth("", repo.smtpUser, repo.smtpPassword, repo.smtpServer), //authentication
		repo.smtpFrom, //from
		email.To,      //List of toos
		[]byte(msg))   //Msg in byte form

	return err

}

/**
Parse the html template
*/
func parseTemplate(templateName string, data interface{}) (string, error) {
	t, err := template.ParseFiles(templateName)
	if err != nil {
		return "", err
	}
	buffer := new(bytes.Buffer)
	if err = t.Execute(buffer, data); err != nil {
		return "", err
	}
	body := buffer.String()
	return body, nil
}

/**
Get all of the news
*/
func (repo *SmtpSender) SendEmailHtml(email *HeaderInfo, templateName string, data interface{}) error {
	body, err := parseTemplate(templateName, data)

	//If it is an error
	if err != nil {
		return err
	}

	//Now build the entire body
	entireBody := "To: " + email.To[0] + "\r\nSubject: " + email.Subject + "\r\n" + MIMEHTML + "\r\n" + body

	msg := []byte(entireBody)

	//Send it
	err = smtp.SendMail(
		repo.smtpServer+repo.smtpPort,                                         //smtp address
		smtp.PlainAuth("", repo.smtpUser, repo.smtpPassword, repo.smtpServer), //authentication
		repo.smtpFrom, //from
		email.To,      //List of toos
		[]byte(msg))   //Msg in byte form

	return err
}

/**
Get all of the news
*/
func (repo *SmtpSender) SendEmailJson(email *HeaderInfo, data interface{}) error {

	//Get the json string
	jsonByte, _ := json.Marshal(data)

	//Now build the entire body
	entireBody := "To: " + email.To[0] + "\r\nSubject: " + email.Subject + "\r\n" + MIMETEXT + "\r\n" + string(jsonByte)

	msg := []byte(entireBody)

	//Send it
	err := smtp.SendMail(
		repo.smtpServer+repo.smtpPort,                                         //smtp address
		smtp.PlainAuth("", repo.smtpUser, repo.smtpPassword, repo.smtpServer), //authentication
		repo.smtpFrom, //from
		email.To,      //List of toos
		[]byte(msg))   //Msg in byte form

	return err
}
