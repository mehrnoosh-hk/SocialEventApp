package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-chi/chi/v5"
	"log"
	"math/rand"
	"net/http"
)

func main() {
	r := chi.NewRouter()
	r.Route("/events", func(r chi.Router) {
		r.Get("/", ListEvents)
		r.Post("/", CreateEvent)
		r.Route("/{eventID}", func(r chi.Router) {
			r.Get("/", GetEvent)
			r.Put("/", UpdateEvent)
			r.Delete("/", DeleteEvent)
		})
	})
	err := http.ListenAndServe(":3333", r)
	if err != nil {
		log.Fatal(err)
	}
}

func ListEvents(w http.ResponseWriter, r *http.Request) {
	marshalled, _ := json.Marshal(eventsDB)
	w.Write(marshalled)
}

func CreateEvent(w http.ResponseWriter, r *http.Request) {
	var e Event
	err := json.NewDecoder(r.Body).Decode(&e)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	id, _ := dbNewEvent(&e)
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(id))
}

func GetEvent(w http.ResponseWriter, r *http.Request) {
	if eventID := chi.URLParam(r, "eventID"); eventID != "" {
		event, err := dbGetEventByID(eventID)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("No event with this ID exist"))
		}
		marshaled, _ := json.Marshal(event)
		w.Write([]byte(marshaled))
	}
}

func UpdateEvent(w http.ResponseWriter, r *http.Request) {
	// Remember it should not be a pointer
	var eventUpdate Event
	err := json.NewDecoder(r.Body).Decode(&eventUpdate)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("error occurred during converting json"))
		return
	}
	if eventID := chi.URLParam(r, "eventID"); eventID != "" {
		event, err := dbUpdateEventByID(eventID, &eventUpdate)
		if err != nil {
			w.Write([]byte(err.Error()))
			return
		}
		w.WriteHeader(http.StatusCreated)
		eventMarshaled, _ := json.Marshal(event)
		w.Write(eventMarshaled)
	}
}

func DeleteEvent(w http.ResponseWriter, r *http.Request) {}

//--
// Data model objects and persistence mocks:
//--

type Event struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	PublicEvent bool   `json:"public_event"`
}

var eventsDB = []*Event{
	{ID: "01", Title: "BirthDay", PublicEvent: false},
	{ID: "02", Title: "NewYear", PublicEvent: true},
	{ID: "03", Title: "MovieNight", PublicEvent: true},
}

func dbNewEvent(event *Event) (string, error) {
	event.ID = fmt.Sprintf("%d", rand.Intn(100)+3)
	eventsDB = append(eventsDB, event)
	return event.ID, nil
}

// dbGetEventByID retrieve an event from database based on a given id and returns error if there is
// no event with such id exists
func dbGetEventByID(id string) (*Event, error) {
	for _, event := range eventsDB {
		if event.ID == id {
			return event, nil
		}
	}
	return nil, errors.New("event does not found")
}

func dbRemoveEventByID(id string) (*Event, error) {
	for i, event := range eventsDB {
		if event.ID == id {
			eventsDB = append(eventsDB[:i], eventsDB[i+1:]...)
			return event, nil
		}
	}
	return nil, errors.New("there is no event with such id")
}

func dbUpdateEventByID(id string, eventUpdate *Event) (*Event, error) {
	for _, event := range eventsDB {
		if event.ID == id {
			event.Title = eventUpdate.Title
			event.PublicEvent = eventUpdate.PublicEvent
			return event, nil
		}
	}
	return nil, errors.New("there is no event with such id")
}
