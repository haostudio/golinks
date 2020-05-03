package middlewares

import (
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	"go.opencensus.io/trace"
)

// Trace requests with opencesnsus.
func Trace(ctx *gin.Context) {
	logger := GetLogger(ctx)

	var span *trace.Span
	reqCtx, span := trace.StartSpan(
		ctx.Request.Context(), ctx.Request.URL.String())
	defer span.End()

	// add attributes
	span.AddAttributes(trace.StringAttribute("type", "http_request"))
	span.AddAttributes(trace.StringAttribute("logger_id", logger.GetID()))
	span.AddAttributes(trace.StringAttribute("client_ip", ctx.ClientIP()))
	span.AddAttributes(trace.StringAttribute("http_method", ctx.Request.Method))
	for k, v := range ctx.Request.Header {
		span.AddAttributes(
			trace.StringAttribute(fmt.Sprintf("header.%s", k), strings.Join(v, ",")),
		)
	}

	// set ctx back to request.
	ctx.Request = ctx.Request.WithContext(reqCtx)

	ctx.Next()

	// get response code
	span.AddAttributes(
		trace.Int64Attribute("status_code", int64(ctx.Writer.Status())),
	)
}
