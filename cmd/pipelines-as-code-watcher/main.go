package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/openshift-pipelines/pipelines-as-code/pkg/params"
	"github.com/openshift-pipelines/pipelines-as-code/pkg/reconciler"

	"knative.dev/pkg/injection/sharedmain"
	"knative.dev/pkg/signals"
)

const globalProbesPort = "8080"

func main() {
	probesPort := globalProbesPort
	envProbePort := os.Getenv("PAC_WATCHER_PORT")
	if envProbePort != "" {
		probesPort = envProbePort
	}

	// set up client/informer overrides for kcp
	ctx := signals.NewContext()
	run := params.New()
	err := run.Clients.NewClients(ctx, &run.Info)
	if err != nil {
		log.Fatal("failed to init clients : ", err)
	}
	run.Informers.Clients = run.Clients
	run.Informers.NewInformers(ctx)

	mux := http.NewServeMux()
	mux.HandleFunc("/live", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		_, _ = fmt.Fprint(w, "ok")
	})
	go func() {
		// start the web server on port and accept requests
		log.Printf("Readiness and health check server listening on port %s", probesPort)
		// timeout values same as default one from triggers eventlistener
		// https://github.com/tektoncd/triggers/blame/b5b0ee1249402187d1ceff68e0b9d4e49f2ee957/pkg/sink/initialization.go#L48-L52
		srv := &http.Server{
			ReadTimeout:  5 * time.Second,
			WriteTimeout: 40 * time.Second,
			Addr:         ":" + probesPort,
			Handler:      mux,
		}
		log.Fatal(srv.ListenAndServe())
	}()

	sharedmain.MainWithContext(ctx, "pac-watcher", reconciler.NewController())
}
