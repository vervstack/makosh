package makosh_resolver

import (
	"net/http"
	"os"
	"strings"

	errors "github.com/Red-Sock/trace-errors"

	"github.com/godverv/makosh/internal/interceptors"
)

const (
	MakoshURL    = "MAKOSH_URL"
	MakoshSecret = "MAKOSH_SECRET"
)

const Header = "Grpc-Metadata-" + interceptors.AuthHeader

var (
	ErrNoServiceDiscoveryURL = errors.New("no " + MakoshURL + " is defined. " +
		"Specify it as environment variable or as and argument for constructor ")
	ErrNoMakoshSecret = errors.New("no " + MakoshSecret + " is defined. " +
		"Specify it as environment variable or as and argument for constructor ")
)

type ResolverBuilder interface {
	NewResolver(target string) (EndpointsResolver, error)
}

type RemoteResolverBuilder struct {
	remoteServiceDiscoveryURL string
	secret                    string
	isPublic                  bool

	protoc string
}

func NewBuilder(opts ...opt) (*RemoteResolverBuilder, error) {
	rsd := &RemoteResolverBuilder{
		remoteServiceDiscoveryURL: os.Getenv(MakoshURL),
		secret:                    os.Getenv(MakoshSecret),
		isPublic:                  false,
	}

	for _, o := range opts {
		o(rsd)
	}

	if rsd.remoteServiceDiscoveryURL == "" {
		return nil, errors.Wrap(ErrNoServiceDiscoveryURL, "Need url to connect to makosh")
	}
	if rsd.secret == "" && !rsd.isPublic {
		return nil, errors.Wrap(ErrNoMakoshSecret, "Need secret to connect to address resolver or "+
			"service discovery to be public")
	}

	if rsd.remoteServiceDiscoveryURL[len(rsd.remoteServiceDiscoveryURL)-1] != '/' {
		rsd.remoteServiceDiscoveryURL += "/"
	}

	protocEndIdx := strings.Index(rsd.remoteServiceDiscoveryURL, "://")
	if protocEndIdx == -1 {
		rsd.protoc = "http"
	} else {
		rsd.protoc = rsd.remoteServiceDiscoveryURL[:protocEndIdx]
	}

	return rsd, nil
}

func (r *RemoteResolverBuilder) NewResolver(target string) (EndpointsResolver, error) {
	if target == "makosh" {
		return EndpointsResolver{
			Resolver: NewStaticResolver(r.remoteServiceDiscoveryURL[:len(r.remoteServiceDiscoveryURL)-1]),
		}, nil
	}

	req, err := http.NewRequest(
		http.MethodGet,
		r.protoc+"://"+r.remoteServiceDiscoveryURL+"api/endpoints/"+target,
		nil)
	if err != nil {
		return EndpointsResolver{}, errors.Wrap(err, "error creating http request")
	}

	if !r.isPublic {
		req.Header.Set(Header, r.secret)
	}

	makoshResolver := NewMakoshResolver(req)

	return EndpointsResolver{
		Resolver: makoshResolver,
	}, nil
}
