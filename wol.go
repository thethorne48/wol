package main

import (
	"fmt"
	"net/http"
)

var port int

func main() {
	port = 8080
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		var mac = r.URL.Query().Get("mac")
		fmt.Fprintf(w, "Broadcasting WOL request")
		if err := SendMagicPacket(mac); err != nil {
			panic(err)
		}
	})

	println(fmt.Sprintf("Starting HTTP Server on port %d", port))
	if err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil); err != nil {
		panic(err)
	}
}
