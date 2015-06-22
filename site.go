package main

import "github.com/gorilla/sessions"

type site struct {
	p     *persistence
	Store sessions.Store
}

func newSite(p *persistence, store sessions.Store) *site {
	s := site{
		p:     p,
		Store: store,
	}
	return &s
}

func (s *site) CreateUser(username, password string) (*user, error) {
	return s.p.CreateUser(username, password)
}

func (s *site) GetUser(username string) (*user, error) {
	return s.p.GetUser(username)
}

func (s *site) CreateFact(title, details, source_name, source_url, tags string, user *user) *Fact {
	return s.p.CreateFact(title, details, source_name, source_url, tags, user)
}

func (s *site) ListFacts(offset, count int) ([]Fact, error) {
	return s.p.ListFacts(offset, count)
}

func (s *site) TotalFactsCount() int {
	return s.p.TotalFactsCount()
}

func (s *site) GetFactById(id string) (*Fact, error) {
	return s.p.GetFactById(id)
}

func (s *site) ListTags() ([]Tag, error) {
	return s.p.ListTags()
}

func (s *site) GetTagBySlug(slug string) Tag {
	return Tag{Slug: slug, Name: slug}
}

func (s *site) ListFactsByTag(tag Tag) ([]Fact, error) {
	return s.p.ListFactsByTag(tag)
}

func (s *site) ImportFact(uuid, title, details, source_name, source_url, added, username string, tags []string) {
	s.p.ImportFact(uuid, title, details, source_name, source_url, added, username, tags)
}
