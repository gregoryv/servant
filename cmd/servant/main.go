package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gregoryv/servant"
	"github.com/gregoryv/servant/htsec"
)

func main() {
	sys := servant.NewSystem()
	sec := htsec.NewSecure()
	sec.Use(github)
	srv := http.Server{
		Addr:    ":8100",
		Handler: servant.NewRouter(sys, sec),
	}

	fmt.Println("listen", srv.Addr)
	if err := srv.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
