package resolver

import (
	"context"
	"io"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	errors "github.com/Red-Sock/trace-errors"
	"google.golang.org/grpc/resolver"
)

const (
	MakoshSourceURL = "MAKOSH_SRC_URL"
)

var ErrNoMakoshSourceURL = errors.New("no " + MakoshSourceURL + "is defined")

type Resolver struct {
	ctx    context.Context
	cancel func()

	target     string
	clientConn resolver.ClientConn

	srcUrl string

	wg sync.WaitGroup
}

func (r *Resolver) ResolveNow(opts resolver.ResolveNowOptions) {

}

func (r *Resolver) Close() {

}

func (r *Resolver) watch() {
	defer r.wg.Done()

	r.lookup(r.target)
	ticker := time.NewTicker(5 * time.Second)
	for {
		select {
		case <-r.ctx.Done():
			return
		case <-ticker.C:
			r.lookup(r.target)
		}
	}
}

func (r *Resolver) lookup(target string) {
	resp, err := http.Get(r.target)
	if err != nil {
		log.Println(err)
		return
	}
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return
	}
	addrs := strings.Split(string(data), " ")
	log.Println("resolver made request:", addrs)
	newAddrs := make([]resolver.Address, 0, len(addrs))
	for _, a := range addrs {
		newAddrs = append(newAddrs, resolver.Address{Addr: a})
	}
	// Обновляем адреса в ClientConn
	err = r.clientConn.UpdateState(resolver.State{
		Addresses:     newAddrs,
		ServiceConfig: defaultServiceConfig,
	})
	if err != nil {
		log.Println(err)
		return
	}
}
