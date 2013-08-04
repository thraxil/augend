package main

import (
	"fmt"
	"github.com/tpjg/goriakpbc"
	"strings"
)

type Tag struct {
	Name       string
	Facts      riak.Many
	riak.Model `riak:"augend.tag"`
}

func (t *Tag) Resolve(count int) (err error) {
	return nil
}

func (t Tag) Url() string {
	return "/tag/" + t.Name + "/"
}

func (t Tag) ListFacts() []Fact {
	fl := make([]Fact, t.Facts.Len())
	for i, f := range t.Facts {
		var lfact Fact
		f.Get(&lfact)
		fl[i] = lfact
	}
	return fl
}

func normalizeTag(t string) string {
	t = strings.Trim(t, " \n\t,-")
	t = strings.ToLower(t)
	return t
}

func getOrCreateTag(t string) *Tag {
	t = normalizeTag(t)
	if t == "" {
		return nil
	}
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
