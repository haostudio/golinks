package authweb

import "github.com/haostudio/golinks/internal/service/golinks/modules/webbase"

// User defines a user data for template
type User struct {
	Email string
}

// PageData defines the data for links.html template.
type PageData struct {
	webbase.Data

	FormInputEmail string
	FormBtnAction  string
	FormInputName  string

	Users []User
	Admin bool
}

// LoginData defines the data for login.html template.
type LoginData struct {
	webbase.Data

	FormInputEmail        string
	FormInputPassword     string
	FormInputCallback     string
	FormInputAction       string
	FormLoginBtnAction    string
	FormRegisterBtnAction string

	Callback string
}
