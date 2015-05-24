package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
)

type JsonFile struct {
	Facts []struct {
		Title      string   `json:"title"`
		Details    string   `json:"details"`
		SourceName string   `json:"source_name"`
		SourceUrl  string   `json:"source_url"`
		Added      string   `json:"added"`
		Tags       []string `json:"tags"`
	} `json:"facts"`
}

func importJsonFile(filename string) {
	fmt.Println(filename)
	data, _ := ioutil.ReadFile(filename)
	var facts JsonFile
	json.Unmarshal(data, &facts)
	for _, f := range facts.Facts {
		fmt.Println(f.Title)
		ImportFact(f.Title, f.Details, f.SourceName, f.SourceUrl, f.Added,
			"anders", f.Tags)
	}
}

type KeyJsonFile struct {
	Keys []string
}

func repairIndex(filename string) {
	fmt.Println(filename)
	data, _ := ioutil.ReadFile(filename)
	var keys KeyJsonFile
	json.Unmarshal(data, &keys)
	for _, k := range keys.Keys {
		fmt.Println(k)
		ImportFactIndexOnly(k)
	}
}

func repairIndices() {
	index := getOrCreateTagIndex()
	n := index.Tags.Len()
	tags := make(TagList, n)

	facts := getOrCreateFactIndex()
	if facts == nil {
		fmt.Println("unable to get/create fact index")
		return
	}
	seen := make(map[string]bool)

	for i, t := range index.Tags {
		var ltag Tag
		t.Get(&ltag)
		tags[i] = ltag
		log.Println("TAG: ", ltag.Name)
		for _, f := range ltag.ListFacts() {
			log.Println("\tFACT: ", f.Title)
			_, ok := seen[f.Key()]
			if ok {
				log.Println("\talready have it")
			} else {
				facts.Facts.Add(&f)
				log.Println("\tadded")
				seen[f.Key()] = true
			}
		}
	}
	facts.SaveAs("fact-index")
	log.Println("done")
}
