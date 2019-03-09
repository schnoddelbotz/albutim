package lib

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"
)

type server struct {
	album Album
}

func Serve(a Album, httpPort string) {
	srv := &server{album: a}

	http.Handle("/assets/", http.StripPrefix("/assets", http.FileServer(_escFS(false))))
	http.Handle("/originals/", http.StripPrefix("/originals", http.FileServer(http.Dir(a.RootPath))))
	if a.NoScaledPreviews {
		http.Handle("/preview/", http.StripPrefix("/preview", http.FileServer(http.Dir(a.RootPath))))
	} else {
		http.HandleFunc("/preview/", srv.previewHandler)
	}
	if a.NoScaledThumbs {
		http.Handle("/thumbs/", http.StripPrefix("/thumbs", http.FileServer(http.Dir(a.RootPath))))
	} else {
		http.HandleFunc("/thumbs/", srv.thumbHandler)
	}
	http.HandleFunc("/albumdata.json", srv.albumDataHandler)
	http.HandleFunc("/", srv.indexHandler)

	log.Printf("Starting Server on port %s\n", httpPort)
	err := http.ListenAndServe(":"+httpPort, nil)
	if err != nil {
		log.Fatalf("Could not start server: %s", err)
	}
}

func (s *server) albumDataHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json")
	myjson, _ := json.Marshal(s.album)
	_, _ = w.Write(myjson)
}

func (s *server) indexHandler(w http.ResponseWriter, r *http.Request) {
	_, err := w.Write(renderIndexTemplate(s.album))
	if err != nil {
		log.Printf("Sending index.html failed: %s", err)
	}
}

func (s *server) thumbHandler(w http.ResponseWriter, r *http.Request) {
	// FIXME: directory traversal....???!?!!
	original := s.album.RootPath + "/" + strings.Replace(r.URL.Path, "/thumbs/", "", 1)
	thumb := s.album.RootPath + r.URL.Path
	if s.serveCached(w, r, thumb) {
		return
	}
	// fixme: scale from preview if present
	scaled, err := getScaled(original, 0, 105 /* FIXME config value */)
	if err != nil {
		http.Error(w, "Preview failed", 500)
		return
	}
	w.Header().Set("Content-type", "image/jpeg")
	_, _ = w.Write(scaled)
	s.album.addCache(thumb, scaled)
}

func (s *server) previewHandler(w http.ResponseWriter, r *http.Request) {
	// FIXME: directory traversal....???!?!!
	original := s.album.RootPath + "/" + strings.Replace(r.URL.Path, "/preview/", "", 1)
	preview := s.album.RootPath + r.URL.Path
	if s.serveCached(w, r, preview) {
		return
	}
	scaled, err := getScaled(original, 0, 700 /* FIXME config value */)
	if err != nil {
		http.Error(w, "Preview failed", 500)
		return
	}
	w.Header().Set("Content-type", "image/jpeg")
	_, _ = w.Write(scaled)
	s.album.addCache(preview, scaled)
}

func (s *server) serveCached(w http.ResponseWriter, r *http.Request, file string) bool {
	if _, err := os.Stat(file); err == nil {
		http.ServeFile(w, r, file)
		return true
	}
	return false
}
