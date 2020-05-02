// Copyright 2019 Reaction Engineering International. All rights reserved.
// Use of this source code is governed by the MIT license in the file LICENSE.txt.

package preferences

//Get the setting group
type Preferences struct {
	//And the value
	Settings *SettingGroup `json:"settings"`

	//We can also old other groups
	Options *OptionGroup `json:"options"`
}
