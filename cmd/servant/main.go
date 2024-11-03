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

	h := servant.New()

	if err := http.ListenAndServe(bind, h); err != nil {
		log.Fatal(err)
	}
}
