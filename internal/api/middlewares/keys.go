package middlewares

// keys for values to save in gin.Context.
const (
	getUserKey = "golinks.middlewares.auth_get_user"
	getOrgKey  = "golinks.middlewares.auth_get_org"
)

// keys for values to save in cookies.
const (
	tokenCookieKey = "GOLINKS_TOKEN"
)
