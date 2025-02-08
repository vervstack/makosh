package makosh_resolver

import (
	"encoding/json"
	"io"
	"net/http"

	errors "github.com/Red-Sock/trace-errors"

	"github.com/godverv/makosh/pkg/makosh_be"
)

type MakoshResolver struct {
	EndpointsContainer

	req *http.Request
}

func NewMakoshResolver(r *http.Request) *MakoshResolver {
	return &MakoshResolver{
		req: r,
	}
}

func (r *MakoshResolver) Resolve() error {
	updaters := r.getUpdaters()
	if len(updaters) == 0 {
		return nil
	}

	err := r.fetchEndpoints()
	if err != nil {
		return errors.Wrap(err)
	}

	err = r.notifySubscribers()
	if err != nil {
		return errors.Wrap(err)
	}

	return nil
}

func (r *MakoshResolver) fetchEndpoints() error {
	r.m.Lock()
	defer r.m.Unlock()

	resp, err := http.DefaultClient.Do(r.req)
	if err != nil {
		return errors.Wrap(err, "error getting response via http request")
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return errors.Wrap(err, "error reading body")
	}

	if resp.StatusCode != http.StatusOK {
		return errors.New(string(body))
	}

	var endpointsResponse makosh_be.ListEndpoints_Response
	err = json.Unmarshal(body, &endpointsResponse)
	if err != nil {
		return errors.Wrap(err, "error parsing list endpoints response")
	}

	r.addrsCache = endpointsResponse.Urls

	return nil
}
