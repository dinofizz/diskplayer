package diskplayer

import (
	"github.com/docker/docker/pkg/mount"
	"github.com/zmb3/spotify"
	"golang.org/x/oauth2"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"os/exec"
)

type IndexPage struct {
	Lsblk []byte
}

type ErrorPage struct {
	Body []byte
}

type DiskplayerServer interface {
	RunRecordServer() error
	RunCallbackServer() (*http.Server, error)
	TokenChannel() chan *oauth2.Token
	Authenticator() *spotify.Authenticator
}

// NewDiskplayerServer returns a new DiskplayerServer instance.
// The arguments are required if the server instance is to be used to obtain a new Spotify auth token.
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
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	http.HandleFunc("/", indexHandler)
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
	devPath := r.FormValue("device_path")

	folder := ConfigValue(RECORD_FOLDER_PATH)
	filename := ConfigValue(RECORD_FILENAME)
	dstPath := folder + "/" + filename

	m, err := mount.Mounted(folder);
	if err != nil {
		returnErrorPage(w, err)
	}

	if m == false {
		err := mount.Mount(devPath, folder, "vfat", "");
		if err != nil {
			returnErrorPage(w, err)
		}
	}

	err = Record(webUrl, dstPath)
	if err != nil {
		returnErrorPage(w, err)
	}

	err = mount.Unmount(folder);
	if err != nil {
		returnErrorPage(w, err)
	}

	http.Redirect(w, r, "/static/success.html", http.StatusFound)
}

// indexHandler handles requests to the server for the root location "/".
// A listing of the attached devices is obtained and applied to the index.html template response.
// An error page is returned if an error occurred.
func indexHandler(w http.ResponseWriter, r *http.Request) {
	cmd := exec.Command("lsblk", "--nodeps")
	stdout, err := cmd.Output()
	if err != nil {
		returnErrorPage(w, err);
	}

	p := &IndexPage{Lsblk: stdout}
	t, _ := template.ParseFiles("./templates/index.html")
	t.Execute(w, p)
}

// returnErrorPage returns an HTML error page, inserting error details into the error.html template.
func returnErrorPage(w http.ResponseWriter, err error) {
	p := &ErrorPage{Body: []byte(err.Error())}
	t, _ := template.ParseFiles("./templates/error.html")
	t.Execute(w, p)
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
