// Copyright 2019 Reaction Engineering International. All rights reserved.
// Use of this source code is governed by the MIT license in the file LICENSE.txt.

package preferences

//go:generate mockgen -destination=../mocks/mock_preferencesRepo.go -package=mocks -mock_names Repo=MockPreferencesRepo github.com/reaction-eng/restlib/preferences Repo

import (
	"github.com/reaction-eng/restlib/users"
)

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
