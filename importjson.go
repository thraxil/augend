package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

type factdump struct {
	Key        string   `json:"key"`
	Title      string   `json:"title"`
	Details    string   `json:"details"`
	SourceName string   `json:"source_name"`
	SourceUrl  string   `json:"source_url"`
	Added      string   `json:"added"`
	User       string   `json:"user"`
	Tags       []string `json:"tags"`
}

func importJsonFile(filename string, s *site) {
	fmt.Println(filename)
	data, _ := ioutil.ReadFile(filename)
	var facts []factdump
	json.Unmarshal(data, &facts)
	for _, f := range facts {
		fmt.Println(f.Title)
		s.ImportFact(f.Key, f.Title, f.Details, f.SourceName, f.SourceUrl, f.Added,
			"anders", f.Tags)
	}
}

type KeyJsonFile struct {
	Keys []string
}

func dumpJSON(filename string) {
	// index := getOrCreateFactIndex()
	// facts := make([]factdump, 0)

	// for _, f := range index.Facts {
	// 	var lfact Fact
	// 	f.Get(&lfact)
	// 	facts = append(facts, factdump{
	// 		lfact.Key(),
	// 		lfact.Title,
	// 		lfact.Details,
	// 		lfact.SourceName,
	// 		lfact.SourceUrl,
	// 		lfact.Added,
	// 		lfact.GetUser().Username,
	// 		lfact.ListTagStrings(),
	// 	})
	// }
	// output, _ := json.Marshal(facts)

	// err := ioutil.WriteFile(filename, output, 0644)
	// if err != nil {
	// 	log.Println("could not write output")
	// } else {
	// 	log.Println("done")
	// }
}
