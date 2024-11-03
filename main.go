package main

import (
	"log"
	"net/http"
)

func main() {
	bind := ":8100"
	debug.Println("listen", bind)

	h := logware(
		AuthLayer(
			Endpoints(),
		),
	)

	if err := http.ListenAndServe(bind, h); err != nil {
		log.Fatal(err)
	}
}
