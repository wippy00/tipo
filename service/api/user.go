package api

import (
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// var test_user_01 = User{Id: 1, Name: "MariOne"}
// var test_user_02 = User{Id: 2, Name: "Gianna"}
// var test_user_03 = User{Id: 3, Name: "Rino"}
// var test_user_04 = User{Id: 4, Name: "Francesco"}
// var test_user_05 = User{Id: 5, Name: "Maria"}
// var test_db = []User{test_user_01, test_user_02, test_user_03, test_user_04, test_user_05}

// func (rt *_router) logIn(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
// 	w.WriteHeader(http.StatusOK)
// 	_, _ = w.Write([]byte("Hello, World!"))
// }

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
