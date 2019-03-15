package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"text/template"
)

func ml5Handler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		log.Print(fmt.Sprintf("GET: %s", r.URL.Path))
		if r.URL.Path == "/" {
			r.URL.Path = "index.html"
		} else if r.URL.Path == "/ml5" {
			log.Print(fmt.Sprintf("Q: %s", r.URL.RawQuery))
			r.URL.Path = "ml5.html"
		}
		lp := filepath.Join("templates", "layout.html")
		fp := filepath.Join("templates", filepath.Clean(r.URL.Path))

		// Return a 404 if the template doesn't exist
		info, err := os.Stat(fp)
		if err != nil {
			if os.IsNotExist(err) {
				http.NotFound(w, r)
				log.Print(err)
				return
			}
		}

		// Return a 404 if the request is for a directory
		if info.IsDir() {
			http.NotFound(w, r)
			return
		}

		tmpl, err := template.ParseFiles(lp, fp)
		tmpl.New("img").Parse(`{{define "img"}}` + r.URL.RawQuery + `{{end}}`)

		if err != nil {
			// Log the detailed error
			log.Println(err.Error())
			// Return a generic "Internal Server Error" message
			http.Error(w, http.StatusText(500), 500)
			return
		}
		if err := tmpl.ExecuteTemplate(w, "layout", nil); err != nil {
			log.Println(err.Error())
			http.Error(w, http.StatusText(500), 500)
		}

	case "POST":
		log.Print(fmt.Sprintf("POST: %s", r.URL.Path))

		file, err := ioutil.TempFile("/tmp", "img.")
		if err != nil {
			log.Print(err)
		}
		log.Print(file.Name())
		f, _, _ := r.FormFile("image")
		defer f.Close()
		io.Copy(file, f)
		log.Print(file.Name())
		w.Write([]byte(fmt.Sprintf(`{"uploadUrl":"/ml5?%s"}`, file.Name())))
	}
}
