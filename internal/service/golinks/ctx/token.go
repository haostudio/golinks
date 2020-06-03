package ctx

import (
	"github.com/gin-gonic/gin"
)

// GetToken returns the token cookie.
func GetToken(ctx *gin.Context) (token string, err error) {
	return ctx.Cookie(tokenCookieKey)
}

// SetToken sets the token cookie
func SetToken(ctx *gin.Context, token string, maxAge int) {
	ctx.SetCookie(tokenCookieKey, token, maxAge, "", "", false, false)
}

// DeleteToken deletes the token cookie
func DeleteToken(ctx *gin.Context) {
	ctx.SetCookie(tokenCookieKey, "", 0, "", "", false, false)
}
