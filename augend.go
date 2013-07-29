package main

import (
	"fmt"
	"github.com/nu7hatch/gouuid"
	"github.com/tpjg/goriakpbc"
	"html/template"
	"net/http"
	"path/filepath"
	"strings"
	"time"
)

type Fact struct {
	Title string ""
	Details string ""
	SourceName string ""
	SourceUrl string ""
	Added string ""
	Tags riak.Many
	riak.Model `riak:"test.augend.fact"`
}

func (f *Fact) Url() string {
	return "/fact/" + f.Key() + "/"
}

func (f *Fact) Resolve(count int) (err error) {
	fmt.Println("resolve fact")
	return nil
}

func (f *Fact) AddTag(t string) {
	tag := getOrCreateTag(t)
	f.Tags.Add(tag)
	f.SaveAs(f.Key())
	tag.Facts.Add(f)
	tag.SaveAs(t)
}

func (f *Fact) HasTags() bool {
	fmt.Println(f.Tags.Len())
	return f.Tags.Len() > 0
}

func (f *Fact) ListTags() []Tag {
	tl := make([]Tag, f.Tags.Len())
	for i, t := range f.Tags {
		var ltag Tag
		t.Get(&ltag)
		tl[i] = ltag
	}
	return tl
}

type FactIndex struct {
	Facts riak.Many
	riak.Model `riak:"test.augend.index"`
}

func getOrCreateTag(t string) *Tag {
	var tag Tag
	err := riak.LoadModel(t, &tag)
	if err != nil {
		fmt.Println(err)
		fmt.Println("creating new tag")
		var ntag Tag
		err := riak.NewModel(t, &ntag)
		if err != nil {
			fmt.Println("could not create new tag")
			fmt.Println(err)
			return nil
		}
		ntag.Name = t
		ntag.SaveAs(t)
		tag_index := getOrCreateTagIndex()
		if tag_index == nil {
			fmt.Println("no tag index!")
			return nil
		}
		tag_index.Tags.Add(&ntag)
		tag_index.SaveAs("tag-index")
		return &ntag
	}
	return &tag
}

func getOrCreateFactIndex() *FactIndex {
	var index FactIndex
	err := riak.LoadModel("fact-index", &index)
	if err != nil {
		fmt.Println(err)
		fmt.Println("creating new fact index")
		err = riak.NewModel("fact-index", &index)
		if err != nil {
			return nil
		}
		return &index
	}
	return &index
}

func NewFact(title, details, source_name, source_url, tags string) *Fact {
	var fact Fact
	u4, err := uuid.NewV4()
	if err != nil {
		fmt.Println("error:", err)
		return nil
	}

	err = riak.NewModel(u4.String(), &fact)
	fact.Title = title
	fact.Details = details
	t := time.Now()
	fact.Added = t.Format(time.RFC3339)

	fact.SaveAs(u4.String())

	index := getOrCreateFactIndex()
	if index == nil {
		fmt.Println("unable to get/create fact index")
		return nil
	}
	index.Facts.Add(&fact)
	index.SaveAs("fact-index")
	for _, t := range strings.Split(tags, ",") {
		fact.AddTag(t)
	}
	return &fact
}

func main () {
	var riak_host = "localhost:10017"
	err := riak.ConnectClient(riak_host)
	if err != nil {
		fmt.Println("error:", err)
		return
	}

//	fmt.Println(index.Facts.Len())
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/fact/", factHandler)
	http.HandleFunc("/tag/", tagHandler)
	http.HandleFunc("/add/", addHandler)
	http.Handle("/media/", http.StripPrefix("/media/",
		http.FileServer(http.Dir("media"))))
	http.ListenAndServe(":9999", nil)
}

type IndexResponse struct {
	Facts []Fact
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	index := getOrCreateFactIndex()
	n := index.Facts.Len()
	facts := make([]Fact, n)
	for i, f := range index.Facts {
		var lfact Fact
		f.Get(&lfact)
		facts[n - 1 - i] = lfact
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

func tagHandler(w http.ResponseWriter, r *http.Request) {
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

type Tag struct {
	Name string
	Facts riak.Many
	riak.Model `riak:"test.augend.tag"`
}

func (t *Tag) Resolve(count int) (err error) {
	return nil
}

func (t *Tag) Url() string {
	return "/tag/" + t.Name + "/"
}

type TagIndex struct {
	Tags riak.Many
	riak.Model `riak:"test.augend.index"`
}

func getOrCreateTagIndex() *TagIndex {
	var index TagIndex
	err := riak.LoadModel("tag-index", &index)
	if err != nil {
		fmt.Println(err)
		fmt.Println("creating new tag index")
		return createTagIndex()
	}
	return &index
}

func createTagIndex() *TagIndex {
	var nindex TagIndex
	err := riak.NewModel("tag-index", &nindex)
	if err != nil {
		fmt.Println("could not create tag index")
		fmt.Println(err)
		return nil
	}
	return &nindex
}
