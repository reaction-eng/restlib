package utils

import (
	"context"
	"errors"
)

const UserKey = "user"
const OrganizationKey = "organization"

var UserNotLoggedInError = errors.New("no_user_logged_in")

func UserFromContext(context context.Context) (userId int, organizationId int, err error) {
	//We have gone through the auth, so we should know the id of the logged in user
	loggedInUserString := context.Value(UserKey)
	organizationIdString := context.Value(OrganizationKey)

	if loggedInUserString == nil || organizationIdString == nil {
		return 0, 0, UserNotLoggedInError
	}

	return loggedInUserString.(int), organizationIdString.(int), nil
}
