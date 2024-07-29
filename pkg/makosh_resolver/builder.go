package makosh_resolver

import (
	"net/http"
	"os"

	errors "github.com/Red-Sock/trace-errors"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/resolver"

	"github.com/godverv/makosh/internal/interceptors"
)

const Schema = "verv"

const (
	MakoshURL    = "MAKOSH_URL"
	MakoshSecret = "MAKOSH_SECRET"
)

const Header = "Grpc-Metadata-" + interceptors.Header

var (
	ErrNoMakoshSourceURL = errors.New("no " + MakoshURL + "is defined")
	ErrNoMakoshSecret    = errors.New("no " + MakoshSecret + "is defined")
)

type connectionType int

const (
	GrpcConnection connectionType = iota
	HttpConnection
)

type Builder struct {
	makoshUrl string
	secret    string
	logger    logrus.StdLogger

	overrides map[string][]string
}

type opt func(b *Builder)

func New(opts ...opt) (*Builder, error) {
	b := &Builder{
		makoshUrl: os.Getenv(MakoshURL),
		secret:    os.Getenv(MakoshSecret),
		logger:    logrus.New(),
		overrides: make(map[string][]string),
	}

	for _, o := range opts {
		o(b)
	}

	if b.makoshUrl == "" {
		return nil, errors.Wrap(ErrNoMakoshSourceURL, "Need url to connect to makosh")
	}
	if b.secret == "" {
		return nil, errors.Wrap(ErrNoMakoshSecret, "Need secret to connect to address resolver")
	}

	return b, nil
}

// Build - for grpc resolving
func (b *Builder) Build(t resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error) {
	return b.newResolver(t.URL.Host, updateGrpc(cc))
}

func (b *Builder) BuildHTTPResolver(targetName string, addressUpdater UpdateAddresses) (*Resolver, error) {
	return b.newResolver(targetName, addressUpdater)
}

func (b *Builder) Scheme() string {
	return Schema
}

func (b *Builder) newResolver(targetName string, addressUpdater UpdateAddresses) (*Resolver, error) {
	r := &Resolver{
		doUpdate:  addressUpdater,
		overrides: b.overrides[targetName],
	}
	url := b.makoshUrl + "/v1/endpoints/" + targetName

	var err error
	r.getAddressesRequest, err = http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, errors.Wrap(err, "error creating http request")
	}

	r.getAddressesRequest.Header.Set(Header, b.secret)

	err = r.Resolve()
	if err != nil {
		return nil, errors.Wrap(err, "error resolving addresses")
	}

	return r, nil
}

func updateGrpc(cc resolver.ClientConn) UpdateAddresses {
	return func(addrs []string) error {
		var state resolver.State

		state.Addresses = make([]resolver.Address, 0, len(addrs))

		for _, a := range addrs {
			state.Addresses = append(state.Addresses, resolver.Address{Addr: a})
		}

		return cc.UpdateState(state)
	}
}
