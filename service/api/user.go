package api

import (
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type User struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
	// Photo image.RGBA `json:"photo"`
}

// var test_user_01 = User{Id: 1, Name: "MariOne"}
// var test_user_02 = User{Id: 2, Name: "Gianna"}
// var test_user_03 = User{Id: 3, Name: "Rino"}
// var test_user_04 = User{Id: 4, Name: "Francesco"}
// var test_user_05 = User{Id: 5, Name: "Maria"}
// var test_db = []User{test_user_01, test_user_02, test_user_03, test_user_04, test_user_05}

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
