package main

import (
	"fmt"
	"github.com/gorilla/sessions"
	"github.com/tpjg/goriakpbc"
	"net/http"
)

var store sessions.Store

func main() {
	var riak_host = "localhost:10017"
	err := riak.ConnectClient(riak_host)
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	store = sessions.NewCookieStore([]byte("secret-key-change-me-in-prod"))

	err = ensureBuckets()
	if err != nil {
		fmt.Println("problem creating buckets. can't start")
		return
	}
	//	fmt.Println(index.Facts.Len())
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/fact/", factHandler)
	http.HandleFunc("/tag/", tagHandler)
	http.HandleFunc("/add/", addHandler)
	http.HandleFunc("/register/", registerHandler)
	http.HandleFunc("/login/", loginHandler)
	http.Handle("/media/", http.StripPrefix("/media/",
		http.FileServer(http.Dir("media"))))
	http.ListenAndServe(":9999", nil)
}
