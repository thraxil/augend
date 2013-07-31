package main

import (
	"fmt"
	"github.com/tpjg/goriakpbc"
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
