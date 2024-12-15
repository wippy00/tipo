package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

func (rt *_router) getConversation(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	_, hasAut, err := checkAuth(r, rt)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	if !hasAut {
		http.Error(w, "not authorized to view conversation", http.StatusUnauthorized)
		return
	}

	var string_id string = ps.ByName("id")
	id, err := strconv.Atoi(string_id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	conversations, err := rt.db.GetConversation(int64(id))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("content-type", "application/json")
	conversationsJSON, err := json.Marshal(conversations)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	_, _ = w.Write(conversationsJSON)
}

func (rt *_router) getConversationOfUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	userId, hasAut, err := checkAuth(r, rt)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	if !hasAut {
		http.Error(w, "not authorized to view conversation", http.StatusUnauthorized)
		return
	}

	conversations, err := rt.db.GetConversationsOfUser(int64(userId))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("content-type", "application/json")
	conversationsJSON, err := json.Marshal(conversations)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	_, _ = w.Write(conversationsJSON)
}

func (rt *_router) updateConversationName(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var respConversation Conversation

	_, hasAut, err := checkAuth(r, rt)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	if !hasAut {
		http.Error(w, "not authorized to change conversation name", http.StatusUnauthorized)
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&respConversation); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var string_id string = ps.ByName("id")
	id, err := strconv.ParseInt(string_id, 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	dbConversation, err := rt.db.UpdateConversationName(id, respConversation.Name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("content-type", "application/json")
	userJSON, err := json.Marshal(NewConversation(dbConversation))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	_, _ = w.Write(userJSON)

}

func (rt *_router) updateConversationPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var respConversation User

	_, hasAut, err := checkAuth(r, rt)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	if !hasAut {
		http.Error(w, "not authorized to change user photo", http.StatusUnauthorized)
		return
	}

	var string_id string = ps.ByName("id")
	id, err := strconv.ParseInt(string_id, 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
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
	respConversation.Id = id
	respConversation.Photo = photo

	dbConversation, err := rt.db.UpdateConversationPhoto(respConversation.Id, respConversation.Photo)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("content-type", "application/json")
	userJSON, err := json.Marshal(NewConversation(dbConversation))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	_, _ = w.Write(userJSON)

}
