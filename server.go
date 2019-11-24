package diskplayer

import (
	"github.com/zmb3/spotify"
	"golang.org/x/oauth2"
	"html/template"
	"log"
	"net/http"
	"net/url"
)

type ErrorPage struct {
	Body []byte
}

type DiskplayerServer interface {
	RunRecordServer() error
	//RunCallbackServer(h http.Handler) (*http.Server, error)
	RunCallbackServer() (*http.Server, error)
	TokenChannel() chan *oauth2.Token
	Authenticator() *spotify.Authenticator
}

func NewDiskplayerServer(a *spotify.Authenticator, ch chan *oauth2.Token) *RealDiskplayerServer {
	h := CallbackHandler{
		ch:   ch,
		auth: a,
	}
	return &RealDiskplayerServer{cbh: h}
}

type RealDiskplayerServer struct {
	cbh CallbackHandler
}

// RunRecordServer creates a web server running on the port defined in the configuration file under the recorder.
// server_port field.
// Files are served directly from the "static" folder.
func (s *RealDiskplayerServer) RunRecordServer() error {
	p := ConfigValue(RECORD_SERVER_PORT)

	http.Handle("/", http.FileServer(http.Dir("./static")))
	http.HandleFunc("/record", recordHandler)
	return http.ListenAndServe(":"+p, nil)
}

// RunCallbackServer creates a web server running on the port defined in the configuration file under the spotify.
// callback_url field.
// A pointer to the server object is returned so that it can be shutdown when no longer needed.
func (s *RealDiskplayerServer) RunCallbackServer() (*http.Server, error) {
	r := ConfigValue(SPOTIFY_CALLBACK_URL)
	u, err := url.Parse(r)
	if err != nil {
		return nil, err
	}

	http.Handle(u.EscapedPath(), s.cbh)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Got request for:", r.URL.String())
	})
	server := &http.Server{Addr: ":" + u.Port(), Handler: nil}
	go server.ListenAndServe()
	return server, nil
}

func (s *RealDiskplayerServer) TokenChannel() chan *oauth2.Token {
	return s.cbh.ch
}

func (s *RealDiskplayerServer) Authenticator() *spotify.Authenticator {
	return s.cbh.auth
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

type CallbackHandler struct {
	ch   chan *oauth2.Token
	auth *spotify.Authenticator
}

// An implementation of the Handler ServeHTTP function for the CallbackHandler struct.
func (h CallbackHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t, err := h.auth.Token(STATE_IDENTIFIER, r)
	if err != nil {
		http.Error(w, "Couldn't get token", http.StatusForbidden)
		log.Fatal(err)
	}
	if st := r.FormValue("state"); st != STATE_IDENTIFIER {
		http.NotFound(w, r)
		log.Fatalf("State mismatch: %s != %s\n", st, STATE_IDENTIFIER)
	}

	h.ch <- t
}
