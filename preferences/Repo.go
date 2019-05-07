package preferences

import (
	"bitbucket.org/reidev/restlib/users"
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
