package users

//a struct to rep user account
type BasicUser struct {
	Id_        int    `json:"id"`
	Email_     string `json:"email"`
	Password_  string `json:"password"`
	Token_     string `json:"token";sql:"-"`
	activated_ bool
}

/**
Add the required setters and getters
*/
func (basic *BasicUser) Id() int {
	return basic.Id_
}
func (basic *BasicUser) SetId(id int) {
	basic.Id_ = id
}
func (basic *BasicUser) Email() string {
	return basic.Email_
}

//func (basic *BasicUser) SetId(id int)  {
//	basic.Id_ = id
//}
func (basic *BasicUser) Password() string {
	return basic.Password_
}
func (basic *BasicUser) SetPassword(pw string) {
	basic.Password_ = pw
}
func (basic *BasicUser) Token() string {
	return basic.Token_
}
func (basic *BasicUser) SetToken(tk string) {
	basic.Token_ = tk
}

func (basic *BasicUser) Activated() bool {
	return basic.activated_
}
