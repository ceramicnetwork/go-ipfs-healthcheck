// Package healthcheck runs a server that responds with the status of the IPFS
// node.
package healthcheck

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/ipfs/go-cid"
	coreiface "github.com/ipfs/interface-go-ipfs-core"
)

type Config struct {
	port string
}

type ServerContext struct {
	ipfs coreiface.CoreAPI
}

type ipfsServerContextKey struct {
	key string
}

type status struct {
	Message string
}

func StartServer(port string, ipfs coreiface.CoreAPI) {
	ctx := ServerContext{ipfs}
	http.HandleFunc("/", createHandler(healthcheckHandler, ctx))
	fmt.Println("Healthcheck server listening on port ", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func createHandler(
	fn func(http.ResponseWriter, *http.Request),
	ctx ServerContext,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_ctx := context.WithValue(
			r.Context(),
			ipfsServerContextKey{"ipfs"},
			ctx.ipfs,
		)
		_r := r.Clone(_ctx)
		fn(w, _r)
	}
}

func healthcheckHandler(w http.ResponseWriter, r *http.Request) {
	// Use CID of empty directory which is pinned on all nodes by default
	// https://github.com/ipfs/go-ipfs/issues/8404#issuecomment-917426813
	c, err := cid.Decode("QmUNLLsPACCz1vLxQVkXqqLX5R1X345qqfHbsf67hvA3Nn")
	if err != nil {
		log.Panic(err)
	}

	var failed bool = false
	var _status *status

	ctx := r.Context()
	ipfs, _ := ctx.Value(ipfsServerContextKey{"ipfs"}).(coreiface.CoreAPI)
	nd, err := ipfs.Dag().Get(ctx, c)
	if err != nil {
		failed = true
	} else {
		_, err = nd.Stat()
		if err != nil {
			failed = true
		}
	}

	if failed {
		w.WriteHeader(http.StatusInternalServerError)
		_status = &status{Message: "Health check failed"}
	} else {
		_status = &status{Message: "Health check succeeded"}
	}

	fmt.Fprintf(w, _status.Message)
}
