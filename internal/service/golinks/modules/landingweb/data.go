package landingweb

import "github.com/haostudio/golinks/internal/service/golinks/modules/webbase"

// Data defines the data for landing.html template.
type Data struct {
	webbase.Data
	AuthEnabled bool
}

// NewData returns landing page data.
func NewData(authEnabled bool) Data {
	return Data{
		Data:        webbase.NewData("Golinks"),
		AuthEnabled: authEnabled,
	}
}
