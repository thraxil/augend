package main

import (
	"html/template"
	"time"

	"github.com/russross/blackfriday"
)

type Fact struct {
	UUID       string ""
	Title      string ""
	Details    string ""
	SourceName string ""
	SourceUrl  string ""
	Added      time.Time
	User       user
	Tags       []Tag
}

func (f *Fact) Url() string {
	return "/fact/" + f.UUID + "/"
}

func (f Fact) RenderAdded() string {
	return f.Added.Format(time.RFC3339)
}

func (f Fact) HasTags() bool {
	return len(f.Tags) > 0
}

func (f Fact) ListTags() []Tag {
	return f.Tags
}

func (f Fact) ListTagStrings() []string {
	tl := make([]string, len(f.Tags))
	for i, t := range f.Tags {
		tl[i] = t.Name
	}
	return tl
}

func (f Fact) HasSourceName() bool {
	return f.SourceName != ""
}

func (f Fact) HasSourceUrl() bool {
	return f.SourceUrl != ""
}

func (f Fact) RenderDetails() template.HTML {
	return template.HTML(string(blackfriday.MarkdownCommon([]byte(f.Details))))
}

// func ImportFact(title, details, source_name, source_url, added, username string,
// 	tags []string) *Fact {

// 	var user User
// 	riak.LoadModel(username, &user)

// 	var fact Fact
// 	u4, err := uuid.NewV4()
// 	if err != nil {
// 		fmt.Println("error:", err)
// 		return nil
// 	}

// 	err = riak.NewModel(u4.String(), &fact)
// 	fact.Title = title
// 	fact.Details = details
// 	fact.SourceName = source_name
// 	fact.SourceUrl = source_url
// 	fact.User.Set(&user)
// 	fact.Added = added

// 	fact.SaveAs(u4.String())

// 	index := getOrCreateFactIndex()
// 	if index == nil {
// 		fmt.Println("unable to get/create fact index")
// 		return nil
// 	}
// 	index.Facts.Add(&fact)
// 	index.SaveAs("fact-index")
// 	for _, t := range tags {
// 		fact.AddTag(t)
// 	}
// 	return &fact

// }
