package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gregoryv/servant"
)

func main() {
	sys := servant.NewSystem()
	srv := http.Server{
		Addr:    ":8100",
		Handler: servant.NewRouter(sys),
	}
	// graceful shutdown support
	shutdownComplete, stop := sys.Shutdown()
	srv.RegisterOnShutdown(stop)
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c,
			syscall.SIGTERM, // needed for systemd
			syscall.SIGHUP, os.Kill, os.Interrupt,
		)
		<-c
		srv.Shutdown(context.Background())

	}()

	log.SetFlags(0)
	log.Println("listen", srv.Addr)
	if err := srv.ListenAndServe(); err != nil {
		if !errors.Is(err, http.ErrServerClosed) {
			log.Fatal(err)
		}
		// otherwise closing down gracefully
	}
	// wait for shutdown to complete
	<-shutdownComplete
}
