package main

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

	_ "github.com/lib/pq"
	"github.com/nu7hatch/gouuid"
)

type persistence struct {
	Database *sql.DB
}

func slugify(s string) string {
	s = strings.Trim(s, " \t\n\r-")
	s = strings.Replace(s, " ", "-", -1)
	s = strings.ToLower(s)
	return s
}

func newPersistence(dbURL string) *persistence {
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal(err)
	}
	return &persistence{Database: db}
}

func (p *persistence) Close() {
	p.Database.Close()
}

func (p *persistence) CreateUser(username, password string) (*user, error) {
	var user user
	user.Username = username
	encpassword := user.SetPassword(password)

	tx, err := p.Database.Begin()
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	stmt, err := tx.Prepare("insert into users (username, password) values ($1, $2)")
	if err != nil {
		log.Fatal(err)
		tx.Rollback()
		return nil, err
	}
	defer stmt.Close()
	_, err = stmt.Exec(username, encpassword)
	tx.Commit()

	u, _ := p.GetUser(username)
	log.Println("created")
	return u, nil

}

func (p *persistence) GetUser(username string) (*user, error) {
	stmt, err := p.Database.Prepare("select password from users where username = $1")
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	defer stmt.Close()

	var password string

	err = stmt.QueryRow(username).Scan(&password)
	if err != nil {
		return nil, err
	}
	return &user{Username: username, Password: []byte(password)}, err
}

func (p *persistence) CreateFact(title, details, source_name, source_url, tags string, user *user) *Fact {
	log.Println("CreateFact()")
	tx, err := p.Database.Begin()
	if err != nil {
		log.Println(err)
		return nil
	}
	stmt, err := tx.Prepare(
		`INSERT into facts (id, title, details, source_name, source_url, owner)
     VALUES            ($1,   $2,    $3,    $4,          $5,         $6)`)
	if err != nil {
		log.Println(err)
		tx.Rollback()
		return nil
	}

	log.Println("prepared")
	u4, err := uuid.NewV4()
	if err != nil {
		fmt.Println("error:", err)
		tx.Rollback()
		return nil
	}
	_, err = stmt.Exec(u4.String(), title, details, source_name, source_url, user.Username)
	if err != nil {
		log.Println(err)
		tx.Rollback()
		return nil
	}

	for _, t := range strings.Split(tags, ",") {
		p.AddTagToFact(tx, u4.String(), t)
	}

	tx.Commit()
	f, _ := p.GetFactById(u4.String())
	return f
}

func (p *persistence) ImportFact(uuid, title, details, source_name, source_url, added, username string, tags []string) {
	tx, err := p.Database.Begin()
	if err != nil {
		log.Println(err)
		return
	}
	stmt, err := tx.Prepare(
		`INSERT into facts (id, title, details, source_name, source_url, owner, added)
     VALUES            ($1,   $2,    $3,    $4,          $5,         $6,    $7)`)
	if err != nil {
		log.Println(err)
		tx.Rollback()
		return
	}

	log.Println("prepared")
	log.Println(uuid, title, details, source_name, source_url, username)
	_, err = stmt.Exec(uuid, title, details, source_name, source_url, username, added)
	if err != nil {
		log.Println(err)
		tx.Rollback()
		return
	}

	for _, t := range tags {
		p.AddTagToFact(tx, uuid, t)
	}

	tx.Commit()
}

func (p *persistence) AddTagToFact(tx *sql.Tx, uuid string, tagName string) (*Tag, error) {
	log.Println("AddTagToFact", uuid, tagName)
	slug := slugify(tagName)
	tag, _ := p.GetOrCreateTag(tx, slug, tagName)

	row := tx.QueryRow(
		`select count(*) from fact_tags where fact_uuid = $1
     and tag_slug = $2`, uuid, tag.Slug)
	var cnt int
	err := row.Scan(&cnt)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	if cnt > 0 {
		log.Println("fact already has tag")
		// already on there
		return tag, nil
	}
	stmt, err := tx.Prepare(
		"insert into fact_tags (fact_uuid, tag_slug) values ($1, $2)")
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer stmt.Close()

	_, err = stmt.Exec(uuid, tag.Slug)
	if err != nil {
		log.Println(err)
		log.Println(uuid, tag.Slug)
		return nil, err
	}
	return tag, nil
}

func (p *persistence) GetOrCreateTag(tx *sql.Tx, slug, name string) (*Tag, error) {
	stmt, err := tx.Prepare("select tagname from tags where slug = $1")
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	defer stmt.Close()
	err = stmt.QueryRow(slug).Scan(&name)
	switch {
	case err == sql.ErrNoRows:
		log.Println("need to create the tag")
		return p.CreateTag(tx, slug, name)
	case err != nil:
		return nil, err
	default:
		log.Println("tag already exists")
		return &Tag{Name: name, Slug: slug}, nil
	}
}

func (p *persistence) CreateTag(tx *sql.Tx, slug, name string) (*Tag, error) {
	// here, we assume that the tag doesn't already exist
	stmt, err := tx.Prepare(
		"insert into tags (slug, tagname) values ($1, $2)")
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	defer stmt.Close()

	_, err = stmt.Exec(slug, name)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return &Tag{Name: name, Slug: slug}, nil
}

func (p *persistence) TotalFactsCount() int {
	row := p.Database.QueryRow(`select count(*) from facts`)
	var cnt int
	err := row.Scan(&cnt)
	if err != nil {
		log.Println(err)
		return 0
	}
	return cnt
}

func (p *persistence) ListFacts(offset, count int) ([]Fact, error) {
	log.Println("listing facts")
	facts := make([]Fact, 0)
	rows, err := p.Database.Query(
		`SELECT id
     FROM facts
     ORDER BY added DESC
     LIMIT $1 OFFSET $2
  `, count, offset)
	if err != nil {
		return facts, err
	}
	for rows.Next() {
		var uuid string
		err := rows.Scan(&uuid)
		if err != nil {
			log.Println(err)
			return facts, err
		}
		fact, err := p.GetFactById(uuid)
		if err != nil {
			log.Println(err)
			return facts, err
		}
		facts = append(facts, *fact)
	}
	err = rows.Err()
	return facts, err
}

func (p *persistence) ListFactsByTag(tag Tag) ([]Fact, error) {
	log.Println("listing facts by tag")
	facts := make([]Fact, 0)
	rows, err := p.Database.Query(
		`SELECT ft.fact_uuid
     FROM fact_tags ft, facts f
     WHERE tag_slug = $1
       AND f.id = ft.fact_uuid
     ORDER BY added DESC`, tag.Slug)
	if err != nil {
		return facts, err
	}
	for rows.Next() {
		var uuid string
		err := rows.Scan(&uuid)
		if err != nil {
			log.Println(err)
			return facts, err
		}
		fact, err := p.GetFactById(uuid)
		if err != nil {
			log.Println(err)
			return facts, err
		}
		facts = append(facts, *fact)
	}
	err = rows.Err()
	return facts, err
}

func (p *persistence) GetFactById(id string) (*Fact, error) {
	var fact Fact
	stmt, err := p.Database.Prepare("select title, details, source_name, source_url, owner, added from facts where id = $1")
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	defer stmt.Close()

	var owner string

	err = stmt.QueryRow(id).Scan(&fact.Title, &fact.Details, &fact.SourceName, &fact.SourceUrl, &owner, &fact.Added)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	user, _ := p.GetUser(owner)
	fact.User = *user
	fact.UUID = id

	tags := make([]Tag, 0)

	rows, err := p.Database.Query(
		`select t.tagname, t.slug from tags t, fact_tags ft
     where t.slug = ft.tag_slug
       and ft.fact_uuid = $1`, id)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	for rows.Next() {
		tag := new(Tag)
		err := rows.Scan(&tag.Name, &tag.Slug)
		if err != nil {
			log.Println(err)
			return &fact, err
		}
		tags = append(tags, *tag)
	}
	fact.Tags = tags

	return &fact, err
}

func (p *persistence) ListTags() ([]Tag, error) {
	tags := make([]Tag, 0)
	rows, err := p.Database.Query(
		`SELECT slug, tagname
     FROM tags
     ORDER BY slug ASC`)
	if err != nil {
		log.Println(err)
		return tags, err
	}
	for rows.Next() {
		tag := new(Tag)
		err := rows.Scan(&tag.Slug, &tag.Name)
		if err != nil {
			log.Println(err)
			return tags, err
		}
		tags = append(tags, *tag)
	}
	err = rows.Err()
	return tags, err
}
