// Copyright 2019 Reaction Engineering International. All rights reserved.
// Use of this source code is governed by the MIT license in the file LICENSE.txt.

package users

import (
	"html/template"
	"net/http"
)

/**
Function used to show user documentation
*/
func (handler *Handler) handleUserDocumentation(w http.ResponseWriter, r *http.Request) {

	//Load int he welcome html
	tmpl, _ := template.New("Users Api").Parse(getUserDocumentation())

	//Show it
	tmpl.Execute(w, nil)

}

/**
Return the hard coded documentation
*/
func getUserDocumentation() string {

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
                <div class="sub header">The user api is used to create, retrieve and login users..</div>
            </div>
        </h2>
        <!-------User Create ------------------>
        <table class="ui celled striped table">
            <thead>
                <tr>
                    <th colspan="3">
                        User Create
                    </th>
                </tr>
            </thead>
            <tbody>
            <tr>
                <td colspan="3">
                    This method is used to create a new user.
                </td>
            </tr>
            <tr>
                <td>URL</td>
                <td colspan="2">/users/</td>
            </tr>
            <tr>
                <td>Method</td>
                <td colspan="2">POST</td>
            </tr>
            <tr>
                <td>Access</td>
                <td colspan="2">Public</td>
            </tr>
            <tr>
                <td>Json Input</td>
                <td colspan="2">
                    User:{<br/>
                        email:string<br/>
                        password:string<br/>

                    }
                </td>
            </tr>
            <tr>
                <td>Json Response (Success)</td>
                <td>Code: 201 </td>
                <td>
                    Response:{<br/>
                    status:true<br/>
                    message:string
                    }
                </td>
            </tr>
            <tr>
                <td>Json Response (Failure)</td>
                <td>Code: 422 </td>
                <td>
                    Response:{<br/>
                    status:false<br/>
                    message:string
                    }
                </td>
            </tr>

            </tbody>
        </table>
        <!-------User Activate ------------------>
        <table class="ui celled striped table">
            <thead>
                <tr>
                    <th colspan="3">
                        User Activate
                    </th>
                </tr>
            </thead>
            <tbody>
            <tr>
                <td colspan="3">
                    Activate a User
                </td>
            </tr>
            <tr>
                <td>URL</td>
                <td colspan="2">/users/activate</td>
            </tr>
            <tr>
                <td>Method</td>
                <td colspan="2">POST</td>
            </tr>
            <tr>
                <td>Access</td>
                <td colspan="2">Public</td>
            </tr>
            <tr>
                <td>Json Input</td>
                <td colspan="2">
                    User:{<br/>
                        email:string<br/>
                        activation_token:string
                    }
                </td>
            </tr>
            <tr>
                <td>Json Response (Success)</td>
                <td>Code: 200 </td>
                <td>
                    Response:{<br/>
                    status:true<br/>
                    message:string
                    }
                </td>
            </tr>
            <tr>
                <td>Json Response (Failure)</td>
                <td>Code: 422 </td>
                <td>
                    Response:{<br/>
                    status:false<br/>
                    message:string
                    }
                </td>
            </tr>

            </tbody>
        </table>
        <!-------Password Reset ------------------>
        <table class="ui celled striped table">
            <thead>
            <tr>
                <th colspan="3">
                    User Update
                </th>
            </tr>
            </thead>
            <tbody>
            <tr>
                <td colspan="3">
                    Method to update the current user with the changed values specified.
                </td>
            </tr>
            <tr>
                <td>URL</td>
                <td colspan="2">/users/</td>
            </tr>
            <tr>
                <td>Method</td>
                <td colspan="2">PUT</td>
            </tr>
            <tr>
                <td>Access</td>
                <td colspan="2">Authorized</td>
            </tr>
            <tr>
                <td>Json Input</td>
                <td colspan="2">
                    User:{<br/>
                    email:string<br/>
                    ...
                    }
                </td>
            </tr>
            <tr>
                <td>Json Response (Success)</td>
                <td>Code: 202 </td>
                <td>
                    User:{<br/>
                    email:string<br/>
                    ...
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
        </table><!-------Password Reset ------------------>
        <table class="ui celled striped table">
            <thead>
            <tr>
                <th colspan="3">
                    Get User
                </th>
            </tr>
            </thead>
            <tbody>
            <tr>
                <td colspan="3">
                    Gets the current data for the logged in user.
                </td>
            </tr>
            <tr>
                <td>URL</td>
                <td colspan="2">/users/</td>
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
                    User:{<br/>
                    email:string<br/>
                    ...
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
        <!-------User Create ------------------>
        <table class="ui celled striped table">
            <thead>
            <tr>
                <th colspan="3">
                    User Login
                </th>
            </tr>
            </thead>
            <tbody>
            <tr>
                <td colspan="3">
                    This method allows the user to login
                </td>
            </tr>
            <tr>
                <td>URL</td>
                <td colspan="2">/users/login</td>
            </tr>
            <tr>
                <td>Method</td>
                <td colspan="2">POST</td>
            </tr>
            <tr>
                <td>Access</td>
                <td colspan="2">Public</td>
            </tr>
            <tr>
                <td>Json Input</td>
                <td colspan="2">
                    User:{<br/>
                    email:string<br/>
                    password:string<br/>
                    }
                </td>
            </tr>
            <tr>
                <td>Json Response (Success)</td>
                <td>Code: 201 </td>
                <td>
                    User:{<br/>
                    email:string<br/>
                    id:int<br/>
                    token:string<br/>
                    }
                </td>
            </tr>
            <tr>
                <td>Json Response (Failure)</td>
                <td>Code: 422 or 403 </td>
                <td>
                    Response:{<br/>
                    message:string
                    }
                </td>
            </tr>

            </tbody>
        </table>
        <!-------Password Change ------------------>
        <table class="ui celled striped table">
            <thead>
            <tr>
                <th colspan="3">
                    Password Change
                </th>
            </tr>
            </thead>
            <tbody>
            <tr>
                <td colspan="3">
                    Method to change the current user's password.
                    The Email is required to match the logged in user.
                </td>
            </tr>
            <tr>
                <td>URL</td>
                <td colspan="2">/users/password/change</td>
            </tr>
            <tr>
                <td>Method</td>
                <td colspan="2">POST</td>
            </tr>
            <tr>
                <td>Access</td>
                <td colspan="2">Authorized</td>
            </tr>
            <tr>
                <td>Json Input</td>
                <td colspan="2">
                    {<br/>
                    email:string<br/>
                    passwordold:string<br/>
                    password:string<br/>
                    }
                </td>
            </tr>
            <tr>
                <td>Json Response (Success)</td>
                <td>Code: 202 </td>
                <td>
                    Response{<br/>
                    status:true<br/>
                    token:string<br/>
                    }
                </td>
            </tr>
            <tr>
                <td>Json Response (Failure)</td>
                <td>Code: 404 </td>
                <td>
                    Response:{<br/>
                    status:false<br/>
                    message:string<br/>
                    }
                </td>
            </tr>

            </tbody>
        </table>
        <!-------Password Reset ------------------>
        <table class="ui celled striped table">
            <thead>
            <tr>
                <th colspan="3">
                    Password Reset Get
                </th>
            </tr>
            </thead>
            <tbody>
            <tr>
                <td colspan="3">
                    Method to request a password reset email.
                </td>
            </tr>
            <tr>
                <td>URL</td>
                <td colspan="2">/users/password/reset</td>
            </tr>
            <tr>
                <td>Method</td>
                <td colspan="2">GET</td>
            </tr>
            <tr>
                <td>Access</td>
                <td colspan="2">Public</td>
            </tr>
            <tr>
                <td>URL Param</td>
                <td colspan="2">
                    email=string
                </td>
            </tr>
            <tr>
                <td>Json Response (Success)</td>
                <td>Code: 201 </td>
                <td>
                    Response{<br/>
                    status:true<br/>
                    token:string<br/>
                    }
                </td>
            </tr>
            <tr>
                <td>Json Response (Failure)</td>
                <td>Code: 404 </td>
                <td>
                    Response:{<br/>
                    status:false<br/>
                    message:string<br/>
                    }
                </td>
            </tr>

            </tbody>
        </table>
        <!-------Password Reset ------------------>
        <table class="ui celled striped table">
            <thead>
            <tr>
                <th colspan="3">
                    Get User Activation Token
                </th>
            </tr>
            </thead>
            <tbody>
            <tr>
                <td colspan="3">
                    Method to request a new user activation token
                </td>
            </tr>
            <tr>
                <td>URL</td>
                <td colspan="2">/users/activation/</td>
            </tr>
            <tr>
                <td>Method</td>
                <td colspan="2">GET</td>
            </tr>
            <tr>
                <td>Access</td>
                <td colspan="2">Public</td>
            </tr>
            <tr>
                <td>URL Param</td>
                <td colspan="2">
                    email=string
                </td>
            </tr>
            <tr>
                <td>Json Response (Success)</td>
                <td>Code: 201 </td>
                <td>
                    Response{<br/>
                    status:true<br/>
                    message:string<br/>
                    }
                </td>
            </tr>
            <tr>
                <td>Json Response (Failure)</td>
                <td>Code: 404 </td>
                <td>
                    Response:{<br/>
                    status:false<br/>
                    message:string<br/>
                    }
                </td>
            </tr>

            </tbody>
        </table>
		<!-------Password Reset ------------------>
        <table class="ui celled striped table">
            <thead>
            <tr>
                <th colspan="3">
                    Password Reset Post
                </th>
            </tr>
            </thead>
            <tbody>
            <tr>
                <td colspan="3">
                    Method to request change a password after a reset
                </td>
            </tr>
            <tr>
                <td>URL</td>
                <td colspan="2">/users/password/reset</td>
            </tr>
            <tr>
                <td>Method</td>
                <td colspan="2">POST</td>
            </tr>
            <tr>
                <td>Access</td>
                <td colspan="2">Public</td>
            </tr>
            <tr>
                <td>Json Input</td>
                <td colspan="2">
                    {<br/>
                        email:string<br/>
						reset_token:string<br/>
						password:string<br/>
                    }
                </td>
            </tr>
            <tr>
                <td>Json Response (Success)</td>
                <td>Code: 201 </td>
                <td>
                    Response{<br/>
                    status:true<br/>
                    token:string<br/>
                    }
                </td>
            </tr>
            <tr>
                <td>Json Response (Failure)</td>
                <td>Code: 400 </td>
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
				<li>user_not_activated</li>
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
