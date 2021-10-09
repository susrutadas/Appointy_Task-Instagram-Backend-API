package main

import (
	"encoding/json"
	"net/http"
	"regexp"
	"sync"
)

var (
	listUserRe   = regexp.MustCompile(`^\/users[\/]*$`)   //List All Users
	getUserRe    = regexp.MustCompile(`^\/users\/(\d+)$`) //Get a user using id
	createUserRe = regexp.MustCompile(`^\/users[\/]*$`)   //Create an User
	listPostRe   = regexp.MustCompile(`^\/posts[\/]*$`)   //List all posts of a user
	getPostRe    = regexp.MustCompile(`^\/posts\/(\d+)$`) //Get a post using id
	createPostRe = regexp.MustCompile(`^\/posts[\/]*$`)   //Create a post
)

//Attributes of Users
type user struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

//Attributes of Posts
type post struct {
	ID              string `json:"id"`
	Caption         string `json:"caption"`
	ImageURL        string `json:"image"`
	PostedTimestamp string `json:"timestamp"`
}

//Storing User data
type datastore struct {
	m map[string]user
	*sync.RWMutex
}

//Storing Post data
type datastore1 struct {
	m1 map[string]post
	*sync.RWMutex
}

//Handler for users
type userHandler struct {
	store *datastore
}

//Handler for posts
type postHandler struct {
	store *datastore1
}

//Endpoints for User
func (h *userHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	switch {
	case r.Method == http.MethodGet && listUserRe.MatchString(r.URL.Path):
		h.List(w, r)
		return
	case r.Method == http.MethodGet && getUserRe.MatchString(r.URL.Path):
		h.Get(w, r)
		return
	case r.Method == http.MethodPost && createUserRe.MatchString(r.URL.Path):
		h.Create(w, r)
		return
	default:
		notFound(w, r)
		return
	}
}

// Endpoints for Post
func (h *postHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	switch {
	case r.Method == http.MethodGet && listPostRe.MatchString(r.URL.Path):
		h.List(w, r)
		return
	case r.Method == http.MethodGet && getPostRe.MatchString(r.URL.Path):
		h.Get(w, r)
		return
	case r.Method == http.MethodPost && createPostRe.MatchString(r.URL.Path):
		h.Create(w, r)
		return
	default:
		notFound(w, r)
		return
	}
}

//List all Users
func (h *userHandler) List(w http.ResponseWriter, r *http.Request) {
	h.store.RLock()
	users := make([]user, 0, len(h.store.m))
	for _, v := range h.store.m {
		users = append(users, v)
	}
	h.store.RUnlock()
	jsonBytes, err := json.Marshal(users)
	if err != nil {
		internalServerError(w, r)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes)
}

//List all Posts
func (h *postHandler) List(w http.ResponseWriter, r *http.Request) {
	h.store.RLock()
	posts := make([]post, 0, len(h.store.m1))
	for _, k := range h.store.m1 {
		posts = append(posts, k)
	}
	h.store.RUnlock()
	jsonBytes, err := json.Marshal(posts)
	if err != nil {
		internalServerError(w, r)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes)
}

//Get a user using id
func (h *userHandler) Get(w http.ResponseWriter, r *http.Request) {
	matches := getUserRe.FindStringSubmatch(r.URL.Path)
	if len(matches) < 2 {
		notFound(w, r)
		return
	}
	h.store.RLock()
	u, ok := h.store.m[matches[1]]
	h.store.RUnlock()
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("user not found"))
		return
	}
	jsonBytes, err := json.Marshal(u)
	if err != nil {
		internalServerError(w, r)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes)
}

//Get a post using id
func (h *postHandler) Get(w http.ResponseWriter, r *http.Request) {
	matches := getPostRe.FindStringSubmatch(r.URL.Path)
	if len(matches) < 2 {
		notFound(w, r)
		return
	}
	h.store.RLock()
	u, ok := h.store.m1[matches[1]]
	h.store.RUnlock()
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("post not found"))
		return
	}
	jsonBytes, err := json.Marshal(u)
	if err != nil {
		internalServerError(w, r)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes)
}

//Create an User
func (h *userHandler) Create(w http.ResponseWriter, r *http.Request) {
	var u user
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		internalServerError(w, r)
		return
	}
	h.store.Lock()
	h.store.m[u.ID] = u
	h.store.Unlock()
	jsonBytes, err := json.Marshal(u)
	if err != nil {
		internalServerError(w, r)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes)
}

//Create a Post
func (h *postHandler) Create(w http.ResponseWriter, r *http.Request) {
	var p post
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		internalServerError(w, r)
		return
	}
	h.store.Lock()
	h.store.m1[p.ID] = p
	h.store.Unlock()
	jsonBytes, err := json.Marshal(p)
	if err != nil {
		internalServerError(w, r)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes)
}

//Error return
func internalServerError(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte("internal server error"))
}

func notFound(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte("not found"))
}

//Main function
func main() {
	mux := http.NewServeMux()
	userH := &userHandler{
		store: &datastore{
			m: map[string]user{
				"1": user{ID: "1", Name: "Susruta", Email: "susrutadas@gmail.com", Password: "helloworld"},
			},
			RWMutex: &sync.RWMutex{},
		},
	}
	postH := &postHandler{
		store: &datastore1{
			m1: map[string]post{
				"1": post{ID: "1", Caption: "Appointy is truly Amazing!!", ImageURL: "https://www.google.com/images/Appointy.png", PostedTimestamp: "2020-10-09T13:49:51.141Z"},
			},
			RWMutex: &sync.RWMutex{},
		},
	}
	mux.Handle("/users", userH)
	mux.Handle("/users/", userH)
	mux.Handle("/posts", postH)
	mux.Handle("/posts/", postH)

	http.ListenAndServe("localhost:8080", mux)
}
