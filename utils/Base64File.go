// Copyright 2019 Reaction Engineering International. All rights reserved.
// Use of this source code is governed by the MIT license in the file LICENSE.txt.

package utils

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"io/ioutil"
	"mime/multipart"
	"regexp"
	"strings"
)

type Base64File struct {
	data []byte

	name string

	//Hold the mime data if known
	mime string
}

func NewBase64FileFromData(name string, data []byte) *Base64File {
	return &Base64File{name: name, data: data}
}

func NewBase64FileFromForm(file multipart.File, fileInfo *multipart.FileHeader) (*Base64File, error) {
	//Store the file name and mime
	// Read entire JPG into byte slice.
	reader := bufio.NewReader(file)
	content, error := ioutil.ReadAll(reader)

	//Now create the info
	b64File := Base64File{
		data: content,
		name: fileInfo.Filename,
		mime: fileInfo.Header.Get("Content-GetType"),
	}

	return &b64File, error
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
		mime: "",
	}

	return &file, err
}

func (file *Base64File) GetDataBytes() []byte {
	return file.data
}

func (file *Base64File) GetEncodedData() string {
	return base64.StdEncoding.EncodeToString(file.data)
}

func (file *Base64File) GetMime() string {
	return file.mime
}

func (file *Base64File) GetDataReader() *bytes.Reader {

	reader := bytes.NewReader(file.data)
	return reader
}

func (file *Base64File) GetName() string {
	return file.name
}
