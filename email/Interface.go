package email

/**
Simple struct for email
*/
type HeaderInfo struct {
	To      []string
	Subject string
}

/**
Simple email message
*/
type Interface interface {

	/**
	Get the specific news istem
	*/
	SendEmail(email *HeaderInfo, body string) error

	/**
	Send html email
	*/
	SendEmailHtml(email *HeaderInfo, templateName string, data interface{}) error

	/**
	Send html email
	*/
	SendEmailJson(email *HeaderInfo, data interface{}) error
}
