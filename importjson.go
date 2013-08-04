package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
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
