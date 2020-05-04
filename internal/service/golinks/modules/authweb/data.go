package authweb

import "github.com/haostudio/golinks/internal/service/golinks/modules/webbase"

// User defines a user data for template
type User struct {
	Email string
}

// PageData defines the data for links.html template.
type PageData struct {
	webbase.Data

	FormInputEmail    string
	FormInputPassword string
	FormBtnAction     string
	FormInputName     string

	Users []User
	Admin bool
}
