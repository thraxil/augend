package main

import (
	"fmt"
	"github.com/tpjg/goriakpbc"
)

type FactIndex struct {
	Facts      riak.Many
	riak.Model `riak:"test.augend.index"`
}

func getOrCreateFactIndex() *FactIndex {
	var index FactIndex
	err := riak.LoadModel("fact-index", &index)
	if err != nil {
		fmt.Println("creating new fact index")
		return createFactIndex()
	}
	return &index
}

func createFactIndex() *FactIndex {
	var index FactIndex
	err := riak.NewModel("fact-index", &index)
	if err != nil {
		fmt.Println("could not create new fact index")
		fmt.Println(err)
		return nil
	}
	index.SaveAs("fact-index")
	return &index
}

type TagIndex struct {
	Tags       riak.Many
	riak.Model `riak:"test.augend.index"`
}

func getOrCreateTagIndex() *TagIndex {
	var index TagIndex
	err := riak.LoadModel("tag-index", &index)
	if err != nil {
		fmt.Println("creating new tag index")
		return createTagIndex()
	}
	return &index
}

func createTagIndex() *TagIndex {
	var nindex TagIndex
	err := riak.NewModel("tag-index", &nindex)
	if err != nil {
		fmt.Println("could not create tag index")
		fmt.Println(err)
		return nil
	}
	nindex.SaveAs("tag-index")
	return &nindex
}

func ensureBuckets() error {
	_, err := riak.NewBucket("test.augend.fact")
	if err != nil {
		fmt.Println("could not get/create fact bucket")
		return err
	}
	_, err = riak.NewBucket("test.augend.index")
	if err != nil {
		fmt.Println("could not get/create fact bucket")
		return err
	}
	return nil
}
