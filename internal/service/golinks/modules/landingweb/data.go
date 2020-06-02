package landingweb

import (
	"github.com/gin-gonic/gin"

	"github.com/haostudio/golinks/internal/service/golinks/modules/webbase"
)

// Data defines the data for landing.html template.
type Data struct {
	webbase.Data
}

// NewData returns landing page data.
func NewData(ctx *gin.Context) Data {
	return Data{
		Data: webbase.NewData("Golinks", ctx),
	}
}
