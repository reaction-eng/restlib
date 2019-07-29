// Copyright 2019 Reaction Engineering International. All rights reserved.
// Use of this source code is governed by the MIT license in the file LICENSE.txt.

package preferences

import (
	"github.com/reaction-eng/restlib/users"
)

/**
Define an interface for roles
*/
type Repo interface {
	/**
	Get the preferences for this repo
	*/
	GetPreferences(user users.User) (*Preferences, error)

	/**
	Update the User pref
	*/
	SetPreferences(user users.User, userSetting *SettingGroup) (*Preferences, error)

	/**
	Allow databases to be closed
	*/
	CleanUp()
}
