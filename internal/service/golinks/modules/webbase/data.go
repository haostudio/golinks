package webbase

import "github.com/haostudio/golinks/internal/version"

// Data defines the base data.
type Data struct {
	Title      string
	BinVersion string
}

// NewData returns a database.
func NewData(title string) Data {
	return Data{
		Title:      title,
		BinVersion: version.Version(),
	}
}

// ErrorPage defines the data for error.html template.
type ErrorPage struct {
	Data
	StatusCode  int
	ErrTitle    string
	Description string
}

// NewErrorPage returns error page data.
func NewErrorPage() ErrorPage {
	return ErrorPage{
		Data: NewData("Golinks - Error"),
	}
}
