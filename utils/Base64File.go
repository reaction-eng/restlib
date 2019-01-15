package utils

import (
	"bytes"
	"encoding/base64"
	"regexp"
	"strings"
)

type Base64File struct {
	data []byte

	name string
}

/**
Decode the base 64 string
*/
func NewBase64File(data string) (*Base64File, error) {
	//Split the data
	//Get the name info

	r, _ := regexp.Compile(";name=.*;")

	//Now get the matching string
	nameString := r.FindString(data)

	//Now remove the start and end
	nameString = strings.Replace(nameString, ";name=", "", 1)
	nameString = strings.Replace(nameString, ";", "", -1)

	//Find the data location
	// First compile the delimiter expression.
	re := regexp.MustCompile(";base64,")

	// Split based on pattern.
	dataString := re.Split(data, -1)[1]

	//Get the data
	dataBytes, err := base64.StdEncoding.DecodeString(dataString)

	//If there is no error
	if err != nil {
		return nil, err
	}

	//Now create the info
	file := Base64File{
		data: dataBytes,
		name: nameString,
	}

	return &file, err
}

func (file *Base64File) GetDataBytes() []byte {
	return file.data
}

func (file *Base64File) GetDataReader() *bytes.Reader {

	reader := bytes.NewReader(file.data)
	return reader
}

func (file *Base64File) GetName() string {
	return file.name
}
