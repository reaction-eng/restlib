// Copyright 2019 Reaction Engineering International. All rights reserved.
// Use of this source code is governed by the MIT license in the file LICENSE.txt.

package roles_test

import (
	"strings"
	"testing"

	"github.com/reaction-eng/restlib/roles"
	"github.com/stretchr/testify/assert"
)

func TestPermissions_AllowedTo(t *testing.T) {
	testCases := []struct {
		permissions     *roles.Permissions
		request         []string
		expectedAllowed bool
	}{
		{
			permissions: &roles.Permissions{
				Permissions: []string{"perm1", "perm2"},
			},
			request:         []string{"perm1"},
			expectedAllowed: true,
		},
		{
			permissions: &roles.Permissions{
				Permissions: []string{"perm1", "perm2"},
			},
			request:         []string{"perm3"},
			expectedAllowed: false,
		},
		{
			permissions: &roles.Permissions{
				Permissions: []string{"perm1", "perm2"},
			},
			request:         []string{"perm1", "perm2"},
			expectedAllowed: true,
		},
		{
			permissions: &roles.Permissions{
				Permissions: []string{"perm1", "perm2"},
			},
			request:         []string{"perm2", "perm1"},
			expectedAllowed: true,
		},
		{
			permissions: &roles.Permissions{
				Permissions: []string{"perm1", "perm2"},
			},
			request:         []string{"perm1", "perm3"},
			expectedAllowed: false,
		},
		{
			permissions: &roles.Permissions{
				Permissions: []string{"perm1", "perm2"},
			},
			request:         []string{},
			expectedAllowed: true,
		},
		{
			permissions: &roles.Permissions{
				Permissions: []string{},
			},
			request:         []string{},
			expectedAllowed: true,
		},
		{
			permissions: &roles.Permissions{
				Permissions: []string{},
			},
			request:         []string{"perm1"},
			expectedAllowed: false,
		},
	}

	for _, testCase := range testCases {
		// arrange
		// act
		allowed := testCase.permissions.AllowedTo(testCase.request...)

		// assert
		assert.Equal(t, testCase.expectedAllowed, allowed, "permissions: "+strings.Join(testCase.permissions.Permissions, ",")+" request: "+strings.Join(testCase.request, ","))
	}
}
