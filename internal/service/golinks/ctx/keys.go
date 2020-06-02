package ctx

// keys for values to save in gin.Context.
const (
	ctxKey = "golinks.middlewares.ctx"

	userKey = "golinks.middlewares.user"
	orgKey  = "golinks.middlewares.org"
)

// keys for values to save in cookies.
const (
	tokenCookieKey = "GOLINKS_TOKEN"
)
