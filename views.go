package main

import (
	"fmt"
	"github.com/tpjg/goriakpbc"
	"html/template"
	"net/http"
	"path/filepath"
	"strings"
)

type IndexResponse struct {
	Facts []Fact
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	index := getOrCreateFactIndex()
	if index == nil {
		fmt.Fprintf(w, "could not retrieve or create main fact index")
		return
	}
	n := index.Facts.Len()
	facts := make([]Fact, n)
	for i, f := range index.Facts {
		var lfact Fact
		f.Get(&lfact)
		facts[n-1-i] = lfact
	}

	ir := IndexResponse{
		Facts: facts,
	}
	pattern := filepath.Join("templates", "index.html")
	tmpl := template.Must(template.ParseGlob(pattern))
	tmpl.Execute(w, ir)
}

type TagIndexResponse struct {
	Tags []Tag
}

type TagResponse struct {
	Tag Tag
}

func tagHandler(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.String(), "/")
	if len(parts) < 3 || parts[2] == "" {
		index := getOrCreateTagIndex()
		n := index.Tags.Len()
		tags := make([]Tag, n)
		for i, t := range index.Tags {
			var ltag Tag
			t.Get(&ltag)
			tags[i] = ltag
		}

		ir := TagIndexResponse{
			Tags: tags,
		}
		pattern := filepath.Join("templates", "tags.html")
		tmpl := template.Must(template.ParseGlob(pattern))
		tmpl.Execute(w, ir)
		return
	}

	id := parts[2]
	var tag Tag
	riak.LoadModel(id, &tag)
	tr := TagResponse{Tag: tag}
	pattern := filepath.Join("templates", "tag.html")
	tmpl := template.Must(template.ParseGlob(pattern))
	tmpl.Execute(w, tr)
}

type FactResponse struct {
	Fact Fact
}

func factHandler(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.String(), "/")
	if len(parts) < 3 {
		http.Error(w, "bad request", 400)
		return
	}
	id := parts[2]
	if id == "" {
		http.Error(w, "bad request", 400)
		return
	}
	var fact Fact
	riak.LoadModel(id, &fact)
	fr := FactResponse{Fact: fact}
	pattern := filepath.Join("templates", "fact.html")
	tmpl := template.Must(template.ParseGlob(pattern))
	tmpl.Execute(w, fr)
}

func addHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		title := r.PostFormValue("title")
		details := r.PostFormValue("details")
		source_name := r.PostFormValue("source_name")
		source_url := r.PostFormValue("source_url")
		tags := r.PostFormValue("tags")
		NewFact(title, details, source_name, source_url, tags)
		http.Redirect(w, r, "/", http.StatusFound)
	} else {
		pattern := filepath.Join("templates", "add.html")
		tmpl := template.Must(template.ParseGlob(pattern))
		tmpl.Execute(w, nil)
	}
}
