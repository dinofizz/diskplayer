package diskplayer

import (
	"html/template"
	"log"
	"net/http"
	"net/url"
)

type ErrorPage struct {
	Body []byte
}

func RunRecordServer() {
	p := GetConfigString(RECORD_SERVER_PORT)

	http.Handle("/", http.FileServer(http.Dir("./static")))
	http.HandleFunc("/record", recordHandler)
	http.ListenAndServe(":"+p, nil)
}

func RunCallbackServer(h http.Handler) *http.Server {
	r := GetConfigString(SPOTIFY_CALLBACK_URL)
	u, err := url.Parse(r)
	HandleError(err)

	http.Handle(u.EscapedPath(), h)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Got request for:", r.URL.String())
	})
	server := &http.Server{Addr: ":" + u.Port(), Handler: nil}
	go server.ListenAndServe()
	return server
}

func recordHandler(w http.ResponseWriter, r *http.Request) {
	webUrl := r.FormValue("web_url")
	log.Println(webUrl)
	e := Record(webUrl)
	if e != nil {
		p := &ErrorPage{Body: []byte(e.Error())}
		t, _ := template.ParseFiles("./static/error.html")
		t.Execute(w, p)
	} else {
		http.Redirect(w, r, "/success.html", http.StatusFound)
	}
}
