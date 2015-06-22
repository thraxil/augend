package main

import "strings"

type Tag struct {
	Name string
	Slug string
}

func (t Tag) Url() string {
	return "/tag/" + t.Slug + "/"
}

// func (t Tag) ListFacts() []Fact {
// 	fl := make([]Fact, t.Facts.Len())
// 	for i, f := range t.Facts {
// 		var lfact Fact
// 		f.Get(&lfact)
// 		fl[i] = lfact
// 	}
// 	return fl
// }

func normalizeTag(t string) string {
	t = strings.Trim(t, " \n\t,-")
	t = strings.ToLower(t)
	return t
}

type TagList []Tag

func (p TagList) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
func (p TagList) Len() int           { return len(p) }
func (p TagList) Less(i, j int) bool { return p[i].Name < p[j].Name }
