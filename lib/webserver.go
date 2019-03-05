package lib

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
)

type server struct {
	albumData Album
}

func Serve(albumRoot string, httpPort string) {
	log.Print("Reading images...")
	var err error
	srv := new(server)
	srv.albumData.Data, err = scanDir("testalbum")
	if err != nil {
		log.Fatalf("Cannot scan '%s': %s", albumRoot, err)
	}

	http.Handle("/assets/", http.StripPrefix("/assets", http.FileServer(_escFS(false))))
	http.Handle("/originals/", http.StripPrefix("/originals", http.FileServer(http.Dir(albumRoot))))
	http.HandleFunc("/thumbs/", srv.thumbHandler)
	http.HandleFunc("/albumdata.json", srv.albumDataHandler)
	http.HandleFunc("/", srv.indexHandler)

	log.Printf("Starting Server on port %s\n", httpPort)
	http.ListenAndServe(":"+httpPort, nil)
}

func (s *server) albumDataHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json")
	json, _ := json.Marshal(s.albumData)
	w.Write(json)
}

func (s *server) indexHandler(w http.ResponseWriter, r *http.Request) {
	w.Write(_escFSMustByte(false, "/index.html"))
}

func (s *server) thumbHandler(w http.ResponseWriter, r *http.Request) {
	original := strings.Replace(r.URL.Path, "/thumbs/", "", 1)
	fmt.Fprintf(w, "HERE YOURS THUMB: %s", original)
}
