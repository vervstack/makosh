package makosh_resolver

import (
	"encoding/json"
	"io"
	"net/http"
	"sync"

	errors "github.com/Red-Sock/trace-errors"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/resolver"

	"github.com/godverv/makosh/pkg/makosh_be"
)

type UpdateAddresses func(addrs []string) error

type Resolver struct {
	getAddressesRequest *http.Request

	m sync.RWMutex

	doUpdate UpdateAddresses

	log logrus.StdLogger

	overrides []string
}

func (r *Resolver) ResolveNow(_ resolver.ResolveNowOptions) {
	err := r.Resolve()
	if err != nil {
		r.log.Println("error resolving connections", err.Error())
	}
}

func (r *Resolver) Close() {}

func (r *Resolver) Resolve() error {
	if len(r.overrides) != 0 {
		return r.doUpdate(r.overrides)
	}

	httpResp, err := http.DefaultClient.Do(r.getAddressesRequest)
	if err != nil {
		return errors.Wrap(err, "error getting ")
	}

	body, err := io.ReadAll(httpResp.Body)
	if err != nil {
		return errors.Wrap(err, "error reading body")
	}

	if httpResp.StatusCode != http.StatusOK {
		return errors.New(string(body))
	}

	var endpointsResponse makosh_be.ListEndpoints_Response
	err = json.Unmarshal(body, &endpointsResponse)
	if err != nil {
		return errors.Wrap(err, "error parsing list endpoints response")
	}

	err = r.doUpdate(endpointsResponse.Urls)
	if err != nil {
		return errors.Wrap(err, "error updating connection state")
	}

	return nil
}
