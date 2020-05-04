package service

import (
	"net"
	"sync"

	"github.com/popodidi/log"
	"github.com/soheilhy/cmux"
)

// NewMux returns a service mux with services.
func NewMux(logger log.Logger, services ...Service) Mux {
	return &mux{
		logger:   logger,
		services: services,
	}
}

type mux struct {
	logger   log.Logger
	services []Service
}

func (m *mux) Append(svc Service) {
	m.services = append(m.services, svc)
}

func (m *mux) Serve(ls net.Listener) error {
	var wg sync.WaitGroup
	mux := cmux.New(ls)
	for _, service := range m.services {
		serviceListener := mux.Match(service.Matchers()...)
		wg.Add(1)
		go func(svc Service, ls net.Listener) {
			defer wg.Done()
			err := svc.Serve(ls)
			if err != nil {
				m.logger.Warn("%s service stopped: %v", svc, err)
			}
		}(service, serviceListener)
	}
	wg.Add(1)
	var err error
	go func() {
		defer wg.Done()
		err = mux.Serve()
	}()
	wg.Wait()
	return err
}
