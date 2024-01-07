package main

import (
	"encoding/json"
	"fmt"
	"github.com/gosimple/slug"
	"go_rest_api/pkg/players"
	"net/http"
	"regexp"
)

var (
	PlayerRe       = regexp.MustCompile(`^/players/*$`)
	PlayerReWithID = regexp.MustCompile(`^/players/([a-z0-9]+(?:-[a-z0-9]+)+)$`)
)

type playerStore interface {
	Add(kitNumber string, player players.Player) error
	Get(kitNumber string) (players.Player, error)
	Update(kitNumber string, player players.Player) error
	List() (map[string]players.Player, error)
	Remove(kitNumber string) error
}

func main() {
	// Create the store and player handler
	store := players.NewMemStore()
	playersHandler := NewPlayersHandler(store)

	// Create a new request multiplexer
	mux := http.NewServeMux()

	// Registering routes and handlers
	mux.Handle("/", &homeHandler{})
	mux.Handle("/players", playersHandler)
	mux.Handle("/players/", playersHandler)

	// Run the server
	http.ListenAndServe(":8080", mux)
}

func InternalServerErrorHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte("500 Internal Server Error"))
}

func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte("404 Not Found`"))
}

type PlayersHandler struct {
	store playerStore
}

func NewPlayersHandler(s playerStore) *PlayersHandler {
	return &PlayersHandler{
		store: s,
	}
}

func (h *PlayersHandler) CreatePlayer(w http.ResponseWriter, r *http.Request) {
	// Player object will be populated from JSON payload
	var player players.Player
	if err := json.NewDecoder(r.Body).Decode(&player); err != nil {
		InternalServerErrorHandler(w, r)
		return
	}
	// Create URL from player kit number
	resourceID := slug.Make(player.FirstName + " " + player.LastName)
	if err := h.store.Add(resourceID, player); err != nil {
		InternalServerErrorHandler(w, r)
		return
	}

	// Set status code to 200
	w.WriteHeader(http.StatusOK)
}
func (h *PlayersHandler) ListPlayers(w http.ResponseWriter, r *http.Request) {
	resources, err := h.store.List()

	jsonBytes, err := json.Marshal(resources)
	if err != nil {
		InternalServerErrorHandler(w, r)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes)
}
func (h *PlayersHandler) GetPlayer(w http.ResponseWriter, r *http.Request) {
	// Extract the resource ID/slug of a recipe
	matches := PlayerReWithID.FindStringSubmatch(r.URL.Path)
	//// Expect matches to be length >= 2 (full string + 1 matching group)
	if len(matches) < 2 {
		InternalServerErrorHandler(w, r)
		return
	}

	//Retrieve player from the store
	player, err := h.store.Get(matches[1])
	if err != nil {
		// Special case of NotFound error
		if err == players.NotFoundErr {
			NotFoundHandler(w, r)
			return
		}

		// Every other error
		InternalServerErrorHandler(w, r)
		return
	}

	// Convert the struct into JSON payload
	jsonBytes, err := json.Marshal(player)
	if err != nil {
		InternalServerErrorHandler(w, r)
		return
	}

	// Write the results
	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes)
}
func (h *PlayersHandler) UpdatePlayer(w http.ResponseWriter, r *http.Request) {
	// Extract the resource ID/slug of a recipe
	matches := PlayerReWithID.FindStringSubmatch(r.URL.Path)
	//// Expect matches to be length >= 2 (full string + 1 matching group)
	if len(matches) < 2 {
		InternalServerErrorHandler(w, r)
		return
	}
	// Update player from the store
	var player players.Player
	if err := json.NewDecoder(r.Body).Decode(&player); err != nil {
		InternalServerErrorHandler(w, r)
		return
	}

	if err := h.store.Update(matches[1], player); err != nil {
		if err == players.NotFoundErr {
			NotFoundHandler(w, r)
			return
		}
		InternalServerErrorHandler(w, r)
		return
	}
	w.WriteHeader(http.StatusOK)
}
func (h *PlayersHandler) DeletePlayer(w http.ResponseWriter, r *http.Request) {
	// Extract the resource ID/slug of a recipe
	matches := PlayerReWithID.FindStringSubmatch(r.URL.Path)
	//// Expect matches to be length >= 2 (full string + 1 matching group)
	if len(matches) < 2 {
		InternalServerErrorHandler(w, r)
		return
	}

	if err := h.store.Remove(matches[1]); err != nil {
		InternalServerErrorHandler(w, r)
		return
	}
	w.WriteHeader(http.StatusOK)
}

type homeHandler struct{}

func (h *homeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("This is my home page"))
}

func (h *PlayersHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch {
	case r.Method == http.MethodPost && PlayerRe.MatchString(r.URL.Path):
		h.CreatePlayer(w, r)
		return
	case r.Method == http.MethodGet && PlayerRe.MatchString(r.URL.Path):
		h.ListPlayers(w, r)
		return
	case r.Method == http.MethodGet && PlayerReWithID.MatchString(r.URL.Path):
		fmt.Println("Inside Get Player route")
		h.GetPlayer(w, r)
		return
	case r.Method == http.MethodPut && PlayerReWithID.MatchString(r.URL.Path):
		h.UpdatePlayer(w, r)
	case r.Method == http.MethodDelete && PlayerReWithID.MatchString(r.URL.Path):
		h.DeletePlayer(w, r)
		return
	default:
		return
	}
}
