package main

import (
	"flag"
	"fmt"
	"html/template"
	"net/http"
	"os"
)

var bind_addr = flag.String("bind_addr", ":8000", "Host/Port to listen on")
var piano_ctl = flag.String("piano_ctl", "/home/pi/.config/pianobar/ctl", "Pianobar fifo")

var actionMap = map[string]string{
	"play":     "P",
	"pause":    "S",
	"next":     "n",
	"voldown":  "(",
	"volreset": "^",
	"volup":    ")",
}

func doAction(action string) error {
	fmt.Println("action:", action)
	f, err := os.OpenFile(*piano_ctl, os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	cmd, ok := actionMap[action]
	if ok {
		f.WriteString(cmd)
	} else {
		fmt.Println("unknown aciton:", action)
	}

	return nil
}

var topTemplate = template.Must(template.ParseFiles("top.html"))

func topHandler(w http.ResponseWriter, r *http.Request) {
	topTemplate.Execute(w, "")
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
