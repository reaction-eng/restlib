// Copyright 2019 Reaction Engineering International. All rights reserved.
// Use of this source code is governed by the MIT license in the file LICENSE.txt.

package utils

import "encoding/json"

//import "strings"
//import "bufio"
import "bytes"
import "fmt"

type JSONstring string

// struct to JSON
func (str JSONstring) MarshalJSON() ([]byte, error) {

	// Marshal the JSON string into a byte array
	b, err := json.Marshal(string(str))

	// Get rid of quotes at beginning and end
	b = bytes.TrimPrefix(b, []byte("\""))
	b = bytes.TrimSuffix(b, []byte("\""))

	// Get rid of all "\" characters
	b = bytes.Replace(b, []byte("\\"), []byte(""), -1)
	//fmt.Println(string(b))

	return b, err

}

// JSON to aliased string object
func (str *JSONstring) UnmarshalJSON(data []byte) error {

	// Make a buffer to help compact the incoming JSON data
	buffer := new(bytes.Buffer)
	if err := json.Compact(buffer, data); err != nil {
		fmt.Println(err)
	}

	// Convert the buffered byte array to a string directly
	s := string(buffer.Bytes())
	//fmt.Println(s)

	// Save it into the JSONstring ptr
	*str = JSONstring(s)

	return nil

}
