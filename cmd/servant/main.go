package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gregoryv/servant"
)

func main() {
	sys := servant.NewSystem()
	srv := http.Server{
		Addr:    ":8100",
		Handler: servant.NewRouter(sys),
	}
	fmt.Println("listen", srv.Addr)
	if err := srv.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
