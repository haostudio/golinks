package landingweb

import "github.com/haostudio/golinks/cmd/golinks/server/modules/webbase"

// Data defines the data for landing.html template.
type Data struct {
	webbase.Data
}

// NewData returns landing page data.
func NewData() Data {
	return Data{
		Data: webbase.NewData("Golinks"),
	}
}
