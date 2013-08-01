package main

import (
	"flag"
	"fmt"
	"github.com/gorilla/sessions"
	"github.com/stvp/go-toml-config"
	"github.com/tpjg/goriakpbc"
	"net/http"
)

var store sessions.Store
var template_dir = "templates"

func main() {
	var configFile string
	flag.StringVar(&configFile, "config", "./dev.conf", "TOML config file")
	flag.Parse()
	var (
		riak_host  = config.String("riak_host", "")
		port       = config.String("port", "9999")
		media_dir  = config.String("media_dir", "media")
		secret_key = config.String("secret_key", "change-me")
		t_dir      = config.String("template_dir", "templates")
	)
	config.Parse(configFile)
	template_dir = *t_dir

	err := riak.ConnectClient(*riak_host)
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	store = sessions.NewCookieStore([]byte(*secret_key))

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
	http.HandleFunc("/logout/", logoutHandler)
	http.Handle("/media/", http.StripPrefix("/media/",
		http.FileServer(http.Dir(*media_dir))))
	http.ListenAndServe(":"+*port, nil)
}
