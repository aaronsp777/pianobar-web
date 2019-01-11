package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
)

var bind_addr = flag.String("bind_addr", ":8000", "Host/Port to listen on")
var piano_ctl = flag.String("piano_ctl", "/home/pi/.config/pianobar/ctl", "Pianobar fifo")

func doAction(a string) error {
	fmt.Println("action:", a)
	f, err := os.OpenFile(*piano_ctl, os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	if a == "play" {
		f.WriteString("P")
	} else if a == "pause" {
		f.WriteString("S")
	} else if a == "next" {
		f.WriteString("n")
	} else {
		fmt.Println("unknown aciton:", a)
	}

	return nil
}

func topHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprintln(w, "<!DOCTYPE html>")
	fmt.Fprintln(w, "<html>")
	fmt.Fprintln(w, "<body>")
	fmt.Fprintln(w, "<meta name='viewport' content='width=device-width, initial-scale=1.0'>")
	fmt.Fprintln(w, "<form action='/action/' method='get'>")
	fmt.Fprintln(w, "  <input type='submit' name='action' value='pause' class='pause'>")
	fmt.Fprintln(w, "  <input type='submit' name='action' value='play' class='play'>")
	fmt.Fprintln(w, "  <input type='submit' name='action' value='next' class='next'>")
	fmt.Fprintln(w, "</form>")
	fmt.Fprintln(w, "</body>")
	fmt.Fprintln(w, "</html>")
}

func actionHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		fmt.Fprintf(w, "ParseForm: %s", err)
		return
	}
	action := r.FormValue("action")
	if action == "" {
		fmt.Fprintf(w, "no action")
		return
	}
	err = doAction(action)
	if err != nil {
		fmt.Fprintf(w, "doAction: %s", err)
		return
	}

	http.Error(w, "No Content", 204)
}

func main() {
	flag.Parse()

	http.HandleFunc("/action/", actionHandler)
	http.HandleFunc("/", topHandler)

	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	http.Handle("/favicon.ico", fs)

	fmt.Println("pianobar-web started, listening on", *bind_addr)
	http.ListenAndServe(*bind_addr, nil)
}
