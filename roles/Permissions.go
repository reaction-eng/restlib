package roles

/**
* This package is used to check for roles
 */
type Permissions struct {
	//Store a list of
	Permissions []string `json:"permissions"`

	//And s list of Roles

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
