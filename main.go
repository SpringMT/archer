// main
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/gorilla/mux"
	elastigo "github.com/mattbaird/elastigo/lib"
	"html/template"
	"log"
	"net/http"
	"os"
)

var (
	host *string = flag.String("host", "localhost", "Elasticsearch Host")
)

var es *elastigo.Conn

type Message struct {
	Token       string  `json:"token"`
	Timestamp   float64 `json:"timestamp"`
	ChannelName string  `json:"channel_name"`
	UserName    string  `json:"user_name"`
	Text        string  `json:"text"`
}

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("views/index.html")
	if err != nil {
		log.Fatal(err)
	}
	t.Execute(w, nil)
}

func ShowChannelHandler(w http.ResponseWriter, r *http.Request) {
	channel := mux.Vars(r)["channel"]
	es := elastigo.NewConn()
	es.Domain = *host
	res, err := es.SearchUri("slack", "messages", map[string]interface{}{"q": "channel_name:" + channel})

	if err != nil {
		log.Println("ES Search Error: ", err)
		w.WriteHeader(http.StatusNotAcceptable)
		return
	}
	if res.Hits.Len() != 0 {
		s := make([]string, res.Hits.Len())
		for _, v := range res.Hits.Hits {
			str, _ := v.Source.MarshalJSON()
			s = append(s, string(str))
		}
		t, err := template.New("channel").Parse("<html><body>{{. | printf `%q` }}</body></html>\n")
		if err != nil {
			log.Fatal(err)
		}
		err = t.ExecuteTemplate(w, "channel", s)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		log.Println("ES NotFound")
		w.WriteHeader(http.StatusNotFound)
		return
	}
}

func PostMessageHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	if r.Method == "POST" {
		var message Message
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&message)
		if err != nil {
			log.Println("JSON Parse Error: ", err)
			w.WriteHeader(http.StatusNotAcceptable)
			return
		}
		fmt.Println(message)
		var res elastigo.BaseResponse
		es := elastigo.NewConn()
		es.Domain = *host
		res, err = es.Index("slack", "messages", "", nil, message)
		fmt.Println(res)
		if err != nil {
			log.Println("ES Save Error: ", err)
			w.WriteHeader(http.StatusNotAcceptable)
			return
		}
		w.WriteHeader(http.StatusCreated)
	}
}

func main() {
	rtr := mux.NewRouter()
	rtr.HandleFunc("/channel/{channel}", ShowChannelHandler)
	rtr.HandleFunc("/message", PostMessageHandler)
	rtr.HandleFunc("/", IndexHandler)
	http.Handle("/", rtr)
	http.ListenAndServe(":8080", nil)
}
