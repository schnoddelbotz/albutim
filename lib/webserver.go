
package lib

import (
	"fmt"
	"net/http"
	"strings"
)

func Serve(albumRoot string, httpPort string) {
	fmt.Printf("Starting Server on port %s\n", httpPort)
	http.Handle("/assets/", http.StripPrefix("/assets", http.FileServer(_escFS(false))))
	http.Handle("/originals/", http.StripPrefix("/originals", http.FileServer(http.Dir(albumRoot))))
	http.HandleFunc("/thumbs/", thumbHandler)
	http.HandleFunc("/albumdata.json", albumDataHandler)
	http.HandleFunc("/", indexHandler)
	http.ListenAndServe(":"+httpPort, nil)
}

func albumDataHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json")
	fmt.Fprintf(w, `{"albumTitle":"FooTitle Bar"}`)
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	w.Write(_escFSMustByte(false, "/index.html"))
}

func thumbHandler(w http.ResponseWriter, r *http.Request) {
	original := strings.Replace(r.URL.Path, "/thumbs/", "", 1)
	fmt.Fprintf(w, "HERE YOURS THUMB: %s", original)
}