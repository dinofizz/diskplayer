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

// RunRecordServer creates a web server running on the port defined in the configuration file under the recorder.
// server_port field.
// Files are served directly from the "static" folder.
func RunRecordServer() {
	p := ConfigValue(RECORD_SERVER_PORT)

	http.Handle("/", http.FileServer(http.Dir("./static")))
	http.HandleFunc("/record", recordHandler)
	http.ListenAndServe(":"+p, nil)
}

// RunCallbackServer creates a web server running on the port defined in the configuration file under the spotify.
// callback_url field.
// A pointer to the server object is returned so that it can be shutdown when no longer needed.
func RunCallbackServer(h http.Handler) *http.Server {
	r := ConfigValue(SPOTIFY_CALLBACK_URL)
	u, err := url.Parse(r)
	if err != nil {
		log.Fatal(err)
	}

	http.Handle(u.EscapedPath(), h)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Got request for:", r.URL.String())
	})
	server := &http.Server{Addr: ":" + u.Port(), Handler: nil}
	go server.ListenAndServe()
	return server
}

// recordHandler handles requests to the server which contain a Spotify web URL to be recorded.
// If the recording is successful, redirection to a success page occurs, otherwise an error page is returned.
func recordHandler(w http.ResponseWriter, r *http.Request) {
	webUrl := r.FormValue("web_url")
	e := Record(webUrl)
	if e != nil {
		p := &ErrorPage{Body: []byte(e.Error())}
		t, _ := template.ParseFiles("./static/error.html")
		t.Execute(w, p)
	} else {
		http.Redirect(w, r, "/success.html", http.StatusFound)
	}
}
