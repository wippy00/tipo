package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

type LoginRequest struct {
	Name string `json:"name"`
}

func (rt *_router) logIn(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var loginReq LoginRequest

	err := json.NewDecoder(r.Body).Decode(&loginReq)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err = checkUserName(loginReq.Name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user, isNew, err := rt.db.LoginUser(loginReq.Name)
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

	if isNew {
		w.WriteHeader(http.StatusCreated)
	} else {
		w.WriteHeader(http.StatusOK)
	}
	_, _ = w.Write(userJSON)

}

func checkAuth(r *http.Request, rt *_router) (int64, bool, error) {
	auth_id_str := r.Header.Get("authorization")
	if auth_id_str == "" {
		return int64(0), false, errors.New("missing auth header")
	}
	atoi, err := strconv.Atoi(auth_id_str)
	if err != nil {
		return int64(atoi), false, err
	}

	auth_id := int64(atoi)

	user, valid, err := rt.db.UserExistById(auth_id)
	if err != nil {
		return user.Id, false, err
	}

	if !valid {
		return user.Id, false, nil
	}

	return user.Id, true, nil
}
