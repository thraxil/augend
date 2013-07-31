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
	Facts    []Fact
	Username string
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
	sess, _ := store.Get(r, "augend")
	username, found := sess.Values["user"]
	if found && username != "" {
		ir.Username = username.(string)
	}
	pattern := filepath.Join("templates", "index.html")
	tmpl := template.Must(template.ParseGlob(pattern))
	tmpl.Execute(w, ir)
}

type TagIndexResponse struct {
	Tags     []Tag
	Username string
}

type TagResponse struct {
	Tag      Tag
	Username string
}

func tagHandler(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.String(), "/")
	sess, _ := store.Get(r, "augend")
	username, found := sess.Values["user"]
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
		if found && username != "" {
			ir.Username = username.(string)
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
	if found && username != "" {
		tr.Username = username.(string)
	}
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
	sess, _ := store.Get(r, "augend")
	username, found := sess.Values["user"]
	if !found || username == "" {
		http.Redirect(w, r, "/login/", http.StatusFound)
		return
	}
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

func registerForm(w http.ResponseWriter, req *http.Request) {
	pattern := filepath.Join("templates", "register.html")
	tmpl := template.Must(template.ParseGlob(pattern))
	tmpl.Execute(w, nil)
}

func registerHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		registerForm(w, r)
		return
	}
	username, password, pass2 := r.FormValue("username"), r.FormValue("password"), r.FormValue("pass2")
	if password != pass2 {
		fmt.Fprintf(w, "passwords don't match")
		return
	}
	user := NewUser(username, password)

	sess, _ := store.Get(r, "augend")
	sess.Values["user"] = user.Username
	sess.Save(r, w)
	http.Redirect(w, r, "/", http.StatusFound)
}

func loginForm(w http.ResponseWriter, req *http.Request) {
	pattern := filepath.Join("templates", "login.html")
	tmpl := template.Must(template.ParseGlob(pattern))
	tmpl.Execute(w, nil)
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		loginForm(w, r)
		return
	}
	username, password := r.FormValue("username"), r.FormValue("password")
	var user User
	err := riak.LoadModel(username, &user)
	if err != nil {
		fmt.Println("couldn't load user:", err)
		fmt.Fprintf(w, "user not found")
		return
	}
	if !user.CheckPassword(password) {
		fmt.Fprintf(w, "login failed")
		return
	}
	// store userid in session
	sess, _ := store.Get(r, "augend")
	sess.Values["user"] = user.Username
	sess.Save(r, w)
	http.Redirect(w, r, "/", http.StatusFound)
}
