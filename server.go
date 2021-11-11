// Package healthcheck runs a server that responds with the status of the IPFS
// node.
package healthcheck

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/ipfs/go-cid"
	coreiface "github.com/ipfs/interface-go-ipfs-core"
)

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
	server := http.Server{
		Addr: ":" + port,
	}

	ctx := ServerContext{ipfs}
	http.HandleFunc("/", createHandler(ctx, healthcheckHandler))

	// Shutdown gracefully

	idleConnsClosed := make(chan struct{})

	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt)
		signal.Notify(sigint, syscall.SIGTERM)
		<-sigint

		fmt.Println("Healthcheck server shutting down...")
		close(idleConnsClosed)
		if err := server.Shutdown(context.Background()); err != nil {
			log.Printf("Healthcheck server error on Shutdown: %+v", err)
		}
	}()

	fmt.Println("Healthcheck server listening on port", port)
	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		log.Printf("Healthcheck server error on ListenAndServe: %+v", err)
	}

	<-idleConnsClosed
}

func createHandler(
	ctx ServerContext,
	fn func(http.ResponseWriter, *http.Request),
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		updatedCtx := context.WithValue(
			r.Context(),
			ipfsServerContextKey{"ipfs"},
			ctx.ipfs,
		)
		updatedRequest := r.Clone(updatedCtx)
		fn(w, updatedRequest)
	}
}

func healthcheckHandler(w http.ResponseWriter, r *http.Request) {
	// Use CID of empty directory which is pinned on all nodes by default
	// https://github.com/ipfs/go-ipfs/issues/8404#issuecomment-917426813
	c, err := cid.Decode("QmUNLLsPACCz1vLxQVkXqqLX5R1X345qqfHbsf67hvA3Nn")
	if err != nil {
		log.Panic(err)
	}

	var failed = false
	var healthCheck *status

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
		healthCheck = &status{Message: "Health check failed"}
	} else {
		healthCheck = &status{Message: "Health check succeeded"}
	}

	_, _ = fmt.Fprintf(w, healthCheck.Message)
}
