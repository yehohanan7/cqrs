package cqrs

import (
	"encoding/json"
	"net/http"
	"time"

	"reflect"

	"fmt"

	"github.com/golang/glog"
	"github.com/gorilla/feeds"
	"github.com/gorilla/mux"
)

func generateFeed(store EventStore) string {
	feed := &feeds.Feed{
		Title:       "event feeds",
		Link:        &feeds.Link{Href: "http://localhost:8000/events"},
		Description: "All events on this service",
	}

	feed.Items = make([]*feeds.Item, 0)

	for _, event := range store.GetAllEvents() {
		feed.Items = append(feed.Items, &feeds.Item{
			Id:          event.Id,
			Title:       event.AggregateName,
			Link:        &feeds.Link{Href: fmt.Sprintf("/events/%s", event.Id)},
			Author:      &feeds.Author{Name: "name", Email: "email"},
			Description: reflect.TypeOf(event.Payload).String(),
			Created:     time.Now(),
		})
	}

	atom, err := feed.ToAtom()
	if err != nil {
		glog.Error("error while building atom feed", err)
	}
	return atom
}

func eventsHandler(store EventStore) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		w.Header().Set("Content-Type", "text/xml")
		w.Write([]byte(generateFeed(store)))
	}
}

func eventHandler(store EventStore) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		vars := mux.Vars(r)
		id := vars["id"]
		json.NewEncoder(w).Encode(store.GetEvent(id))
	}
}

func EventFeed(router *mux.Router, store EventStore) {
	router.HandleFunc("/events", eventsHandler(store)).Methods("GET")
	router.HandleFunc("/events/{id}", eventHandler(store)).Methods("GET")
}
