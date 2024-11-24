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

	if err := http.ListenAndServe(bind, sys); err != nil {
		log.Fatal(err)
	}
}
