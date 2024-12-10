package api

import (
	"encoding/json"
	"net/http"

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
	user, err := rt.db.LoginUser(loginReq.Name)
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
