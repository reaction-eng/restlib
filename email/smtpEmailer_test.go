// Copyright 2019 Reaction Engineering International. All rights reserved.
// Use of this source code is governed by the MIT license in the file LICENSE.txt.

package email_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"testing"
	"text/template"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/reaction-eng/restlib/utils"

	"github.com/reaction-eng/restlib/email"

	"github.com/golang/mock/gomock"
	"github.com/reaction-eng/restlib/mocks"
)

func TestNewSmtpEmailer(t *testing.T) {
	// arrange
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockConfiguration := mocks.NewMockConfiguration(mockCtrl)
	mockConfiguration.EXPECT().GetStringFatal("smtp_server").Return("SMTPSERVER").Times(1)
	mockConfiguration.EXPECT().GetStringFatal("smtp_port").Return("SMTPPORT").Times(1)
	mockConfiguration.EXPECT().GetStringFatal("smtp_user").Return("SMTPUSER").Times(1)
	mockConfiguration.EXPECT().GetStringFatal("smtp_password").Return("SMTPPASSWORD").Times(1)
	mockConfiguration.EXPECT().GetStringFatal("smtp_from").Return("SMTPFROM").Times(1)

	mockSmtpConnection := mocks.NewMockSmtpConnection(mockCtrl)

	// act
	email.NewSmtpEmailer(mockConfiguration, mockSmtpConnection)

	// assert
}

func TestSmtpEmailer_Send(t *testing.T) {
	testCases := []struct {
		header      email.HeaderInfo
		body        string
		attachments map[string][]*utils.Base64File
		error       error
	}{
		{email.HeaderInfo{To: []string{"ToString"}, Bcc: nil, Subject: "test 123", ReplyTo: "reply"}, "Email Body", nil, nil},
		{email.HeaderInfo{To: []string{"ToString"}, Bcc: []string{"bccOne", "bbcTwo"}, Subject: "test 1234"}, "Email Body", nil, nil},
		{email.HeaderInfo{To: []string{"ToString"}, Bcc: []string{"bccOne", "bbcTwo"}, Subject: "test 1234", ReplyTo: ""}, "Email Body", nil, nil},
		{email.HeaderInfo{To: []string{"ToString"}, Bcc: []string{"bccOne", "bbcTwo"}, Subject: "test 1234", ReplyTo: ""}, "Email Body", nil, errors.New("test error")},
		{email.HeaderInfo{To: []string{"ToString"}, Bcc: []string{"bccOne", "bbcTwo"}, Subject: "test 1234", ReplyTo: "Test One"}, "Email Body", map[string][]*utils.Base64File{"testOne": {utils.NewBase64FileFromData("attachment1", []byte{1, 2, 3, 4})}}, nil},
	}

	// arrange
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockConfiguration := mocks.NewMockConfiguration(mockCtrl)
	mockConfiguration.EXPECT().GetStringFatal("smtp_server").Return("SMTPSERVER").Times(len(testCases))
	mockConfiguration.EXPECT().GetStringFatal("smtp_port").Return("SMTPPORT").Times(len(testCases))
	mockConfiguration.EXPECT().GetStringFatal("smtp_user").Return("SMTPUSER").Times(len(testCases))
	mockConfiguration.EXPECT().GetStringFatal("smtp_password").Return("SMTPPASSWORD").Times(len(testCases))
	mockConfiguration.EXPECT().GetStringFatal("smtp_from").Return("SMTPFROM").Times(len(testCases))

	for _, testCase := range testCases {
		// arrange
		mockMail := mocks.NewMockMail(mockCtrl)
		mockMail.EXPECT().To(testCase.header.To).Times(1)
		bccCalled := 0
		if testCase.header.Bcc != nil {
			bccCalled = 1
		}
		mockMail.EXPECT().Bcc(testCase.header.Bcc).Times(bccCalled)
		mockMail.EXPECT().Subject(testCase.header.Subject).Times(1)
		mockMail.EXPECT().From("SMTPFROM").Times(1)
		replyTo := 0
		if len(testCase.header.ReplyTo) > 0 {
			replyTo = 1
		}
		mockMail.EXPECT().ReplyTo(testCase.header.ReplyTo).Times(replyTo)
		mockMail.EXPECT().SetPlain(testCase.body).Times(1)

		for _, values := range testCase.attachments {
			for _, value := range values {
				//Save it to the mail
				mockMail.EXPECT().Attach(value.GetName(), value.GetDataReader()).Times(1)
			}
		}

		mockMail.EXPECT().Send().Times(1).Return(testCase.error)

		mockSmtpConnection := mocks.NewMockSmtpConnection(mockCtrl)
		mockSmtpConnection.EXPECT().New("SMTPSERVER", "SMTPUSER", "SMTPPASSWORD", "SMTPFROM", "SMTPPORT").Times(1).Return(mockMail)

		emailer := email.NewSmtpEmailer(mockConfiguration, mockSmtpConnection)

		// act
		err := emailer.Send(&testCase.header, testCase.body, testCase.attachments)

		assert.Equal(t, testCase.error, err)
	}
}

func TestSmtpEmailer_TemplateString(t *testing.T) {

	referenceTime := time.Unix(1574622126, 0)

	testCases := []struct {
		header         email.HeaderInfo
		templateString string
		data           interface{}
		expected       string
		attachments    map[string][]*utils.Base64File
		error          error
	}{
		{email.HeaderInfo{To: []string{"ToString"}, Bcc: nil, Subject: "test 123", ReplyTo: "reply"},
			"<!DOCTYPE html><html><body><h1>ExampleInfo</h1><h3>{{.Info1}}</h3><h3>Email: <small>{{.Info2}}</small></h3></body></html>",
			struct {
				Info1 string
				Info2 int
			}{"alpha beta", 32},
			"<!DOCTYPE html><html><body><h1>ExampleInfo</h1><h3>alpha beta</h3><h3>Email: <small>32</small></h3></body></html>",
			nil,
			nil},
		{email.HeaderInfo{To: []string{"ToString"}, Bcc: nil, Subject: "test 123", ReplyTo: "reply"},
			"<!DOCTYPE html><html><body><h1>ExampleInfo</h1><h3>{{.Info1}}</h3><h3>Email: <small>{{.Info2}}</small></h3></body></html>",
			struct {
				Info1 string
				Info2 int
			}{"alpha beta", 32},
			"<!DOCTYPE html><html><body><h1>ExampleInfo</h1><h3>alpha beta</h3><h3>Email: <small>32</small></h3></body></html>",
			map[string][]*utils.Base64File{"testOne": {utils.NewBase64FileFromData("attachment1", []byte{1, 2, 3, 4})}},
			nil},
		{email.HeaderInfo{To: []string{"ToString"}, Bcc: nil, Subject: "test 123", ReplyTo: "reply"},
			`<!DOCTYPE html><html><body><h1>ExampleInfo</h1><h3>{{.Info1}}</h3><h3>Email: <small>{{formatInTimeZone .Info2 "America/Denver" "2006-01-02 03:04:05 PM"}}</small></h3></body></html>`,
			struct {
				Info1 string
				Info2 *time.Time
			}{"alpha beta", &referenceTime},
			"<!DOCTYPE html><html><body><h1>ExampleInfo</h1><h3>alpha beta</h3><h3>Email: <small>2019-11-24 12:02:06 PM</small></h3></body></html>",
			map[string][]*utils.Base64File{"testOne": {utils.NewBase64FileFromData("attachment1", []byte{1, 2, 3, 4})}},
			nil},
		{email.HeaderInfo{To: []string{"ToString"}, Bcc: []string{"bccOne", "bbcTwo"}, Subject: "test 1234"},
			"<!DOCTYPE html><html><body><h1>ExampleInfo</h1><h3>{{.Info1}}</h3><h3>Email: <small>{{.Info2}}</small></h3></body></html>",
			struct {
				Info1 string
				Info2 int
			}{"alpha beta", 32},
			"<!DOCTYPE html><html><body><h1>ExampleInfo</h1><h3>alpha beta</h3><h3>Email: <small>32</small></h3></body></html>",
			nil,
			nil},
		{email.HeaderInfo{To: []string{"ToString"}, Bcc: []string{"bccOne", "bbcTwo"}, Subject: "test 1234"},
			"<!DOCTYPE html><html><body><h1>ExampleInfo</h1><h3>{{.Info1}}</h3><h3>Email: <small>{{.Info3}}</small></h3></body></html>",
			struct {
				Info1 string
				Info2 int
			}{"alpha beta", 32},
			"<!DOCTYPE html><html><body><h1>ExampleInfo</h1><h3>alpha beta</h3><h3>Email: <small>",
			nil,
			template.ExecError{}},
	}

	// arrange
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockConfiguration := mocks.NewMockConfiguration(mockCtrl)
	mockConfiguration.EXPECT().GetStringFatal("smtp_server").Return("SMTPSERVER").Times(len(testCases))
	mockConfiguration.EXPECT().GetStringFatal("smtp_port").Return("SMTPPORT").Times(len(testCases))
	mockConfiguration.EXPECT().GetStringFatal("smtp_user").Return("SMTPUSER").Times(len(testCases))
	mockConfiguration.EXPECT().GetStringFatal("smtp_password").Return("SMTPPASSWORD").Times(len(testCases))
	mockConfiguration.EXPECT().GetStringFatal("smtp_from").Return("SMTPFROM").Times(len(testCases))

	for _, testCase := range testCases {
		// arrange
		mockMail := mocks.NewMockMail(mockCtrl)
		mockMail.EXPECT().To(testCase.header.To).Times(1)
		bccCalled := 0
		if testCase.header.Bcc != nil {
			bccCalled = 1
		}
		mockMail.EXPECT().Bcc(testCase.header.Bcc).Times(bccCalled)
		mockMail.EXPECT().Subject(testCase.header.Subject).Times(1)
		mockMail.EXPECT().From("SMTPFROM").Times(1)
		replyTo := 0
		if len(testCase.header.ReplyTo) > 0 {
			replyTo = 1
		}
		mockMail.EXPECT().ReplyTo(testCase.header.ReplyTo).Times(replyTo)

		// create a buffer for writing
		var mockHtmlWriter bytes.Buffer
		mockMail.EXPECT().Html().Times(1).Return(&mockHtmlWriter)

		for _, values := range testCase.attachments {
			for _, value := range values {
				//Save it to the mail
				mockMail.EXPECT().Attach(value.GetName(), value.GetDataReader()).Times(1)
			}
		}
		if testCase.error == nil {
			tryJsonString, _ := json.Marshal(testCase.data)
			mockMail.EXPECT().SetPlain(string(tryJsonString)).Times(1)
			mockMail.EXPECT().Send().Times(1).Return(testCase.error)
		}
		mockSmtpConnection := mocks.NewMockSmtpConnection(mockCtrl)
		mockSmtpConnection.EXPECT().New("SMTPSERVER", "SMTPUSER", "SMTPPASSWORD", "SMTPFROM", "SMTPPORT").Times(1).Return(mockMail)

		emailer := email.NewSmtpEmailer(mockConfiguration, mockSmtpConnection)

		// act
		err := emailer.SendTemplateString(&testCase.header, testCase.templateString, testCase.data, testCase.attachments)

		assert.IsType(t, testCase.error, err)
		assert.Equal(t, testCase.expected, mockHtmlWriter.String())
	}
}

func TestSmtpEmailer_SendTemplateFile(t *testing.T) {

	referenceTime := time.Unix(1574622126, 0)

	testCases := []struct {
		header         email.HeaderInfo
		templateString string
		data           interface{}
		expected       string
		attachments    map[string][]*utils.Base64File
		error          error
	}{
		{email.HeaderInfo{To: []string{"ToString"}, Bcc: nil, Subject: "test 123", ReplyTo: "reply"},
			"<!DOCTYPE html><html><body><h1>ExampleInfo</h1><h3>{{.Info1}}</h3><h3>Email: <small>{{.Info2}}</small></h3></body></html>",
			struct {
				Info1 string
				Info2 int
			}{"alpha beta", 32},
			"<!DOCTYPE html><html><body><h1>ExampleInfo</h1><h3>alpha beta</h3><h3>Email: <small>32</small></h3></body></html>",
			nil,
			nil},
		{email.HeaderInfo{To: []string{"ToString"}, Bcc: nil, Subject: "test 123", ReplyTo: "reply"},
			"<!DOCTYPE html><html><body><h1>ExampleInfo</h1><h3>{{.Info1}}</h3><h3>Email: <small>{{.Info2}}</small></h3></body></html>",
			struct {
				Info1 string
				Info2 int
			}{"alpha beta", 32},
			"<!DOCTYPE html><html><body><h1>ExampleInfo</h1><h3>alpha beta</h3><h3>Email: <small>32</small></h3></body></html>",
			map[string][]*utils.Base64File{"testOne": {utils.NewBase64FileFromData("attachment1", []byte{1, 2, 3, 4})}},
			nil},
		{email.HeaderInfo{To: []string{"ToString"}, Bcc: nil, Subject: "test 123", ReplyTo: "reply"},
			`<!DOCTYPE html><html><body><h1>ExampleInfo</h1><h3>{{.Info1}}</h3><h3>Email: <small>{{formatInTimeZone .Info2 "America/Denver" "2006-01-02 03:04:05 PM"}}</small></h3></body></html>`,
			struct {
				Info1 string
				Info2 *time.Time
			}{"alpha beta", &referenceTime},
			"<!DOCTYPE html><html><body><h1>ExampleInfo</h1><h3>alpha beta</h3><h3>Email: <small>2019-11-24 12:02:06 PM</small></h3></body></html>",
			map[string][]*utils.Base64File{"testOne": {utils.NewBase64FileFromData("attachment1", []byte{1, 2, 3, 4})}},
			nil},
		{email.HeaderInfo{To: []string{"ToString"}, Bcc: []string{"bccOne", "bbcTwo"}, Subject: "test 1234"},
			"<!DOCTYPE html><html><body><h1>ExampleInfo</h1><h3>{{.Info1}}</h3><h3>Email: <small>{{.Info2}}</small></h3></body></html>",
			struct {
				Info1 string
				Info2 int
			}{"alpha beta", 32},
			"<!DOCTYPE html><html><body><h1>ExampleInfo</h1><h3>alpha beta</h3><h3>Email: <small>32</small></h3></body></html>",
			nil,
			nil},
		{email.HeaderInfo{To: []string{"ToString"}, Bcc: []string{"bccOne", "bbcTwo"}, Subject: "test 1234"},
			"<!DOCTYPE html><html><body><h1>ExampleInfo</h1><h3>{{.Info1}}</h3><h3>Email: <small>{{.Info3}}</small></h3></body></html>",
			struct {
				Info1 string
				Info2 int
			}{"alpha beta", 32},
			"<!DOCTYPE html><html><body><h1>ExampleInfo</h1><h3>alpha beta</h3><h3>Email: <small>",
			nil,
			template.ExecError{}},
		{email.HeaderInfo{To: []string{"ToString"}, Bcc: nil, Subject: "test 123", ReplyTo: "reply"},
			`<!DOCTYPE html><html><body><h1>ExampleInfo</h1><h3>{{.Info1}}</h3><h3>Email: <small>{{formatInTimeZone .Info2 "America/Denver" "2006-01-02 03:04:05 PM"}}</small></h3></body></html>`,
			struct {
				Info1 string
				Info2 *time.Time
			}{"alpha beta", &referenceTime},
			"<!DOCTYPE html><html><body><h1>ExampleInfo</h1><h3>alpha beta</h3><h3>Email: <small>2019-11-24 12:02:06 PM</small></h3></body></html>",
			map[string][]*utils.Base64File{"testOne": {utils.NewBase64FileFromData("attachment1", []byte{1, 2, 3, 4})}},
			errors.New("send error")},
	}

	// arrange
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockConfiguration := mocks.NewMockConfiguration(mockCtrl)
	mockConfiguration.EXPECT().GetStringFatal("smtp_server").Return("SMTPSERVER").Times(len(testCases))
	mockConfiguration.EXPECT().GetStringFatal("smtp_port").Return("SMTPPORT").Times(len(testCases))
	mockConfiguration.EXPECT().GetStringFatal("smtp_user").Return("SMTPUSER").Times(len(testCases))
	mockConfiguration.EXPECT().GetStringFatal("smtp_password").Return("SMTPPASSWORD").Times(len(testCases))
	mockConfiguration.EXPECT().GetStringFatal("smtp_from").Return("SMTPFROM").Times(len(testCases))

	for _, testCase := range testCases {
		// arrange
		mockMail := mocks.NewMockMail(mockCtrl)
		mockMail.EXPECT().To(testCase.header.To).Times(1)
		bccCalled := 0
		if testCase.header.Bcc != nil {
			bccCalled = 1
		}
		mockMail.EXPECT().Bcc(testCase.header.Bcc).Times(bccCalled)
		mockMail.EXPECT().Subject(testCase.header.Subject).Times(1)
		mockMail.EXPECT().From("SMTPFROM").Times(1)
		replyTo := 0
		if len(testCase.header.ReplyTo) > 0 {
			replyTo = 1
		}
		mockMail.EXPECT().ReplyTo(testCase.header.ReplyTo).Times(replyTo)

		// create a buffer for writing
		var mockHtmlWriter bytes.Buffer
		mockMail.EXPECT().Html().Times(1).Return(&mockHtmlWriter)

		for _, values := range testCase.attachments {
			for _, value := range values {
				//Save it to the mail
				mockMail.EXPECT().Attach(value.GetName(), value.GetDataReader()).Times(1)
			}
		}
		if _, isType := testCase.error.(template.ExecError); !isType {
			tryJsonString, _ := json.Marshal(testCase.data)
			mockMail.EXPECT().SetPlain(string(tryJsonString)).Times(1)
			mockMail.EXPECT().Send().Times(1).Return(testCase.error)
		}
		mockSmtpConnection := mocks.NewMockSmtpConnection(mockCtrl)
		mockSmtpConnection.EXPECT().New("SMTPSERVER", "SMTPUSER", "SMTPPASSWORD", "SMTPFROM", "SMTPPORT").Times(1).Return(mockMail)

		emailer := email.NewSmtpEmailer(mockConfiguration, mockSmtpConnection)

		// Save the html template to a temp file
		tmpFile, tempFileError := ioutil.TempFile(os.TempDir(), "prefix-")
		assert.Nil(t, tempFileError)
		if _, tempFileError = tmpFile.Write([]byte(testCase.templateString)); tempFileError != nil {
			assert.Nil(t, tempFileError)
		}
		if tempFileError = tmpFile.Close(); tempFileError != nil {
			assert.Nil(t, tempFileError)
		}

		// act
		err := emailer.SendTemplateFile(&testCase.header, tmpFile.Name(), testCase.data, testCase.attachments)

		assert.IsType(t, testCase.error, err)
		assert.Equal(t, testCase.expected, mockHtmlWriter.String())
		os.Remove(tmpFile.Name())
	}
}

func buildMockTable(mockCtrl *gomock.Controller) *mocks.MockTableInfo {
	mockChild := mocks.NewMockTableInfo(mockCtrl)
	mockChild.EXPECT().GetTitle().Return("Example Child").Times(1)
	mockChild.EXPECT().IsNode().Return(false).Times(1).Times(1)
	mockChild.EXPECT().GetValue().Return("example child Value").Times(1)
	mockChild.EXPECT().GetChildren().Times(0)

	mockTable := mocks.NewMockTableInfo(mockCtrl)
	mockTable.EXPECT().GetTitle().Return("Example Table").Times(2)
	mockTable.EXPECT().IsNode().Times(0)
	mockTable.EXPECT().GetValue().Times(0)
	mockTable.EXPECT().GetChildren().Return([]email.TableInfo{mockChild})

	return mockTable
}

func TestSmtpEmailer_SendTable(t *testing.T) {
	// arrange
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	testCases := []struct {
		header      email.HeaderInfo
		table       email.TableInfo
		expected    string
		attachments map[string][]*utils.Base64File
		error       error
	}{
		{email.HeaderInfo{To: []string{"ToString"}, Bcc: nil, Subject: "test 123", ReplyTo: "reply"},
			buildMockTable(mockCtrl),
			"<!DOCTYPE html><html lang=\"en\"><head><meta charset=\"UTF-8\"><style>.header{  background-color: #aed957;  font-size: 18px;  text-align: center;  font-weight: bold;}.title{  text-align: left;  font-size: 15px;font-weight: bold;  background-color: gray;}.content{    background-color: white;}</style></head><body><h1>Example Table</h1><table width=\"99%\" border=\"0\" cellpadding=\"1\" cellspacing=\"0\" bgcolor=\"#EAEAEA\"><tr class=\"header\"><td>Example Table</td></tr><tr class=\"title\"><td><strong>Example Child</strong></td></tr><tr><td>example child Value</td></tr></table></body></html>",
			nil,
			nil},
		{email.HeaderInfo{To: []string{"ToString"}, Bcc: nil, Subject: "test 123"},
			buildMockTable(mockCtrl),
			"<!DOCTYPE html><html lang=\"en\"><head><meta charset=\"UTF-8\"><style>.header{  background-color: #aed957;  font-size: 18px;  text-align: center;  font-weight: bold;}.title{  text-align: left;  font-size: 15px;font-weight: bold;  background-color: gray;}.content{    background-color: white;}</style></head><body><h1>Example Table</h1><table width=\"99%\" border=\"0\" cellpadding=\"1\" cellspacing=\"0\" bgcolor=\"#EAEAEA\"><tr class=\"header\"><td>Example Table</td></tr><tr class=\"title\"><td><strong>Example Child</strong></td></tr><tr><td>example child Value</td></tr></table></body></html>",
			map[string][]*utils.Base64File{"testOne": {utils.NewBase64FileFromData("attachment1", []byte{1, 2, 3, 4})}},
			nil},
		{email.HeaderInfo{To: []string{"ToString"}, Bcc: []string{"bccOne", "bbcTwo"}, Subject: "test 123"},
			buildMockTable(mockCtrl),
			"<!DOCTYPE html><html lang=\"en\"><head><meta charset=\"UTF-8\"><style>.header{  background-color: #aed957;  font-size: 18px;  text-align: center;  font-weight: bold;}.title{  text-align: left;  font-size: 15px;font-weight: bold;  background-color: gray;}.content{    background-color: white;}</style></head><body><h1>Example Table</h1><table width=\"99%\" border=\"0\" cellpadding=\"1\" cellspacing=\"0\" bgcolor=\"#EAEAEA\"><tr class=\"header\"><td>Example Table</td></tr><tr class=\"title\"><td><strong>Example Child</strong></td></tr><tr><td>example child Value</td></tr></table></body></html>",
			map[string][]*utils.Base64File{"testOne": {utils.NewBase64FileFromData("attachment1", []byte{1, 2, 3, 4})}},
			nil},
		{email.HeaderInfo{To: []string{"ToString"}, Bcc: []string{"bccOne", "bbcTwo"}, Subject: "test 123"},
			buildMockTable(mockCtrl),
			"<!DOCTYPE html><html lang=\"en\"><head><meta charset=\"UTF-8\"><style>.header{  background-color: #aed957;  font-size: 18px;  text-align: center;  font-weight: bold;}.title{  text-align: left;  font-size: 15px;font-weight: bold;  background-color: gray;}.content{    background-color: white;}</style></head><body><h1>Example Table</h1><table width=\"99%\" border=\"0\" cellpadding=\"1\" cellspacing=\"0\" bgcolor=\"#EAEAEA\"><tr class=\"header\"><td>Example Table</td></tr><tr class=\"title\"><td><strong>Example Child</strong></td></tr><tr><td>example child Value</td></tr></table></body></html>",
			map[string][]*utils.Base64File{"testOne": {utils.NewBase64FileFromData("attachment1", []byte{1, 2, 3, 4})}},
			errors.New("send error")},
	}

	mockConfiguration := mocks.NewMockConfiguration(mockCtrl)
	mockConfiguration.EXPECT().GetStringFatal("smtp_server").Return("SMTPSERVER").Times(len(testCases))
	mockConfiguration.EXPECT().GetStringFatal("smtp_port").Return("SMTPPORT").Times(len(testCases))
	mockConfiguration.EXPECT().GetStringFatal("smtp_user").Return("SMTPUSER").Times(len(testCases))
	mockConfiguration.EXPECT().GetStringFatal("smtp_password").Return("SMTPPASSWORD").Times(len(testCases))
	mockConfiguration.EXPECT().GetStringFatal("smtp_from").Return("SMTPFROM").Times(len(testCases))

	for _, testCase := range testCases {
		// arrange
		mockMail := mocks.NewMockMail(mockCtrl)
		mockMail.EXPECT().To(testCase.header.To).Times(1)
		bccCalled := 0
		if testCase.header.Bcc != nil {
			bccCalled = 1
		}
		mockMail.EXPECT().Bcc(testCase.header.Bcc).Times(bccCalled)
		mockMail.EXPECT().Subject(testCase.header.Subject).Times(1)
		mockMail.EXPECT().From("SMTPFROM").Times(1)
		replyTo := 0
		if len(testCase.header.ReplyTo) > 0 {
			replyTo = 1
		}
		mockMail.EXPECT().ReplyTo(testCase.header.ReplyTo).Times(replyTo)

		// create a buffer for writing
		var mockHtmlWriter bytes.Buffer
		mockMail.EXPECT().Html().Times(1).Return(&mockHtmlWriter)

		for _, values := range testCase.attachments {
			for _, value := range values {
				//Save it to the mail
				mockMail.EXPECT().Attach(value.GetName(), value.GetDataReader()).Times(1)
			}
		}
		mockMail.EXPECT().SetPlain("HTML Email Required").Times(1)
		mockMail.EXPECT().Send().Times(1).Return(testCase.error)

		mockSmtpConnection := mocks.NewMockSmtpConnection(mockCtrl)
		mockSmtpConnection.EXPECT().New("SMTPSERVER", "SMTPUSER", "SMTPPASSWORD", "SMTPFROM", "SMTPPORT").Times(1).Return(mockMail)

		emailer := email.NewSmtpEmailer(mockConfiguration, mockSmtpConnection)

		// act
		err := emailer.SendTable(&testCase.header, testCase.table, testCase.attachments)

		assert.IsType(t, testCase.error, err)
		assert.Equal(t, testCase.expected, mockHtmlWriter.String())
	}
}
