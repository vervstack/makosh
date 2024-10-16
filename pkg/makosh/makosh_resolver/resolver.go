package makosh_resolver

import (
	"encoding/json"
	stderrs "errors"
	"io"
	"net/http"
	"sync"

	errors "github.com/Red-Sock/trace-errors"

	"github.com/godverv/makosh/pkg/makosh_be"
)

type UpdateAddresses func(addrs []string) error

type EndpointsResolver interface {
	AddUpdateCallbacks(callbacks ...UpdateAddresses)
	Resolve() error
	GetUpdaters() []UpdateAddresses
}

type MakoshResolver struct {
	// updaters - contains callbacks for places where this resolver have to update addresses
	// when Resolve is called
	updaters []UpdateAddresses
	m        sync.Mutex

	req *http.Request
}

func NewMakoshResolver(r *http.Request) *MakoshResolver {
	return &MakoshResolver{
		req: r,
	}
}

func (r *MakoshResolver) AddUpdateCallbacks(updaters ...UpdateAddresses) {
	r.m.Lock()
	r.updaters = append(r.updaters, updaters...)
	r.m.Unlock()
}

func (r *MakoshResolver) GetUpdaters() []UpdateAddresses {
	return r.getUpdaters()
}

func (r *MakoshResolver) Resolve() error {
	updaters := r.getUpdaters()
	if len(updaters) == 0 {
		return nil
	}

	r.m.Lock()
	defer r.m.Unlock()

	ep, err := r.fetchEndpoints()
	if err != nil {
		return errors.Wrap(err)
	}

	var updateErrs []error
	for _, doUpdate := range updaters {
		err = doUpdate(ep)
		if err != nil {
			updateErrs = append(updateErrs, err)
		}
	}

	if len(updateErrs) == 0 {
		return nil
	}

	return stderrs.Join(updateErrs...)
}

func (r *MakoshResolver) fetchEndpoints() ([]string, error) {
	resp, err := http.DefaultClient.Do(r.req)
	if err != nil {
		return nil, errors.Wrap(err, "error getting ")
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "error reading body")
	}

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(string(body))
	}

	var endpointsResponse makosh_be.ListEndpoints_Response
	err = json.Unmarshal(body, &endpointsResponse)
	if err != nil {
		return nil, errors.Wrap(err, "error parsing list endpoints response")
	}

	return endpointsResponse.Urls, nil
}

func (r *MakoshResolver) getUpdaters() []UpdateAddresses {
	r.m.Lock()
	out := make([]UpdateAddresses, len(r.updaters))
	copy(out, r.updaters)
	r.m.Unlock()

	return out
}
