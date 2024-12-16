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

	auth_id, hasAut, err := checkAuth(r, rt)
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

	dbConversation, err := rt.db.UpdateConversationName(id, auth_id, respConversation.Name)
	if err != nil && err.Error() == "auth user is not in this group" {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	if err != nil && err.Error() == "conversation is not a group" {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
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

	auth_id, hasAut, err := checkAuth(r, rt)
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

	dbConversation, err := rt.db.UpdateConversationPhoto(respConversation.Id, auth_id, respConversation.Photo)
	if err != nil && err.Error() == "auth user is not in this group" {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	if err != nil && err.Error() == "conversation is not a group" {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
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

func (rt *_router) addUserToConversation(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	auth_id, hasAut, err := checkAuth(r, rt)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	if !hasAut {
		http.Error(w, "not authorized to add user to conversation", http.StatusUnauthorized)
		return
	}

	var string_conversation_id string = ps.ByName("conversation_id")
	conversation_id, err := strconv.ParseInt(string_conversation_id, 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var string_user_id string = ps.ByName("user_id")
	user_id, err := strconv.ParseInt(string_user_id, 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	dbConversation, err := rt.db.AddUserToConversation(conversation_id, auth_id, user_id)
	if err != nil && err.Error() == "auth user is not in this group" {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	if err != nil && err.Error() == "conversation is not a group" {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err != nil && err.Error() == "user is already in conversation" {
		http.Error(w, err.Error(), http.StatusConflict)
		return
	}
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

func (rt *_router) removeUserFromConversation(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	auth_id, hasAut, err := checkAuth(r, rt)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	if !hasAut {
		http.Error(w, "not authorized to remove user from this conversation", http.StatusUnauthorized)
		return
	}

	var string_conversation_id string = ps.ByName("conversation_id")
	conversation_id, err := strconv.ParseInt(string_conversation_id, 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = rt.db.RemoveUserFromConversation(conversation_id, auth_id, auth_id)
	if err != nil && err.Error() == "auth user is not in this conversations" {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	if err != nil && err.Error() == "conversation is not a group" {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err != nil && err.Error() == "user is already in conversation" {
		http.Error(w, err.Error(), http.StatusConflict)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
