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

type FactIndex struct {
	Facts      riak.Many
	riak.Model `riak:"test.augend.index"`
}

func ensureBuckets() error {
	_, err := riak.NewBucket("test.augend.fact")
	if err != nil {
		fmt.Println("could not get/create fact bucket")
		return err
	}
	_, err = riak.NewBucket("test.augend.index")
	if err != nil {
		fmt.Println("could not get/create fact bucket")
		return err
	}
	return nil
}

func getOrCreateTag(t string) *Tag {
	var tag Tag
	err := riak.LoadModel(t, &tag)
	if err != nil {
		fmt.Println("creating new tag")
		return createTag(t)
	}
	return &tag
}

func createTag(t string) *Tag {
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

func getOrCreateFactIndex() *FactIndex {
	var index FactIndex
	err := riak.LoadModel("fact-index", &index)
	if err != nil {
		fmt.Println("creating new fact index")
		return createFactIndex()
	}
	return &index
}

func createFactIndex() *FactIndex {
	var index FactIndex
	err := riak.NewModel("fact-index", &index)
	if err != nil {
		fmt.Println("could not create new fact index")
		fmt.Println(err)
		return nil
	}
	index.SaveAs("fact-index")
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
	fact.SourceName = source_name
	fact.SourceUrl = source_url
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

func main() {
	var riak_host = "localhost:10017"
	err := riak.ConnectClient(riak_host)
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	err = ensureBuckets()
	if err != nil {
		fmt.Println("problem creating buckets. can't start")
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

type Tag struct {
	Name       string
	Facts      riak.Many
	riak.Model `riak:"test.augend.tag"`
}

func (t *Tag) Resolve(count int) (err error) {
	return nil
}

func (t *Tag) Url() string {
	return "/tag/" + t.Name + "/"
}

func (t *Tag) ListFacts() []Fact {
	fmt.Println("ListFacts()")
	fl := make([]Fact, t.Facts.Len())
	fmt.Println(t.Facts.Len())
	for i, f := range t.Facts {
		var lfact Fact
		f.Get(&lfact)
		fl[i] = lfact
	}
	return fl
}

type TagIndex struct {
	Tags       riak.Many
	riak.Model `riak:"test.augend.index"`
}

func getOrCreateTagIndex() *TagIndex {
	var index TagIndex
	err := riak.LoadModel("tag-index", &index)
	if err != nil {
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
	nindex.SaveAs("tag-index")
	return &nindex
}
