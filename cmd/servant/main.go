package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gregoryv/servant"
	"github.com/gregoryv/servant/htsec"
	"github.com/gregoryv/servant/htsec/github"
	"github.com/gregoryv/servant/htsec/google"
)

func main() {
	sys := servant.NewSystem()

	guard := htsec.NewGuard(
		github.Default(),
		google.Default(),
	)
	sys.SetSecurity(guard)

	srv := http.Server{
		Addr:    ":8100",
		Handler: servant.NewRouter(sys),
	}

	fmt.Println("listen", srv.Addr)
	if err := srv.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
