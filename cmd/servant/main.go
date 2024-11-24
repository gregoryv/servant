package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gregoryv/servant"
)

func main() {
	bind := ":8100"
	fmt.Println("listen", bind)

	sys := servant.NewSystem()

	h := sys.Handler()
	if err := http.ListenAndServe(bind, h); err != nil {
		log.Fatal(err)
	}
}
