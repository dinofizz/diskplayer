package server

import (
	"github.com/dinofizz/diskplayer/internal/recorder"
	"html/template"
	"log"
	"net/http"
)

type ErrorPage struct {
	Body []byte
}

func RunRecordServer() {
	http.Handle("/", http.FileServer(http.Dir("./static")))
	http.HandleFunc("/record", recordHandler)
	http.ListenAndServe(":3000", nil)
}

func RunCallbackServer(callback func(w http.ResponseWriter, r *http.Request)) *http.Server {
	http.HandleFunc("/callback", callback)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Got request for:", r.URL.String())
	})
	server := &http.Server{Addr: ":8080", Handler: nil}
	go server.ListenAndServe()
	return server
}

func recordHandler(w http.ResponseWriter, r *http.Request) {
	webUrl := r.FormValue("web_url")
	log.Println(webUrl)
	e := recorder.Record(webUrl)
	if e != nil {
		p := &ErrorPage{Body: []byte(e.Error())}
		t, _ := template.ParseFiles("./static/error.html")
		t.Execute(w, p)
	} else {
		http.Redirect(w, r, "/success.html", http.StatusFound)
	}
}
