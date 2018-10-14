package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"contrib.go.opencensus.io/exporter/stackdriver"
	"go.opencensus.io/trace"
	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/memcache"
)

var spannerService *SpannerService

type Hoge struct {
	CreatedAt time.Time `json:"createdAt"`
}

func main() {
	log.Println("Start Main")

	exporter, err := stackdriver.NewExporter(stackdriver.Options{
		ProjectID: "sinmetal-go",
	})
	if err != nil {
		panic(err)
	}
	trace.RegisterExporter(exporter)

	ctx := context.Background()
	ctx, span := trace.StartSpan(ctx, "/main")
	defer span.End()

	sc, err := NewSpannerClient(ctx, "projects/gcpug-public-spanner/instances/merpay-sponsored-instance/databases/sinmetal")
	if err != nil {
		panic(err)
	}

	spannerService = NewSpannerService(sc)

	http.HandleFunc("/datastore", datastoreHandler)
	http.HandleFunc("/memcache", memcacheHandler)
	http.HandleFunc("/spanner", SpannerSimpleQueryHandler)

	appengine.Main()
}

const Kind = "HelloGAEGo111"

func datastoreHandler(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)

	o := Hoge{
		CreatedAt: time.Now(),
	}
	k, err := datastore.Put(ctx, datastore.NewKey(ctx, Kind, "", time.Now().Unix(), nil), &o)
	if err != nil {
		log.Fatalf("failed Datastore.put(), err=%+v", err)
	}
	log.Printf("%Key=+v\n", k)

	var list []Hoge
	q := datastore.NewQuery(Kind)
	q.GetAll(ctx, &list)
	log.Println(list)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(list); err != nil {
		log.Fatalln(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func memcacheHandler(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)

	if err := memcache.Add(ctx, &memcache.Item{
		Key:   "Hoge",
		Value: []byte("VALUE"),
	}); err != nil {
		msg := fmt.Sprintf("failed Memcache.Add(), err=%+v", err)
		log.Fatal(msg)
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}

	w.Write([]byte("DONE"))
}
