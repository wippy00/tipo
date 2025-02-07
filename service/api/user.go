package api

import (
	"encoding/json"
	"net/http"
	"strconv"

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
func (rt *_router) getUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var string_id string = ps.ByName("id")
	id, err := strconv.Atoi(string_id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user, err := rt.db.GetUser(int64(id))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("content-type", "application/json")
	userJSON, err := json.Marshal(user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	_, _ = w.Write(userJSON)
}

func (rt *_router) updateUserName(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var respUser User

	id, hasAut, err := checkAuth(r, rt)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	if !hasAut {
		http.Error(w, "not authorized to change user name", http.StatusUnauthorized)
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&respUser); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err = checkUserName(respUser.Name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	dbUser, err := rt.db.UpdateUserName(id, respUser.Name)
	if err != nil && err.Error() == "user already exists" {
		http.Error(w, err.Error(), http.StatusConflict)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("content-type", "application/json")
	userJSON, err := json.Marshal(NewUser(dbUser))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	_, _ = w.Write(userJSON)

}

func (rt *_router) updateUserPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var respUser User

	id, hasAut, err := checkAuth(r, rt)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	if !hasAut {
		http.Error(w, "not authorized to change user photo", http.StatusUnauthorized)
		return
	}

	photo_multipart, handler, err := r.FormFile("photo")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	photo, err := validateFile(photo_multipart, handler, err)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	respUser.Id = id
	respUser.Photo = photo

	dbUser, err := rt.db.UpdateUserPhoto(respUser.Id, respUser.Photo)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("content-type", "application/json")
	userJSON, err := json.Marshal(NewUser(dbUser))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	_, _ = w.Write(userJSON)

}
