package api

import (
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// GetUsers returns all users
func (rt *_router) getUsers(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	users, err := rt.db.GetUsers()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("content-type", "application/json")
	usersJSON, err := json.Marshal(users)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	_, _ = w.Write(usersJSON)
}

// GetUser by id returns a single user
// func (rt *_router) getUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
// 	var string_id string = ps.ByName("id")
// 	id, err := strconv.Atoi(string_id)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusBadRequest)
// 		return
// 	}

// 	user, err := rt.db.GetUser(id)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}

// 	w.Header().Set("content-type", "application/json")
// 	userJSON, err := json.Marshal(user)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}
// 	_, _ = w.Write(userJSON)
// }

// // AddUser adds a new user
// func (rt *_router) addUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
// 	var user User
// 	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
// 		http.Error(w, err.Error(), http.StatusBadRequest)
// 		return
// 	}

// 	if err := rt.db.AddUser(NewUser(user)); err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}

// 	w.WriteHeader(http.StatusCreated)
// }
