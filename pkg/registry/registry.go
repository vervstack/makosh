package registry

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"

	"go.vervstack.ru/makosh/pkg/makosh_be"
)

var ErrNotFound = errors.New("not found")

type Registry struct {
	m    sync.RWMutex
	data map[string]*makosh_be.Endpoint
}

func New() *Registry {
	return &Registry{
		data: make(map[string]*makosh_be.Endpoint),
	}
}

func ToStatusError(err error) error {
	if err == nil {
		return nil
	}

	if errors.Is(err, ErrNotFound) {
		return status.Error(codes.NotFound, err.Error())
	}

	return status.Error(codes.Internal, err.Error())
}

func (r *Registry) Upsert(_ context.Context, endpoints []*makosh_be.Endpoint) error {
	r.m.Lock()
	defer r.m.Unlock()

	for _, endpoint := range endpoints {
		if endpoint == nil {
			continue
		}

		stored := r.data[endpoint.GetServiceName()]
		if stored == nil {
			stored = &makosh_be.Endpoint{ServiceName: endpoint.GetServiceName()}
			r.data[endpoint.GetServiceName()] = stored
		}

		isScopedOnly := len(endpoint.GetAddrs()) == 0 && len(endpoint.GetAddresses()) > 0
		if !isScopedOnly {
			stored.Addrs = append([]string(nil), endpoint.GetAddrs()...)
		}
		stored.Addresses = replaceScopes(stored.GetAddresses(), endpoint.GetAddresses())
	}

	return nil
}

func (r *Registry) List(
	_ context.Context, serviceName string, scope makosh_be.AddressScope,
) (*makosh_be.ListEndpoints_Response, error) {
	r.m.RLock()
	stored := r.data[serviceName]
	var endpoint *makosh_be.Endpoint
	if stored != nil {
		endpoint = proto.Clone(stored).(*makosh_be.Endpoint)
	}
	r.m.RUnlock()

	if endpoint == nil {
		return nil, fmt.Errorf("service %s: %w", serviceName, ErrNotFound)
	}

	response := &makosh_be.ListEndpoints_Response{}

	seen := make(map[string]struct{})
	appendUrl := func(url string) {
		_, isSeen := seen[url]
		if isSeen {
			return
		}
		seen[url] = struct{}{}
		response.Urls = append(response.Urls, url)
	}

	if scope == makosh_be.AddressScope_ADDRESS_SCOPE_UNSPECIFIED {
		for _, addr := range endpoint.GetAddrs() {
			appendUrl(addr)
		}
	}

	for _, address := range endpoint.GetAddresses() {
		isMatching := scope == makosh_be.AddressScope_ADDRESS_SCOPE_UNSPECIFIED || address.GetScope() == scope
		if !isMatching {
			continue
		}

		response.Addresses = append(response.Addresses, address)
		appendUrl(address.GetAddr())
	}

	return response, nil
}

func (r *Registry) Delete(_ context.Context, serviceName string, scope makosh_be.AddressScope) error {
	r.m.Lock()
	defer r.m.Unlock()

	if scope == makosh_be.AddressScope_ADDRESS_SCOPE_UNSPECIFIED {
		delete(r.data, serviceName)
		return nil
	}

	stored := r.data[serviceName]
	if stored == nil {
		return nil
	}

	kept := make([]*makosh_be.Address, 0, len(stored.GetAddresses()))
	for _, address := range stored.GetAddresses() {
		if address.GetScope() != scope {
			kept = append(kept, address)
		}
	}
	stored.Addresses = kept

	return nil
}

func replaceScopes(stored, incoming []*makosh_be.Address) []*makosh_be.Address {
	replaced := make(map[makosh_be.AddressScope]struct{}, len(incoming))
	for _, address := range incoming {
		replaced[address.GetScope()] = struct{}{}
	}

	result := make([]*makosh_be.Address, 0, len(stored)+len(incoming))
	for _, address := range stored {
		_, isReplaced := replaced[address.GetScope()]
		if !isReplaced {
			result = append(result, address)
		}
	}

	for _, address := range incoming {
		result = append(result, proto.Clone(address).(*makosh_be.Address))
	}

	return result
}
