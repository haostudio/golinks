package middlewares

// keys for values to save in gin.Context.
const (
	authManagerKey = "golinks.middlewares.auth_manager"
	userKey        = "golinks.middlewares.auth_user"
	orgKey         = "golinks.middlewares.auth_org"
)

// keys for values to save in cookies.
const (
	tokenCookieKey = "GOLINKS_TOKEN"
)
