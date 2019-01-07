package main

import (
	"flag"
	"fmt"
	"net/http"
)

var bind_addr = flag.String("bind_addr", ":8000", "Host/Port to listen on")

func main() {
	flag.Parse()

	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	http.Handle("/favicon.ico", fs)

	fmt.Println("Webserver started, listening on", *bind_addr)
	http.ListenAndServe(*bind_addr, nil)
}
