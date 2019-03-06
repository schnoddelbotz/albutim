package lib

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
)

type server struct {
	album Album
}

func Serve(a Album, httpPort string) {
	srv := &server{album: a}

	http.Handle("/assets/", http.StripPrefix("/assets", http.FileServer(_escFS(false))))
	http.Handle("/originals/", http.StripPrefix("/originals", http.FileServer(http.Dir(a.RootPath))))
	http.Handle("/thumbs/", http.StripPrefix("/thumbs", http.FileServer(http.Dir(a.RootPath))))
	http.Handle("/preview/", http.StripPrefix("/preview", http.FileServer(http.Dir(a.RootPath))))
	//http.HandleFunc("/thumbs/", srv.thumbHandler)
	http.HandleFunc("/albumdata.json", srv.albumDataHandler)
	http.HandleFunc("/", srv.indexHandler)

	log.Printf("Starting Server on port %s\n", httpPort)
	http.ListenAndServe(":"+httpPort, nil)
}

func (s *server) albumDataHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json")
	json, _ := json.Marshal(s.album)
	w.Write(json)
}

func (s *server) indexHandler(w http.ResponseWriter, r *http.Request) {
	w.Write(renderIndexTemplate(s.album))
}

func (s *server) thumbHandler(w http.ResponseWriter, r *http.Request) {
	original := strings.Replace(r.URL.Path, "/thumbs/", "", 1)
	fmt.Fprintf(w, "HERE YOURS THUMB: %s", original)
}
