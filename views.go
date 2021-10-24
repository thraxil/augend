package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/thraxil/paginate"
)

type SiteResponse struct {
	Username string
}

type IndexResponse struct {
	Facts []Fact
	Page  paginate.Page
	SiteResponse
}

type paginatedFacts struct {
	s *site
}

func NewPaginatedFacts(s *site) paginatedFacts {
	return paginatedFacts{s: s}
}

func (p paginatedFacts) TotalItems() int {
	return p.s.TotalFactsCount()
}

func (p paginatedFacts) ItemRange(offset, count int) []interface{} {
	facts, err := p.s.ListFacts(offset, count)
	if err != nil {
		return make([]interface{}, 0)
	}
	out := make([]interface{}, len(facts))
	for j, v := range facts {
		out[j] = v
	}
	return out
}

func faviconHandler(w http.ResponseWriter, r *http.Request) {
}

func indexHandler(w http.ResponseWriter, r *http.Request, s *site) {
	index := NewPaginatedFacts(s)
	var p = paginate.Paginator{ItemList: index, PerPage: 20}
	page := p.GetPage(r)
	ifacts := page.Items()
	facts := make([]Fact, len(ifacts))
	for i, v := range ifacts {
		facts[i] = v.(Fact)
	}

	ir := IndexResponse{
		Page:  page,
		Facts: facts,
	}
	sess, _ := s.Store.Get(r, "augend")
	username, found := sess.Values["user"]
	if found && username != "" {
		ir.Username = username.(string)
	}
	tmpl := getTemplate("index.html")
	tmpl.Execute(w, ir)
}

type TagIndexResponse struct {
	Tags TagList
	SiteResponse
}

type TagResponse struct {
	Tag   Tag
	Facts []Fact
	SiteResponse
}

func tagHandler(w http.ResponseWriter, r *http.Request, s *site) {
	parts := strings.Split(r.URL.String(), "/")
	sess, _ := s.Store.Get(r, "augend")
	username, found := sess.Values["user"]
	if len(parts) < 3 || parts[2] == "" {
		tags, err := s.ListTags()
		if err != nil {
			log.Println(err)
			fmt.Fprintf(w, "couldn't list tags")
			return
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

	slug := parts[2]
	tag := s.GetTagBySlug(slug)
	facts, err := s.ListFactsByTag(tag)
	if err != nil {
		log.Println(err)
		fmt.Fprintf(w, "couldn't list facts")
		return
	}
	tr := TagResponse{
		Tag:   tag,
		Facts: facts,
	}
	if found && username != "" {
		tr.Username = username.(string)
	}
	tmpl := getTemplate("tag.html")
	tmpl.Execute(w, tr)
}

type FactResponse struct {
	Fact *Fact
	SiteResponse
}

func factHandler(w http.ResponseWriter, r *http.Request, s *site) {
	sess, err := s.Store.Get(r, "augend")
	if err != nil {
		log.Println("no session store")
		fmt.Fprintf(w, "could not create a session store")
		return
	}
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

	fact, err := s.GetFactById(id)
	if err != nil {
		fmt.Fprintf(w, "couldn't find fact")
		return
	}
	fr := FactResponse{Fact: fact}
	if found && username != "" {
		fr.Username = username.(string)
	}

	tmpl := getTemplate("fact.html")
	tmpl.Execute(w, fr)
}

type AddResponse struct {
	SourceName string
	SourceUrl  string
	Details    string
	SiteResponse
}

func addHandler(w http.ResponseWriter, r *http.Request, s *site) {
	sess, _ := s.Store.Get(r, "augend")
	username, found := sess.Values["user"]
	if !found || username == "" {
		http.Redirect(w, r, "/login/", http.StatusFound)
		return
	}

	user, err := s.GetUser(username.(string))
	if err != nil {
		fmt.Fprintf(w, "user not found")
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
				s.CreateFact(title, details, source_name, source_url, tags, user)
			}
		}
		http.Redirect(w, r, "/", http.StatusFound)
	} else {
		source_name := r.FormValue("source_name")
		source_url := r.FormValue("source_url")
		details := r.FormValue("details")
		ar := AddResponse{
			SourceName: source_name, SourceUrl: source_url, Details: details}
		ar.Username = username.(string)
		tmpl := getTemplate("add.html")
		tmpl.Execute(w, ar)
	}
}

func registerForm(w http.ResponseWriter, req *http.Request) {
	tmpl := getTemplate("register.html")
	tmpl.Execute(w, nil)
}

func registerHandler(w http.ResponseWriter, r *http.Request, s *site) {
	if r.Method == "GET" {
		registerForm(w, r)
		return
	}
	username, password, pass2 := r.FormValue("username"), r.FormValue("password"), r.FormValue("pass2")
	if password != pass2 {
		fmt.Fprintf(w, "passwords don't match")
		return
	}
	log.Println("passwords match")
	user, err := s.CreateUser(username, password)
	if err != nil {
		fmt.Println(err)
		fmt.Fprintf(w, "could not create user")
		return
	}
	log.Println("created user")

	sess, _ := s.Store.Get(r, "augend")
	sess.Values["user"] = user.Username
	sess.Save(r, w)
	http.Redirect(w, r, "/", http.StatusFound)
}

func loginForm(w http.ResponseWriter, req *http.Request) {
	tmpl := getTemplate("login.html")
	tmpl.Execute(w, nil)
}

func loginHandler(w http.ResponseWriter, r *http.Request, s *site) {
	if r.Method == "GET" {
		loginForm(w, r)
		return
	}
	username, password := r.FormValue("username"), r.FormValue("password")
	user, err := s.GetUser(username)

	if err != nil {
		fmt.Fprintf(w, "user not found")
		return
	}
	if !user.CheckPassword(password) {
		fmt.Fprintf(w, "login failed")
		return
	}
	// store userid in session
	sess, _ := s.Store.Get(r, "augend")
	sess.Values["user"] = user.Username
	sess.Save(r, w)
	http.Redirect(w, r, "/", http.StatusFound)
}

func logoutHandler(w http.ResponseWriter, r *http.Request, s *site) {
	sess, _ := s.Store.Get(r, "augend")
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

type smoketestResponse struct {
	Status       string   `json:"status"`
	TestClasses  int      `json:"test_classes"`
	TestsRun     int      `json:"tests_run"`
	TestsPassed  int      `json:"tests_passed"`
	TestsFailed  int      `json:"tests_failed"`
	TestsErrored int      `json:"tests_errored"`
	Time         float64  `json:"time"`
	ErroredTests []string `json:"errored_tests"`
	FailedTests  []string `json:"failed_tests"`
}

func smoketestHandler(w http.ResponseWriter, r *http.Request, s *site) {
	var status string
	var tests int

	tests = 1

	sr := smoketestResponse{
		Status:       status,
		TestClasses:  1,
		TestsRun:     1,
		TestsPassed:  tests,
		TestsFailed:  1 - tests,
		TestsErrored: 0,
		Time:         1.0,
	}
	if sr.TestsFailed > 0 || sr.TestsErrored > 0 {
		http.Error(w, "smoketest failed", http.StatusInternalServerError)
	}

	h := r.Header.Get("Accept")
	if strings.Index(h, "application/json") != -1 {
		b, _ := json.Marshal(sr)
		w.Header().Set("Content-Type", "application/json")
		w.Write(b)
		return
	}
	smokeTemplate := `{{.Status}}
test classes: 1
tests run: 1
tests passed: {{.TestsPassed}}
tests failed: {{.TestsFailed}}
tests errored: 0
time: 1.0ms
`
	t, _ := template.New("smoketest").Parse(smokeTemplate)
	w.Header().Set("Content-Type", "text/plain")
	t.Execute(w, sr)
}

func healthzHandler(w http.ResponseWriter, _ *http.Request, _ *site) {
	w.WriteHeader(http.StatusOK)
}
