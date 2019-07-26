// Copyright 2019 Reaction Engineering International. All rights reserved.
// Use of this source code is governed by the MIT license in the file LICENSE.txt.

package roles

import (
	"html/template"
	"net/http"
)

/**
Function used to show user documentation
*/
func (handler *Handler) handlePermissionsDocumentation(w http.ResponseWriter, r *http.Request) {

	//Load int he welcome html
	tmpl, _ := template.New("Permissions Api").Parse(getDocumentation())

	//Show it
	tmpl.Execute(w, nil)

}

/**
Return the hard coded documentation
*/
func getDocumentation() string {

	return `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/semantic-ui/2.4.1/semantic.min.css">

</head>
<body>
    <div class="ui container">
        <h2 class="ui header">
            <i class="user icon"></i>
            <div class="content">
                User Api
                <div class="sub header">The user to get the allowed permissions for the logged in user</div>
            </div>
        </h2>
        <!-------User Permission Get ------------------>
        <table class="ui celled striped table">
            <thead>
            <tr>
                <th colspan="3">
                    Get the User Permissions
                </th>
            </tr>
            </thead>
            <tbody>
            <tr>
                <td colspan="3">
                    Get the User Permissions for the current logged in user.
                </td>
            </tr>
            <tr>
                <td>URL</td>
                <td colspan="2">/users/permissions</td>
            </tr>
            <tr>
                <td>Method</td>
                <td colspan="2">Get</td>
            </tr>
            <tr>
                <td>Access</td>
                <td colspan="2">Authorized</td>
            </tr>
            <tr>
                <td>Json Response (Success)</td>
                <td>Code: 200 </td>
                <td>
                    Permission:{<br/>
                    permissions:[string, etc]<br/>
                    }
                </td>
            </tr>
            <tr>
                <td>Json Response (Failure)</td>
                <td>Code: 403 or 422 </td>
                <td>
                    Response:{<br/>
                    status:false<br/>
                    message:string<br/>
                    }
                </td>
            </tr>

            </tbody>
        </table>
        
        <!-- The list of possible erros -->
        <div class="ui segment">
        	<h3 class="ui header">
            	<div class="content">
             	   Possible Error Messages
           	 	</div>
        	</h3>
        	<ul>
        		<li>create_user_added: when a new user is created in the database</li>
        		<li>validate_missing_email: The user did not include the email</li>
        		<li>validate_password_insufficient: password does not meet requirements</li>
        		<li>validate_email_in_use: the email is already in use</li>
        		<li>login_invalid_password: invalid password</li>
        		<li>auth_missing_token: </li>
        		<li>login_user_id_not_found</li>
        		<li>login_email_not_found</li>
        		<li>password_change_success</li>
        		<li>password_change_missing_email</li>
        		<li>password_change_request_received</li>
        		<li>password_change_forbidden</li>
        	</ul>
  			<p></p>
		</div>
    </div>
</body>
</html>
	`

}
