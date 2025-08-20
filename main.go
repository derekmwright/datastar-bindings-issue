package main

import (
	"bytes"
	"fmt"
	"html/template"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/starfederation/datastar-go/datastar"
)

var tpl = template.Must(template.ParseGlob("templates/*.gohtml"))

func handleRoot(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	if err := tpl.ExecuteTemplate(w, "index.gohtml", nil); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func handleLanding(w http.ResponseWriter, r *http.Request) {
	buf := &bytes.Buffer{}

	if err := tpl.ExecuteTemplate(buf, "landing.gohtml", nil); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	sse := datastar.NewSSE(w, r)
	if err := sse.PatchElements(buf.String()); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	// Clear out all form bindings
	type forms struct {
		Name     *string `json:"name"`
		Location *string `json:"location"`
		Date     *string `json:"date"`
	}

	empty := forms{}
	if err := sse.MarshalAndPatchSignals(empty); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func handleForm1(w http.ResponseWriter, r *http.Request) {
	if err := tpl.ExecuteTemplate(w, "form1.gohtml", nil); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func handleForm2(w http.ResponseWriter, r *http.Request) {
	if err := tpl.ExecuteTemplate(w, "form2.gohtml", nil); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func handleForm3(w http.ResponseWriter, r *http.Request) {
	if err := tpl.ExecuteTemplate(w, "form3.gohtml", nil); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func handleSubmit(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	fmt.Println(r.PostForm.Get("name"))
	fmt.Println(r.PostForm.Get("location"))
	fmt.Println(r.PostForm.Get("date"))

	http.Redirect(w, r, "/", http.StatusFound)
}

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/", handleRoot)
	r.Get("/landing", handleLanding)
	r.Get("/form1", handleForm1)
	r.Get("/form2", handleForm2)
	r.Get("/form3", handleForm3)
	r.Post("/submit", handleSubmit)

	if err := http.ListenAndServe(":8088", r); err != nil {
		panic(err)
	}
}
