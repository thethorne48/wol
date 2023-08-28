package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/thethorne48/wol/magic"
)

var port int

func main() {
	p := os.Getenv("PORT")
	if p != "" {
		var err error
		port, err = strconv.Atoi(p)
		if err != nil {
			panic(err)
		}
	} else {
		port = 8080
	}
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		var mac = r.URL.Query().Get("mac")
		fmt.Fprintf(w, "Broadcasting WOL request to %s", mac)
		if err := magic.SendMagicPacket(mac); err != nil {
			panic(err)
		}
	})

	println(fmt.Sprintf("Starting HTTP Server on port %d", port))
	if err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil); err != nil {
		panic(err)
	}
}
