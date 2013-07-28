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
	riak.Model `riak:"test.augend.fact"`
}

func (f *Fact) Url() string {
	return "/fact/" + f.Key() + "/"
}

func (f *Fact) Resolve(count int) (err error) {
	return nil
}

type FactIndex struct {
	Facts riak.Many
	riak.Model `riak:"test.augend.index"`
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

func NewFact(title, details, source_name, source_url string) *Fact {
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
		NewFact(title, details, source_name, source_url)
		http.Redirect(w, r, "/", http.StatusFound)
	} else {
		pattern := filepath.Join("templates", "add.html")
		tmpl := template.Must(template.ParseGlob(pattern))
		tmpl.Execute(w, nil)
	}
}
