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
        </table><!-------Password Reset ------------------>
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
                <td>Code: 401 </td>
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
                    Password Reset
                </th>
            </tr>
            </thead>
            <tbody>
            <tr>
                <td colspan="3">
                    Method to request a password reset email.
                    The Email is required to match the logged in user.

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
    </div>
</body>
</html>

	`

}
