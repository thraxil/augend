package main

import (
	"fmt"
	"github.com/tpjg/goriakpbc"
)

type Tag struct {
	Name       string
	Facts      riak.Many
	riak.Model `riak:"test.augend.tag"`
}

func (t *Tag) Resolve(count int) (err error) {
	return nil
}

func (t Tag) Url() string {
	return "/tag/" + t.Name + "/"
}

func (t Tag) ListFacts() []Fact {
	fl := make([]Fact, t.Facts.Len())
	fmt.Println(t.Facts.Len())
	for i, f := range t.Facts {
		var lfact Fact
		f.Get(&lfact)
		fl[i] = lfact
	}
	return fl
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
