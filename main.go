package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
)

var (
	key   = []byte("super-secret-key")
	store = sessions.NewCookieStore(key)
)

func secret(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "cookie-name")

	// check if the user is authenticated
	if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}

	fmt.Fprintln(w, "the cake is a lie")

}

func login(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "cookie-name")
	var body struct {
		Email    string
		Password string
	}

	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	fmt.Printf("user  %v", body)

	if body.Password == "home" {
		session.Values["authenticated"] = true
		session.Save(r, w)
	}
}

func logout(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "cookie-name")

	session.Values["authenticated"] = false
	session.Save(r, w)
}

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/l", login)
	router.HandleFunc("/s", secret)
	router.HandleFunc("/lo", logout)

	http.ListenAndServe(":3000", router)
}
