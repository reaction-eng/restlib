// Copyright 2019 Reaction Engineering International. All rights reserved.
// Use of this source code is governed by the MIT license in the file LICENSE.txt.

package roles

/**
* This package is used to check for roles
 */
type Permissions struct {
	//Store a list of
	Permissions []string `json:"permissions"`
}

/**
Check to see if the user has permission to do something

*/
func (perm *Permissions) AllowedTo(tasks ...string) bool {
	//March over each task
	for _, task := range tasks {
		//See if it is in my list of permissions
		if !contains(perm.Permissions, task) {
			return false
		}
	}

	//I Guess we can
	return true

}

/**
Write a little support function for
*/
func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
