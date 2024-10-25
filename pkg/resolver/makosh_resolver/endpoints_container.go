package makosh_resolver

import (
	stderrs "errors"
	"sync"

	errors "github.com/Red-Sock/trace-errors"
)

type Resolver interface {
	//AddSubscribers - adds callback which will be called whenever addresses is updated
	AddSubscribers(callbacks ...DoUpdateAddress)
	//GetSubscribers - returns subscribers added with AddUpdateCallbacks
	GetSubscribers() []DoUpdateAddress

	//Resolve - does something in order to update addresses
	// and notify subscribers (added with AddUpdateCallbacks)
	Resolve() error

	//SetAddrs - sets addresses manually and notifies subscribers
	SetAddrs(addrs ...string) error
	//GetAddrs - returns current addresses from last resolving (or resolves and returns)
	GetAddrs() []string
}

type EndpointsResolver struct {
	Resolver
}

type DoUpdateAddress func(addrs []string) error

type EndpointsContainer struct {
	// updaters - contains callbacks for places where this resolver have to update addresses
	// when Resolve is called
	m          sync.Mutex
	updaters   []DoUpdateAddress
	addrsCache []string
}

func (r *EndpointsContainer) AddSubscribers(updaters ...DoUpdateAddress) {
	r.m.Lock()
	r.updaters = append(r.updaters, updaters...)
	r.m.Unlock()
}

func (r *EndpointsContainer) GetSubscribers() []DoUpdateAddress {
	return r.getUpdaters()
}

func (r *EndpointsContainer) SetAddrs(addrs ...string) error {
	r.m.Lock()
	r.addrsCache = addrs
	r.m.Unlock()

	err := r.notifySubscribers()
	if err != nil {
		return errors.Wrap(err)
	}

	return nil
}

func (r *EndpointsContainer) GetAddrs() []string {
	r.m.Lock()
	addrs := make([]string, len(r.addrsCache))
	copy(addrs, r.addrsCache)
	r.m.Unlock()

	return addrs
}

func (r *EndpointsContainer) notifySubscribers() error {
	r.m.Lock()
	defer r.m.Unlock()

	var updateErrs []error
	for _, doUpdate := range r.updaters {
		err := doUpdate(r.addrsCache)
		if err != nil {
			updateErrs = append(updateErrs, err)
		}
	}

	if len(updateErrs) == 0 {
		return nil
	}

	return stderrs.Join(updateErrs...)
}

func (r *EndpointsContainer) getUpdaters() []DoUpdateAddress {
	r.m.Lock()
	out := make([]DoUpdateAddress, len(r.updaters))
	copy(out, r.updaters)
	r.m.Unlock()

	return out
}
