package service

import (
	"fmt"
	"net"

	"github.com/soheilhy/cmux"
)

// Mux defines a service multiplexer.
type Mux interface {
	Append(svc Service)
	Serve(ls net.Listener) error
}

// Service defines a service interface.
type Service interface {
	fmt.Stringer
	Matchers() []cmux.Matcher
	Serve(ls net.Listener) error
}
