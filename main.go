package main

import (
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"sync"
)

type templateHandler struct {
	once     sync.Once
	filename string
	templ    *template.Template
}

func (t *templateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	//makes sure only execute once per request
	t.once.Do(func() {
		//find the template file and parse it
		t.templ = template.Must(template.ParseFiles(filepath.Join("templates", t.filename)))
	})
	//return the parsed template to the browser
	t.templ.Execute(w, r)
}

func main() {

	r := newRoom()

	http.Handle("/", &templateHandler{filename: "chat.html"})
	http.Handle("/room", r)

	//Start room logic on a thread
	go r.run()

	//start web server
	log.Println("Starting go-chat server on port 8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}

}
