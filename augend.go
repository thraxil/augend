package main

import (
	"flag"
	"fmt"
	"github.com/gorilla/sessions"
	"github.com/peterbourgon/g2s"
	"github.com/stvp/go-toml-config"
	"github.com/tpjg/goriakpbc"
	"net/http"
	"time"
)

var store sessions.Store
var template_dir = "templates"
var statsd g2s.Statter

func main() {
	var configFile string
	var importjson string
	var keyjson string
	flag.StringVar(&configFile, "config", "./dev.conf", "TOML config file")
	flag.StringVar(&importjson, "importjson", "", "json file to import")
	flag.StringVar(&keyjson, "keyjson", "", "json file with keys to repair index")
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

	if importjson != "" {
		fmt.Println("importing JSON file")
		importJsonFile(importjson)
	}
	if keyjson != "" {
		fmt.Println("importing Key JSON file and repairing index")
		repairIndex(keyjson)
	}
	statsd, _ = g2s.Dial("udp", "127.0.0.1:8125")
	//	fmt.Println(index.Facts.Len())
	http.HandleFunc("/favicon.ico", faviconHandler)
	http.HandleFunc("/", makeHandler(indexHandler))
	http.HandleFunc("/fact/", makeHandler(factHandler))
	http.HandleFunc("/tag/", makeHandler(tagHandler))
	http.HandleFunc("/add/", makeHandler(addHandler))
	http.HandleFunc("/register/", makeHandler(registerHandler))
	http.HandleFunc("/login/", makeHandler(loginHandler))
	http.HandleFunc("/logout/", makeHandler(logoutHandler))
	http.Handle("/media/", http.StripPrefix("/media/",
		http.FileServer(http.Dir(*media_dir))))
	http.ListenAndServe(":"+*port, nil)
}

func makeHandler(f func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		t0 := time.Now()
		f(w, r)
		t1 := time.Now()
		statsd.Counter(1.0, "augend.response.200", 1)
		statsd.Timing(1.0, "augend.view.GET", t1.Sub(t0))
	}
}
