
package lib

import (
	"fmt"
	"net/http"
)

func DoStuff() {
	print("Doing stuff\n")
	http.Handle("/", http.FileServer(_escFS(false)))
	http.HandleFunc("/albumdata.js", albumdataHandler)
	http.ListenAndServe(":3000", nil)
}

func albumdataHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "foo=1;")
}