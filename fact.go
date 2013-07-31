package main

import (
	"fmt"
	"github.com/nu7hatch/gouuid"
	"github.com/tpjg/goriakpbc"
	"strings"
	"time"
)

type Fact struct {
	Title      string ""
	Details    string ""
	SourceName string ""
	SourceUrl  string ""
	Added      string ""
	Tags       riak.Many
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

func (f Fact) HasTags() bool {
	return f.Tags.Len() > 0
}

func (f Fact) ListTags() []Tag {
	tl := make([]Tag, f.Tags.Len())
	for i, t := range f.Tags {
		var ltag Tag
		t.Get(&ltag)
		tl[i] = ltag
	}
	return tl
}

func (f Fact) HasSource() bool {
	return f.SourceName != ""
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
