package main

import (
	"errors"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"util-pipe/internal/dbg"
	"util-pipe/internal/utils"

	"github.com/go-chi/chi"
)

func InitWiki(r *chi.Mux, path string) error {
	if !utils.IsDir(path) {
		return errors.New("not include")
	}
	return filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		ext := filepath.Ext(info.Name())
		if info.IsDir() || ext != ".html" {
			return nil
		}
		link := strings.ReplaceAll(path[len(config.WikiPath):len(path)-len(ext)], "\\", "/")
		if regexp.MustCompile(`[^a-zA-Z0-9_\-\/]`).MatchString(link) {
			log.Println("page does not meet the requirements wiki:", link)
			return nil
		}
		dbg.Log.Printf("init wiki page: %s - %v\n", path, link)
		config.wiki.ListLinks = append(config.wiki.ListLinks, "/wiki"+link)
		pageWiki := template.Must(template.ParseFiles(path))
		r.HandleFunc("/wiki"+link, func(w http.ResponseWriter, r *http.Request) { pageWiki.Execute(w, nil) })
		return nil
	})
}
