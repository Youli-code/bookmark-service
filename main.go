package main

import (
	"encoding/json"
	"net/http"
	"strconv"
	"sync"
)

type Bookmark struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
	URL   string `json:"url"`
}

var (
	bookmarks = make(map[int]Bookmark)
	idCounter = 1
	mu        sync.Mutex
)

func createBookmark(w http.ResponseWriter, r *http.Request) {
	var b Bookmark

	err := json.NewDecoder(r.Body).Decode(&b)
	if err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	mu.Lock()
	b.ID = idCounter
	bookmarks[idCounter] = b
	idCounter++
	mu.Unlock()

	json.NewEncoder(w).Encode(b)
}

func listBookmarks(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	defer mu.Unlock()

	list := []Bookmark{}

	for _, b := range bookmarks {
		list = append(list, b)
	}

	json.NewEncoder(w).Encode(list)
}

func deleteBookmark(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	mu.Lock()
	delete(bookmarks, id)
	mu.Unlock()

	w.WriteHeader(http.StatusNoContent)
}

func main() {
	http.HandleFunc("/bookmarks", listBookmarks)
	http.HandleFunc("/create", createBookmark)
	http.HandleFunc("/delete", deleteBookmark)

	http.ListenAndServe(":8080", nil)
}
