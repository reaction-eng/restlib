// Copyright 2019 Reaction Engineering International. All rights reserved.
// Use of this source code is governed by the MIT license in the file LICENSE.txt.

package google

import (
	"time"
)

type gFile struct {
	//Keep a boolean if it is a file
	Id string `json:"Id"`

	//Hold the item
	Name string `json:"name"`

	//Hold if we should hide the item
	HideListing bool `json:"hideListing"`

	//Keep a date if useful
	Date *time.Time `json:"date"`
}

//Hold a base type document
func (file gFile) GetId() string {
	return file.Id
}
func (file gFile) GetName() string {
	return file.Name
}
func (file gFile) GetDate() *time.Time {
	return file.Date
}
