package main

import (
	"log"
	"net/http"
)

func main() {
	http.Handle("/", http.FileServer(http.Dir("./static")))
	http.HandleFunc("/record", recordHandler)
	http.ListenAndServe(":3000", nil)
}

func recordHandler(w http.ResponseWriter, r *http.Request) {
	web_url := r.FormValue("web_url")
	log.Println(web_url)
	http.Redirect(w, r, "/success.html", http.StatusFound)
}
