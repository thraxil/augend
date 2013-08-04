package main

import (
	"fmt"
	"github.com/tpjg/goriakpbc"
	"html/template"
	"net/http"
	"path/filepath"
	"strings"
	"strconv"
)

type SiteResponse struct {
	Username string
}

type IndexResponse struct {
	Facts []Fact
	Page  Page
	SiteResponse
}

func (f FactIndex) TotalItems() int {
	return f.Facts.Len()
}

func (f FactIndex) ItemRange(offset, count int) []Fact {
	facts_on_page := minint(count, (f.TotalItems() - offset))
	facts := make([]Fact, facts_on_page)
	for i := 0; i < facts_on_page; i++ {
		var lfact Fact
		f.Facts[offset+i].Get(&lfact)
		facts[facts_on_page-1-i] = lfact
	}
	return facts
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	index := getOrCreateFactIndex()
	if index == nil {
		fmt.Fprintf(w, "could not retrieve or create main fact index")
		return
	}
	var p = Paginator{ItemList: index, PerPage: 20}
	page := p.GetPage(r)
	facts := page.Items()
	ir := IndexResponse{
		Page:  page,
		Facts: facts,
	}
	sess, _ := store.Get(r, "augend")
	username, found := sess.Values["user"]
	if found && username != "" {
		ir.Username = username.(string)
	}
	tmpl := getTemplate("index.html")
	tmpl.Execute(w, ir)
}

type TagIndexResponse struct {
	Tags []Tag
	SiteResponse
}

type TagResponse struct {
	Tag Tag
	SiteResponse
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
		tmpl := getTemplate("tags.html")
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
	tmpl := getTemplate("tag.html")
	tmpl.Execute(w, tr)
}

type FactResponse struct {
	Fact Fact
	SiteResponse
}

func factHandler(w http.ResponseWriter, r *http.Request) {
	sess, _ := store.Get(r, "augend")
	username, found := sess.Values["user"]
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
	if found && username != "" {
		fr.Username = username.(string)
	}

	tmpl := getTemplate("fact.html")
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
		// call once to make sure the form is initialized
		// before we access r.Form directly
		r.PostFormValue("title0")

		for k, _ := range r.Form {
			if strings.HasPrefix(k, "title") {
				idx, err := strconv.Atoi(k[5:])
				if err != nil {
					continue
				}
				title := r.PostFormValue("title" + strconv.Itoa(idx))
				if title == "" {
					// no title? don't bother
					continue
				}
				details := r.PostFormValue("details" + strconv.Itoa(idx))
				source_name := r.PostFormValue("source_name" + strconv.Itoa(idx))
				source_url := r.PostFormValue("source_url" + strconv.Itoa(idx))
				tags := r.PostFormValue("tags" + strconv.Itoa(idx))
				var user User
				riak.LoadModel(username.(string), &user)
				NewFact(title, details, source_name, source_url, tags, user)
			}
		}
		http.Redirect(w, r, "/", http.StatusFound)
	} else {
		tmpl := getTemplate("add.html")
		ctx := SiteResponse{Username: username.(string)}
		tmpl.Execute(w, ctx)
	}
}

func registerForm(w http.ResponseWriter, req *http.Request) {
	tmpl := getTemplate("register.html")
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
	tmpl := getTemplate("login.html")
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

func logoutHandler(w http.ResponseWriter, r *http.Request) {
	sess, _ := store.Get(r, "augend")
	delete(sess.Values, "user")
	sess.Save(r, w)
	http.Redirect(w, r, "/", http.StatusFound)
}

func getTemplate(filename string) *template.Template {
	var t = template.New("base.html")
	return template.Must(t.ParseFiles(
		filepath.Join(template_dir, "base.html"),
		filepath.Join(template_dir, filename),
	))
}
